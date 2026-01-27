import {
  Home,
  Users,
  Package,
  UserCircle,
  Ticket,
  Settings,
  Shield,
  UserPlus
} from "lucide-react";

export const SIDEBAR_ITEMS = [
  /* =====================
     DASHBOARD
  ====================== */
  {
    label: "Dashboard",
    path: "/dashboard",
    icon: Home,
    roles: ["admin", "support"],
  },
  {
    label: "Dashboard",
    path: "/customer/dashboard",
    icon: Home,
    roles: ["customer"],
  },

  /* =====================
     ADMIN
  ====================== */
  {
    label: "Admin Users",
    path: "/admin/users",
    icon: Shield,
    roles: ["admin"],
  },
  {
    label: "Products",
    path: "/admin/products",
    icon: Package,
    roles: ["admin"],
  },
  // {
  //   label: "Create Product",
  //   path: "/admin/products/new",
  //   icon: Package,
  //   roles: ["admin"],
  // },
  {
    label: "AMCs",
    path: "/admin/amc",
    icon: Users,
    roles: ["admin"],
  },

  /* =====================
     SUPPORT
  ====================== */
  {
    label: "My Tickets",
    path: "/support/tickets",
    icon: Ticket,
    roles: ["support"],
  },
 
{
  label: "Assign Products",
  path: "/admin/assign-products",
  roles: ["admin"],
  icon: UserPlus,
},

  /* =====================
     CUSTOMER
  ====================== */
  {
    label: "My Tickets",
    path: "/customer/tickets",
    icon: Ticket,
    roles: ["customer"],
  },
  {
    label: "My AMC",
    path: "/customer/amc",
    icon: UserCircle,
    roles: ["customer"],
  },
  {
  label: "Create Ticket",
  path: "/tickets/new",
  icon: Ticket,
  roles: ["customer"],
},

  /* =====================
     SETTINGS
  ====================== */
  {
    label: "Settings",
    path: "/settings",
    icon: Settings,
    roles: ["admin"],
  },
];
