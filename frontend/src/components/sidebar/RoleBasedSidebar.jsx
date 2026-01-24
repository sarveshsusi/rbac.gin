import { NavLink } from "react-router-dom";
import { useAuth } from "../../auth/useAuth";
import { SIDEBAR_ITEMS } from "../../config/sidebarConfig";

export default function Sidebar() {
  const { user } = useAuth();

  // auth not ready or logged out
  if (!user) return null;

  return (
    <aside className="w-64 min-h-screen border-r bg-white p-4">
      <h2 className="mb-6 text-lg font-bold">App</h2>

      <nav className="space-y-1">
        {SIDEBAR_ITEMS
          .filter(item => item.roles.includes(user.role))
          .map(item => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) =>
                `flex items-center gap-3 rounded px-3 py-2 text-sm
                 ${isActive
                   ? "bg-blue-600 text-white"
                   : "text-gray-700 hover:bg-gray-100"}`
              }
            >
              <item.icon className="h-4 w-4" />
              {item.label}
            </NavLink>
          ))}
      </nav>
    </aside>
  );
}
