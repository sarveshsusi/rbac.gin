import { useEffect, useState } from "react";
import {
  getCategories,
  createCategory,
  getBrandsByCategory,
  createBrand,
  getModelsByBrand,
  createModel,
} from "../api/adminLookup.api";
import { createProduct } from "../api/product.api";

export default function AdminCreateProductAllInOne() {
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

  /* LOAD CATEGORIES */
  useEffect(() => {
    getCategories()
      .then((res) => setCategories(res.data || []))
      .catch(() => setCategories([]));
  }, []);

  /* LOAD BRANDS */
  useEffect(() => {
    if (!categoryId) {
      setBrands([]);
      setBrandId("");
      return;
    }

    getBrandsByCategory(categoryId)
      .then((res) => setBrands(res.data || []))
      .catch(() => setBrands([]));
  }, [categoryId]);

  /* LOAD MODELS */
  useEffect(() => {
    if (!brandId) {
      setModels([]);
      setModelId("");
      return;
    }

    getModelsByBrand(brandId)
      .then((res) => setModels(res.data || []))
      .catch(() => setModels([]));
  }, [brandId]);

  /* CREATE CATEGORY */
  const addCategory = async () => {
    if (!newCategory.trim()) return;

    const res = await createCategory(newCategory.trim());
    setCategories((prev) => [...prev, res.data]);
    setCategoryId(res.data.id);
    setNewCategory("");
  };

  /* CREATE BRAND */
  const addBrand = async () => {
    if (!newBrand.trim() || !categoryId) return;

    const res = await createBrand(newBrand.trim(), categoryId);
    setBrands((prev) => [...prev, res.data]);
    setBrandId(res.data.id);
    setNewBrand("");
  };

  /* CREATE MODEL */
  const addModel = async () => {
    if (!newModel.trim() || !brandId) return;

    const res = await createModel(newModel.trim(), brandId);
    setModels((prev) => [...prev, res.data]);
    setModelId(res.data.id);
    setNewModel("");
  };

  /* CREATE PRODUCT */
  const submit = async (e) => {
    e.preventDefault();

    if (!categoryId || !brandId || !modelId || !productName.trim()) {
      alert("Complete all steps");
      return;
    }

    setLoading(true);
    try {
      await createProduct({
        name: productName.trim(),
        category_id: categoryId,
        brand_id: brandId,
        model_id: modelId,
      });

      alert("Product created successfully");
      setProductName("");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form
      onSubmit={submit}
      className="bg-white p-6 rounded-xl shadow max-w-xl space-y-4"
    >
      <h1 className="text-2xl font-semibold">Create Product (All in One)</h1>

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

      <input
        className="border px-3 py-2 rounded"
        placeholder="New category"
        value={newCategory}
        onChange={(e) => setNewCategory(e.target.value)}
      />
      <button type="button" onClick={addCategory}>+ Add Category</button>

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

      <input
        className="border px-3 py-2 rounded"
        placeholder="New brand"
        value={newBrand}
        onChange={(e) => setNewBrand(e.target.value)}
        disabled={!categoryId}
      />
      <button type="button" onClick={addBrand} disabled={!categoryId}>
        + Add Brand
      </button>

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

      <input
        className="border px-3 py-2 rounded"
        placeholder="New model"
        value={newModel}
        onChange={(e) => setNewModel(e.target.value)}
        disabled={!brandId}
      />
      <button type="button" onClick={addModel} disabled={!brandId}>
        + Add Model
      </button>

      {/* PRODUCT */}
      <input
        className="w-full border rounded px-3 py-2"
        placeholder="Product name"
        value={productName}
        onChange={(e) => setProductName(e.target.value)}
        disabled={!modelId}
      />

      <button
        disabled={loading}
        className="bg-blue-600 text-white w-full py-2 rounded"
      >
        {loading ? "Creating..." : "Create Product"}
      </button>
    </form>
  );
}
