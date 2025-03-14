import { createApp } from "vue";
import App from "./App.vue";
import { createRouter, createWebHistory } from "vue-router";
import HomePage from "./components/HomePage.vue";
import ImageList from "./components/ImageList.vue";
import DeviceImages from "./components/DeviceImages.vue";
import UserLogin from "./components/UserLogin.vue";
import MonitorPlatformPage from "./components/MonitorPlatformPage.vue";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import { isLoggedIn } from "./auth";

const routes = [
  { path: "/", redirect: "/login" }, // 默认重定向到登录页面
  { path: "/login", component: UserLogin },
  { path: "/home", component: HomePage, meta: { requiresAuth: true } },
  { path: "/devices", component: ImageList, meta: { requiresAuth: true } },
  {
    path: "/devices/:deviceId",
    component: DeviceImages,
    meta: { requiresAuth: true },
  },
  {
    path: "/monitor",
    component: MonitorPlatformPage,
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  if (to.matched.some((record) => record.meta.requiresAuth) && !isLoggedIn()) {
    next("/login");
  } else {
    next();
  }
});

const app = createApp(App);
app.use(router);
app.use(ElementPlus);
app.mount("#app");
