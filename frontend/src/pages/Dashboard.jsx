import { useAuth } from "../auth/useAuth";
import usePageTitle from "../components/hooks/usePageTitle";

export default function Dashboard() {
   usePageTitle("Dashboard • RBAC App");
  const { user, logout } = useAuth();

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold">Dashboard</h1>
      <p>Email: {user.email}</p>
      <p>Role: {user.role}</p>
      <button onClick={logout} className="mt-4 bg-red-600 text-white px-4 py-2">
        Logout
      </button>
    </div>
  );
}
