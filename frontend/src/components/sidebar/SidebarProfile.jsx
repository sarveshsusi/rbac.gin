import { User, LogOut } from "lucide-react";

export function SidebarProfile({ user, logout, open }) {
  return (
    <div className="mt-5 m-1">
      <div className={`rounded-xl ${open ? 'bg-slate-100' : 'bg-white'} p-3 space-y-3`}>
        {/* PROFILE INFO */}
        <div className="flex items-start gap-3">
          <User className={`h-5 w-5 ${open ? 'text-slate-500' : 'text-slate-500'} mt-0.5 shrink-0`} />

          <div
            className={`
              transition-all duration-200 justify-between items-center
              ${
                open
                  ? "opacity-100 translate-x-0"
                  : "opacity-0 -translate-x-2 pointer-events-none"
              }
            `}
          >
            <p className="text-sm font-medium text-slate-700">
              {user.name}
            </p>

            <p className="text-xs text-slate-500 lowercase">
              {user.email}
            </p>

            {/* <p className="text-xs text-slate-500 capitalize">
              {user.role}
            </p> */}
          </div>
        </div>

        {/* LOGOUT BUTTON */}
        <button
          onClick={logout}
          className={`
            w-full flex items-center justify-center gap-2
            rounded-lg px-3 py-2.5
            text-sm font-medium
            text-red-600
            hover:bg-red-100
            transition
            ${open ? "opacity-100" : "opacity-0 pointer-events-none"}
          `}
        >
          <LogOut size={16} />
          Logout
        </button>
      </div>
    </div>
  );
}
