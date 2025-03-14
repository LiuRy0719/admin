export function isLoggedIn() {
  return !!localStorage.getItem("token");
}

export function login(token) {
  localStorage.setItem("token", token);
  localStorage.setItem("isLoggedIn", "true");
}

export function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("isLoggedIn");
}
