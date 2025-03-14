<template>
  <el-container style="height: 100vh">
    <el-aside width="220px" class="sidebar" v-if="isLoggedIn">
      <div class="logo">
        <img src="@/assets/logo.png" alt="Logo" />
      </div>
      <el-menu
        :default-active="activeIndex"
        class="el-menu-vertical-demo"
        @select="handleSelect"
      >
        <el-menu-item index="1" @click="navigate('/home')" icon="el-icon-house"
          >首页</el-menu-item
        >
        <el-menu-item
          index="2"
          @click="navigate('/devices')"
          icon="el-icon-s-management"
          >设备管理</el-menu-item
        >
        <el-menu-item
          index="3"
          @click="navigate('/monitor')"
          icon="el-icon-video-camera"
          >监控平台</el-menu-item
        >
        <el-menu-item index="4" @click="logout" icon="el-icon-switch-button"
          >登出</el-menu-item
        >
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header" v-if="isLoggedIn">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item
            v-for="(item, index) in breadcrumbs"
            :key="index"
            >{{ item }}</el-breadcrumb-item
          >
        </el-breadcrumb>
      </el-header>

      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import { isLoggedIn, logout } from "./auth";

export default {
  data() {
    return {
      activeIndex: "1",
      breadcrumbs: ["首页"],
    };
  },
  computed: {
    isLoggedIn() {
      return isLoggedIn();
    },
  },
  methods: {
    handleSelect(key) {
      if (key === "1") {
        this.$router.push("/home");
      } else if (key === "2") {
        this.$router.push("/devices");
      } else if (key === "3") {
        this.$router.push("/monitor");
      }
    },
    updateBreadcrumbs(route) {
      const paths = route.path.split("/").filter((p) => p);
      this.breadcrumbs = ["首页", ...paths];
    },
    navigate(path) {
      this.$router.push(path);
    },
    logout() {
      logout();
      this.$router.push("/login");
      location.reload();
    },
  },
  watch: {
    $route(to) {
      this.updateBreadcrumbs(to);
    },
  },
  created() {
    this.updateBreadcrumbs(this.$route);
  },
};
</script>

<style>
.sidebar {
  background-color: #2c3e50;
  color: #ecf0f1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 20px 0;
}

.logo {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 60px;
  margin-bottom: 20px;
}

.logo img {
  height: 40px;
}

.el-menu-vertical-demo {
  border-right: none;
}

.el-menu-item {
  color: #ecf0f1;
}

.el-menu-item:hover {
  background-color: #34495e;
}

.el-menu-item.is-active {
  background-color: #1abc9c;
}

.header {
  background-color: #ecf0f1;
  color: #2c3e50;
  padding: 0 20px;
  display: flex;
  align-items: center;
  height: 60px;
}
</style>
