// src/libs/authService.ts
import api, { APIError } from "./api";
import { csrfService } from "./csrfService";
import axios, { AxiosError } from "axios";

// Updated API response types to match the backend
interface AuthResponse {
  message: string;
  code: string;
  user?: {
    username: string;
    email?: string;
  };
}

interface AuthCheckResponse {
  authenticated: boolean;
  code: string;
}

// Error codes aligned with the backend
export enum AuthErrorCode {
  // Authentication errors
  INVALID_CREDENTIALS = "invalid_credentials",
  INVALID_INPUT = "invalid_input",
  USERNAME_EXISTS = "username_exists",
  EMAIL_EXISTS = "email_exists",

  // Security errors
  INVALID_CSRF_TOKEN = "INVALID_CSRF_TOKEN",
  UNAUTHORIZED = "unauthorized",

  // Technical errors
  NETWORK_ERROR = "network_error",
  SERVER_ERROR = "server_error",

  // Other errors
  LOGIN_ERROR = "login_error",
  SIGNUP_ERROR = "signup_error",
  LOGOUT_ERROR = "logout_error",
  REFRESH_ERROR = "refresh_token_error",
}

export class AuthError extends Error {
  constructor(
    public code: AuthErrorCode,
    message: string,
    public originalError?: unknown
  ) {
    super(message);
    this.name = "AuthError";
  }
}

export const authService = {
  async signup(
    username: string,
    email: string,
    password: string
  ): Promise<AuthResponse> {
    try {
      const response = await api.post<AuthResponse>("/auth/signup", {
        username,
        email,
        password,
      });

      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        const errorCode = err.response?.data?.code;

        switch (errorCode) {
          case "username_exists":
            throw new AuthError(
              AuthErrorCode.USERNAME_EXISTS,
              "Username already taken"
            );
          case "email_exists":
            throw new AuthError(
              AuthErrorCode.EMAIL_EXISTS,
              "Email already registered"
            );
          case "invalid_input":
            throw new AuthError(
              AuthErrorCode.INVALID_INPUT,
              err.response?.data?.error || "Invalid input data"
            );
          default:
            throw new AuthError(
              AuthErrorCode.SIGNUP_ERROR,
              "Failed to create account"
            );
        }
      }
      throw new AuthError(
        AuthErrorCode.NETWORK_ERROR,
        "Network error during signup"
      );
    }
  },

  async login(username: string, password: string): Promise<AuthResponse> {
    try {
      const response = await api.post<AuthResponse>("/auth/login", {
        username,
        password,
      });

      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        if (err.response?.status === 401) {
          throw new AuthError(
            AuthErrorCode.INVALID_CREDENTIALS,
            "Invalid username or password"
          );
        }

        throw new AuthError(
          AuthErrorCode.LOGIN_ERROR,
          err.response?.data?.error || "Login failed"
        );
      }
      throw new AuthError(
        AuthErrorCode.NETWORK_ERROR,
        "Network error during login"
      );
    }
  },

  async logout(): Promise<void> {
    try {
      await api.post<AuthResponse>("/auth/logout");
      csrfService.clearToken();
    } catch (error) {
      // Always clear the token in case of error
      csrfService.clearToken();

      if (axios.isAxiosError(error)) {
        throw new AuthError(
          AuthErrorCode.LOGOUT_ERROR,
          "Failed to logout properly"
        );
      }
      throw new AuthError(
        AuthErrorCode.NETWORK_ERROR,
        "Network error during logout"
      );
    }
  },

  async refreshToken(): Promise<AuthResponse> {
    try {
      const response = await api.post<AuthResponse>("/auth/refresh");
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        if (err.response?.status === 401) {
          csrfService.clearToken();
          throw new AuthError(AuthErrorCode.UNAUTHORIZED, "Session expired");
        }
      }
      throw new AuthError(
        AuthErrorCode.REFRESH_ERROR,
        "Failed to refresh session"
      );
    }
  },

  async isAuthenticated(): Promise<boolean> {
    try {
      const response = await api.get<AuthCheckResponse>("/auth/check");
      return response.data.authenticated;
    } catch (error) {
      return false;
    }
  },
};
