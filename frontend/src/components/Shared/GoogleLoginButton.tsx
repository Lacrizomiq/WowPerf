// src/components/auth/GoogleLoginButton.tsx

"use client";

import React from "react";
import { useAuth } from "@/providers/AuthContext";

interface GoogleLoginButtonProps {
  /** Texte du bouton - par défaut selon le contexte */
  text?: string;
  /** Variante du bouton */
  variant?: "signin" | "signup" | "continue";
  /** Si le bouton est disabled */
  disabled?: boolean;
  /** Classe CSS additionnelle */
  className?: string;
  /** Callback en cas d'erreur */
  onError?: (error: Error) => void;
}

const GoogleLoginButton: React.FC<GoogleLoginButtonProps> = ({
  text,
  variant = "signin",
  disabled = false,
  className = "",
  onError,
}) => {
  const { loginWithGoogle } = useAuth();
  const [isLoading, setIsLoading] = React.useState(false);

  // Textes par défaut selon les guidelines Google
  const getDefaultText = () => {
    switch (variant) {
      case "signup":
        return "Sign up with Google";
      case "continue":
        return "Continue with Google";
      case "signin":
      default:
        return "Sign in with Google";
    }
  };

  const handleGoogleLogin = async () => {
    if (disabled || isLoading) return;

    try {
      setIsLoading(true);
      await loginWithGoogle();
      // La redirection est gérée automatiquement par le backend
    } catch (error) {
      console.error("Google login failed:", error);
      setIsLoading(false);

      if (onError && error instanceof Error) {
        onError(error);
      }
    }
  };

  const buttonText = text || getDefaultText();
  const isDisabled = disabled || isLoading;

  return (
    <button
      type="button"
      onClick={handleGoogleLogin}
      disabled={isDisabled}
      className={`
        w-full flex items-center justify-center px-4 py-3
        bg-white border border-gray-300 rounded-md
        text-gray-700 text-sm font-medium
        hover:bg-gray-50 hover:border-gray-400
        focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
        disabled:opacity-50 disabled:cursor-not-allowed
        transition-all duration-200
        ${className}
      `}
      aria-label={buttonText}
    >
      {/* Google Icon SVG - Officiel */}
      <svg
        className="w-5 h-5 mr-3"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          fill="#4285F4"
          d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
        />
        <path
          fill="#34A853"
          d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
        />
        <path
          fill="#FBBC05"
          d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
        />
        <path
          fill="#EA4335"
          d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
        />
      </svg>

      {/* Loading Spinner ou Texte */}
      {isLoading ? (
        <span className="flex items-center">
          <svg
            className="animate-spin -ml-1 mr-2 h-4 w-4 text-gray-700"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            ></circle>
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          Redirecting...
        </span>
      ) : (
        buttonText
      )}
    </button>
  );
};

export default GoogleLoginButton;
