import { Outlet, useLocation } from "react-router-dom";
import Sidebar from "../components/sidebar/Sidebar";
import { SidebarProvider } from "../components/sidebar/SidebarProvider";
import { useSidebar } from "../components/sidebar/useSidebar";
import MobileBottomNav from "../components/sidebar/MobileBottomNav";
import Topbar from "./TopBar";

function LayoutContent() {
  const { open } = useSidebar();
  const location = useLocation();

  const pageTitleMap = {
    "/admin/users": "Users",
    "/dashboard": "Dashboard",
  };

  const title = pageTitleMap[location.pathname] ?? "RBAC Portal";

  return (
    <div className="min-h-screen flex bg-[#F4F7FB]">
      {/* SIDEBAR */}
      <Sidebar />

      {/* MAIN WRAPPER */}
      <div
        className={`
          flex-1 flex flex-col min-h-screen
          transition-all duration-300 ease-in-out

          ${open ? "lg:ml-64" : "lg:ml-16"}
          ml-0
        `}
      >
        {/* TOP HEADER */}
        <Topbar title={title} />

        {/* PAGE CONTENT */}
        <main
          className="
            flex-1
            px-4 md:px-6 py-4
            pb-24
          "
        >
          <Outlet />
        </main>
      </div>

      {/* MOBILE BOTTOM NAV */}
      <MobileBottomNav />
    </div>
  );
}

export default function AppLayout() {
  return (
    <SidebarProvider>
      <LayoutContent />
    </SidebarProvider>
  );
}
