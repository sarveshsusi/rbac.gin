import api from "./axios";

export const login = (data) => api.post("/auth/login", data);
export const logout = () => api.post("/logout");
export const createUser = (data) => api.post("/admin/users", data);

// 🔥 MUST use the SAME axios instance
export const profile = () => api.get("/profile");
export const resetPassword = (data) =>
  api.post("/auth/reset-password", data);

// ADD BELOW existing exports

export const verify2FA = (code, twoFAToken) =>
  api.post(
    "/auth/verify-2fa",
    { code },
    {
      headers: {
        "X-2FA-Token": twoFAToken,
      },
    }
  );
