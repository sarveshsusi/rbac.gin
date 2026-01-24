import { Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "./auth/useAuth";
import RequireAuth from "./auth/RequireAuth";
import RequireGuest from "./auth/RequireGuest";

import LoginPage from "./pages/Login-page";
import Dashboard from "./pages/Dashboard";
import AdminUsers from "./pages/AdminUsers";
import AppLayout from "./layout/AppLayout";

function App() {
  const { loading } = useAuth();
  if (loading) return null;

  return (
    <Routes>
      {/* Guest-only */}
      <Route
        path="/login"
        element={
          <RequireGuest>
            <LoginPage />
          </RequireGuest>
        }
      />

      {/* Protected */}
      <Route element={<RequireAuth />}>
        <Route element={<AppLayout />}>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/admin/users" element={<AdminUsers />} />
        </Route>
      </Route>

      {/* Fallback */}
      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  );
}

export default App;
