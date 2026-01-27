import { useEffect, useState } from "react";
import { createTicket } from "../api/ticket.api";
import api from "../api/axios";

export default function CreateTicket() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [priority, setPriority] = useState("low");
  const [productId, setProductId] = useState("");
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(false);

  /* =========================
     LOAD CUSTOMER PRODUCTS
  ========================= */
  useEffect(() => {
    api
      .get("/customer/products")
      .then((res) => setProducts(res.data))
      .catch(() => setProducts([]));
  }, []);

  /* =========================
     SUBMIT
  ========================= */
  const submit = async (e) => {
    e.preventDefault();

    if (!productId) {
      alert("Please select a product");
      return;
    }

    setLoading(true);
    try {
      await createTicket({
        title: title.trim(),
        description: description.trim(),
        priority,          // ✅ lowercase enum
        product_id: productId,
      });

      alert("Ticket created successfully");

      setTitle("");
      setDescription("");
      setPriority("low");
      setProductId("");
    } catch (err) {
      alert(err.response?.data?.error || "Ticket creation failed");
    } finally {
      setLoading(false);
    }
  };

  /* =========================
     UI
  ========================= */
  return (
    <form
      onSubmit={submit}
      className="bg-white p-6 rounded-xl shadow max-w-xl"
    >
      <h2 className="text-xl font-semibold mb-4">Create Ticket</h2>

      {/* TITLE */}
      <input
        className="w-full border rounded px-3 py-2 mb-3"
        placeholder="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        required
      />

      {/* DESCRIPTION */}
      <textarea
        className="w-full border rounded px-3 py-2 mb-3"
        placeholder="Description"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
      />

      {/* PRODUCT */}
      <select
        className="w-full border rounded px-3 py-2 mb-3"
        value={productId}
        onChange={(e) => setProductId(e.target.value)}
        required
      >
        <option value="">Select Product</option>
        {products.map((p) => (
          <option key={p.id} value={p.id}>
            {p.name}
          </option>
        ))}
      </select>

      {/* PRIORITY */}
      <select
        className="w-full border rounded px-3 py-2 mb-4"
        value={priority}
        onChange={(e) => setPriority(e.target.value)}
      >
        <option value="low">Low</option>
        <option value="medium">Medium</option>
        <option value="high">High</option>
        <option value="critical">Critical</option>
      </select>

      <button
        disabled={loading}
        className="bg-blue-600 text-white px-4 py-2 rounded disabled:opacity-60"
      >
        {loading ? "Creating..." : "Submit"}
      </button>
    </form>
  );
}
