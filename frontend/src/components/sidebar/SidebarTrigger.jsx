import { Menu } from "lucide-react";
import { useSidebar } from "./useSidebar";

export function SidebarTrigger() {
  const { open, setOpen } = useSidebar();

  return (
    <button
      onClick={() => setOpen(!open)}
      className="lg:hidden rounded-md border bg-white p-2"
      aria-label="Toggle sidebar"
    >
      <Menu className="h-5 w-5" />
    </button>
  );
}
