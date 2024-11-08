import api from "./api";
import axios, { AxiosError } from "axios";

interface AuthError {
  error: string;
  message?: string;
}

export const authService = {
  async signup(username: string, email: string, password: string) {
    try {
      const response = await api.post<{ message: string }>("/auth/signup", {
        username,
        email,
        password,
      });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<AuthError>;
        throw new Error(err.response?.data?.message || "Signup failed");
      }
      throw error;
    }
  },

  async login(username: string, password: string) {
    try {
      const response = await api.post<{ message: string }>("/auth/login", {
        username,
        password,
      });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response?.status === 401) {
          throw new Error("Invalid credentials");
        }
        const err = error as AxiosError<AuthError>;
        throw new Error(err.response?.data?.message || "Login failed");
      }
      throw new Error("An error occurred during login");
    }
  },

  async logout() {
    try {
      // La vérification CSRF sera gérée par l'interceptor
      await api.post("/auth/logout");
      // Ne pas rediriger ici, laisser le contexte gérer la redirection
      return true;
    } catch (error) {
      console.error("Error during logout:", error);
      throw error;
    }
  },

  async refreshToken() {
    try {
      const response = await api.post<{ message: string }>("/auth/refresh");
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<AuthError>;
        throw new Error(err.response?.data?.message || "Token refresh failed");
      }
      throw error;
    }
  },

  async isAuthenticated() {
    try {
      const response = await api.get<{ authenticated: boolean }>("/auth/check");
      return response.data.authenticated;
    } catch (error) {
      return false;
    }
  },

  async initiateOAuthLogin() {
    try {
      const response = await api.get<{ url: string }>("/auth/battle-net/login");
      return response.data.url;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<AuthError>;
        throw new Error(
          err.response?.data?.message || "OAuth initiation failed"
        );
      }
      throw error;
    }
  },

  async handleOAuthCallback(code: string, state: string) {
    try {
      const response = await api.get("/auth/battle-net/callback", {
        params: { code, state },
      });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<AuthError>;
        throw new Error(err.response?.data?.message || "OAuth callback failed");
      }
      throw error;
    }
  },
};
