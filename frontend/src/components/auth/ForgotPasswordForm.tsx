"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import api from "@/libs/api";

const ForgotPasswordForm = () => {
  const [email, setEmail] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isSubmitting) return;

    setIsSubmitting(true);
    setError("");

    try {
      await api.post("/auth/forgot-password", { email });
      setSuccess(true);
    } catch (err: any) {
      setError(
        err.response?.data?.error ||
          "An error occurred while processing your request. Please try again."
      );
    } finally {
      setIsSubmitting(false);
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
                If an account exists with the email {email}, you will receive
                password reset instructions shortly.
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
    <form onSubmit={handleSubmit} className="space-y-6">
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
          value={email}
          onChange={(e) => {
            setEmail(e.target.value);
            setError("");
          }}
          required
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border ${
            error ? "border-red-500" : "border-gray-600"
          } rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500`}
          placeholder="Enter your email address"
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
