import api from "./axios";

export const getUsers = (page = 1) =>
  api.get(`/admin/users?page=${page}`);
