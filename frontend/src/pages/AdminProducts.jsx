import { useEffect, useState } from "react";
import api from "../api/axios";
import { createProduct } from "../api/product.api";

export default function AdminCreateProduct() {
  const [categories, setCategories] = useState([]);
  const [brands, setBrands] = useState([]);
  const [models, setModels] = useState([]);

  const [categoryId, setCategoryId] = useState("");
  const [brandId, setBrandId] = useState("");
  const [modelId, setModelId] = useState("");
  const [name, setName] = useState("");

  const [loading, setLoading] = useState(false);

  /* =========================
     LOAD CATEGORIES
  ========================= */
  useEffect(() => {
    api.get("/categories").then((res) => {
      setCategories(res.data);
    });
  }, []);

  /* =========================
     LOAD BRANDS (CATEGORY)
  ========================= */
  useEffect(() => {
    if (!categoryId) {
      setBrands([]);
      setBrandId("");
      return;
    }

    api
      .get(`/categories/${categoryId}/brands`)
      .then((res) => setBrands(res.data));
  }, [categoryId]);

  /* =========================
     LOAD MODELS (BRAND)
  ========================= */
  useEffect(() => {
    if (!brandId) {
      setModels([]);
      setModelId("");
      return;
    }

    api
      .get(`/brands/${brandId}/models`)
      .then((res) => setModels(res.data));
  }, [brandId]);

  /* =========================
     SUBMIT
  ========================= */
  const submit = async (e) => {
    e.preventDefault();

    if (!categoryId || !brandId || !modelId || !name.trim()) {
      alert("Please fill all fields");
      return;
    }

    setLoading(true);
    try {
      await createProduct({
        name: name.trim(),
        category_id: categoryId,
        brand_id: brandId,
        model_id: modelId,
      });

      alert("Product created successfully");

      // reset
      setName("");
      setCategoryId("");
      setBrandId("");
      setModelId("");
      setBrands([]);
      setModels([]);
    } catch (err) {
      alert(err.response?.data?.error || "Failed to create product");
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
      className="bg-white p-6 rounded-xl shadow max-w-xl space-y-4"
    >
      <h1 className="text-2xl font-semibold">Create Product</h1>

      {/* CATEGORY */}
      <select
        className="w-full border rounded px-3 py-2"
        value={categoryId}
        onChange={(e) => setCategoryId(e.target.value)}
      >
        <option value="">Select Category</option>
        {categories.map((c) => (
          <option key={c.id} value={c.id}>
            {c.name}
          </option>
        ))}
      </select>

      {/* BRAND */}
      <select
        className="w-full border rounded px-3 py-2"
        value={brandId}
        onChange={(e) => setBrandId(e.target.value)}
        disabled={!categoryId}
      >
        <option value="">Select Brand</option>
        {brands.map((b) => (
          <option key={b.id} value={b.id}>
            {b.name}
          </option>
        ))}
      </select>

      {/* PRODUCT NAME */}
      <input
        className="w-full border rounded px-3 py-2"
        placeholder="Product name"
        value={name}
        onChange={(e) => setName(e.target.value)}
        disabled={!brandId}
      />

      {/* MODEL */}
      <select
        className="w-full border rounded px-3 py-2"
        value={modelId}
        onChange={(e) => setModelId(e.target.value)}
        disabled={!brandId}
      >
        <option value="">Select Model</option>
        {models.map((m) => (
          <option key={m.id} value={m.id}>
            {m.name}
          </option>
        ))}
      </select>

      <button
        disabled={loading}
        className="bg-blue-600 text-white px-4 py-2 rounded disabled:opacity-60"
      >
        {loading ? "Creating..." : "Create Product"}
      </button>
    </form>
  );
}
