import { useEffect, useState } from "react";
import { createTicket } from "../api/ticket.api";
import api from "../api/axios";

export default function CreateTicket() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
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
      // Find the selected product to see if we have AMC data (optional optimization)
      const selectedProduct = products.find(p => p.id === productId);
      const amcId = selectedProduct?.amc_id; // assuming property exists, else undefined

      await createTicket({
        title: title.trim(),
        description: description.trim(),
        product_id: productId,
        amc_id: amcId
      });

      alert("Ticket created successfully");
      setTitle("");
      setDescription("");
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
        placeholder="Issue Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        required
      />

      {/* DESCRIPTION */}
      <textarea
        className="w-full border rounded px-3 py-2 mb-3 h-24"
        placeholder="Detailed Description of the problem..."
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        required
      />

      {/* PRODUCT */}
      <label className="block text-sm text-gray-600 mb-1">Select Product</label>
      <select
        className="w-full border rounded px-3 py-2 mb-3"
        value={productId}
        onChange={(e) => setProductId(e.target.value)}
        required
      >
        <option value="">-- Choose Product --</option>
        {products.map((p) => (
          <option key={p.id} value={p.id}>
            {p.name}
          </option>
        ))}
      </select>

      <button
        disabled={loading}
        className="bg-blue-600 text-white px-4 py-2 rounded disabled:opacity-60 w-full font-medium"
      >
        {loading ? "Creating Ticket..." : "Submit Ticket"}
      </button>
    </form>
  );
}
