// src/libs/api.ts
import axios, { AxiosError } from "axios";
import { csrfService } from "./csrfService";

export interface APIError {
  error: string;
  code: string;
  details?: string;
}

export interface ApiResponse {
  message: string;
  code: string;
  error?: string;
}

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "https://test.wowperf.com/api";

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
    "X-Requested-With": "XMLHttpRequest",
  },
});

// Interceptor for requests
api.interceptors.request.use(
  async (config) => {
    // Skip CSRF for OAuth routes
    if (config.headers["X-Skip-CSRF"]) {
      return config;
    }

    if (
      config.method &&
      config.method.toLowerCase() !== "get" &&
      config.url &&
      csrfService.isProtectedRoute(config.url, config.method)
    ) {
      try {
        const token = await csrfService.getToken();
        if (token) {
          config.headers["X-CSRF-Token"] = token;
        }
      } catch (error) {
        console.error("Error setting CSRF token:", error);
      }
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Interceptor for responses
api.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    console.error("Response error:", {
      message: error.message,
      status: error?.response?.status,
      data: error?.response?.data,
      config: {
        url: error?.config?.url,
        method: error?.config?.method,
        data: error?.config?.data,
      },
    });

    const originalRequest = error.config;

    // Specific handling of CSRF errors
    if (
      csrfService.shouldResetToken(error) &&
      !originalRequest._retry &&
      csrfService.isProtectedRoute(originalRequest.url, originalRequest.method)
    ) {
      originalRequest._retry = true;

      try {
        // Force the token refresh
        const token = await csrfService.getToken(true);
        if (token) {
          originalRequest.headers["X-CSRF-Token"] = token;
          return api(originalRequest);
        }
      } catch (refreshError) {
        console.error("Failed to refresh CSRF token:", refreshError);
        csrfService.clearToken();
        return Promise.reject({
          ...error,
          message: "CSRF token refresh failed",
        });
      }
    }

    // Handling other errors
    if (error.response?.status === 401) {
      // Session expired or not authenticated
      csrfService.clearToken();
    }

    return Promise.reject(error);
  }
);

export default api;

// Utility functions for CSRF management
export const resetCSRFToken = () => csrfService.clearToken();
export const preloadCSRFToken = () => csrfService.getToken();

// Utility function to check if a route is protected
export const isProtectedRoute = (url: string, method: string) =>
  csrfService.isProtectedRoute(url, method);
