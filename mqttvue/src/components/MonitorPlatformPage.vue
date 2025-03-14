<template>
  <div class="monitor-platform">
    <h2>实时监控平台</h2>
    <div class="image-container">
      <img :src="imageSrc" alt="实时图片" v-if="imageSrc" />
    </div>
  </div>
</template>

<script>
export default {
  name: "MonitorPlatformPage",
  data() {
    return {
      imageSrc: null,
      socket: null,
    };
  },
  methods: {
    connectWebSocket() {
      this.socket = new WebSocket("ws://localhost:8199/ws");
      this.socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        this.imageSrc = "data:image/jpeg;base64," + data.image;
      };
      this.socket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
      this.socket.onclose = () => {
        console.log("WebSocket connection closed");
      };
    },
  },
  mounted() {
    this.connectWebSocket();
  },
  beforeUnmount() {
    if (this.socket) {
      this.socket.close();
    }
  },
};
</script>

<style>
.monitor-platform {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.image-container {
  width: 80%;
  height: 80vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
  border: 1px solid #ddd;
}

.image-container img {
  max-width: 100%;
  max-height: 100%;
}
</style>
