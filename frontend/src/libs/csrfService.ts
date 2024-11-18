// src/libs/csrfService.ts
import axios, { AxiosInstance } from "axios";
import https from "https";

const CSRF_PROTECTED_ROUTES = [
  "/auth/login",
  "/auth/signup",
  "/auth/logout",
  "/user/email",
  "/user/password",
  "/user/username",
  "/user/account",
];

class CSRFService {
  private static instance: CSRFService;
  private token: string | null = null;
  private tokenExpiryTime: number | null = null;
  private readonly TOKEN_LIFETIME = 60 * 60 * 1000;
  private readonly csrfAxios: AxiosInstance;

  private constructor() {
    this.csrfAxios = axios.create({
      baseURL: process.env.NEXT_PUBLIC_API_URL,
      withCredentials: true,
      headers: {
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
    return (
      method.toLowerCase() !== "get" &&
      CSRF_PROTECTED_ROUTES.some((route) => url.includes(route))
    );
  }

  async getToken(forceRefresh = false): Promise<string | null> {
    if (forceRefresh || !this.token || this.isTokenExpired()) {
      try {
        const response = await this.csrfAxios.get("/auth/csrf-token", {
          headers: {
            Origin: process.env.NEXT_PUBLIC_APP_URL,
          },
        });
        this.token = response.data.token;
        this.tokenExpiryTime = Date.now() + this.TOKEN_LIFETIME;
      } catch (error) {
        console.error("Failed to fetch CSRF token:", error);
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
}

export const csrfService = CSRFService.getInstance();
