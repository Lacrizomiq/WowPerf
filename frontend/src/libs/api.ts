import axios, {
  AxiosError,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from "axios";

// Type definitions
interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

interface CSRFResponse {
  token: string;
}

interface APIErrorResponse {
  error: string;
  code: string;
  details?: string;
  debug?: {
    origin?: string;
    referer?: string;
    hasToken?: boolean;
  };
}

// Custom error type
export class APIError extends Error {
  constructor(
    public code: string,
    message: string,
    public status?: number,
    public details?: string
  ) {
    super(message);
    this.name = "APIError";
  }
}

// API client instance
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
    "X-Requested-With": "XMLHttpRequest",
  },
  xsrfCookieName: "_csrf",
  xsrfHeaderName: "X-CSRF-Token",
});

// CSRF token management
let csrfToken: string | null = null;
let tokenPromise: Promise<string> | null = null;

const getCSRFToken = async (): Promise<string> => {
  if (csrfToken) {
    return csrfToken;
  }

  if (tokenPromise) {
    return tokenPromise;
  }

  tokenPromise = api
    .get<CSRFResponse>("/api/csrf-token")
    .then((response: AxiosResponse<CSRFResponse>) => {
      if (!response.data.token) {
        throw new Error("No CSRF token received from server");
      }
      csrfToken = response.data.token;
      console.log("🔑 CSRF token fetched:", csrfToken);
      tokenPromise = null;
      return csrfToken;
    })
    .catch((error: AxiosError) => {
      console.error("Failed to fetch CSRF token:", error);
      tokenPromise = null;
      throw error;
    });

  return tokenPromise;
};

// Request interceptor
api.interceptors.request.use(
  async (config: CustomAxiosRequestConfig) => {
    const method = config.method?.toLowerCase();

    // Skip CSRF for safe methods
    if (["get", "head", "options"].includes(method || "")) {
      return config;
    }

    try {
      const token = await getCSRFToken();

      // Ensure headers object exists and is properly typed
      config.headers = config.headers || {};
      config.headers["X-CSRF-Token"] = token;
      config.headers["X-Requested-With"] = "XMLHttpRequest";
      config.headers["Origin"] = window.location.origin;

      console.log("📨 Request configuration:", {
        url: config.url,
        method: config.method,
        headers: {
          "X-CSRF-Token": token,
          "Content-Type": config.headers["Content-Type"],
          Origin: config.headers["Origin"],
        },
      });

      return config;
    } catch (error) {
      console.error("Failed to prepare request:", error);
      return Promise.reject(error);
    }
  },
  (error: any) => {
    console.error("Request preparation failed:", error);
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response: AxiosResponse) => {
    console.log("✅ Response:", {
      status: response.status,
      headers: response.headers,
      cookies: document.cookie,
    });
    return response;
  },
  async (error: AxiosError<APIErrorResponse>) => {
    console.error("❌ Response error:", {
      status: error.response?.status,
      data: error.response?.data,
      headers: error.response?.headers,
      cookies: document.cookie,
    });

    const config = error.config as CustomAxiosRequestConfig;

    // Handle CSRF errors with retry
    if (
      error.response?.status === 403 &&
      error.response?.data?.code === "INVALID_CSRF_TOKEN" &&
      !config._retry
    ) {
      config._retry = true;
      csrfToken = null; // Reset token
      return api(config);
    }

    // Transform error to custom APIError
    const apiError = new APIError(
      error.response?.data?.code || "UNKNOWN_ERROR",
      error.response?.data?.error || "An unexpected error occurred",
      error.response?.status,
      error.response?.data?.details
    );

    return Promise.reject(apiError);
  }
);

// Utility functions
export const resetCSRFToken = (): void => {
  csrfToken = null;
  tokenPromise = null;
};

export const getApiConfig = () => ({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,
});

// Debug helper for development
if (process.env.NODE_ENV === "development") {
  api.interceptors.request.use((request) => {
    console.log("🔄 Request:", {
      url: request.url,
      method: request.method,
      headers: request.headers,
    });
    return request;
  });
}

export default api;
