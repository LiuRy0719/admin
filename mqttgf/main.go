package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"
)

var (
   broker     = "mqtt://112.6.224.25:20042"
   topic      = "image_topic/#"  // 使用通配符订阅所有设备的图片
   clientID   = "go_mqtt_client"
   username   = "your_username"
   password   = "your_password"
   jwtKey     = []byte("my_secret_key") // JWT 密钥
   devices    = []string{}  // 设备列表
   upgrader   = websocket.Upgrader{
       CheckOrigin: func(r *http.Request) bool {
           return true
       },
   }
   clients    = make(map[*websocket.Conn]bool)
   broadcast  = make(chan ImageMessage)
)

type ImageMessage struct {
   DeviceID string `json:"device_id"`
   Image    []byte `json:"image"`
}

type Claims struct {
   Username string `json:"username"`
   jwt.StandardClaims
}

func archiveAndClearImages(r *ghttp.Request) {
	deviceID := r.Get("device").String()
	deviceDir := filepath.Join("images", deviceID)
	zipFileName := fmt.Sprintf("%s_%d.zip", deviceID, time.Now().Unix())
	zipFilePath := filepath.Join("images", zipFileName)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		r.Response.WriteJson(g.Map{"error": err.Error()})
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(deviceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		zipEntry, err := zipWriter.Create(info.Name())
		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return err
		}

		return os.Remove(path)
	})
	if err != nil {
		r.Response.WriteJson(g.Map{"error": err.Error()})
		return
	}

	r.Response.Header().Set("Content-Disposition", "attachment; filename="+zipFileName)
	r.Response.Header().Set("Content-Type", "application/zip")
	r.Response.ServeFile(zipFilePath)

	// 删除设备目录
	os.Remove(deviceDir)
}
func connectHandler(client mqtt.Client) {
   fmt.Println("Connected to MQTT Broker")
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {
   topicParts := strings.Split(msg.Topic(), "/")
   if len(topicParts) < 2 {
       fmt.Println("Invalid topic format")
       return
   }
   deviceID := topicParts[1]

   deviceDir := filepath.Join("images", deviceID)
   err := os.MkdirAll(deviceDir, os.ModePerm)
   if err != nil {
       fmt.Println("Error creating directory:", err)
       return
   }

   fileName := fmt.Sprintf("%d.jpg", time.Now().Unix())
   filePath := filepath.Join(deviceDir, fileName)
   file, err := os.Create(filePath)
   if err != nil {
       fmt.Println("Error creating file:", err)
       return
   }
   defer file.Close()
   file.Write(msg.Payload())

   // 发送消息到 WebSocket 客户端
   broadcast <- ImageMessage{DeviceID: deviceID, Image: msg.Payload()}
}

func handleConnections(r *ghttp.Request) {
   ws, err := upgrader.Upgrade(r.Response.Writer, r.Request, nil)
   if err != nil {
       fmt.Println(err)
       return
   }
   defer ws.Close()

   clients[ws] = true

   for {
       var msg ImageMessage
       err := ws.ReadJSON(&msg)
       if err != nil {
           fmt.Println(err)
           delete(clients, ws)
           break
       }
   }
}

func handleMessages() {
   for {
       msg := <-broadcast
       for client := range clients {
           err := client.WriteJSON(msg)
           if err != nil {
               fmt.Println(err)
               client.Close()
               delete(clients, client)
           }
       }
   }
}

func loginHandler(r *ghttp.Request) {
   var creds struct {
       Username string `json:"username"`
       Password string `json:"password"`
   }

   if err := r.Parse(&creds); err != nil {
       r.Response.WriteHeader(http.StatusBadRequest)
       return
   }

   // 这里应该进行用户身份验证
   if creds.Username != "admin" || creds.Password != "password" {
       r.Response.WriteHeader(http.StatusUnauthorized)
       return
   }

   expirationTime := time.Now().Add(5 * time.Minute)
   claims := &Claims{
       Username: creds.Username,
       StandardClaims: jwt.StandardClaims{
           ExpiresAt: expirationTime.Unix(),
       },
   }

   token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
   tokenString, err := token.SignedString(jwtKey)
   if err != nil {
       r.Response.WriteHeader(http.StatusInternalServerError)
       return
   }

   r.Response.WriteJson(g.Map{"token": tokenString})
}

func main() {
   // 创建保存图像的目录
   os.Mkdir("images", os.ModePerm)

   // MQTT客户端配置
   opts := mqtt.NewClientOptions()
   opts.AddBroker(broker)
   opts.SetClientID(clientID)
   opts.SetUsername(username)
   opts.SetPassword(password)
   opts.SetDefaultPublishHandler(messageHandler)
   opts.OnConnect = connectHandler

   // 创建并连接MQTT客户端
   client := mqtt.NewClient(opts)
   if token := client.Connect(); token.Wait() && token.Error() != nil {
       panic(token.Error())
   }

   // 订阅主题
   if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
       panic(token.Error())
   }

   // 启动WebSocket处理协程
   go handleMessages()

   // 启动HTTP服务器
   s := g.Server()
   
   // 添加允许跨域的中间件
   s.Use(func(r *ghttp.Request) {
       r.Response.CORSDefault()
       r.Middleware.Next()
   })

   s.BindHandler("/login", loginHandler)
   s.BindHandler("/ws", handleConnections)

   s.BindHandler("/images", func(r *ghttp.Request) {
       devices, err := os.ReadDir("images")
       if err != nil {
           r.Response.WriteJson(g.Map{"error": err.Error()})
           return
       }

       var deviceList []string
       for _, device := range devices {
           if device.IsDir() {
               deviceList = append(deviceList, device.Name())
           }
       }
       r.Response.WriteJson(g.Map{"devices": deviceList})
   })

   s.BindHandler("/images/{device}", func(r *ghttp.Request) {
       device := r.Get("device").String()
       images, err := os.ReadDir(filepath.Join("images", device))
       if err != nil {
           r.Response.WriteJson(g.Map{"error": err.Error()})
           return
       }

       var imageList []string
       for _, img := range images {
           if !img.IsDir() {
               imageList = append(imageList, img.Name())
           }
       }
       r.Response.WriteJson(g.Map{"images": imageList})
   })
   
   s.BindHandler("/images/{device}/{name}", func(r *ghttp.Request) {
       device := r.Get("device").String()
       name := r.Get("name").String()
       r.Response.ServeFile(filepath.Join("images", device, name))
   })
   
   s.BindHandler("/delete/{device}/{name}", func(r *ghttp.Request) {
       device := r.Get("device").String()
       name := r.Get("name").String()
       err := os.Remove(filepath.Join("images", device, name))
       if err != nil {
           r.Response.WriteJson(g.Map{"error": err.Error()})
           return
       }
       r.Response.WriteJson(g.Map{"message": "Image deleted successfully"})
   })

   // 设备管理接口
   s.BindHandler("/devices", func(r *ghttp.Request) {
       r.Response.WriteJson(g.Map{"devices": devices})
   })
   
   s.BindHandler("/devices/add", func(r *ghttp.Request) {
       device := r.Get("device").String()
       devices = append(devices, device)
       r.Response.WriteJson(g.Map{"message": "Device added successfully"})
   })
   
   s.BindHandler("/devices/delete", func(r *ghttp.Request) {
       device := r.Get("device").String()
       for i, d := range devices {
           if d == device {
               devices = append(devices[:i], devices[i+1:]...)
               break
           }
       }
       r.Response.WriteJson(g.Map{"message": "Device deleted successfully"})
   })
   s.BindHandler("/archive/{device}", archiveAndClearImages)
   s.SetPort(8199)
   s.Run()
}