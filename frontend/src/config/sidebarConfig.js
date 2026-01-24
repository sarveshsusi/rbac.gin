import {
  Home,
  Users,
  Package,
  UserCircle,
  Ticket,
  Settings,
  Shield,
} from "lucide-react";

export const SIDEBAR_ITEMS = [
  // ===== Common =====
  {
    label: "Dashboard",
    path: "/dashboard",
    icon: Home,
    roles: ["admin", "user", "support"],
  },

  // ===== Admin =====
  {
    label: "Admin Users",
    path: "/admin/users",
    icon: Shield,
    roles: ["admin"],
  },

  {
    label: "Users",
    path: "/users",
    icon: Users,
    roles: ["admin"],
  },

  {
    label: "Products",
    path: "/products",
    icon: Package,
    roles: ["admin"],
  },

  // ===== Support / Operations =====
  {
    label: "Customers",
    path: "/customers",
    icon: UserCircle,
    roles: ["admin", "support"],
  },

  {
    label: "Tickets",
    path: "/tickets",
    icon: Ticket,
    roles: ["admin", "support"],
  },

  // ===== Settings =====
  {
    label: "Settings",
    path: "/settings",
    icon: Settings,
    roles: ["admin"],
  },
];
