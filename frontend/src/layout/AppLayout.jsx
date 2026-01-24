import { Outlet } from "react-router-dom";
import Sidebar from "../components/sidebar/Sidebar";
import { SidebarProvider } from "../components/sidebar/SidebarProvider";
import { useSidebar } from "../components/sidebar/useSidebar";
import MobileBottomNav from "../components/sidebar/MobileBottomNav";

function LayoutContent() {
  const { open } = useSidebar();

  return (
    <div className="min-h-screen flex">
      {/* SIDEBAR */}
      <Sidebar />

      {/* MAIN CONTENT */}
      <main
        className={`
          flex-1 min-h-screen
          transition-all duration-300 ease-in-out

          /* Desktop offset */
          ${open ? "lg:ml-64" : "lg:ml-16"}

          /* Mobile */
          ml-0

          px-6 py-4
          pb-20
        `}
      >
        <Outlet />
      </main>

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
