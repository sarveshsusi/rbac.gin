import api from "./axios";

/* =========================
   ADMIN
========================= */

export const createProduct = (data) =>
  api.post("/admin/products", data);

export const getProducts = () =>
  api.get("/admin/products");


/* =========================
   CUSTOMER
========================= */
export const getCustomerProducts = () =>
  api.get("/customer/products");

export const assignProductToCustomer = (customerId, productId) =>
  api.post(`/admin/customers/${customerId}/products`, {
    product_id: productId,
  });
