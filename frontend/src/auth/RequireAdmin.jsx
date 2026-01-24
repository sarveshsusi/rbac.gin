import { Navigate } from "react-router-dom";
import { useAuth } from "./useAuth";

export default function RequireAdmin({ children }) {
  const { user, loading } = useAuth();

  if (loading) return null;

  if (!user || user.role !== "admin") {
    return <Navigate to="/" replace />;
  }

  return children;
}
