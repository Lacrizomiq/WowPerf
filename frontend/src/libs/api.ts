import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";
import { authService } from "./authService";

interface CSRFResponse {
  csrf_token: string;
}

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

const refreshCSRFToken = async (): Promise<string> => {
  try {
    const response = await api.get<CSRFResponse>("/auth/csrf");
    const newToken = response.data.csrf_token;
    document
      .querySelector('meta[name="csrf-token"]')
      ?.setAttribute("content", newToken);
    return newToken;
  } catch (error) {
    console.error("Failed to refresh CSRF token:", error);
    throw error;
  }
};

api.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    const token = authService.getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    let csrfToken = document
      .querySelector('meta[name="csrf-token"]')
      ?.getAttribute("content");
    if (!csrfToken) {
      csrfToken = await refreshCSRFToken();
    }
    if (csrfToken) {
      config.headers["X-CSRF-Token"] = csrfToken;
    }

    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    if (
      error.response &&
      error.response.status === 403 &&
      error.response.data === "CSRF token mismatch"
    ) {
      const originalRequest = error.config;
      if (!originalRequest) {
        return Promise.reject(error);
      }

      try {
        const newToken = await refreshCSRFToken();

        if (originalRequest.headers) {
          originalRequest.headers["X-CSRF-Token"] = newToken;
        }
        return api(originalRequest);
      } catch (refreshError) {
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default api;
