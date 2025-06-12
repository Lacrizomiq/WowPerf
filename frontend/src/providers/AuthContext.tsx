// src/contexts/AuthContext.tsx
"use client";

import React, {
  createContext,
  useState,
  useContext,
  useEffect,
  useCallback,
} from "react";
import {
  authService,
  AuthError,
  AuthErrorCode,
  AuthMethod,
} from "@/libs/authService";
import { useRouter } from "next/navigation";
import { resetCSRFToken, preloadCSRFToken } from "@/libs/api";
import { usePathname } from "next/navigation";
import { googleAuthService } from "@/libs/googleAuthService";

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: UserData | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  signup: (
    username: string,
    email: string,
    password: string,
    captchaToken?: string
  ) => Promise<void>;
  checkAuth: () => Promise<void>;
  loginWithGoogle: () => Promise<void>;
}

interface UserData {
  username: string;
  email?: string;
  authMethod?: AuthMethod;
  hasGoogleLinked?: boolean;
}

interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: UserData | null;
}

const initialState: AuthState = {
  isAuthenticated: false,
  isLoading: true,
  user: null,
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [state, setState] = useState<AuthState>(initialState);
  const router = useRouter();

  // Using a state to track CSRF initialization only if necessary
  const [csrfInitialized, setCsrfInitialized] = useState(false);

  const updateState = useCallback((updates: Partial<AuthState>) => {
    setState((prev) => ({ ...prev, ...updates }));
  }, []);

  // Checking authentication
  const checkAuth = useCallback(async () => {
    try {
      const isAuth = await authService.isAuthenticated();
      updateState({
        isAuthenticated: isAuth,
        isLoading: false,
      });

      // If authenticated, preload the CSRF token
      if (isAuth && !csrfInitialized) {
        await preloadCSRFToken();
        setCsrfInitialized(true);
      }
    } catch (error) {
      updateState({
        isAuthenticated: false,
        isLoading: false,
        user: null,
      });
    }
  }, [updateState, csrfInitialized]);

  // Effect for initial authentication check
  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  const handleAuthError = useCallback(
    async (error: unknown): Promise<string> => {
      if (error instanceof AuthError) {
        switch (error.code) {
          case AuthErrorCode.INVALID_CSRF_TOKEN:
            try {
              // Attempt to refresh the CSRF token
              await preloadCSRFToken();
              return "Please try again. Security token refreshed.";
            } catch {
              resetCSRFToken();
              return "Security verification failed. Please try again.";
            }
          case AuthErrorCode.INVALID_CREDENTIALS:
            return "Invalid username or password";
          case AuthErrorCode.USERNAME_EXISTS:
            return "Username already exists";
          case AuthErrorCode.EMAIL_EXISTS:
            return "Email already exists";
          case AuthErrorCode.CAPTCHA_REQUIRED:
            return "Please complete the captcha verification";
          case AuthErrorCode.CAPTCHA_INVALID:
            return "Captcha verification failed. Please try again.";
          case AuthErrorCode.NETWORK_ERROR:
            return "Connection error. Please check your internet connection.";
          case AuthErrorCode.UNAUTHORIZED:
            updateState({
              isAuthenticated: false,
              user: null,
            });
            resetCSRFToken();
            router.push("/login");
            return "Session expired. Please log in again.";
          default:
            return error.message || "An unexpected error occurred";
        }
      }
      return "An unexpected error occurred";
    },
    [router, updateState]
  );

  const login = useCallback(
    async (email: string, password: string) => {
      try {
        console.log("Starting login process...");
        const response = await authService.login(email, password);
        console.log("Login successful:", response);

        updateState({
          isAuthenticated: true,
          user: response.user,
          isLoading: false,
        });

        console.log("State updated after login");

        try {
          await preloadCSRFToken();
          setCsrfInitialized(true);
          console.log("CSRF token preloaded");
        } catch (csrfError) {
          console.error("Failed to preload CSRF token:", csrfError);
        }

        // Add a small delay before redirect
        await new Promise((resolve) => setTimeout(resolve, 100));
        console.log("Redirecting to profile...");
        router.push("/profile");
      } catch (error) {
        console.error("Login process failed:", error);
        const errorMessage = await handleAuthError(error);
        throw new Error(errorMessage);
      }
    },
    [router, handleAuthError, updateState]
  );

  const logout = useCallback(async () => {
    try {
      await authService.logout();
      updateState({
        isAuthenticated: false,
        user: null,
      });
      resetCSRFToken();
      setCsrfInitialized(false);
      router.push("/login");
    } catch (error) {
      console.error("Logout failed:", error);
      // Cleaning the state even in case of error
      updateState({
        isAuthenticated: false,
        user: null,
      });
      resetCSRFToken();
      setCsrfInitialized(false);
      router.push("/login");
    }
  }, [router, updateState]);

  const signup = useCallback(
    async (
      username: string,
      email: string,
      password: string,
      captchaToken?: string
    ) => {
      try {
        const signupResponse = await authService.signup(
          username,
          email,
          password,
          captchaToken
        );

        updateState({
          isAuthenticated: true,
          user: signupResponse.user,
        });

        // Preload the CSRF token
        await preloadCSRFToken();
        setCsrfInitialized(true);

        // Redirect
        router.push("/profile");
      } catch (error) {
        const errorMessage = await handleAuthError(error);
        throw new Error(errorMessage);
      }
    },
    [router, handleAuthError, updateState]
  );

  const loginWithGoogle = useCallback(async () => {
    try {
      await googleAuthService.initiateGoogleLogin();
      // Le backend redirige automatiquement vers Google
    } catch (error) {
      console.error("Failed to initiate Google login:", error);
      const errorMessage = await handleAuthError(error);
      throw new Error(errorMessage);
    }
  }, [handleAuthError]);

  const value = {
    isAuthenticated: state.isAuthenticated,
    isLoading: state.isLoading,
    user: state.user,
    login,
    logout,
    signup,
    checkAuth,
    loginWithGoogle,
  };

  if (state.isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        Loading...
      </div>
    );
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// Hook to use the authentication context
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

export const useRequireAuth = () => {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      // Store the current page to redirect after login if necessary
      router.push(`/login?redirect=${encodeURIComponent(pathname)}`);
    }
  }, [isAuthenticated, isLoading, router, pathname]);

  return { isAuthenticated, isLoading };
};
