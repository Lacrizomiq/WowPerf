import api from "./api";

interface LoginResponse {
  token: string;
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

const AUTH_TOKEN_KEY = "auth_token";

export const authService = {
  async signup(data: SignupData) {
    await api.post("/auth/signup", data);
  },

  async login(data: LoginData) {
    const response = await api.post<LoginResponse>("/auth/login", data);
    const token = response.data.token;
    localStorage.setItem(AUTH_TOKEN_KEY, token);
    return token;
  },

  async logout(): Promise<void> {
    await api.post("/auth/logout");
    localStorage.removeItem(AUTH_TOKEN_KEY);
  },

  getToken(): string | null {
    return localStorage.getItem(AUTH_TOKEN_KEY);
  },

  isAuthenticated(): boolean {
    return !!this.getToken();
  },
};
