<template>
  <div id="device-manager">
    <h1>设备管理</h1>
    <input v-model="newDevice" placeholder="输入设备名称" />
    <button @click="addDevice">添加设备</button>
    <ul>
      <li v-for="device in devices" :key="device">
        {{ device }}
        <button @click="deleteDevice(device)">删除</button>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  data() {
    return {
      devices: [],
      newDevice: "",
    };
  },
  created() {
    this.fetchDevices();
  },
  methods: {
    fetchDevices() {
      fetch("http://localhost:8199/devices")
        .then((response) => response.json())
        .then((data) => {
          this.devices = data.devices;
        })
        .catch((error) => {
          console.error("Error fetching devices:", error);
        });
    },
    addDevice() {
      if (this.newDevice.trim() === "") {
        alert("设备名称不能为空");
        return;
      }
      fetch(`http://localhost:8199/devices/add?device=${this.newDevice}`, {
        method: "POST",
      })
        .then((response) => response.json())
        .then((data) => {
          if (data.message) {
            alert(data.message);
            this.fetchDevices();
            this.newDevice = "";
          } else {
            alert("添加失败");
          }
        })
        .catch((error) => {
          console.error("Error adding device:", error);
        });
    },
    deleteDevice(device) {
      fetch(`http://localhost:8199/devices/delete?device=${device}`, {
        method: "POST",
      })
        .then((response) => response.json())
        .then((data) => {
          if (data.message) {
            alert(data.message);
            this.fetchDevices();
          } else {
            alert("删除失败");
          }
        })
        .catch((error) => {
          console.error("Error deleting device:", error);
        });
    },
  },
};
</script>

<style>
#device-manager {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
