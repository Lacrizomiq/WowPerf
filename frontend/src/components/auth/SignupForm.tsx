"use client";

import React, { useState, useRef } from "react";
import { useAuth } from "@/providers/AuthContext";
import Link from "next/link";
import { AuthError, AuthErrorCode } from "@/libs/authService";
import { AxiosError } from "axios";
import { useForm } from "react-hook-form";
import HCaptcha from "@hcaptcha/react-hcaptcha";
import GoogleLoginButton from "@/components/Shared/GoogleLoginButton";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Alert, AlertDescription } from "@/components/ui/alert";

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

  // States for HCAPTCHA
  const [captchaToken, setCaptchaToken] = useState<string>("");
  const captchaRef = useRef<HCaptcha>(null);

  // Vérifier si hCaptcha est activé (présence de la clé dans l'env)
  const captchaEnabled = !!process.env.NEXT_PUBLIC_HCAPTCHA_SITE_KEY;

  const siteKey = process.env.NEXT_PUBLIC_HCAPTCHA_SITE_KEY;
  console.log("HCAPTCHA Debug:", {
    preview: siteKey ? `[${siteKey.substring(0, 5)}...]` : "[undefined]",
    length: siteKey?.length || 0,
    hasQuotes: siteKey
      ? {
          startsWithQuote: siteKey.startsWith('"'),
          endsWithQuote: siteKey.endsWith('"'),
          firstChar: `[${siteKey.charAt(0)}]`,
          lastChar: `[${siteKey.charAt(siteKey.length - 1)}]`,
        }
      : null,
    captchaEnabled,
  });

  const onCaptchaVerify = (token: string) => {
    setCaptchaToken(token);
  };

  const onCaptchaExpire = () => {
    setCaptchaToken("");
  };

  const onCaptchaError = (err: string) => {
    console.error("[CAPTCHA] Error:", err);
    setCaptchaToken("");
  };

  const resetCaptcha = () => {
    if (captchaRef.current) {
      captchaRef.current.resetCaptcha();
      setCaptchaToken("");
    }
  };

  // Form submission handler using React Hook Form
  const onSubmit = async (data: SignupFormData) => {
    try {
      await signup(data.username, data.email, data.password, captchaToken);
      // Redirection is handled in AuthContext
    } catch (err) {
      console.error("Signup error:", err);

      if (captchaEnabled) {
        resetCaptcha();
      }

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
          case AuthErrorCode.CAPTCHA_REQUIRED:
            setError("root", {
              type: "custom",
              message: "Please complete the captcha verification",
            });
            break;
          case AuthErrorCode.CAPTCHA_INVALID:
            setError("root", {
              type: "custom",
              message: "Captcha verification failed. Please try again.",
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

  const isFormValid = captchaEnabled ? captchaToken !== "" : true;

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
          Register an account
        </h1>
        <p className="text-muted-foreground text-sm">
          Join WoW Perf to track your performance
        </p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <div className="space-y-2">
          <label
            htmlFor="username"
            className="block text-sm font-medium text-foreground"
          >
            Username
          </label>
          <Input
            id="username"
            type="text"
            disabled={isSubmitting}
            placeholder="Enter your username"
            className={`
              border-2 bg-slate-800/40 text-white placeholder:text-slate-400
              focus:border-primary focus:bg-slate-600/50
              ${errors.username ? "border-destructive" : ""}
              `}
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
            <p className="text-lg text-red-600" role="alert">
              {errors.username.message}
            </p>
          )}
          <p className="text-xs text-muted-foreground">
            Must be between 3 and 50 characters
          </p>
        </div>

        <div className="space-y-2">
          <label
            htmlFor="email"
            className="block text-sm font-medium text-foreground"
          >
            Email
          </label>
          <Input
            id="email"
            type="email"
            disabled={isSubmitting}
            placeholder="Enter your email"
            className={`
              border-2 bg-slate-800/40 text-white placeholder:text-slate-400
              focus:border-primary focus:bg-slate-600/50
              ${errors.email ? "border-destructive" : ""}
              `}
            {...register("email", {
              required: "Email is required",
              pattern: {
                value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                message: "Please enter a valid email address",
              },
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
              ${errors.password ? "border-destructive" : ""}
            `}
            {...register("password", {
              required: "Password is required",
              minLength: {
                value: 8,
                message: "Password must be at least 8 characters long",
              },
            })}
          />
          {errors.password && (
            <p className="text-lg text-red-600" role="alert">
              {errors.password.message}
            </p>
          )}
          <p className="text-xs text-muted-foreground">
            Must be at least 8 characters long
          </p>
        </div>

        {/* HCAPTCHA */}
        {captchaEnabled && (
          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">
              Security Verification
            </label>
            <div className="flex justify-center">
              <HCaptcha
                ref={captchaRef}
                sitekey={process.env.NEXT_PUBLIC_HCAPTCHA_SITE_KEY!}
                onVerify={onCaptchaVerify}
                onExpire={onCaptchaExpire}
                onError={onCaptchaError}
                theme="dark"
              />
            </div>
            {captchaEnabled && !captchaToken && (
              <p className="text-xs text-muted-foreground text-center">
                Please complete the security verification above
              </p>
            )}
          </div>
        )}

        {errors.root && (
          <Alert variant="destructive">
            <AlertDescription>{errors.root.message}</AlertDescription>
          </Alert>
        )}

        <div>
          <Button
            type="submit"
            disabled={isSubmitting || !isFormValid}
            className="w-full"
          >
            {isSubmitting ? "Creating account..." : "Sign Up"}
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
            variant="signup"
            disabled={isSubmitting}
            onError={(error: any) => {
              setError("root", {
                type: "custom",
                message: error.message || "Google sign-up failed",
              });
            }}
          />
        </div>

        <div className="flex items-center justify-center">
          <Link
            href="/login"
            className="text-sm text-white hover:text-purple-600 transition-colors duration-200"
          >
            Already have an account? Sign in
          </Link>
        </div>
      </form>
    </div>
  );
};

export default SignupForm;
