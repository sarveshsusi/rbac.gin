import api from "./axios";

export const getAdminAMCs = () =>
  api.get("/admin/amc");

export const getCustomerAMCs = () =>
  api.get("/customer/amc");
