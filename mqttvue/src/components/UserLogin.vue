<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h2>登录</h2>
      </div>
      <el-form :model="form" @submit.prevent="login" class="login-form">
        <el-form-item>
          <el-input
            v-model="form.username"
            placeholder="用户名"
            size="large"
            prefix-icon="el-icon-user"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            size="large"
            prefix-icon="el-icon-lock"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            @click="login"
            size="large"
            class="login-button"
            >登录</el-button
          >
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import { login } from "../auth";

export default {
  name: "UserLogin",
  data() {
    return {
      form: {
        username: "",
        password: "",
      },
    };
  },
  methods: {
    login() {
      fetch("http://localhost:8199/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(this.form),
      })
        .then((response) => response.json())
        .then((data) => {
          if (data.token) {
            login(data.token);
            localStorage.setItem("isLoggedIn", "true");
            this.$router.push("/home");
            location.reload();
          } else {
            this.$message.error("登录失败");
          }
        })
        .catch(() => {
          this.$message.error("登录失败");
        });
    },
  },
  created() {
    if (localStorage.getItem("isLoggedIn") === "true") {
      this.$router.push("/home");
    }
  },
};
</script>

<style>
html,
body {
  margin: 0;
  height: 100%;
  overflow: hidden;
}

.login-container {
  height: 95vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-image: url("@/assets/background.jpg");
  background-size: cover;
  background-repeat: no-repeat;
  background-position: center;
}

.login-box {
  width: 400px;
  padding: 20px;
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 20px;
}

.login-header h2 {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.login-form {
  display: flex;
  flex-direction: column;
}

.login-button {
  width: 100%;
}
</style>
