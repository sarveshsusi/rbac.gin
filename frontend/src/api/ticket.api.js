import api from "./axios";

/* =====================
   CUSTOMER
===================== */
export const createTicket = (data) =>
  api.post("/customer/tickets", data);

export const getCustomerTickets = () =>
  api.get("/customer/tickets");

export const submitFeedback = (ticketId, data) =>
  api.post(`/customer/tickets/${ticketId}/feedback`, data);

/* =====================
   ADMIN
===================== */
export const getAdminTickets = () =>
  api.get("/admin/tickets");

export const adminCreateTicket = (data) =>
  api.post("/admin/tickets", data);

export const assignTicket = (ticketId, data) =>
  api.post(`/admin/tickets/${ticketId}/assign`, data); 
  // data: { engineer_id, priority, support_mode, service_call_type }

/* =====================
   SUPPORT
===================== */
export const getSupportTickets = () =>
  api.get("/support/tickets");

export const startTicket = (ticketId) =>
  api.post(`/support/tickets/${ticketId}/start`);

export const closeTicket = (ticketId, formData) =>
  api.post(`/support/tickets/${ticketId}/close`, formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
