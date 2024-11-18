import api from "./api";
import axios, { AxiosError } from "axios";

// Response Types
interface AuthResponse {
  message: string;
  user?: {
    username: string;
    email?: string;
  };
  csrf_token?: string;
}

interface LoginResponse extends AuthResponse {
  user: {
    username: string;
  };
}

interface SignupResponse extends AuthResponse {
  code?: string;
}

interface AuthCheckResponse {
  authenticated: boolean;
}

interface APIError {
  error: string;
  code: string;
  details?: string;
}

// Error Handling
export enum AuthErrorCode {
  INVALID_CREDENTIALS = "invalid_credentials",
  INVALID_INPUT = "invalid_input",
  USERNAME_EXISTS = "username_exists",
  EMAIL_EXISTS = "email_exists",
  SIGNUP_ERROR = "signup_error",
  LOGIN_ERROR = "login_error",
  NETWORK_ERROR = "network_error",
  OAUTH_ERROR = "oauth_error",
  UNKNOWN_ERROR = "unknown_error",
  CSRF_ERROR = "csrf_error",
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

// Auth Service Implementation
export const authService = {
  async signup(
    username: string,
    email: string,
    password: string
  ): Promise<SignupResponse> {
    try {
      const response = await api.post<SignupResponse>("/auth/signup", {
        username,
        email,
        password,
      });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        const data = err.response?.data;

        if (data?.code === "INVALID_CSRF_TOKEN") {
          throw new AuthError(
            AuthErrorCode.CSRF_ERROR,
            "Security verification failed"
          );
        }

        switch (data?.code) {
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
              data.error || "Invalid input"
            );
          default:
            throw new AuthError(
              AuthErrorCode.UNKNOWN_ERROR,
              data?.error || "Signup failed"
            );
        }
      }
      throw new AuthError(AuthErrorCode.NETWORK_ERROR, "Network error");
    }
  },

  async login(username: string, password: string): Promise<LoginResponse> {
    try {
      const response = await api.post<LoginResponse>("/auth/login", {
        username,
        password,
      });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        const data = err.response?.data;

        if (data?.code === "INVALID_CSRF_TOKEN") {
          throw new AuthError(
            AuthErrorCode.CSRF_ERROR,
            "Security verification failed"
          );
        }

        if (err.response?.status === 401) {
          throw new AuthError(
            AuthErrorCode.INVALID_CREDENTIALS,
            "Invalid username or password"
          );
        }

        throw new AuthError(
          AuthErrorCode.LOGIN_ERROR,
          data?.error || "Login failed"
        );
      }
      throw new AuthError(AuthErrorCode.NETWORK_ERROR, "Network error");
    }
  },

  async logout(): Promise<void> {
    try {
      await api.post<AuthResponse>("/auth/logout");
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        if (err.response?.data?.code === "INVALID_CSRF_TOKEN") {
          throw new AuthError(
            AuthErrorCode.CSRF_ERROR,
            "Security verification failed"
          );
        }
      }
      throw new AuthError(AuthErrorCode.UNKNOWN_ERROR, "Logout failed");
    }
  },

  async isAuthenticated(): Promise<boolean> {
    try {
      const response = await api.get<AuthCheckResponse>("/auth/check");
      return response.data.authenticated;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;
        if (err.response?.status === 401) {
          return false;
        }
        if (err.response?.data?.code === "INVALID_CSRF_TOKEN") {
          throw new AuthError(
            AuthErrorCode.CSRF_ERROR,
            "Security verification failed"
          );
        }
      }
      return false;
    }
  },
};
