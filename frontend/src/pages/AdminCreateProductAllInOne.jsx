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
import { useToast } from "../components/toast/usetoast";

export default function AdminCreateProductAllInOne() {
  const { showToast } = useToast();

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
  const [selectionLocked, setSelectionLocked] = useState(false);

  // 🔥 Loading states for skeletons
  const [loadingCategories, setLoadingCategories] = useState(true);
  const [loadingBrands, setLoadingBrands] = useState(false);
  const [loadingModels, setLoadingModels] = useState(false);

  /* =======================
     LOAD CATEGORIES
  ======================= */
useEffect(() => {
  setLoadingCategories(true);
  getCategories()
    .then((r) => setCategories(r.data || []))
    .catch(() => showToast("Failed to load categories", "error"))
    .finally(() => setLoadingCategories(false));
}, [showToast]);


  /* =======================
     LOAD BRANDS
  ======================= */
useEffect(() => {
  if (!categoryId) {
    setBrands([]);
    return;
  }

  setLoadingBrands(true);
  getBrandsByCategory(categoryId)
    .then((r) => setBrands(r.data || []))
    .catch(() => showToast("Failed to load brands", "error"))
    .finally(() => setLoadingBrands(false));
}, [categoryId, showToast]);


  /* =======================
     LOAD MODELS
  ======================= */
useEffect(() => {
  if (!brandId) {
    setModels([]);
    return;
  }

  setLoadingModels(true);
  getModelsByBrand(brandId)
    .then((r) => setModels(r.data || []))
    .catch(() => showToast("Failed to load models", "error"))
    .finally(() => setLoadingModels(false));
}, [brandId, showToast]);


  /* =======================
     CREATE LOOKUPS
  ======================= */
  const addCategory = async () => {
    if (!newCategory.trim()) {
      showToast("Category name required", "error");
      return;
    }
    const res = await createCategory(newCategory.trim());
    setCategories((p) => [...p, res.data]);
    setCategoryId(res.data.id);
    setNewCategory("");
    showToast("Category created");
  };

  const addBrand = async () => {
    if (!categoryId) {
      showToast("Select category first", "info");
      return;
    }
    const res = await createBrand(newBrand.trim(), categoryId);
    setBrands((p) => [...p, res.data]);
    setBrandId(res.data.id);
    setNewBrand("");
    showToast("Brand created");
  };

  const addModel = async () => {
    if (!brandId) {
      showToast("Select brand first", "info");
      return;
    }
    const res = await createModel(newModel.trim(), brandId);
    setModels((p) => [...p, res.data]);
    setModelId(res.data.id);
    setNewModel("");
    showToast("Model created");
  };

  /* =======================
     SUBMIT PRODUCT
  ======================= */
  const submit = async (e) => {
    e.preventDefault();

    if (!selectionLocked || !productName.trim()) {
      showToast("Enter product name", "error");
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

      showToast("Product created successfully 🎉", "success");

      // 🔥 RESET FLOW
      setProductName("");
      setCategoryId("");
      setBrandId("");
      setModelId("");
      setBrands([]);
      setModels([]);
      setSelectionLocked(false);
    } catch {
      showToast("Failed to create product", "error");
    } finally {
      setLoading(false);
    }
  };

  const progress = selectionLocked ? 100 : 50;

  return (
    <div className="max-w-6xl mx-auto px-4 py-8 space-y-8">
      {/* PROGRESS BAR */}
      <div className="h-2 bg-slate-200 rounded-full">
        <div
          className="h-full bg-blue-600 transition-all duration-500"
          style={{ width: `${progress}%` }}
        />
      </div>

      <Stepper activeStep={selectionLocked ? 2 : 1} />

      {/* STEP 1 */}
      <EnterpriseCard title="Product Selection">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Select
            value={categoryId}
            onChange={setCategoryId}
            options={categories}
            placeholder="Category"
            loading={loadingCategories}
          />
          <Select
            value={brandId}
            onChange={setBrandId}
            options={brands}
            placeholder="Brand"
            disabled={!categoryId}
            loading={loadingBrands}
          />
          <Select
            value={modelId}
            onChange={setModelId}
            options={models}
            placeholder="Model"
            disabled={!brandId}
            loading={loadingModels}
          />

          <button
            disabled={!categoryId || !brandId || !modelId}
            onClick={() => {
              setSelectionLocked(true);
              showToast("Selection locked", "info");
            }}
            className="bg-blue-600 text-white rounded-xl px-6 py-3"
          >
            Select
          </button>
        </div>

        {!selectionLocked && (
          <div className="mt-6 space-y-4">
            <InlineAdd value={newCategory} onChange={setNewCategory} onAdd={addCategory} placeholder="New category" />
            <InlineAdd value={newBrand} onChange={setNewBrand} onAdd={addBrand} placeholder="New brand" disabled={!categoryId} />
            <InlineAdd value={newModel} onChange={setNewModel} onAdd={addModel} placeholder="New model" disabled={!brandId} />
          </div>
        )}
      </EnterpriseCard>

      {/* STEP 2 */}
      <EnterpriseCard title="Create Product">
        <form onSubmit={submit} className="space-y-4">
          <input
            disabled={!selectionLocked}
            value={productName}
            onChange={(e) => setProductName(e.target.value)}
            placeholder="Product name"
            className="enterprise-input"
          />
          <button
            disabled={loading}
            className="bg-blue-600 text-white px-10 py-3 rounded-xl w-full md:w-auto"
          >
            {loading ? "Creating..." : "Create Product"}
          </button>
        </form>
      </EnterpriseCard>
    </div>
  );
}

/* =======================
   UI COMPONENTS
======================= */

const Stepper = ({ activeStep }) => (
  <div className="flex items-center gap-6">
    <Step active={activeStep >= 1} label="Select" />
    <div className="flex-1 h-px bg-slate-300" />
    <Step active={activeStep >= 2} label="Create" />
  </div>
);

const Step = ({ active, label }) => (
  <div className="flex items-center gap-2">
    <div className={`w-8 h-8 rounded-full flex items-center justify-center ${active ? "bg-blue-600 text-white" : "bg-slate-300"}`}>
      ✓
    </div>
    <span className="text-sm">{label}</span>
  </div>
);

const EnterpriseCard = ({ title, children }) => (
  <div className="bg-white rounded-2xl border border-slate-300 p-6 space-y-6">
    <h2 className="font-semibold">{title}</h2>
    {children}
  </div>
);

const Select = ({ value, onChange, options, placeholder, disabled, loading }) => (
  <select
    disabled={disabled || loading}
    value={value}
    onChange={(e) => onChange(e.target.value)}
    className="enterprise-input"
  >
    <option value="">
      {loading ? "Loading..." : placeholder}
    </option>
    {!loading &&
      options.map((o) => (
        <option key={o.id} value={o.id}>
          {o.name}
        </option>
      ))}
  </select>
);

const InlineAdd = ({ value, onChange, onAdd, placeholder, disabled }) => (
  <div className="flex gap-3">
    <input
      disabled={disabled}
      value={value}
      onChange={(e) => onChange(e.target.value)}
      placeholder={placeholder}
      className="enterprise-input flex-1"
    />
    <button
      onClick={onAdd}
      disabled={disabled}
      className="text-blue-600 text-sm"
    >
      + Add
    </button>
  </div>
);
