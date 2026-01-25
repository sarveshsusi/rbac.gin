import { useRef } from "react";
import {
  Sidebar as RadixSidebar,
  SidebarContent,
  SidebarHeader,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarFooter,
} from "../ui/sidebar";

import { NavLink } from "react-router-dom";
import { Menu, X, User, LogOut } from "lucide-react";
import { useSidebar } from "./useSidebar";
import { useAuth } from "../../auth/useAuth";
import { SIDEBAR_ITEMS } from "../../config/sidebarConfig";

export default function Sidebar() {
  const { open, setOpen } = useSidebar();
  const { user, logout } = useAuth();
  const hoverTimeout = useRef(null);

  if (!user) return null;

  const role = user.role.toLowerCase();
  const allowedItems = SIDEBAR_ITEMS.filter((item) =>
    item.roles.includes(role)
  );

  return (
    <>
      {/* ================= SIDEBAR ================= */}
      <RadixSidebar
        className={`
          fixed inset-y-0 left-0 z-50
          flex flex-col h-screen          /* ✅ KEY FIX */
          w-64 p-1
          bg-gradient-to-b from-white to-slate-50
          shadow-[4px_0_24px_rgba(0,0,0,0.06)]
          transition-all duration-300 ease-in-out

          /* MOBILE */
          ${open ? "translate-x-0" : "-translate-x-full"}

          /* DESKTOP */
          lg:translate-x-0
          ${open ? "lg:w-64" : "lg:w-16"}
        `}
        onMouseEnter={() => {
          if (window.innerWidth < 1024) return;
          clearTimeout(hoverTimeout.current);
          hoverTimeout.current = setTimeout(() => setOpen(true), 120);
        }}
        onMouseLeave={() => {
          if (window.innerWidth < 1024) return;
          clearTimeout(hoverTimeout.current);
          hoverTimeout.current = setTimeout(() => setOpen(false), 120);
        }}
      >
        {/* ================= HEADER ================= */}
        <SidebarHeader>
          <div className="relative flex h-14 items-center justify-between px-4">
            {/* Brand */}
            <span
              className={`
                text-sm font-semibold tracking-wide
                whitespace-nowrap transition-all duration-200
                ${open ? "opacity-100 translate-x-0" : "opacity-0 -translate-x-2"}
              `}
            >
              RBAC App
            </span>

            {/* Collapsed logo */}
            <span
              className={`
                absolute left-4 text-sm font-semibold
                transition-opacity duration-200
                ${open ? "opacity-0" : "opacity-100"}
              `}
            >
              R
            </span>

            {/* Toggle */}
            <button
              onClick={() => setOpen(!open)}
              className="rounded-md p-1.5 hover:bg-slate-100"
            >
              {open ? <X size={18} /> : <Menu size={16} />}
            </button>
          </div>
        </SidebarHeader>

        {/* ================= MENU ================= */}
        <SidebarContent className="flex-1 pt-2">
          <SidebarGroup>
            <SidebarGroupContent>
              <SidebarMenu className="space-y-1">
                {allowedItems.map((item) => (
                  <SidebarMenuItem key={item.path}>
                    <NavLink
                      to={item.path}
                      className={({ isActive }) =>
                        `
                        group flex items-center gap-3
                        rounded-xl px-3 py-2.5 text-sm
                        transition-all duration-200
                        ${
                          isActive
                            ? "bg-blue-600 text-white shadow-md"
                            : "text-slate-600 hover:bg-slate-100"
                        }
                        `
                      }
                    >
                      <item.icon className="h-4 w-4 shrink-0" />

                      <span
                        className={`
                          whitespace-nowrap transition-all duration-200
                          ${
                            open
                              ? "opacity-100 translate-x-0"
                              : "opacity-0 -translate-x-2 pointer-events-none"
                          }
                        `}
                      >
                        {item.label}
                      </span>
                    </NavLink>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>

        {/* ================= FOOTER (BOTTOM FIXED) ================= */}
        <SidebarFooter className="mt-auto px-3 pb-3">
          <div className="rounded-xl bg-slate-100 p-3 mt-2 space-y-3">
            {/* Profile */}
            <div className="flex items-center gap-3 mb-2">
              <User className="h-5 w-5 text-slate-500" />

              <div
                className={`
                  transition-all duration-200
                  ${
                    open
                      ? "opacity-100 translate-x-0"
                      : "opacity-0 -translate-x-2 pointer-events-none"
                  }
                `}
              >
                <p className="text-sm font-medium text-slate-700">
                  {user.name}
                </p>
               
                <p className="text-xs text-slate-500 capitalize">
                  {user.email}
                </p>
                 <p className="text-xs text-slate-500 capitalize">
                  {user.role}
                </p>
              </div>
            </div>

            {/* Logout */}
            <button
              onClick={logout}
              className={`
                flex justify-center items-center gap-2
                w-full rounded-lg px-3 py-3
                text-sm font-medium
                text-red-600 hover:bg-red-100
                transition-all
                ${
                  open
                    ? "opacity-100"
                    : "opacity-0 pointer-events-none"
                }
              `}
            >
              <LogOut size={16} />
              Logout
            </button>
          </div>
        </SidebarFooter>
      </RadixSidebar>

      {/* ================= MOBILE OVERLAY ================= */}
      {open && (
        <div
          className="fixed inset-0 z-40 bg-black/40 lg:hidden"
          onClick={() => setOpen(false)}
        />
      )}
    </>
  );
}
