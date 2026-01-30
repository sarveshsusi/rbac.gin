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
    "/dashboard": "Dashboard",
    "/customer/dashboard": "Dashboard",

    "/admin/users": "Admin Users",
    "/admin/products": "Products",
    "/admin/amc": "AMC Contracts",
    "/admin/assign-products": "Assign Products",

    "/support/tickets": "My Tickets",

    "/customer/tickets": "My Tickets",
    "/customer/amc": "My AMC",

    "/tickets/new": "Create Ticket",
    "/admin/tickets": "Ticket Assignment",
    "/admin/tickets/new": "Create Ticket (Admin)",
  };


  const title =
    pageTitleMap[location.pathname] ?? "RBAC Portal";

  return (
    <div className="min-h-screen flex bg-[#F4F7FB]">
      {/* SIDEBAR */}
      <Sidebar />

      {/* MAIN CONTENT */}
      <div
        className={`
          flex-1 flex flex-col min-h-screen
          transition-all duration-300 ease-in-out
          ${open ? "lg:ml-64" : "lg:ml-16"}
          ml-0
        `}
      >
        {/* TOPBAR */}
        <Topbar title={title} />

        {/* CONTENT */}
        <main className="flex-1 px-4 md:px-6 py-4 pb-24">
          <Outlet />
        </main>
      </div>

      {/* MOBILE NAV */}
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
