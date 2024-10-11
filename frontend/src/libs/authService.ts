import api from "./api";

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
      throw error;
    }
  },

  async logout() {
    try {
      await api.post("/auth/logout");
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
    } catch (error) {
      throw error;
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
};
