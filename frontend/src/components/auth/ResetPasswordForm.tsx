"use client";

import React, { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import api from "@/libs/api";
import { toast } from "react-hot-toast";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertTriangle, Loader2 } from "lucide-react";

type ResetPasswordFormData = {
  password: string;
  confirmPassword: string;
};

const ResetPasswordForm = () => {
  const [isValidToken, setIsValidToken] = React.useState<boolean | null>(null);
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams?.get("token");

  const {
    register,
    handleSubmit,
    watch,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<ResetPasswordFormData>({
    defaultValues: {
      password: "",
      confirmPassword: "",
    },
  });

  // Watch password for confirmation validation
  const password = watch("password");

  useEffect(() => {
    const validateToken = async () => {
      if (!token) {
        setIsValidToken(false);
        return;
      }

      try {
        await api.get(`/auth/validate-reset-token?token=${token}`);
        setIsValidToken(true);
      } catch (err) {
        setIsValidToken(false);
      }
    };

    validateToken();
  }, [token]);

  const onSubmit = async (data: ResetPasswordFormData) => {
    const toastId = toast.loading("Resetting password...");

    try {
      await api.post("/auth/reset-password", {
        token: token,
        new_password: data.password,
      });

      toast.success("Password reset successfully!", {
        id: toastId,
        duration: 2000,
      });

      setTimeout(() => {
        router.push("/login?reset=success");
      }, 2000);
    } catch (err: any) {
      const errorMessage =
        err.response?.data?.error ||
        "An error occurred while resetting your password";

      setError("root", {
        type: "manual",
        message: errorMessage,
      });

      toast.error(errorMessage, {
        id: toastId,
      });
    }
  };

  // Loading state
  if (isValidToken === null) {
    return (
      <div className="flex items-center justify-center space-x-2">
        <Loader2 className="h-5 w-5 animate-spin text-primary" />
        <span className="text-muted-foreground">Validating reset token...</span>
      </div>
    );
  }

  // Invalid token state
  if (isValidToken === false) {
    return (
      <div className="space-y-6">
        <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-6">
          <div className="flex items-start space-x-3">
            <div className="flex-shrink-0">
              <AlertTriangle className="h-6 w-6 text-destructive" />
            </div>
            <div className="flex-1">
              <h3 className="text-lg font-medium text-destructive mb-2">
                Invalid or Expired Reset Link
              </h3>
              <p className="text-sm text-destructive/80 mb-4">
                This password reset link is invalid or has expired. Please
                request a new one.
              </p>
              <Button
                onClick={() => router.push("/forgot-password")}
                variant="outline"
                className="border-destructive/30 text-destructive hover:bg-destructive/10"
              >
                Request New Reset Link
              </Button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="space-y-2">
        <label
          htmlFor="password"
          className="block text-sm font-medium text-foreground"
        >
          New Password
        </label>
        <Input
          id="password"
          type="password"
          disabled={isSubmitting}
          placeholder="Enter your new password"
          className={errors.password || errors.root ? "border-destructive" : ""}
          {...register("password", {
            required: "Password is required",
            minLength: {
              value: 8,
              message: "Password must be at least 8 characters long",
            },
            maxLength: {
              value: 32,
              message: "Password must be less than 32 characters",
            },
          })}
        />
        {errors.password && (
          <p className="text-sm text-destructive" role="alert">
            {errors.password.message}
          </p>
        )}
        <p className="text-xs text-muted-foreground">
          Must be between 8 and 32 characters long
        </p>
      </div>

      <div className="space-y-2">
        <label
          htmlFor="confirmPassword"
          className="block text-sm font-medium text-foreground"
        >
          Confirm New Password
        </label>
        <Input
          id="confirmPassword"
          type="password"
          disabled={isSubmitting}
          placeholder="Confirm your new password"
          className={
            errors.confirmPassword || errors.root ? "border-destructive" : ""
          }
          {...register("confirmPassword", {
            required: "Please confirm your password",
            validate: (value) =>
              value === password || "The passwords do not match",
          })}
        />
        {errors.confirmPassword && (
          <p className="text-sm text-destructive" role="alert">
            {errors.confirmPassword.message}
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
          {isSubmitting ? (
            <span className="flex items-center">
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Resetting...
            </span>
          ) : (
            "Reset Password"
          )}
        </Button>
      </div>

      <div className="text-center">
        <Button
          type="button"
          variant="ghost"
          onClick={() => router.push("/login")}
          className="text-primary hover:text-primary/80"
        >
          Back to login
        </Button>
      </div>
    </form>
  );
};

export default ResetPasswordForm;
