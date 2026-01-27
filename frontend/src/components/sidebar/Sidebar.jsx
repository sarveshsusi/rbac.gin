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
import { Menu, X } from "lucide-react";
import { useSidebar } from "./useSidebar";
import { useAuth } from "../../auth/useAuth";
import { SIDEBAR_ITEMS } from "../../config/sidebarConfig";
import { SidebarProfile } from "./SidebarProfile";

export default function Sidebar() {
  const img1 = "https://ik.imagekit.io/cj2po8igx/EMERD/ChatGPT_Image_Jan_26__2026__09_32_12_PM-removebg-preview.png"
  const img2 ="https://ik.imagekit.io/cj2po8igx/EMERD/ChatGPT_Image_Jan_26__2026__09_33_14_PM-removebg-preview.png"
  
  const { open, setOpen } = useSidebar();
  const { user, logout } = useAuth();
  const hoverTimeout = useRef(null);

  if (!user) return null;

  const role = user.role.toLowerCase();
  const allowedItems = SIDEBAR_ITEMS.filter((item) =>
    item.roles.includes(role),
  );

  const handleNavClick = () => {
    if (window.innerWidth < 1024) {
      setOpen(false);
    }
  };

  return (
    <>
      {/* ================= SIDEBAR ================= */}
      <RadixSidebar
        className={`
          fixed inset-y-0 left-0 z-50
          flex flex-col h-screen 
          w-64 p-1
          bg-linear-to-b from-white to-white
          shadow-[4px_0_24px_rgba(0,0,0,0.06)]
          transition-all duration-300 ease-in-out

          /* MOBILE */
          ${open ? "translate-x-0" : "-translate-x-full"}
          pb-[calc(env(safe-area-inset-bottom)+24px)]

          /* DESKTOP */
          lg:translate-x-0 border-r border-slate-200
          lg:pb-0
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
        <div className="border-b border-slate-200"> 
          <div className="relative flex h-14 items-center justify-between px-4">
            <span
              className={`
                text-lg font-extrabold tracking-wide
                whitespace-nowrap transition-all duration-200 

                ${open ? "opacity-95 translate-x-0" : "opacity-0 -translate-x-2"}
              `}
            >
              <img src={img1} alt="" className="" />
              {/* <span className="text-blue-600 text-2xl">E</span>MERD */}
            </span>

          <span
  className={`
    absolute left-4 text-lg font-semibold text-blue-600
    transition-all duration-300 ease-out
    will-change-transform will-change-opacity object-contain w-full h-full
    ${open
      ? "opacity-0 -translate-x-3 scale-110 blur-sm tracking-widest"
      : "opacity-100 translate-x-0 scale-100 blur-0 tracking-normal"}
  `}
>
  
  <img src={img2} alt="" className=""/>
</span>


            <button
              onClick={() => setOpen(!open)}
              className="lg:hidden rounded-md hover:bg-slate-100"
            >
              {open ? <X size={18} /> : <Menu size={15} />}
            </button>
          </div>
        </div>

        {/* ================= MENU ================= */}
        <SidebarContent className="border border-slate-200 flex-1 pt-2">
          <SidebarGroup>
            <SidebarGroupContent>
              <SidebarMenu className="space-y-1">
                {allowedItems.map((item) => (
                  <SidebarMenuItem key={item.path}>
                    <NavLink
                      to={item.path}
                      onClick={handleNavClick}
                      className={({ isActive }) =>
                        `
                        group flex items-center gap-3 mb-2
                        rounded-lg px-3 py-2.5 text-sm p-5
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
                          whitespace-nowrap transition-all duration-200 text-md
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

        {/* ================= FOOTER ================= */}
        <SidebarFooter
          className="
            mt-auto px-3 pb-3
            border-t border-slate-200/60
            lg:border-none
          "
        >
          <SidebarProfile user={user} logout={logout} open={open} />
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
