import React from "react";
import * as Collapsible from "@radix-ui/react-collapsible";

// Base container
export function Sidebar({ children, className, ...props }) {
  return (
    <Collapsible.Root {...props} className={className}>
      {children}
    </Collapsible.Root>
  );
}

export function SidebarHeader({ children }) {
  return <div className="sticky top-0 border-b bg-white">{children}</div>;
}

export function SidebarFooter({ children }) {
  return <div className="border-t bg-white">{children}</div>;
}

export function SidebarContent({ children }) {
  return (
    <div className="overflow-y-auto px-1 py-2">{children}</div>
  );
}

export function SidebarGroup({ children }) {
  return <div className="mb-4">{children}</div>;
}

export function SidebarGroupLabel({ children }) {
  return (
    <div className="px-3 py-1 text-xs font-semibold uppercase text-gray-500">
      {children}
    </div>
  );
}

export function SidebarGroupContent({ children }) {
  return <div className="space-y-1">{children}</div>;
}

export function SidebarMenu({ children }) {
  return <nav>{children}</nav>;
}

export function SidebarMenuItem({ children }) {
  return <div>{children}</div>;
}
