import { Search, MoreHorizontal } from "lucide-react";

const USERS = [
  {
    id: 1,
    name: "Jane Smith",
    email: "janesmith@gmail.com",
    role: "Customer",
    status: "Pending",
    date: "Nov 05, 2023",
  },
  {
    id: 2,
    name: "Marcus Wright",
    email: "marcus@enterprise.com",
    role: "Admin",
    status: "Active",
    date: "Oct 12, 2023",
  },
  {
    id: 3,
    name: "Aisha Khan",
    email: "a.khan@design.io",
    role: "Manager",
    status: "Inactive",
    date: "Sep 28, 2023",
  },
];

const STATUS_STYLE = {
  Active: "bg-green-100 text-green-700",
  Pending: "bg-yellow-100 text-yellow-700",
  Inactive: "bg-slate-200 text-slate-700",
};

export default function UsersTable() {
  return (
    <div className="mt-8 w-full rounded-2xl bg-white border border-slate-200 shadow-[0_10px_40px_rgba(15,23,42,0.08)] overflow-hidden">

      {/* SEARCH */}
      <div className="flex items-center gap-4 p-4 sm:p-5 border-b">
        <div className="relative w-full max-w-md">
          <Search className="absolute left-3 top-3.5 h-4 w-4 text-slate-400" />
          <input
            className="w-full rounded-xl bg-slate-50 px-4 py-3 pl-9 text-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
            placeholder="Search users..."
          />
        </div>
      </div>

      {/* ================= MOBILE CARDS ================= */}
      <div className="lg:hidden divide-y">
        {USERS.map((u) => (
          <div key={u.id} className="p-4 flex items-start justify-between">
            <div className="flex gap-3">
              <div className="h-10 w-10 rounded-full bg-blue-600 text-white font-semibold flex items-center justify-center">
                {u.name[0]}
                {u.name.split(" ")[1][0]}
              </div>

              <div>
                <p className="font-medium text-slate-900">{u.name}</p>
                <p className="text-xs text-slate-500">{u.email}</p>

                <div className="mt-2 flex flex-wrap items-center gap-2">
                  <span className="text-xs text-slate-600">{u.role}</span>
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${STATUS_STYLE[u.status]}`}>
                    {u.status}
                  </span>
                </div>

                <p className="mt-1 text-xs text-slate-400">
                  Created {u.date}
                </p>
              </div>
            </div>

            <MoreHorizontal className="h-5 w-5 text-slate-400 mt-1" />
          </div>
        ))}
      </div>

      {/* ================= DESKTOP TABLE ================= */}
      <div className="hidden lg:block overflow-x-auto">
        <table className="min-w-[720px] w-full text-sm">
          <thead className="bg-slate-50 text-slate-500">
            <tr>
              <th className="px-6 py-4 text-left">User</th>
              <th className="px-6 py-4 text-left">Role</th>
              <th className="px-6 py-4 text-left">Status</th>
              <th className="px-6 py-4 text-left">Created At</th>
              <th className="px-6 py-4 text-right">Actions</th>
            </tr>
          </thead>

          <tbody>
            {USERS.map((u) => (
              <tr key={u.id} className="border-t hover:bg-slate-50">
                <td className="px-6 py-4">
                  <div className="flex items-center gap-3">
                    <div className="h-10 w-10 rounded-full bg-blue-600 text-white font-semibold flex items-center justify-center">
                      {u.name[0]}
                      {u.name.split(" ")[1][0]}
                    </div>
                    <div>
                      <p className="font-medium text-slate-900">{u.name}</p>
                      <p className="text-xs text-slate-500">{u.email}</p>
                    </div>
                  </div>
                </td>

                <td className="px-6 py-4">{u.role}</td>

                <td className="px-6 py-4">
                  <span className={`rounded-full px-3 py-1 text-xs font-medium ${STATUS_STYLE[u.status]}`}>
                    {u.status}
                  </span>
                </td>

                <td className="px-6 py-4 text-slate-600">{u.date}</td>

                <td className="px-6 py-4 text-right">
                  <MoreHorizontal className="h-5 w-5 text-slate-400 cursor-pointer" />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* FOOTER */}
      <div className="flex flex-col sm:flex-row gap-3 items-center justify-between p-4 sm:p-5 text-sm text-slate-500 border-t">
        <span>Showing 1 to 3 of 12 users</span>
        <div className="flex gap-2">
          <button className="h-8 w-8 rounded-lg border">1</button>
          <button className="h-8 w-8 rounded-lg border bg-blue-600 text-white">2</button>
          <button className="h-8 w-8 rounded-lg border">3</button>
        </div>
      </div>
    </div>
  );
}
