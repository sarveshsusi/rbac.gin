import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import { AuthProvider } from "./auth/AuthProvider";
import "./index.css";

import { registerSW } from "virtual:pwa-register";
import { ToastProvider } from "./components/toast/ToastProvider";

registerSW();

ReactDOM.createRoot(document.getElementById("root")).render(
  
    <AuthProvider>
      <BrowserRouter>
      <ToastProvider>
        <App />
        </ToastProvider>
      </BrowserRouter>
    </AuthProvider>
);
