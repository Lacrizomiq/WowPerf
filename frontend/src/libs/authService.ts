import api, { APIError, resetCSRFToken } from "./api";
import axios, { AxiosError } from "axios";

// Precise types for API responses
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

// Possible errors
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
  ): Promise<SignupResponse> {
    try {
      console.log("AuthService: Starting signup request");
      const response = await api.post<SignupResponse>("/auth/signup", {
        username,
        email,
        password,
      });

      console.log("AuthService: Signup successful");
      return response.data;
    } catch (error) {
      console.error("AuthService: Signup error:", error);

      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        switch (err.response?.data?.code) {
          case "username_exists":
            throw new AuthError(
              AuthErrorCode.USERNAME_EXISTS,
              "Username already exists"
            );
          case "email_exists":
            throw new AuthError(
              AuthErrorCode.EMAIL_EXISTS,
              "Email already exists"
            );
          case "invalid_input":
            throw new AuthError(
              AuthErrorCode.INVALID_INPUT,
              "Invalid signup data"
            );
          default:
            throw new AuthError(
              AuthErrorCode.SIGNUP_ERROR,
              err.response?.data?.error || "Signup failed",
              err
            );
        }
      }

      throw new AuthError(
        AuthErrorCode.UNKNOWN_ERROR,
        "An unexpected error occurred during signup",
        error
      );
    }
  },

  async login(username: string, password: string): Promise<LoginResponse> {
    try {
      const response = await api.post<LoginResponse>("/auth/login", {
        username,
        password,
      });

      // Store user information if needed
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
          err.response?.data?.error || "Login failed",
          err
        );
      }

      throw new AuthError(
        AuthErrorCode.NETWORK_ERROR,
        "Network error during login",
        error
      );
    }
  },

  async logout(): Promise<void> {
    try {
      await api.post<AuthResponse>("/auth/logout");
      // Reset the CSRF token after logout
      resetCSRFToken();
    } catch (error) {
      console.error("Logout error:", error);
      // Reset the CSRF token even in case of error
      resetCSRFToken();
      throw error;
    }
  },

  async refreshToken(): Promise<AuthResponse> {
    try {
      const response = await api.post<AuthResponse>("/auth/refresh");
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        // If the refresh token is invalid or expired
        if (err.response?.status === 401) {
          resetCSRFToken(); // Reset CSRF token
          throw new AuthError(
            AuthErrorCode.INVALID_CREDENTIALS,
            "Session expired"
          );
        }

        throw new AuthError(
          AuthErrorCode.UNKNOWN_ERROR,
          err.response?.data?.error || "Token refresh failed",
          err
        );
      }

      throw new AuthError(
        AuthErrorCode.NETWORK_ERROR,
        "Network error during token refresh",
        error
      );
    }
  },

  async isAuthenticated(): Promise<boolean> {
    try {
      const response = await api.get<AuthCheckResponse>("/auth/check");
      return response.data.authenticated;
    } catch (error) {
      // In case of error, consider the user not authenticated
      return false;
    }
  },
};
