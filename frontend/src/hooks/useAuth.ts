import { useState, useEffect, useCallback } from "react";
import { authService } from "@/libs/authService";

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    setIsAuthenticated(authService.isAuthenticated());
  }, []);

  const login = useCallback(async (username: string, password: string) => {
    try {
      await authService.login(username, password);
      setIsAuthenticated(true);
    } catch (error) {
      throw error;
    }
  }, []);

  const logout = useCallback(async () => {
    try {
      await authService.logout();
      setIsAuthenticated(false);
    } catch (error) {
      throw error;
    }
  }, []);

  const signup = useCallback(
    async (username: string, email: string, password: string) => {
      try {
        const response = await authService.signup(username, email, password);
        setIsAuthenticated(true);
        return response;
      } catch (error) {
        throw error;
      }
    },
    []
  );

  return { isAuthenticated, login, logout, signup };
}
