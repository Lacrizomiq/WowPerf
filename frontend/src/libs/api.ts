import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";

// Types definitions
interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

interface CSRFResponse {
  token: string;
}

interface CSRFErrorResponse {
  error: string;
  code: string;
}

// Create axios instance with base configuration
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
});

// CSRF token management
let csrfToken: string | null = null;
let tokenPromise: Promise<string> | null = null;

// Get CSRF token with caching
const getCSRFToken = async (): Promise<string> => {
  // Return existing token if available
  if (csrfToken) {
    return csrfToken;
  }

  // Return pending request if exists
  if (tokenPromise) {
    return tokenPromise;
  }

  // Create new token request
  tokenPromise = api
    .get<CSRFResponse>("/api/csrf-token")
    .then((response) => {
      if (!response.data.token) {
        throw new Error("No token received from server");
      }
      csrfToken = response.data.token;
      tokenPromise = null;
      return response.data.token;
    })
    .catch((error) => {
      tokenPromise = null;
      throw error;
    });

  return tokenPromise;
};

// Request interceptor to add CSRF token
api.interceptors.request.use(async (config: CustomAxiosRequestConfig) => {
  const method = config.method?.toLowerCase();

  // Skip CSRF for safe methods
  if (["get", "head", "options"].includes(method || "")) {
    return config;
  }

  try {
    const token = await getCSRFToken();
    // Ensure headers object exists
    config.headers = config.headers || {};
    config.headers["X-CSRF-Token"] = token;
    return config;
  } catch (error) {
    console.error("Failed to fetch CSRF token:", error);
    return Promise.reject(error);
  }
});

// Response interceptor to handle CSRF errors
api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<CSRFErrorResponse>) => {
    const config = error.config as CustomAxiosRequestConfig;

    // Handle CSRF validation failures
    if (
      error.response?.status === 403 &&
      error.response?.data?.code === "INVALID_CSRF_TOKEN" &&
      !config._retry
    ) {
      config._retry = true;
      csrfToken = null; // Reset token
      return api(config); // Retry request
    }

    return Promise.reject(error);
  }
);

// Utility functions
export const resetCSRFToken = (): void => {
  csrfToken = null;
  tokenPromise = null;
};

// Debug logging for development
if (process.env.NODE_ENV === "development") {
  api.interceptors.request.use((request) => {
    console.log("🔄 Request:", {
      url: request.url,
      method: request.method,
      headers: request.headers,
    });
    return request;
  });

  api.interceptors.response.use(
    (response) => {
      console.log("✅ Response:", {
        status: response.status,
        headers: response.headers,
        data: response.data,
      });
      return response;
    },
    (error) => {
      console.log("❌ Response Error:", {
        status: error.response?.status,
        data: error.response?.data,
      });
      return Promise.reject(error);
    }
  );
}

export default api;
