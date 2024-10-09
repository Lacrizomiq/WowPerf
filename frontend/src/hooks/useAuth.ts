"use client";

import { useState, useEffect } from "react";
import { authService } from "@/libs/authService";

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    setIsLoading(true);
    const authStatus = await authService.isAuthenticated();
    setIsAuthenticated(authStatus);
    setIsLoading(false);
  };

  const login = async (username: string, password: string) => {
    try {
      await authService.login({ username, password });
      setIsAuthenticated(true);
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  };

  const signup = async (username: string, email: string, password: string) => {
    try {
      await authService.signup({ username, email, password });
    } catch (error) {
      console.error("Signup failed:", error);
      throw error;
    }
  };

  const logout = async () => {
    try {
      await authService.logout();
      setIsAuthenticated(false);
    } catch (error) {
      console.error("Logout failed:", error);
    }
  };

  return { isAuthenticated, isLoading, login, signup, logout };
}
