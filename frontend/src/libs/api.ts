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
}

// API error
interface APIError {
  error: string;
  code: string;
  details?: string;
}

// Configuration base on the environment
const API_URL = process.env.NEXT_PUBLIC_API_URL;
const APP_URL = process.env.NEXT_PUBLIC_APP_URL;
const isLocalEnv = process.env.NODE_ENV === "development";

// Conditionnal logger
const logDebug = (message: string, data?: any) => {
  if (isLocalEnv) {
    console.log(`🔧 [API Debug] ${message}`, data || "");
  }
};

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

  // api.ts - dans la classe CSRFTokenManager
  private async fetchToken(): Promise<string> {
    try {
      const csrfAxios = axios.create({
        baseURL: API_URL,
        withCredentials: true,
        headers: {
          Accept: "application/json",
          "X-Requested-With": "XMLHttpRequest",
        },
        ...(isLocalEnv && {
          httpsAgent: new https.Agent({
            rejectUnauthorized: false,
            keepAlive: true,
          }),
        }),
      });

      const response = await csrfAxios.get<CSRFResponse>("/api/csrf-token");
      const headerToken = response.headers["x-csrf-token"];
      const dataToken = response.data.token;
      const token = headerToken || dataToken;

      if (!token) {
        throw new Error("No CSRF token in response");
      }

      this.token = token as string; // Cast explicite pour TypeScript
      this.tokenExpiryTime = Date.now() + this.TOKEN_LIFETIME;

      logDebug("CSRF Token fetched", {
        token: token.substring(0, 10) + "...",
        expiresIn: this.TOKEN_LIFETIME,
      });

      return this.token;
    } catch (error) {
      logDebug("Failed to fetch CSRF token:", error);
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
    logDebug("CSRF Token cleared");
  }
}

const csrfManager = CSRFTokenManager.getInstance();

// Create axios instance avec la configuration appropriée
const api = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
    "X-Requested-With": "XMLHttpRequest",
  },
  ...(isLocalEnv && {
    httpsAgent: new https.Agent({
      rejectUnauthorized: false,
      keepAlive: true,
    }),
  }),
});

// Initial log of the local config
if (isLocalEnv) {
  logDebug("API Configuration:", {
    apiUrl: API_URL,
    appUrl: APP_URL,
    environment: process.env.NODE_ENV,
  });
}

// Request interceptor
api.interceptors.request.use(
  async (config: CustomAxiosRequestConfig) => {
    if (
      config.method?.toLowerCase() === "get" ||
      config.method?.toLowerCase() === "head" ||
      config.method?.toLowerCase() === "options" ||
      config.url === "/api/csrf-token"
    ) {
      return config;
    }

    try {
      logDebug("Request:", {
        url: config.url,
        method: config.method,
      });

      const token = await csrfManager.getToken();
      if (!token) {
        throw new Error("No CSRF token available");
      }

      const headers = new AxiosHeaders(config.headers);
      headers.set("X-CSRF-Token", token);
      headers.set("X-Requested-With", "XMLHttpRequest");
      config.headers = headers;

      logDebug("Request headers set:", {
        csrf: token,
        url: config.url,
      });

      return config;
    } catch (error) {
      logDebug("Request Error:", error);
      throw error;
    }
  },
  (error) => {
    logDebug("Request Interceptor Error:", error);
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => {
    logDebug("Response Success:", {
      url: response.config.url,
      status: response.status,
    });
    return response;
  },
  async (error: AxiosError<APIError>) => {
    const originalRequest = error.config as CustomAxiosRequestConfig;

    if (!originalRequest) {
      return Promise.reject(error);
    }

    logDebug("Response Error:", {
      url: error.config?.url,
      status: error.response?.status,
      data: error.response?.data,
    });

    // CSRF error handling
    if (
      error.response?.status === 403 &&
      error.response.data?.code === "INVALID_CSRF_TOKEN" &&
      !originalRequest._csrfRetry
    ) {
      logDebug("CSRF validation failed, retrying with new token");
      originalRequest._csrfRetry = true;

      try {
        const newToken = await csrfManager.getToken(true);
        const headers = new AxiosHeaders(originalRequest.headers);
        headers.set("X-CSRF-Token", newToken);
        originalRequest.headers = headers;

        return api(originalRequest);
      } catch (retryError) {
        logDebug("Failed to retry request with new CSRF token:", retryError);
        csrfManager.clearToken();
        return Promise.reject(retryError);
      }
    }

    // Authentication error handling
    if (error.response?.status === 401) {
      logDebug("Authentication error detected");
      csrfManager.clearToken();
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

export const getApiConfig = () => ({
  apiUrl: API_URL,
  appUrl: APP_URL,
  isLocal: isLocalEnv,
});

// Types for better usage
export type { APIError, CSRFResponse };
