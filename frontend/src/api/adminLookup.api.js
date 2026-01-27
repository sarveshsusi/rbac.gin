import api from "./axios";

/* CATEGORY */
export const getCategories = () =>
  api.get("/admin/categories");

export const createCategory = (name) =>
  api.post("/admin/categories", { name });

/* BRAND */
export const getBrandsByCategory = (categoryId) =>
  api.get(`/admin/categories/${categoryId}/brands`);

export const createBrand = (name, categoryId) =>
  api.post("/admin/brands", {
    name,
    category_id: categoryId,
  });


/* MODEL */
export const getModelsByBrand = (brandId) =>
  api.get(`/admin/brands/${brandId}/models`);

export const createModel = (name, brandId) =>
  api.post("/admin/models", { name, brand_id: brandId });
