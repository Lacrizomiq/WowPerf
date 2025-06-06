"use client";

import React, { useEffect } from "react";
import Link from "next/link";
import { useAuth } from "@/providers/AuthContext";
import { AuthError, AuthErrorCode } from "@/libs/authService";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import GoogleLoginButton from "@/components/Shared/GoogleLoginButton";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Alert, AlertDescription } from "@/components/ui/alert";

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
    <div className="w-full space-y-6">
      {/* Header */}
      <div className="flex flex-col items-center space-y-2 mb-8">
        <div className="bg-primary text-primary-foreground p-2 rounded-lg mb-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="h-6 w-6"
          >
            <rect width="18" height="18" x="3" y="3" rx="2" />
            <path d="M3 9h18" />
            <path d="M9 21V9" />
          </svg>
        </div>
        <h1 className="text-3xl font-bold text-foreground">
          Sign in to your account
        </h1>
        <p className="text-muted-foreground text-sm">
          Welcome back! Please sign in to continue
        </p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <div className="space-y-2">
          <label
            htmlFor="email"
            className="block text-sm font-medium text-foreground"
          >
            Email
          </label>
          <Input
            id="email"
            type="text"
            disabled={isSubmitting}
            placeholder="Enter your email"
            className={`
              border-2 bg-slate-800/40 text-white placeholder:text-slate-400
              focus:border-primary focus:bg-slate-600/50
              ${
                errors.email || errors.root
                  ? "border-red-500"
                  : "border-slate-600"
              }
            `}
            {...register("email", {
              required: "Email is required",
              onChange: handleInputChange,
            })}
          />
          {errors.email && (
            <p className="text-lg text-red-600" role="alert">
              {errors.email.message}
            </p>
          )}
        </div>

        <div className="space-y-2">
          <label
            htmlFor="password"
            className="block text-sm font-medium text-foreground"
          >
            Password
          </label>
          <Input
            id="password"
            type="password"
            disabled={isSubmitting}
            placeholder="Enter your password"
            className={`
              border-2 bg-slate-800/40 text-white placeholder:text-slate-400
              focus:border-primary focus:bg-slate-600/50
              ${errors.password || errors.root ? "border-destructive" : ""}
            `}
            {...register("password", {
              required: "Password is required",
              onChange: handleInputChange,
            })}
          />
          {errors.password && (
            <p className="text-lg text-red-600" role="alert">
              {errors.password.message}
            </p>
          )}
        </div>

        {errors.root && (
          <Alert variant="destructive">
            <AlertDescription>{errors.root.message}</AlertDescription>
          </Alert>
        )}

        <div>
          <Button type="submit" disabled={isSubmitting} className="w-full">
            {isSubmitting ? "Logging in..." : "Log In"}
          </Button>
        </div>

        {/* Séparateur */}
        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <Separator className="w-full border" />
          </div>
          <div className="relative flex justify-center text-xs">
            <span className="bg-slate-800 px-2 text-white">OR</span>
          </div>
        </div>

        {/* Bouton Google */}
        <div>
          <GoogleLoginButton
            variant="signin"
            disabled={isSubmitting}
            onError={(error: any) => {
              setError("root", {
                type: "manual",
                message: error.message || "Google sign-in failed",
              });
            }}
          />
        </div>

        <div className="flex items-center justify-between">
          <Link
            href="/forgot-password"
            className="text-sm text-white hover:text-purple-600 transition-colors duration-200"
          >
            Forgot password?
          </Link>
          <Link
            href="/signup"
            className="text-sm text-white hover:text-purple-600 transition-colors duration-200"
          >
            Not a user? Sign up
          </Link>
        </div>
      </form>
    </div>
  );
};

export default LoginForm;
