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
      window.location.href = "/login"; // Redirect after logout
    } catch (error) {
      console.error("Error during logout:", error);
    }
  },

  async refreshToken() {
    try {
      const response = await api.post("/auth/refresh");
      return response.data;
    } catch (error) {
      console.error("Error refreshing token:", error);
      throw error;
    }
  },

  async isAuthenticated() {
    try {
      await api.get("/auth/check");
      return true;
    } catch (error) {
      return false;
    }
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
      return response.data;
    } catch (error) {
      console.error("Error handling OAuth callback:", error);
      throw error;
    }
  },
};
