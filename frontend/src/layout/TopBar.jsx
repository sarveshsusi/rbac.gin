import { Bell, Moon } from "lucide-react";

export default function Topbar({ title }) {
  return (
    <header
      className="
        sticky top-0 z-40
        h-14
        flex items-center justify-between
        px-4 md:px-6

        bg-white/80 backdrop-blur
        border-b border-slate-200/60

        /* ✅ SAFE AREA */
        pt-[env(safe-area-inset-top)]
      "
    >
      {/* LEFT */}
      <h1 className="text-lg font-semibold text-slate-900">
        {title}
      </h1>

      {/* RIGHT */}
      <div className="flex items-center gap-2">
        <button className="h-9 w-9 rounded-lg flex items-center justify-center text-slate-600 hover:bg-slate-100 transition">
          <Bell size={18} />
        </button>

        <button className="h-9 w-9 rounded-lg flex items-center justify-center text-slate-600 hover:bg-slate-100 transition">
          <Moon size={18} />
        </button>
      </div>
    </header>
  );
}
