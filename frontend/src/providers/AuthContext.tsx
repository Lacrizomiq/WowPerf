"use client";

import React, {
  createContext,
  useState,
  useContext,
  useEffect,
  useCallback,
} from "react";
import { authService } from "@/libs/authService";
import { useRouter } from "next/navigation";

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
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
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  const checkAuth = useCallback(async () => {
    try {
      const isAuth = await authService.isAuthenticated();
      setIsAuthenticated(isAuth);
    } catch (error) {
      setIsAuthenticated(false);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  const login = useCallback(
    async (username: string, password: string) => {
      try {
        await authService.login(username, password);
        setIsAuthenticated(true);
        router.push("/dashboard"); // Redirection after successful login
      } catch (error) {
        throw error;
      }
    },
    [router]
  );

  const logout = useCallback(async () => {
    try {
      await authService.logout();
      setIsAuthenticated(false);
      router.push("/login"); // Handle redirection here rather than in the service
    } catch (error) {
      console.error("Logout failed:", error);
      throw error;
    }
  }, [router]);

  const signup = useCallback(
    async (username: string, email: string, password: string) => {
      try {
        await authService.signup(username, email, password);
        await login(username, password); // Auto login after signup
      } catch (error) {
        throw error;
      }
    },
    [login]
  );

  const initiateOAuthLogin = useCallback(async () => {
    try {
      const url = await authService.initiateOAuthLogin();
      window.location.href = url;
    } catch (error) {
      console.error("OAuth initiation failed:", error);
      throw error;
    }
  }, []);

  const handleOAuthCallback = useCallback(
    async (code: string, state: string) => {
      try {
        await authService.handleOAuthCallback(code, state);
        setIsAuthenticated(true);
        router.push("/dashboard");
      } catch (error) {
        console.error("OAuth callback failed:", error);
        router.push("/login?error=oauth_failed");
      }
    },
    [router]
  );

  const value = {
    isAuthenticated,
    isLoading,
    login,
    logout,
    signup,
    initiateOAuthLogin,
    handleOAuthCallback,
  };

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
