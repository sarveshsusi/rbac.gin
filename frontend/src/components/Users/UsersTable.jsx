import { useEffect, useState, useMemo, useRef } from "react";
import { Search, MoreHorizontal, Filter } from "lucide-react";
import { getUsers } from "../../api/user.api";

/* =========================
   STATUS COLORS
========================= */
const STATUS_STYLE = {
  active: "bg-green-100 text-green-700",
  pending: "bg-yellow-100 text-yellow-700",
  inactive: "bg-slate-200 text-slate-700",
};

/* =========================
   HELPERS
========================= */
const formatDate = (date) => {
  if (!date) return "—";

  const iso = date
    .replace(" ", "T")
    .replace(/\.\d+/, "")
    .replace("+00", "Z");

  const d = new Date(iso);
  if (isNaN(d.getTime())) return "—";

  return d.toLocaleDateString("en-US", {
    month: "short",
    day: "2-digit",
    year: "numeric",
  });
};

const getInitials = (name = "") => {
  const p = name.trim().split(" ");
  return p.length > 1
    ? p[0][0].toUpperCase() + p[1][0].toUpperCase()
    : p[0]?.[0]?.toUpperCase();
};

/* =========================
   SKELETON
========================= */
const SkeletonRow = () => (
  <tr className="animate-pulse">
    <td className="px-6 py-4">
      <div className="flex items-center gap-3">
        <div className="h-10 w-10 rounded-full bg-slate-200" />
        <div className="space-y-2">
          <div className="h-3 w-32 bg-slate-200 rounded" />
          <div className="h-3 w-48 bg-slate-200 rounded" />
        </div>
      </div>
    </td>
    <td className="px-6 py-4"><div className="h-3 w-16 bg-slate-200 rounded" /></td>
    <td className="px-6 py-4"><div className="h-5 w-20 bg-slate-200 rounded-full" /></td>
    <td className="px-6 py-4"><div className="h-3 w-24 bg-slate-200 rounded" /></td>
    <td className="px-6 py-4 text-right"><div className="h-4 w-4 bg-slate-200 rounded" /></td>
  </tr>
);

export default function UsersTable() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);

  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);

  const [search, setSearch] = useState("");
  const [roleFilter, setRoleFilter] = useState("all");
  const [filterOpen, setFilterOpen] = useState(false);

  const filterRef = useRef(null);
  const pageSize = 3;

  /* =========================
     FETCH USERS
  ========================= */
  useEffect(() => {
    let alive = true;

    const fetchUsers = async () => {
      setLoading(true);
      try {
        const res = await getUsers(page);
        if (!alive) return;

        const mapped = (res.data.users || []).map((u) => ({
          ...u,
          role: u.role?.toLowerCase(),
          status: (u.status || "pending").toLowerCase(),
          created_at: u.created_at || u.createdAt || null,
        }));

        setUsers(mapped);
        setTotal(res.data.total || 0);
      } catch (e) {
        console.error(e);
      } finally {
        setLoading(false);
      }
    };

    fetchUsers();
    return () => (alive = false);
  }, [page]);

  /* =========================
     CLOSE FILTER ON OUTSIDE CLICK
  ========================= */
  useEffect(() => {
    const close = (e) => {
      if (filterRef.current && !filterRef.current.contains(e.target)) {
        setFilterOpen(false);
      }
    };
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, []);

  /* =========================
     FILTER
  ========================= */
  const filteredUsers = useMemo(() => {
    const q = search.toLowerCase();
    return users.filter(
      (u) =>
        (u.name?.toLowerCase().includes(q) ||
          u.email?.toLowerCase().includes(q)) &&
        (roleFilter === "all" || u.role === roleFilter)
    );
  }, [users, search, roleFilter]);

  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="mt-8 rounded-2xl bg-white border border-none shadow overflow-hidden">

      {/* SEARCH + FILTER */}
      <div className="flex items-center justify-between gap-3 p-4 border-b border-none bg-white">
        <div className="relative w-full max-w-md">
          <Search className="absolute left-3 top-3.5 h-4 w-4 text-slate-400" />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full rounded-xl bg-slate-50 px-4 py-3 pl-9 text-sm"
            placeholder="Search users..."
          />
        </div>

        <div className="relative" ref={filterRef}>
          <button
            onClick={() => setFilterOpen(!filterOpen)}
            className="flex items-center gap-2 rounded-xl border px-4 py-2 text-sm"
          >
            <Filter className="h-4 w-4 border-none"  />
            {roleFilter === "all" ? "Filter" : roleFilter}
          </button>

          {filterOpen && (
            <div className="absolute right-0 mt-2 w-44 rounded-xl border bg-white shadow-lg z-50">
              {["all", "support", "customer"].map((r) => (
                <button
                  key={r}
                  onClick={() => {
                    setRoleFilter(r);
                    setPage(1);
                    setFilterOpen(false);
                  }}
                  className="block w-full px-4 py-2 text-left capitalize border-none hover:bg-slate-50"
                >
                  {r === "all" ? "All Roles" : r}
                </button>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* DESKTOP TABLE */}
  
      <div className="hidden lg:block overflow-x-auto">
        <table className="w-full text-sm border border-none">
          <thead className="bg-slate-50 border-b border-slate-200 ">
            <tr>
              <th className="px-6 py-4 text-left">User</th>
              <th className="px-6 py-4 text-left">Role</th>
              <th className="px-6 py-4 text-left">Status</th>
              <th className="px-6 py-4 text-left">Created At</th>
              <th className="px-6 py-4 text-right">Actions</th>
            </tr>
          </thead>

          <tbody>
            {loading
              ? Array.from({ length: 3 }).map((_, i) => <SkeletonRow key={i} />)
              : filteredUsers.map((u) => (
                  <tr key={u.id} className="border-slate-300 hover:bg-slate-50">
                    <td className="px-6 py-4 flex gap-3 items-center">
                      <div className="h-10 w-10 rounded-full bg-blue-600 text-white flex items-center justify-center">
                        {getInitials(u.name)}
                      </div>
                      <div>
                        <p className="font-medium">{u.name}</p>
                        <p className="text-xs text-slate-500">{u.email}</p>
                      </div>
                    </td>
                    <td className="px-6 py-4 capitalize">{u.role}</td>
                    <td className="px-6 py-4">
                      <span className={`px-3 py-1 rounded-full text-xs ${STATUS_STYLE[u.status]}`}>
                        {u.status}
                      </span>
                    </td>
                    <td className="px-6 py-4">{formatDate(u.created_at)}</td>
                    <td className="px-6 py-4 text-right">
                      <MoreHorizontal className="h-5 w-5 text-slate-400" />
                    </td>
                  </tr>
                ))}
          </tbody>
        </table>
      </div>

      {/* MOBILE CARDS */}
      <div className="lg:hidden divide-y border-t border-slate-200">
        {loading
          ? Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="p-4 animate-pulse h-20 bg-slate-100" />
            ))
          : filteredUsers.map((u) => (
              <div key={u.id} className="p-4 border-slate-300 flex justify-between">
                <div>
                  <p className="font-medium">{u.name}</p>
                  <p className="text-xs">{u.email}</p>
                   <p className="text-xs">{u.role}</p>
                  <p className="text-xs mt-1">{formatDate(u.created_at)}</p>
                </div>
                <MoreHorizontal className="h-5 w-5 text-slate-400" />
              </div>
            ))}
      </div>

      {/* PAGINATION */}
      <div className="flex justify-between p-4 text-sm border-t">
        <span>
          Showing {(page - 1) * pageSize + 1}–
          {Math.min(page * pageSize, total)} of {total}
        </span>
        <div className="flex gap-2">
          {Array.from({ length: totalPages }).map((_, i) => (
            <button
              key={i}
              onClick={() => setPage(i + 1)}
              className={`h-8 w-8 rounded border ${
                page === i + 1 ? "bg-blue-600 text-white" : ""
              }`}
            >
              {i + 1}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
