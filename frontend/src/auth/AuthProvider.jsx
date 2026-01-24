import { useEffect, useState } from "react";
import { AuthContext } from "./AuthContext";
import {
  login as loginApi,
  logout as logoutApi,
  profile,
} from "../api/auth.api";
import { setAccessToken } from "../api/axios";
import { authChannel } from "./authChannel";

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // 🔄 REHYDRATE ON REFRESH
  useEffect(() => {
    const hydrate = async () => {
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

  // 🔥 MULTI-TAB SYNC (LOGIN + LOGOUT)
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

  const login = async (email, password) => {
    const res = await loginApi({ email, password });
    setAccessToken(res.data.access_token);
    setUser(res.data.user);

    // 🔥 broadcast login
    authChannel.postMessage({ type: "LOGIN" });
  };

  const logout = async () => {
    try {
      await logoutApi();
    } finally {
      setAccessToken(null);
      setUser(null);

      // 🔥 broadcast logout
      authChannel.postMessage({ type: "LOGOUT" });
    }
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
