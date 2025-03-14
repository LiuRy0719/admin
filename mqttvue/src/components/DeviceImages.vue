<template>
  <div>
    <h2>{{ deviceId }} 的图像</h2>
    <el-button type="primary" @click="archiveImages">一键归档</el-button>
    <div class="image-list">
      <el-table :data="images" style="width: 100%">
        <el-table-column prop="name" label="图片名称" width="180">
          <template v-slot="scope">
            <el-button type="text" @click="viewImage(scope.row.name)">
              {{ scope.row.name }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <el-dialog v-model="dialogVisible" width="80%">
      <img :src="currentImageSrc" alt="图片" style="width: 100%" />
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: "DeviceImages",
  data() {
    return {
      images: [],
      deviceId: this.$route.params.deviceId,
      dialogVisible: false,
      currentImageSrc: "",
    };
  },
  methods: {
    fetchImages() {
      fetch(`http://localhost:8199/images/${this.deviceId}`)
        .then((response) => response.json())
        .then(
          (data) =>
            (this.images = data.images.map((image) => ({ name: image })))
        )
        .catch((error) => console.error("Error fetching images:", error));
    },
    archiveImages() {
      fetch(`http://localhost:8199/archive/${this.deviceId}`, {
        method: "POST",
      })
        .then((response) => response.blob())
        .then((blob) => {
          const url = window.URL.createObjectURL(blob);
          const a = document.createElement("a");
          a.href = url;
          a.download = `${this.deviceId}.zip`;
          document.body.appendChild(a);
          a.click();
          a.remove();
          this.fetchImages(); // 更新图片列表
        })
        .catch((error) => console.error("Error archiving images:", error));
    },
    viewImage(imageName) {
      this.currentImageSrc = `http://localhost:8199/images/${this.deviceId}/${imageName}`;
      this.dialogVisible = true;
    },
  },
  created() {
    this.fetchImages();
  },
};
</script>

<style>
.image-list {
  margin-top: 20px;
}

.el-table .el-button {
  padding: 0;
  margin: 0;
}
</style>
