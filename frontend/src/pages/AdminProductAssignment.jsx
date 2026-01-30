import { useEffect, useMemo, useState } from "react";
import api from "../api/axios";
import { assignProductToCustomer } from "../api/product.api";
import { useToast } from "../components/toast/usetoast";
import { Pencil } from "lucide-react"; // optional icon lib (or replace with ✏️)

export default function AdminProductAssignment() {
  const { showToast } = useToast();

  const [users, setUsers] = useState([]);
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);

  const [search, setSearch] = useState("");
  const [pending, setPending] = useState({});

  /* =======================
     LOAD DATA
  ======================= */
  useEffect(() => {
    (async () => {
      try {
        const [u, p] = await Promise.all([
          api.get("/admin/users?role=customer"), // Only fetch customers
          api.get("/admin/products"),
        ]);

        setUsers(u.data?.users || u.data?.data || u.data || []);
        setProducts(p.data?.products || p.data?.data || p.data || []);
      } catch {
        showToast("Failed to load data", "error");
      } finally {
        setLoading(false);
      }
    })();
  }, [showToast]);

  /* =======================
     FILTER
  ======================= */
  const filteredUsers = useMemo(() => {
    return users.filter(
      (u) =>
        u.email?.toLowerCase().includes(search.toLowerCase()) ||
        u.name?.toLowerCase().includes(search.toLowerCase())
    );
  }, [users, search]);

  /* =======================
     ASSIGN LOGIC
  ======================= */
  const addProduct = (userId, product) => {
    setPending((p) => {
      const list = p[userId] || [];
      if (list.some((x) => x.id === product.id)) return p;
      return { ...p, [userId]: [...list, product] };
    });
  };

  const removeProduct = (userId, productId) => {
    setPending((p) => ({
      ...p,
      [userId]: (p[userId] || []).filter((x) => x.id !== productId),
    }));
  };

  const saveAll = async () => {
    try {
      const tasks = [];
      for (const userId in pending) {
        for (const prod of pending[userId]) {
          tasks.push(assignProductToCustomer(userId, prod.id));
        }
      }

      if (!tasks.length) {
        showToast("No changes to save", "info");
        return;
      }

      await Promise.all(tasks);
      showToast("Changes saved successfully!", "success");
      setPending({});
    } catch {
      showToast("Failed to save changes", "error");
    }
  };

  if (loading) return <Skeleton />;

  return (
    <div className="max-w-7xl mx-auto px-4 py-8 space-y-6">
      {/* HEADER */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">
            Assign Products to Customers
          </h1>
          <p className="text-sm text-slate-500">
            Manage which products are accessible to each customer account.
          </p>
        </div>

        <button
          onClick={saveAll}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-xl shadow transition"
        >
          💾 Save All Changes
        </button>
      </div>

      {/* FILTER BAR */}
      <div className="flex flex-col md:flex-row gap-4">
        <input
          className="enterprise-input md:max-w-sm"
          placeholder="Search customers..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      {/* TABLE */}
      <div className="bg-white border rounded-2xl overflow-hidden">
        {/* TABLE HEADER */}
        <div className="hidden md:grid grid-cols-[2fr_3fr_80px] gap-4 px-6 py-4 bg-slate-50 text-xs font-semibold text-slate-500">
          <div>CUSTOMER INFO</div>
          <div>ASSIGNED PRODUCTS</div>
          <div className="text-right">ACTIONS</div>
        </div>

        {/* ROWS */}
        {filteredUsers.map((u) => (
          <Row
            key={u.id}
            user={u}
            products={products}
            assigned={pending[u.id] || []}
            onAdd={addProduct}
            onRemove={removeProduct}
          />
        ))}

        {/* FOOTER */}
        <div className="px-6 py-4 text-sm text-slate-500">
          Showing 1 to {filteredUsers.length} of {users.length} customers
        </div>
      </div>
    </div>
  );
}

/* =======================
   ROW
======================= */

function Row({ user, products, assigned, onAdd, onRemove }) {
  const [open, setOpen] = useState(false);

  const initials =
    (user.name || user.email || "U")
      .split(" ")
      .map((w) => w[0])
      .join("")
      .slice(0, 2)
      .toUpperCase();

  return (
    <div className="border-t px-4 md:px-6 py-4 grid grid-cols-1 md:grid-cols-[2fr_3fr_80px] gap-4 items-center">
      {/* CUSTOMER */}
      <div className="flex items-center gap-3">
        <div className="w-10 h-10 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-semibold">
          {initials}
        </div>
        <div>
          <div className="font-medium">{user.name || "Customer"}</div>
          <div className="text-sm text-slate-500">{user.email}</div>
        </div>
      </div>

      {/* PRODUCTS */}
      <div className="flex flex-wrap gap-2 items-center">
        {assigned.map((p) => (
          <span
            key={p.id}
            className="flex items-center gap-2 bg-slate-100 px-3 py-1 rounded-full text-sm"
          >
            {p.name}
            <button
              onClick={() => onRemove(user.id, p.id)}
              className="text-slate-400 hover:text-red-500"
            >
              ×
            </button>
          </span>
        ))}

        <button
          onClick={() => setOpen((v) => !v)}
          className="text-blue-600 text-sm font-medium"
        >
          + Add Product
        </button>

        {open && (
          <select
            className="enterprise-input max-w-xs"
            defaultValue=""
            onChange={(e) => {
              const prod = products.find((p) => p.id === e.target.value);
              if (prod) onAdd(user.id, prod);
              setOpen(false);
            }}
          >
            <option value="">Select product</option>
            {products.map((p) => (
              <option key={p.id} value={p.id}>
                {p.name}
              </option>
            ))}
          </select>
        )}
      </div>

      {/* ACTION */}
      <div className="flex justify-end">
        <button className="text-slate-400 hover:text-slate-600">
          <Pencil size={18} />
        </button>
      </div>
    </div>
  );
}

/* =======================
   SKELETON
======================= */

function Skeleton() {
  return (
    <div className="p-6 space-y-4">
      {[1, 2, 3].map((i) => (
        <div
          key={i}
          className="h-20 bg-slate-100 rounded-xl animate-pulse"
        />
      ))}
    </div>
  );
}
