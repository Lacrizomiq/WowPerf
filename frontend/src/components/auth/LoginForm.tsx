"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { useAuth } from "@/providers/AuthContext";
import { AuthError, AuthErrorCode } from "@/libs/authService";
import { useRouter } from "next/navigation";

const LoginForm: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { login, isAuthenticated } = useAuth();
  const router = useRouter();

  // Effect to redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      router.push("/profile");
    }
  }, [isAuthenticated, router]);

  const clearError = () => {
    if (error) {
      setError("");
    }
  };

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
    clearError();
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
    clearError();
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (isSubmitting) {
      return;
    }

    // Basic validation
    if (!email.trim() || !password.trim()) {
      setError("Please fill in all fields");
      return;
    }

    setIsSubmitting(true);
    setError("");

    try {
      await login(email, password);
      // The redirection is handled by the useEffect
    } catch (err) {
      console.error("Login error:", err);

      if (err instanceof AuthError) {
        switch (err.code) {
          case AuthErrorCode.INVALID_CREDENTIALS:
            setError("Invalid email or password");
            break;
          case AuthErrorCode.NETWORK_ERROR:
            setError(
              "Network error. Please check your connection and try again"
            );
            break;
          case AuthErrorCode.SERVER_ERROR:
            setError("Server error. Please try again later");
            break;
          case AuthErrorCode.INVALID_INPUT:
            setError("Please check your input and try again");
            break;
          case AuthErrorCode.LOGIN_ERROR:
            setError(err.message || "Login failed. Please try again");
            break;
          default:
            setError("An unexpected error occurred. Please try again");
        }
      } else if (err instanceof Error) {
        setError(err.message || "An unexpected error occurred");
      } else {
        setError("An unexpected error occurred");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label
          htmlFor="email"
          className="block text-sm font-medium text-gray-300 mb-2"
        >
          Email
        </label>
        <input
          id="email"
          type="text"
          value={email}
          onChange={handleEmailChange}
          required
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border ${
            error ? "border-red-500" : "border-gray-600"
          } rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500`}
          placeholder="Enter your email"
        />
      </div>

      <div>
        <label
          htmlFor="password"
          className="block text-sm font-medium text-gray-300 mb-2"
        >
          Password
        </label>
        <input
          id="password"
          type="password"
          value={password}
          onChange={handlePasswordChange}
          required
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border ${
            error ? "border-red-500" : "border-gray-600"
          } rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500`}
          placeholder="Enter your password"
        />
      </div>

      {error && (
        <div
          className="p-3 bg-red-100 border border-red-400 text-red-700 rounded relative"
          role="alert"
        >
          <span className="block sm:inline">{error}</span>
        </div>
      )}

      <div>
        <button
          type="submit"
          disabled={isSubmitting}
          className={`w-full px-4 py-2 bg-gradient-blue text-white rounded-md transition-all duration-200 ease-in-out
            ${
              isSubmitting
                ? "opacity-50 cursor-not-allowed"
                : "hover:bg-blue-700 active:bg-blue-800"
            } focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-800`}
        >
          {isSubmitting ? (
            <span className="flex items-center justify-center">
              <span className="mr-2">Logging in...</span>
            </span>
          ) : (
            "Log In"
          )}
        </button>
      </div>

      <div className="flex items-center justify-between">
        <Link
          href="/forgot-password"
          className="text-sm text-blue-400 hover:text-blue-300 transition-colors duration-200"
        >
          Forgot password?
        </Link>
        <Link
          href="/signup"
          className="text-sm text-blue-400 hover:text-blue-300 transition-colors duration-200"
        >
          Not a user? Sign up
        </Link>
      </div>
    </form>
  );
};

export default LoginForm;
