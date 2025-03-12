"use client";

import React from "react";
import { useRouter } from "next/navigation";
import api from "@/libs/api";
import { useForm } from "react-hook-form";

type ForgotPasswordFormData = {
  email: string;
};

const ForgotPasswordForm = () => {
  const router = useRouter();
  const [success, setSuccess] = React.useState(false);
  const [submittedEmail, setSubmittedEmail] = React.useState("");

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<ForgotPasswordFormData>({
    defaultValues: {
      email: "",
    },
  });

  const onSubmit = async (data: ForgotPasswordFormData) => {
    try {
      await api.post("/auth/forgot-password", { email: data.email });
      setSubmittedEmail(data.email);
      setSuccess(true);
    } catch (err: any) {
      setError("root", {
        type: "manual",
        message:
          err.response?.data?.error ||
          "An error occurred while processing your request. Please try again.",
      });
    }
  };

  if (success) {
    return (
      <div className="rounded-md bg-green-50 p-4">
        <div className="flex">
          <div className="ml-3">
            <h3 className="text-sm font-medium text-green-800">
              Password Reset Email Sent
            </h3>
            <div className="mt-2 text-sm text-green-700">
              <p>
                If an account exists with the email {submittedEmail}, you will
                receive password reset instructions shortly.
              </p>
            </div>
            <div className="mt-4">
              <button
                type="button"
                onClick={() => router.push("/login")}
                className="text-sm font-medium text-green-800 hover:text-green-700"
              >
                Return to login
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
          htmlFor="email"
          className="block text-sm font-medium text-gray-300 mb-2"
        >
          Email Address
        </label>
        <input
          id="email"
          type="email"
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border ${
            errors.email || errors.root ? "border-red-500" : "border-gray-600"
          } rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500`}
          placeholder="Enter your email address"
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
          {isSubmitting ? "Sending..." : "Send Reset Instructions"}
        </button>
      </div>

      <div className="text-center mt-4">
        <button
          type="button"
          onClick={() => router.push("/login")}
          className="text-sm text-blue-400 hover:text-blue-300 transition-colors duration-200"
        >
          Back to login
        </button>
      </div>
    </form>
  );
};

export default ForgotPasswordForm;
