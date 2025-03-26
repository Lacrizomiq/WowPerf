"use client";

import React, { useEffect } from "react";
import Link from "next/link";
import { useAuth } from "@/providers/AuthContext";
import { AuthError, AuthErrorCode } from "@/libs/authService";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";

// Define the form data type
type LoginFormData = {
  email: string;
  password: string;
};

const LoginForm: React.FC = () => {
  // Initialize React Hook Form with validation
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
    clearErrors,
  } = useForm<LoginFormData>({
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const { login, isAuthenticated } = useAuth();
  const router = useRouter();

  // Effect to redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      router.push("/profile");
    }
  }, [isAuthenticated, router]);

  // Form submission handler using React Hook Form
  const onSubmit = async (data: LoginFormData) => {
    try {
      await login(data.email, data.password);
      // Redirection is handled by the useEffect
    } catch (err) {
      console.error("Login error:", err);

      if (err instanceof AuthError) {
        switch (err.code) {
          case AuthErrorCode.INVALID_CREDENTIALS:
            // Set a form-level error for invalid credentials
            setError("root", {
              type: "manual",
              message: "Invalid email or password",
            });
            break;
          case AuthErrorCode.NETWORK_ERROR:
            setError("root", {
              type: "manual",
              message:
                "Network error. Please check your connection and try again",
            });
            break;
          case AuthErrorCode.SERVER_ERROR:
            setError("root", {
              type: "manual",
              message: "Server error. Please try again later",
            });
            break;
          case AuthErrorCode.INVALID_INPUT:
            setError("root", {
              type: "manual",
              message: "Please check your input and try again",
            });
            break;
          case AuthErrorCode.LOGIN_ERROR:
            setError("root", {
              type: "manual",
              message: err.message || "Login failed. Please try again",
            });
            break;
          default:
            setError("root", {
              type: "manual",
              message: "An unexpected error occurred. Please try again",
            });
        }
      } else if (err instanceof Error) {
        setError("root", {
          type: "manual",
          message: err.message || "An unexpected error occurred",
        });
      } else {
        setError("root", {
          type: "manual",
          message: "An unexpected error occurred",
        });
      }
    }
  };

  // Clear errors when user changes input
  const handleInputChange = () => {
    if (errors.root) {
      clearErrors("root");
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
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
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border ${
            errors.email || errors.root ? "border-red-500" : "border-gray-600"
          } rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500`}
          placeholder="Enter your email"
          {...register("email", {
            required: "Email is required",
            onChange: handleInputChange,
          })}
        />
        {errors.email && (
          <p className="mt-1 text-sm text-red-500" role="alert">
            {errors.email.message}
          </p>
        )}
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
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border ${
            errors.password || errors.root
              ? "border-red-500"
              : "border-gray-600"
          } rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500`}
          placeholder="Enter your password"
          {...register("password", {
            required: "Password is required",
            onChange: handleInputChange,
          })}
        />
        {errors.password && (
          <p className="mt-1 text-sm text-red-500" role="alert">
            {errors.password.message}
          </p>
        )}
      </div>

      {errors.root && (
        <div
          className="p-3 bg-red-100 border border-red-400 text-red-700 rounded relative"
          role="alert"
        >
          <span className="block sm:inline">{errors.root.message}</span>
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
