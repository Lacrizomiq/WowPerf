"use client";

import React, {
  createContext,
  useState,
  useContext,
  useEffect,
  useCallback,
} from "react";
import { authService, AuthError, AuthErrorCode } from "@/libs/authService";
import { useRouter } from "next/navigation";
import axios, { AxiosError } from "axios";
import { resetCSRFToken, preloadCSRFToken } from "@/libs/api";

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: UserData | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  signup: (username: string, email: string, password: string) => Promise<void>;
}

interface CSRFErrorResponse {
  code: string;
  error?: string;
  details?: string;
}

interface UserData {
  username: string;
  email?: string;
  battlenet_id?: string;
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
  const [csrfInitialized, setCsrfInitialized] = useState(false);

  const updateState = useCallback((updates: Partial<AuthState>) => {
    setState((prev) => ({ ...prev, ...updates }));
  }, []);

  const checkAuth = useCallback(async () => {
    try {
      const isAuth = await authService.isAuthenticated();
      updateState({
        isAuthenticated: isAuth,
        isLoading: false,
      });
    } catch (error) {
      // Ne pas rediriger, juste mettre à jour l'état
      updateState({
        isAuthenticated: false,
        isLoading: false,
      });
    }
  }, [updateState]);

  // Initialize CSRF only for protected routes
  useEffect(() => {
    const initializeCSRF = async () => {
      if (!csrfInitialized) {
        try {
          await preloadCSRFToken();
          setCsrfInitialized(true);
        } catch (error) {
          console.error("Failed to initialize CSRF token:", error);
          setTimeout(initializeCSRF, 2000);
        }
      }
    };

    initializeCSRF();
  }, [csrfInitialized]);

  // Check authentication after CSRF initialization
  useEffect(() => {
    if (csrfInitialized) {
      checkAuth();
    }
  }, [checkAuth, csrfInitialized]);

  const handleAuthError = useCallback(
    async (error: unknown): Promise<string> => {
      if (error instanceof AuthError) {
        switch (error.code) {
          case AuthErrorCode.CSRF_ERROR:
            try {
              await preloadCSRFToken();
              return "Please try again";
            } catch (e) {
              resetCSRFToken();
              return "Security verification failed. Please try again.";
            }
          case AuthErrorCode.INVALID_CREDENTIALS:
            return "Invalid username or password";
          case AuthErrorCode.USERNAME_EXISTS:
            return "Username already exists";
          case AuthErrorCode.EMAIL_EXISTS:
            return "Email already exists";
          case AuthErrorCode.NETWORK_ERROR:
            return "Network error, please try again";
          default:
            return error.message;
        }
      }

      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<CSRFErrorResponse>;
        if (err.response?.status === 401) {
          updateState({
            isAuthenticated: false,
            user: null,
          });
          resetCSRFToken();
          return "Session expired";
        }
      }

      return "An unexpected error occurred";
    },
    [updateState]
  );

  const login = useCallback(
    async (username: string, password: string) => {
      try {
        const response = await authService.login(username, password);
        updateState({
          isAuthenticated: true,
          user: {
            username: response.user.username,
          },
        });
        router.push("/profile");
      } catch (error) {
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
    } catch (error) {
      console.error("Logout failed:", error);
      updateState({
        isAuthenticated: false,
        user: null,
      });
      resetCSRFToken();
    }
  }, [updateState]);

  const signup = useCallback(
    async (username: string, email: string, password: string) => {
      try {
        await authService.signup(username, email, password);
        // Auto-login after successful signup
        await login(username, password);
      } catch (error) {
        const errorMessage = await handleAuthError(error);
        throw new Error(errorMessage);
      }
    },
    [login, handleAuthError]
  );

  const value = {
    isAuthenticated: state.isAuthenticated,
    isLoading: state.isLoading || !csrfInitialized,
    user: state.user,
    login,
    logout,
    signup,
  };

  if (state.isLoading || !csrfInitialized) {
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

// Utility hook for redirection if not authenticated
