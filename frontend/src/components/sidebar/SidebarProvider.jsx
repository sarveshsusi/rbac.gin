import { useState } from "react";
import { SidebarContext } from "./SidebarContext";

export function SidebarProvider({ children }) {
  const [open, setOpen] = useState(false); // mobile-first

  return (
    <SidebarContext.Provider value={{ open, setOpen }}>
      {children}
    </SidebarContext.Provider>
  );
}
