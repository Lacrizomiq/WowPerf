import api from "./api";
import axios, { AxiosError } from "axios";

export const authService = {
  async signup(username: string, email: string, password: string) {
    try {
      const response = await api.post("/auth/signup", {
        username,
        email,
        password,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  async login(username: string, password: string) {
    try {
      const response = await api.post("/auth/login", { username, password });
      const { access_token, refresh_token } = response.data;
      localStorage.setItem("accessToken", access_token);
      localStorage.setItem("refreshToken", refresh_token);
      return response.data;
    } catch (error) {
      console.error("Error logging in:", error);
      if (axios.isAxiosError(error)) {
        if (error.response && error.response.status === 401) {
          throw new Error("Invalid credentials");
        }
      }
      throw new Error("An error occurred during login");
    }
  },

  async logout() {
    try {
      await api.post("/auth/logout");
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
    } catch (error) {
      console.error("Error during logout:", error);
    }
  },

  async refreshToken() {
    const refreshToken = localStorage.getItem("refreshToken");
    if (!refreshToken) {
      throw new Error("No refresh token available");
    }
    try {
      const response = await api.post("/auth/refresh", {
        refresh_token: refreshToken,
      });
      const { access_token } = response.data;
      localStorage.setItem("accessToken", access_token);
      return access_token;
    } catch (error) {
      throw error;
    }
  },

  isAuthenticated() {
    return !!localStorage.getItem("accessToken");
  },

  async initiateOAuthLogin() {
    try {
      const response = await api.get("/auth/battle-net/login");
      window.location.href = response.data.url;
    } catch (error) {
      console.error("Error initiating OAuth login:", error);
      throw error;
    }
  },

  async handleOAuthCallback(code: string, state: string) {
    try {
      const response = await api.get("/auth/battle-net/callback", {
        params: { code, state },
      });
      const { access_token, refresh_token } = response.data;
      localStorage.setItem("accessToken", access_token);
      localStorage.setItem("refreshToken", refresh_token);
      return response.data;
    } catch (error) {
      console.error("Error handling OAuth callback:", error);
      throw error;
    }
  },
};
