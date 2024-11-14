import axios, {
  AxiosError,
  InternalAxiosRequestConfig,
  AxiosHeaders,
} from "axios";
import https from "https";
// Custom Axios request config
interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
  _csrfRetry?: boolean;
}

// CSRF response
interface CSRFResponse {
  token: string;
  header: string;
}

// API error
interface APIError {
  error: string;
  code: string;
  details?: string;
}

// Create custom HTTPS agent
const httpsAgent = new https.Agent({
  rejectUnauthorized: false, // Important for self-signed certificates
});

class CSRFTokenManager {
  private static instance: CSRFTokenManager;
  private token: string | null = null;
  private tokenPromise: Promise<string> | null = null;
  private tokenExpiryTime: number | null = null;
  private readonly TOKEN_LIFETIME = 60 * 60 * 1000; // 1 hour in milliseconds

  private constructor() {}

  static getInstance(): CSRFTokenManager {
    if (!CSRFTokenManager.instance) {
      CSRFTokenManager.instance = new CSRFTokenManager();
    }
    return CSRFTokenManager.instance;
  }

  private isTokenExpired(): boolean {
    return !this.tokenExpiryTime || Date.now() > this.tokenExpiryTime;
  }

  private async fetchToken(): Promise<string> {
    try {
      // Create a specific Axios instance for the CSRF request
      const csrfAxios = axios.create({
        baseURL: process.env.NEXT_PUBLIC_API_URL,
        withCredentials: true,
        headers: {
          Accept: "application/json",
          "X-Requested-With": "XMLHttpRequest",
        },
        ...(process.env.NODE_ENV === "development" && {
          httpsAgent: new https.Agent({
            rejectUnauthorized: false,
            keepAlive: true,
          }),
        }),
      });

      const response = await csrfAxios.get<CSRFResponse>("/api/csrf-token");

      if (!response.data.token) {
        throw new Error("No CSRF token in response");
      }

      this.token = response.data.token;
      this.tokenExpiryTime = Date.now() + this.TOKEN_LIFETIME;
      return this.token;
    } catch (error) {
      console.error("Failed to fetch CSRF token:", error);
      this.token = null;
      this.tokenPromise = null;
      this.tokenExpiryTime = null;
      throw error;
    }
  }

  async getToken(forceRefresh = false): Promise<string> {
    if (forceRefresh || this.isTokenExpired()) {
      this.token = null;
      this.tokenPromise = null;
      this.tokenExpiryTime = null;
    }

    if (this.token && !this.isTokenExpired()) {
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
    this.tokenExpiryTime = null;
  }
}

const csrfManager = CSRFTokenManager.getInstance();

// Create axios instance
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true, // Important for cookies
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
    "X-Requested-With": "XMLHttpRequest",
  },
  ...(process.env.NODE_ENV === "development" && {
    httpsAgent: new https.Agent({
      rejectUnauthorized: false,
      keepAlive: true, // Keep the connection alive
    }),
  }),
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
      // Log request details in development
      if (process.env.NODE_ENV === "development") {
        console.log("Request details:", {
          url: config.url,
          method: config.method,
          headers: config.headers,
          withCredentials: config.withCredentials,
        });
      }

      const token = await csrfManager.getToken();
      if (!token) {
        throw new Error("No CSRF token available");
      }

      const headers = new AxiosHeaders(config.headers);
      headers.set("X-CSRF-Token", token);
      headers.set("X-Requested-With", "XMLHttpRequest");
      config.headers = headers;

      // Log final request configuration in development
      if (process.env.NODE_ENV === "development") {
        console.log("Final request config:", {
          headers: config.headers,
          withCredentials: config.withCredentials,
        });
      }

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
  (response) => {
    if (process.env.NODE_ENV === "development") {
      console.log(`✅ Response success:`, {
        status: response.status,
        data: response.data,
      });
    }
    return response;
  },
  async (error: AxiosError<APIError>) => {
    const originalRequest = error.config as CustomAxiosRequestConfig;

    if (!originalRequest) {
      return Promise.reject(error);
    }

    // Detailed error log
    if (process.env.NODE_ENV === "development") {
      console.log(
        `❌ API Error for ${originalRequest.method} ${originalRequest.url}:`,
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
    }

    // CSRF error handling
    if (
      error.response?.status === 403 &&
      error.response.data?.code === "INVALID_CSRF_TOKEN" &&
      !originalRequest._csrfRetry
    ) {
      console.log("CSRF validation failed, retrying with new token");
      originalRequest._csrfRetry = true;

      try {
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
      csrfManager.clearToken();
      // I can add here a redirection to the login page or other actions
    }

    return Promise.reject(error);
  }
);

export default api;

// Utility functions
export const resetCSRFToken = () => {
  csrfManager.clearToken();
};

export const preloadCSRFToken = () => {
  return csrfManager.getToken();
};

// Types for better usage
export type { APIError, CSRFResponse };
