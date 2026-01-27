import { useEffect, useState } from "react";
import api from "../api/axios";
import { createProduct } from "../api/product.api";

export default function AdminCreateProductAllInOne() {
  /* =========================
     STATE
  ========================= */
  const [categories, setCategories] = useState([]);
  const [brands, setBrands] = useState([]);
  const [models, setModels] = useState([]);

  const [categoryId, setCategoryId] = useState("");
  const [brandId, setBrandId] = useState("");
  const [modelId, setModelId] = useState("");

  const [newCategory, setNewCategory] = useState("");
  const [newBrand, setNewBrand] = useState("");
  const [newModel, setNewModel] = useState("");

  const [productName, setProductName] = useState("");
  const [loading, setLoading] = useState(false);

  /* =========================
     LOAD CATEGORIES
  ========================= */
  useEffect(() => {
    api.get("/categories").then((res) => {
      setCategories(res.data || []);
    });
  }, []);

  /* =========================
     LOAD BRANDS (CATEGORY)
  ========================= */
  useEffect(() => {
    if (!categoryId) return;

    api
      .get(`/categories/${categoryId}/brands`)
      .then((res) => setBrands(res.data || []));
  }, [categoryId]);

  /* =========================
     LOAD MODELS (BRAND)
  ========================= */
  useEffect(() => {
    if (!brandId) return;

    api
      .get(`/brands/${brandId}/models`)
      .then((res) => setModels(res.data || []));
  }, [brandId]);

  /* =========================
     CREATE LOOKUPS INLINE
  ========================= */
  const createCategory = async () => {
    if (!newCategory.trim()) return;
    const res = await api.post("/categories", { name: newCategory });
    setCategories([...categories, res.data]);
    setCategoryId(res.data.id);
    setNewCategory("");
  };

  const createBrand = async () => {
    if (!newBrand.trim() || !categoryId) return;
    const res = await api.post("/brands", {
      name: newBrand,
      category_id: categoryId,
    });
    setBrands([...brands, res.data]);
    setBrandId(res.data.id);
    setNewBrand("");
  };

  const createModel = async () => {
    if (!newModel.trim() || !brandId) return;
    const res = await api.post("/models", {
      name: newModel,
      brand_id: brandId,
    });
    setModels([...models, res.data]);
    setModelId(res.data.id);
    setNewModel("");
  };

  /* =========================
     CREATE PRODUCT
  ========================= */
  const submit = async (e) => {
    e.preventDefault();

    if (!categoryId || !brandId || !modelId || !productName.trim()) {
      alert("Please complete all steps");
      return;
    }

    setLoading(true);
    try {
      await createProduct({
        name: productName,
        category_id: categoryId,
        brand_id: brandId,
        model_id: modelId,
      });

      alert("Product created successfully");

      // reset only product name
      setProductName("");
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
      <h1 className="text-2xl font-semibold">
        Create Product (All-in-One)
      </h1>

      {/* CATEGORY */}
      <div>
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

        <div className="flex gap-2 mt-2">
          <input
            className="border rounded px-3 py-1 flex-1"
            placeholder="Add new category"
            value={newCategory}
            onChange={(e) => setNewCategory(e.target.value)}
          />
          <button
            type="button"
            onClick={createCategory}
            className="px-3 py-1 bg-gray-100 rounded"
          >
            + Add
          </button>
        </div>
      </div>

      {/* BRAND */}
      <div>
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

        <div className="flex gap-2 mt-2">
          <input
            className="border rounded px-3 py-1 flex-1"
            placeholder="Add new brand"
            value={newBrand}
            onChange={(e) => setNewBrand(e.target.value)}
            disabled={!categoryId}
          />
          <button
            type="button"
            onClick={createBrand}
            className="px-3 py-1 bg-gray-100 rounded"
            disabled={!categoryId}
          >
            + Add
          </button>
        </div>
      </div>

      {/* MODEL */}
      <div>
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

        <div className="flex gap-2 mt-2">
          <input
            className="border rounded px-3 py-1 flex-1"
            placeholder="Add new model"
            value={newModel}
            onChange={(e) => setNewModel(e.target.value)}
            disabled={!brandId}
          />
          <button
            type="button"
            onClick={createModel}
            className="px-3 py-1 bg-gray-100 rounded"
            disabled={!brandId}
          >
            + Add
          </button>
        </div>
      </div>

      {/* PRODUCT NAME */}
      <input
        className="w-full border rounded px-3 py-2"
        placeholder="Product name"
        value={productName}
        onChange={(e) => setProductName(e.target.value)}
        disabled={!modelId}
      />

      <button
        disabled={loading}
        className="bg-blue-600 text-white px-4 py-2 rounded w-full disabled:opacity-60"
      >
        {loading ? "Creating..." : "Create Product"}
      </button>
    </form>
  );
}
