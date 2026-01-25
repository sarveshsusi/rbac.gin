import { NavLink } from "react-router-dom";
import { Menu } from "lucide-react";
import { useSidebar } from "./useSidebar";
import { useAuth } from "../../auth/useAuth";
import { SIDEBAR_ITEMS } from "../../config/sidebarConfig";

export default function MobileBottomNav() {
  const { setOpen } = useSidebar();
  const { user } = useAuth();

  if (!user) return null;

  const role = user.role.toLowerCase();

  // Pick ONLY first 3 items for mobile
  const primaryItems = SIDEBAR_ITEMS
    .filter(item => item.roles.includes(role))
    .slice(0, 3);

  return (
    <nav
      className="
        fixed bottom-0 left-0 right-0 z-50
        bg-white border-t border-gray-200
        flex justify-around items-center
        h-15 
        lg:hidden
      "
    >
      {primaryItems.map(item => (
        <NavLink
          key={item.path}
          to={item.path}
          className={({ isActive }) =>
            `
            flex flex-col items-center gap-1 text-xs
            ${
              isActive
                ? "text-blue-600"
                : "text-gray-500"
            }
            `
          }
        >
          <item.icon className="h-5 w-5" />
          <span>{item.label}</span>
        </NavLink>
      ))}

      {/* MENU BUTTON (4th ITEM) */}
      <button
        onClick={() => setOpen(true)}
        className="flex flex-col items-center gap-1 text-xs text-gray-500"
      >
        <Menu className="h-5 w-5" />
        <span>Menu</span>
      </button>
    </nav>
  );
}
