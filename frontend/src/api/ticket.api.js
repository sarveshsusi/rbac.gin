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

export const assignTicket = (ticketId, engineerId) =>
  api.post(`/admin/tickets/${ticketId}/assign`, {
    engineer_id: engineerId,
  });

export const closeTicket = (ticketId) =>
  api.post(`/admin/tickets/${ticketId}/close`);

/* =====================
   SUPPORT
===================== */
export const getSupportTickets = () =>
  api.get("/support/tickets");

export const resolveTicket = (ticketId) =>
  api.post(`/support/tickets/${ticketId}/resolve`);
