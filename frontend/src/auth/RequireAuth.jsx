import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "./useAuth";

export default function RequireAuth() {
  const { user, loading } = useAuth();

  if (loading) return null; // ⛔ WAIT

  return user ? <Outlet /> : <Navigate to="/login" replace />;
}
