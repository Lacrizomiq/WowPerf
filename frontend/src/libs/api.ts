// src/libs/api.ts
import axios from "axios";
import { csrfService } from "./csrfService";

export interface APIError {
  error: string;
  code: string;
  details?: string;
}

export interface ApiResponse {
  message: string;
  error?: string;
  code?: string;
}

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
    "X-Requested-With": "XMLHttpRequest",
  },
});

api.interceptors.request.use(
  async (config) => {
    if (config.url && config.method) {
      if (csrfService.isProtectedRoute(config.url, config.method)) {
        const token = await csrfService.getToken();
        if (token) {
          config.headers["X-CSRF-Token"] = token;
        }
      }
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (
      error.response?.status === 403 &&
      error.response.data?.code === "INVALID_CSRF_TOKEN" &&
      !originalRequest._retry &&
      csrfService.isProtectedRoute(originalRequest.url, originalRequest.method)
    ) {
      originalRequest._retry = true;
      const token = await csrfService.getToken(true);
      if (token) {
        originalRequest.headers["X-CSRF-Token"] = token;
        return api(originalRequest);
      }
    }

    return Promise.reject(error);
  }
);

export default api;

export const resetCSRFToken = () => csrfService.clearToken();
export const preloadCSRFToken = () => csrfService.getToken();
