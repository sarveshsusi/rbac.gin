import { useEffect, useState } from "react";
import { AuthContext } from "./AuthContext";
import {
  login as loginApi,
  logout as logoutApi,
  profile,
  verify2FA as verify2FAApi,
} from "../api/auth.api";
import { setAccessToken } from "../api/axios";
import { authChannel } from "./authChannel";

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // 🔄 Rehydrate on refresh
  useEffect(() => {
    const hydrate = async () => {
      if (!sessionStorage.getItem("access_token")) {
        setLoading(false);
        return;
      }

      try {
        const res = await profile();
        setUser(res.data);
      } catch {
        setUser(null);
      } finally {
        setLoading(false);
      }
    };

    hydrate();
  }, []);

  // 🔁 Multi-tab sync
  useEffect(() => {
    authChannel.onmessage = async (event) => {
      if (event.data?.type === "LOGOUT") {
        setAccessToken(null);
        setUser(null);
      }

      if (event.data?.type === "LOGIN") {
        try {
          const res = await profile();
          setUser(res.data);
        } catch {
          setUser(null);
        }
      }
    };

    return () => authChannel.close();
  }, []);

  // 🔐 LOGIN
  const login = async (email, password) => {
  const res = await loginApi({ email, password });

  // 🔐 2FA REQUIRED
  if (res.data?.two_fa_required) {
    return {
      status: "2FA",
      token: res.data.two_fa_token,
    };
  }

  // ✅ NORMAL LOGIN
  setAccessToken(res.data.access_token);
  setUser(res.data.user);
  authChannel.postMessage({ type: "LOGIN" });

  return {
    status: "OK",
  };
};


  // 🔐 VERIFY 2FA
  const verify2FA = async (code, twoFAToken) => {
    const res = await verify2FAApi(code, twoFAToken);

    setAccessToken(res.data.access_token);
    setUser(res.data.user);
    authChannel.postMessage({ type: "LOGIN" });

    return { success: true };
  };

  const logout = async () => {
    try {
      await logoutApi();
    } finally {
      setAccessToken(null);
      setUser(null);
      authChannel.postMessage({ type: "LOGOUT" });
    }
  };

  return (
    <AuthContext.Provider
      value={{ user, loading, login, verify2FA, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
}
