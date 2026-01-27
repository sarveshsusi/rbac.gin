import { useEffect, useState } from "react";
import api from "../api/axios";
import { assignProductToCustomer } from "../api/product.api";

export default function AdminProductAssignment() {
  const [users, setUsers] = useState([]);
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const [usersRes, productsRes] = await Promise.all([
          api.get("/admin/users"),
          api.get("/admin/products"),
        ]);

        // ✅ FIX: normalize users response
        const usersData = Array.isArray(usersRes.data)
          ? usersRes.data
          : usersRes.data.users || usersRes.data.data || [];

        const productsData = Array.isArray(productsRes.data)
          ? productsRes.data
          : productsRes.data.products || productsRes.data.data || [];

        setUsers(usersData);
        setProducts(productsData);
      } catch (err) {
        console.error("Failed to load data", err);
        setUsers([]);
        setProducts([]);
      } finally {
        setLoading(false);
      }
    };

    load();
  }, []);

  const assign = async (userId, productId) => {
    if (!productId) return;

    try {
      await assignProductToCustomer(userId, productId);
      alert("Product assigned");
    } catch (err) {
      alert(err.response?.data?.error || "Assignment failed");
    }
  };

  if (loading) {
    return <div className="text-gray-500">Loading...</div>;
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Assign Products</h1>

      {users.length === 0 && (
        <div className="text-gray-500">No users found</div>
      )}

      {users.map((u) => (
        <div
          key={u.id}
          className="bg-white border rounded-xl p-4 space-y-2"
        >
          <div className="font-medium">{u.email}</div>

          <select
            onChange={(e) => assign(u.id, e.target.value)}
            className="border rounded px-3 py-2 w-full"
            defaultValue=""
          >
            <option value="">Assign product</option>

            {products.map((p) => (
              <option key={p.id} value={p.id}>
                {p.name}
              </option>
            ))}
          </select>
        </div>
      ))}
    </div>
  );
}
