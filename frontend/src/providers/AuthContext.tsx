"use client";

import React, { createContext, useState, useContext, useEffect } from "react";
import { authService } from "@/libs/authService";

interface AuthContextType {
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  signup: (username: string, email: string, password: string) => Promise<void>;
  initiateOAuthLogin: () => Promise<void>;
  handleOAuthCallback: (code: string, state: string) => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    authService.isAuthenticated().then(setIsAuthenticated);
  }, []);

  const login = async (username: string, password: string) => {
    await authService.login(username, password);
    setIsAuthenticated(true);
  };

  const logout = async () => {
    await authService.logout();
    setIsAuthenticated(false);
  };

  const signup = async (username: string, email: string, password: string) => {
    await authService.signup(username, email, password);
  };

  const initiateOAuthLogin = async () => {
    await authService.initiateOAuthLogin();
  };

  const handleOAuthCallback = async (code: string, state: string) => {
    await authService.handleOAuthCallback(code, state);
    setIsAuthenticated(true);
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        login,
        logout,
        signup,
        initiateOAuthLogin,
        handleOAuthCallback,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
