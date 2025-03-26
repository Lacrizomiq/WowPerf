"use client";

import React, { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import api from "@/libs/api";
import { toast } from "react-hot-toast";
import { useForm } from "react-hook-form";

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

  if (isValidToken === null) {
    return (
      <div className="text-center">
        <div className="text-gray-300">Validating reset token...</div>
      </div>
    );
  }

  if (isValidToken === false) {
    return (
      <div className="rounded-md bg-red-50 p-4">
        <div className="flex">
          <div className="ml-3">
            <h3 className="text-sm font-medium text-red-800">
              Invalid or Expired Reset Link
            </h3>
            <div className="mt-2 text-sm text-red-700">
              <p>
                This password reset link is invalid or has expired. Please
                request a new one.
              </p>
            </div>
            <div className="mt-4">
              <button
                type="button"
                onClick={() => router.push("/forgot-password")}
                className="text-sm font-medium text-red-800 hover:text-red-700"
              >
                Request New Reset Link
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div>
        <label
          htmlFor="password"
          className="block text-sm font-medium text-gray-300 mb-2"
        >
          New Password
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
          placeholder="Enter your new password"
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
          <p className="mt-1 text-sm text-red-500" role="alert">
            {errors.password.message}
          </p>
        )}
      </div>

      <div>
        <label
          htmlFor="confirmPassword"
          className="block text-sm font-medium text-gray-300 mb-2"
        >
          Confirm New Password
        </label>
        <input
          id="confirmPassword"
          type="password"
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border ${
            errors.confirmPassword || errors.root
              ? "border-red-500"
              : "border-gray-600"
          } rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500`}
          placeholder="Confirm your new password"
          {...register("confirmPassword", {
            required: "Please confirm your password",
            validate: (value) =>
              value === password || "The passwords do not match",
          })}
        />
        {errors.confirmPassword && (
          <p className="mt-1 text-sm text-red-500" role="alert">
            {errors.confirmPassword.message}
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
            } focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2`}
        >
          {isSubmitting ? "Resetting..." : "Reset Password"}
        </button>
      </div>
    </form>
  );
};

export default ResetPasswordForm;
