// feedback.api.js
import api from "./axios";

export const submitFeedback = (ticketId, rating, comment) =>
  api.post(`/customer/tickets/${ticketId}/feedback`, {
    rating,
    comment,
  });
