import api from "./api";

interface LoginResponse {
  message: string;
}

interface SignupData {
  username: string;
  email: string;
  password: string;
}

interface LoginData {
  username: string;
  password: string;
}

export const authService = {
  async signup(data: SignupData) {
    const csrfToken = localStorage.getItem("csrfToken");
    await api.post("/auth/signup", data, {
      headers: { "X-CSRF-Token": csrfToken },
    });
  },

  async login(data: LoginData) {
    const csrfToken = localStorage.getItem("csrfToken");
    const response = await api.post<LoginResponse>("/auth/login", data, {
      headers: { "X-CSRF-Token": csrfToken },
    });
    return response.data.message;
  },

  async logout(): Promise<void> {
    const csrfToken = localStorage.getItem("csrfToken");
    await api.post(
      "/auth/logout",
      {},
      {
        headers: { "X-CSRF-Token": csrfToken },
      }
    );
  },

  async isAuthenticated(): Promise<boolean> {
    try {
      const csrfToken = localStorage.getItem("csrfToken");
      const response = await api.get("/auth/check", {
        headers: { "X-CSRF-Token": csrfToken },
      });
      return response.status === 200;
    } catch (error) {
      return false;
    }
  },

  async getCSRFToken(): Promise<string> {
    const response = await api.get<{ csrf_token: string }>("/csrf-token");
    return response.data.csrf_token;
  },
};
