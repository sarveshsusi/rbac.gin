import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080/api/v1",
  withCredentials: true, // REQUIRED for refresh cookie
});

/* =========================
   ACCESS TOKEN MEMORY
========================= */
let accessToken = sessionStorage.getItem("access_token");
let isRefreshing = false;
let refreshQueue = [];

/* =========================
   TOKEN SETTER
========================= */
export const setAccessToken = (token) => {
  accessToken = token;

  if (token) {
    sessionStorage.setItem("access_token", token);
  } else {
    sessionStorage.removeItem("access_token");
  }
};

/* =========================
   QUEUE HANDLER
========================= */
const processQueue = (token, error) => {
  refreshQueue.forEach(({ resolve, reject }) => {
    if (error) reject(error);
    else resolve(token);
  });
  refreshQueue = [];
};

/* =========================
   REQUEST INTERCEPTOR
========================= */
api.interceptors.request.use((config) => {
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`;
  }
  return config;
});

/* =========================
   RESPONSE INTERCEPTOR
========================= */
api.interceptors.response.use(
  (res) => res,
  async (err) => {
    const original = err.config;

    // Ignore non-401 errors
    if (err.response?.status !== 401) {
      return Promise.reject(err);
    }

    // ❌ Never intercept refresh itself
    if (original.url?.includes("/auth/refresh")) {
      setAccessToken(null);
      return Promise.reject(err);
    }

    // ❌ Prevent infinite loops
    if (original._retry) {
      return Promise.reject(err);
    }

    // 🔄 Queue requests while refreshing
    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        refreshQueue.push({
          resolve: (token) => {
            original.headers.Authorization = `Bearer ${token}`;
            resolve(api(original));
          },
          reject,
        });
      });
    }

    original._retry = true;
    isRefreshing = true;

    try {
      const res = await api.post("/auth/refresh");

      const newToken = res.data.access_token;
      setAccessToken(newToken);

      processQueue(newToken, null);

      original.headers.Authorization = `Bearer ${newToken}`;
      return api(original);
    } catch (refreshErr) {
      setAccessToken(null);
      processQueue(null, refreshErr);
      return Promise.reject(refreshErr);
    } finally {
      isRefreshing = false;
    }
  }
);

export default api;
