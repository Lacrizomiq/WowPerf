import axios, {
  AxiosError,
  InternalAxiosRequestConfig,
  AxiosResponse,
} from "axios";
import { authService } from "@/libs/authService";

interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

// Ajout des interfaces pour le typage des réponses
interface CSRFResponse {
  csrf_token: string;
}

interface CSRFErrorResponse {
  error: string;
}

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,
});

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: unknown) => void;
  reject: (reason?: unknown) => void;
}> = [];

// CSRF token storage avec type explicite
let csrfToken: string | null = null;

const processQueue = (error: unknown | null, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

// Function to get CSRF token with proper typing
const getCSRFToken = async (): Promise<string> => {
  if (csrfToken) return csrfToken;

  try {
    const response = await api.get<CSRFResponse>("/api/csrf-token");
    if (!response.data.csrf_token) {
      throw new Error("No CSRF token in response");
    }
    csrfToken = response.data.csrf_token;
    return response.data.csrf_token;
  } catch (error) {
    console.error("Failed to fetch CSRF token:", error);
    throw error;
  }
};

// Request interceptor
api.interceptors.request.use(
  async (config: CustomAxiosRequestConfig) => {
    // Don't need CSRF token for these requests
    if (
      config.method?.toLowerCase() === "get" ||
      config.method?.toLowerCase() === "head" ||
      config.method?.toLowerCase() === "options" ||
      config.url === "/api/csrf-token"
    ) {
      return config;
    }

    try {
      const token = await getCSRFToken();
      if (config.headers) {
        config.headers["X-CSRF-Token"] = token;
      }
    } catch (error) {
      console.error("Error setting CSRF token:", error);
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor with proper error typing
api.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error: AxiosError<CSRFErrorResponse>) => {
    const originalRequest = error.config as CustomAxiosRequestConfig;
    if (!originalRequest) {
      return Promise.reject(error);
    }

    // Handle CSRF errors
    if (
      error.response?.status === 403 &&
      error.response?.data?.error?.includes("CSRF")
    ) {
      csrfToken = null; // Reset CSRF token
      try {
        const token = await getCSRFToken();
        if (originalRequest.headers) {
          originalRequest.headers["X-CSRF-Token"] = token;
        }
        return api(originalRequest);
      } catch (csrfError) {
        return Promise.reject(csrfError);
      }
    }

    // Handle authentication errors
    if (error.response?.status === 401) {
      if (originalRequest.url === "/auth/refresh") {
        await authService.logout();
        window.location.href = "/login?expired=true";
        return Promise.reject(error);
      }

      if (originalRequest._retry) {
        return Promise.reject(error);
      }

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(() => api(originalRequest));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        await authService.refreshToken();
        processQueue(null);
        return api(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        await authService.logout();
        window.location.href = "/login?expired=true";
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

export default api;
