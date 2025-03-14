<template>
  <div id="device-list">
    <el-row gutter="20">
      <el-col :span="24">
        <h1>设备列表</h1>
        <div v-for="device in devices" :key="device" class="device-item">
          <el-card class="box-card">
            <router-link :to="`/devices/${device}`">{{ device }}</router-link>
          </el-card>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { reactive, toRefs } from "vue";

export default {
  setup() {
    const state = reactive({
      devices: [],
    });

    const fetchDevices = () => {
      fetch("http://localhost:8199/images")
        .then((response) => response.json())
        .then((data) => {
          state.devices = data.devices;
        })
        .catch((error) => {
          console.error("Error fetching devices:", error);
        });
    };

    fetchDevices();

    return {
      ...toRefs(state),
      fetchDevices,
    };
  },
};
</script>

<style>
#device-list {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}

.device-item {
  margin-bottom: 20px; /* 增加设备号之间的间距 */
}

.el-card {
  background-color: #ffffff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}
</style>
