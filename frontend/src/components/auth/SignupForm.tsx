"use client";

import React from "react";
import { useAuth } from "@/providers/AuthContext";
import Link from "next/link";
import { AuthError, AuthErrorCode } from "@/libs/authService";
import { AxiosError } from "axios";
import { useForm } from "react-hook-form";

// Define the form data type
type SignupFormData = {
  username: string;
  email: string;
  password: string;
};

const SignupForm: React.FC = () => {
  // Initialize React Hook Form with default values and validation
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<SignupFormData>({
    defaultValues: {
      username: "",
      email: "",
      password: "",
    },
  });

  const { signup } = useAuth();

  // Form submission handler using React Hook Form
  const onSubmit = async (data: SignupFormData) => {
    try {
      await signup(data.username, data.email, data.password);
      // Redirection is handled in AuthContext
    } catch (err) {
      console.error("Signup error:", err);

      if (err instanceof AuthError) {
        switch (err.code) {
          case AuthErrorCode.USERNAME_EXISTS:
            setError("username", {
              type: "custom",
              message: "This username is already taken",
            });
            break;
          case AuthErrorCode.EMAIL_EXISTS:
            setError("email", {
              type: "custom",
              message: "This email is already registered",
            });
            break;
          case AuthErrorCode.INVALID_INPUT:
            // Check if the error message contains specific details
            if (err.message.includes("username")) {
              setError("username", { type: "custom", message: err.message });
            } else if (err.message.includes("email")) {
              setError("email", { type: "custom", message: err.message });
            } else if (err.message.includes("password")) {
              setError("password", { type: "custom", message: err.message });
            } else {
              setError("root", {
                type: "custom",
                message: "Please check your input and try again",
              });
            }
            break;
          case AuthErrorCode.SIGNUP_ERROR:
            setError("root", { type: "custom", message: err.message });
            break;
          case AuthErrorCode.NETWORK_ERROR:
            setError("root", {
              type: "custom",
              message: "Network error, please try again",
            });
            break;
          case AuthErrorCode.SERVER_ERROR:
            setError("root", {
              type: "custom",
              message: "Server error, please try again later",
            });
            break;
          default:
            setError("root", {
              type: "custom",
              message: err.message || "An unexpected error occurred",
            });
        }
      } else if (err instanceof AxiosError) {
        setError("root", {
          type: "custom",
          message: "Server error. Please try again later.",
        });
      } else {
        setError("root", {
          type: "custom",
          message: "An unexpected error occurred",
        });
      }
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div>
        <label
          htmlFor="username"
          className="block text-sm font-medium text-gray-300 mb-2"
        >
          Username
        </label>
        <input
          id="username"
          type="text"
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 ${
            errors.username ? "border-red-500" : "border-gray-600"
          }`}
          {...register("username", {
            required: "Username is required",
            minLength: {
              value: 3,
              message: "Username must be at least 3 characters long",
            },
            maxLength: {
              value: 50,
              message: "Username must be less than 50 characters",
            },
          })}
        />
        {errors.username && (
          <p className="mt-1 text-sm text-red-500" role="alert">
            {errors.username.message}
          </p>
        )}
        <p className="mt-1 text-sm text-gray-400">
          Must be between 3 and 50 characters
        </p>
      </div>

      <div>
        <label
          htmlFor="email"
          className="block text-sm font-medium text-gray-300 mb-2"
        >
          Email
        </label>
        <input
          id="email"
          type="email"
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 ${
            errors.email ? "border-red-500" : "border-gray-600"
          }`}
          {...register("email", {
            required: "Email is required",
            pattern: {
              value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
              message: "Please enter a valid email address",
            },
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
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 ${
            errors.password ? "border-red-500" : "border-gray-600"
          }`}
          {...register("password", {
            required: "Password is required",
            minLength: {
              value: 8,
              message: "Password must be at least 8 characters long",
            },
          })}
        />
        {errors.password && (
          <p className="mt-1 text-sm text-red-500" role="alert">
            {errors.password.message}
          </p>
        )}
        <p className="mt-1 text-sm text-gray-400">
          Must be at least 8 characters long
        </p>
      </div>

      {errors.root && (
        <div
          className="p-3 bg-red-100 border border-red-400 text-red-700 rounded"
          role="alert"
        >
          {errors.root.message}
        </div>
      )}

      <div>
        <button
          type="submit"
          disabled={isSubmitting}
          className={`w-full px-4 py-2 bg-gradient-blue text-white rounded-md ${
            isSubmitting ? "opacity-50 cursor-not-allowed" : "hover:bg-blue-700"
          } focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-800`}
        >
          {isSubmitting ? "Creating account..." : "Sign Up"}
        </button>
      </div>

      <div className="flex items-center justify-center">
        <Link
          href="/login"
          className="text-sm text-blue-400 hover:text-blue-300"
        >
          Already have an account? Sign in
        </Link>
      </div>
    </form>
  );
};

export default SignupForm;
