// src/libs/csrfService.ts
import axios, { AxiosInstance } from "axios";
import https from "https";

// Routes requiring CSRF protection (non-GET methods only)
const CSRF_PROTECTED_ROUTES = [
  // Auth routes protected by CSRF
  "/auth/logout",
  "/auth/refresh",

  // User routes protected by CSRF
  "/user/email",
  "/user/password",
  "/user/username",
  "/user/account",
];

class CSRFService {
  private static instance: CSRFService;
  private token: string | null = null;
  private tokenExpiryTime: number | null = null;
  private readonly TOKEN_LIFETIME = 3600 * 1000; // 1 hour, consistent with backend
  private readonly csrfAxios: AxiosInstance;

  private constructor() {
    this.csrfAxios = axios.create({
      baseURL: process.env.NEXT_PUBLIC_API_URL,
      withCredentials: true,
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        "X-Requested-With": "XMLHttpRequest",
      },
      ...(process.env.NODE_ENV === "development" && {
        httpsAgent: new https.Agent({
          rejectUnauthorized: false,
        }),
      }),
    });
  }

  static getInstance(): CSRFService {
    if (!CSRFService.instance) {
      CSRFService.instance = new CSRFService();
    }
    return CSRFService.instance;
  }

  isProtectedRoute(url: string, method: string): boolean {
    // CSRF protection is only required for non-GET methods
    if (method.toLowerCase() === "get") {
      return false;
    }

    // Check if the URL matches a protected route
    return CSRF_PROTECTED_ROUTES.some(
      (route) => url.endsWith(route) || url.includes(`${route}/`)
    );
  }

  async getToken(forceRefresh = false): Promise<string | null> {
    // If we force a refresh or if the token is expired/missing
    if (forceRefresh || !this.token || this.isTokenExpired()) {
      try {
        const response = await this.csrfAxios.get("/api/csrf-token", {
          headers: {
            Origin: process.env.NEXT_PUBLIC_APP_URL,
            "X-Requested-With": "XMLHttpRequest",
          },
        });

        if (response.data.token) {
          this.token = response.data.token;
          this.tokenExpiryTime = Date.now() + this.TOKEN_LIFETIME;
          return this.token;
        }

        console.warn("No CSRF token in response");
        return null;
      } catch (error) {
        console.error("Failed to fetch CSRF token:", error);
        this.clearToken();
        return null;
      }
    }

    return this.token;
  }

  private isTokenExpired(): boolean {
    return !this.tokenExpiryTime || Date.now() > this.tokenExpiryTime;
  }

  clearToken(): void {
    this.token = null;
    this.tokenExpiryTime = null;
  }

  // New method to check if a route requires a token reset
  shouldResetToken(error: any): boolean {
    return (
      error?.response?.data?.code === "INVALID_CSRF_TOKEN" ||
      error?.response?.status === 403
    );
  }
}

export const csrfService = CSRFService.getInstance();
