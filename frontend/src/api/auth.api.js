import api from "./axios";

export const login = (data) => api.post("/auth/login", data);
export const logout = () => api.post("/logout");
export const createUser = (data) => api.post("/admin/users", data);

// 🔥 MUST use the SAME axios instance
export const profile = () => api.get("/profile");
