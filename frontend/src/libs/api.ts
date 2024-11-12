import axios, {
  AxiosError,
  InternalAxiosRequestConfig,
  AxiosHeaders,
} from "axios";

// Custom Axios request config
interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
  _csrfRetry?: boolean;
}

// CSRF response
interface CSRFResponse {
  token: string; // CSRF token
  header: string; // Header name
}

// API error
interface APIError {
  error: string;
  code: string;
  details?: string;
}

class CSRFTokenManager {
  private static instance: CSRFTokenManager;
  private token: string | null = null;
  private tokenPromise: Promise<string> | null = null;

  private constructor() {}

  static getInstance(): CSRFTokenManager {
    if (!CSRFTokenManager.instance) {
      CSRFTokenManager.instance = new CSRFTokenManager();
    }
    return CSRFTokenManager.instance;
  }

  private async fetchToken(): Promise<string> {
    try {
      const response = await axios.get<CSRFResponse>(
        `${process.env.NEXT_PUBLIC_API_URL}/api/csrf-token`,
        {
          withCredentials: true,
        }
      );

      if (!response.data.token) {
        throw new Error("No CSRF token in response");
      }

      this.token = response.data.token;
      return this.token;
    } catch (error) {
      console.error("Failed to fetch CSRF token:", error);
      this.token = null;
      this.tokenPromise = null;
      throw error;
    }
  }

  async getToken(forceRefresh = false): Promise<string> {
    if (forceRefresh) {
      this.token = null;
      this.tokenPromise = null;
    }

    if (this.token) {
      return this.token;
    }

    if (this.tokenPromise) {
      return this.tokenPromise;
    }

    this.tokenPromise = this.fetchToken();
    const token = await this.tokenPromise;
    this.tokenPromise = null;
    return token;
  }

  clearToken(): void {
    this.token = null;
    this.tokenPromise = null;
  }
}

const csrfManager = CSRFTokenManager.getInstance();

// Create axios instance
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
});

// Request interceptor
api.interceptors.request.use(
  async (config: CustomAxiosRequestConfig) => {
    // Skip CSRF token for safe methods and token requests
    if (
      config.method?.toLowerCase() === "get" ||
      config.method?.toLowerCase() === "head" ||
      config.method?.toLowerCase() === "options" ||
      config.url === "/api/csrf-token"
    ) {
      return config;
    }

    try {
      // Log request details
      console.log("Request details:", {
        url: config.url,
        method: config.method,
        headers: config.headers,
        withCredentials: config.withCredentials,
      });

      const token = await csrfManager.getToken();
      console.log("CSRF Token details:", {
        exists: !!token,
        length: token?.length,
      });

      if (!token) {
        throw new Error("No CSRF token available");
      }

      const headers = new AxiosHeaders(config.headers);
      headers.set("X-CSRF-Token", token);
      config.headers = headers;

      // Log final request configuration
      console.log("Final request config:", {
        headers: config.headers,
        withCredentials: config.withCredentials,
      });

      return config;
    } catch (error) {
      console.error("Error in request interceptor:", error);
      throw error;
    }
  },
  (error) => {
    console.error("Request interceptor error:", error);
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<APIError>) => {
    const originalRequest = error.config as CustomAxiosRequestConfig;

    if (!originalRequest) {
      return Promise.reject(error);
    }

    // Detailed error log
    console.log(
      `API Error for ${originalRequest.method} ${originalRequest.url}:`,
      {
        status: error.response?.status,
        data: error.response?.data,
        config: {
          headers: originalRequest.headers,
          method: originalRequest.method,
          url: originalRequest.url,
        },
      }
    );

    // CSRF error handling
    if (
      error.response?.status === 403 &&
      error.response.data?.code === "INVALID_CSRF_TOKEN" &&
      !originalRequest._csrfRetry
    ) {
      console.log("CSRF validation failed, retrying with new token");
      originalRequest._csrfRetry = true;

      try {
        // Force the token refresh
        const newToken = await csrfManager.getToken(true);
        const headers = new AxiosHeaders(originalRequest.headers);
        headers.set("X-CSRF-Token", newToken);
        originalRequest.headers = headers;

        return api(originalRequest);
      } catch (retryError) {
        console.error(
          "Failed to retry request with new CSRF token:",
          retryError
        );
        csrfManager.clearToken();
        return Promise.reject(retryError);
      }
    }

    // Authentication error handling
    if (error.response?.status === 401) {
      // Clear CSRF token on authentication errors
      csrfManager.clearToken();
    }

    return Promise.reject(error);
  }
);

export default api;

// Utility function to reset the CSRF token
export const resetCSRFToken = () => {
  csrfManager.clearToken();
};

// Utility function to pre-load a CSRF token
export const preloadCSRFToken = () => {
  return csrfManager.getToken();
};

// Types for better usage
export type { APIError, CSRFResponse };
