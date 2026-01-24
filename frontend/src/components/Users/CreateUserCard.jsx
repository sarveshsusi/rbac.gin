import { useState } from "react";
import { Eye, EyeOff, Info, UserPlus } from "lucide-react";

export default function CreateUserCard() {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("Admin");
  const [showPassword, setShowPassword] = useState(false);

  return (
    <div className=" w-full bg-white border border-slate-200
    shadow-[0_10px_40px_rgba(15,23,42,0.08)]
    rounded-none px-4 py-4
    sm:rounded-2xl sm:px-6 sm:py-6 sm:mx-4
    lg:mx-auto lg:max-w-5xl lg:px-8 lg:py-8">

      {/* HEADER */}
      <div>
        <h2 className="text-lg sm:text-xl font-semibold text-slate-900">
          Create New User
        </h2>
        <p className="mt-1 text-sm text-slate-500">
          Add a new team member and assign their role.
        </p>
      </div>

      {/* FORM */}
      <div className="mt-6 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">

        {/* Username */}
        <div>
          <label className="block text-sm font-medium text-slate-700">
            Username
          </label>
          <input
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="e.g. janesmith"
            className="mt-1 w-full rounded-xl border border-slate-300
                       px-4 py-3 text-sm
                       focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        {/* Email */}
        <div>
          <label className="block text-sm font-medium text-slate-700">
            Email Address
          </label>
          <input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="jane@example.com"
            className="mt-1 w-full rounded-xl border border-slate-300
                       px-4 py-3 text-sm
                       focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        {/* Password */}
        <div className="relative sm:col-span-2 lg:col-span-1">
          <label className="block text-sm font-medium text-slate-700">
            Password
          </label>
          <input
            type={showPassword ? "text" : "password"}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
            className="mt-1 w-full rounded-xl border border-slate-300
                       px-4 py-3 pr-12 text-sm
                       focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
          <button
            type="button"
            onClick={() => setShowPassword(v => !v)}
            className="absolute right-4 top-[42px]
                       text-slate-400 hover:text-slate-600"
          >
            {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
          </button>
        </div>
      </div>

      {/* ROLE + PERMISSIONS */}
      <div className="mt-6 grid grid-cols-1 gap-4 lg:grid-cols-2">
        <div>
          <label className="block text-sm font-medium text-slate-700">
            Role Selection
          </label>
          <select
            value={role}
            onChange={(e) => setRole(e.target.value)}
            className="mt-1 w-full rounded-xl border border-slate-300
                       px-4 py-3 text-sm
                       focus:ring-2 focus:ring-blue-500 focus:outline-none"
          >
            <option>Admin</option>
            <option>Support</option>
            <option>Customer</option>
          </select>
        </div>

        <div className="flex items-start gap-3 rounded-xl bg-blue-50
                        px-4 py-3 text-sm text-blue-700">
          <Info className="mt-0.5 h-4 w-4 shrink-0" />
          <p>
            <span className="font-semibold">Permissions:</span> Full ticket
            management, access to customer logs, and chat support capabilities.
          </p>
        </div>
      </div>

      {/* ACTION */}
      <div className="mt-8 flex justify-end">
        <button
          className="inline-flex w-full sm:w-auto items-center justify-center gap-2
                     rounded-xl bg-blue-600 px-6 py-3
                     text-sm font-semibold text-white
                     shadow-lg shadow-blue-600/30
                     hover:bg-blue-700 transition"
        >
          <UserPlus size={16} />
          Create User
        </button>
      </div>
    </div>
  );
}
