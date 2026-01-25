import { Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "./auth/useAuth";
import RequireAuth from "./auth/RequireAuth";
import RequireGuest from "./auth/RequireGuest";

import LoginPage from "./pages/Login-page";
import VerifyOTP from "./pages/VerifyOTP"; // ✅ ADD
import ResetPassword from "./pages/ResetPassword";
import ForgotPasswordPage from "./pages/ForgotPasswordPage";

import Dashboard from "./pages/Dashboard";
import AdminUsers from "./pages/AdminUsers";
import AppLayout from "./layout/AppLayout";

function App() {
  const { loading } = useAuth();
  if (loading) return null;

  return (
    <Routes>
      {/* =====================
          PUBLIC / GUEST ROUTES
      ====================== */}
      <Route
        path="/login"
        element={
          <RequireGuest>
            <LoginPage />
          </RequireGuest>
        }
      />

      {/* ✅ 2FA OTP (GUEST ONLY) */}
      <Route
        path="/verify-otp"
        element={
          <RequireGuest>
            <VerifyOTP />
          </RequireGuest>
        }
      />

      {/* ✅ PASSWORD RESET (PUBLIC, NO GUARDS) */}
      <Route path="/reset-password" element={<ResetPassword />} />
      <Route path="/forgot-password" element={<ForgotPasswordPage />} />

      {/* =====================
          PROTECTED ROUTES
      ====================== */}
      <Route element={<RequireAuth />}>
        <Route element={<AppLayout />}>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/admin/users" element={<AdminUsers />} />
          
        </Route>
      </Route>
       <Route element={<RequireAuth />}>
        <Route element={<AppLayout />}>
          <Route path="/customer/dashboard" element={<Dashboard />} />     
        </Route>
      </Route>

      {/* =====================
          FALLBACK
      ====================== */}
      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  );
}

export default App;
