import { useState, useEffect, useCallback } from "react";
import { authService } from "@/libs/authService";

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);


  useEffect(() => {
    const checkAuth = async () => {
      try {
        const isAuth = await authService.isAuthenticated();
        setIsAuthenticated(isAuth);
      } catch (error) {
        setIsAuthenticated(false);
      } finally {
        setIsLoading(false);
      }
    };
    checkAuth();
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

  return { isAuthenticated, isLoading, login, logout, signup };
}
