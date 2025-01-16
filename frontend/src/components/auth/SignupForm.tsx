"use client";

import React, { useState } from "react";
import { useAuth } from "@/providers/AuthContext";
import Link from "next/link";
import { AuthError, AuthErrorCode } from "@/libs/authService";
import { AxiosError } from "axios";

const SignupForm: React.FC = () => {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [usernameError, setUsernameError] = useState("");
  const [emailError, setEmailError] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { signup } = useAuth();

  // Validation functions for the password
  const validatePassword = (pass: string): boolean => {
    const hasMinLength = pass.length >= 8;
    if (!hasMinLength) {
      setPasswordError("Password must be at least 8 characters long");
      return false;
    }
    setPasswordError("");
    return true;
  };

  // Validation function for the email
  const validateEmail = (email: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    const isValid = emailRegex.test(email);
    setEmailError(isValid ? "" : "Please enter a valid email address");
    return isValid;
  };

  // Validation function for the username
  const validateUsername = (username: string): boolean => {
    if (username.length < 3) {
      setUsernameError("Username must be at least 3 characters long");
      return false;
    }
    if (username.length > 50) {
      setUsernameError("Username must be less than 50 characters");
      return false;
    }
    setUsernameError("");
    return true;
  };

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (isSubmitting) {
      return;
    }

    // Reset all errors
    setError("");
    setUsernameError("");
    setEmailError("");
    setPasswordError("");

    // Validate all fields
    const isUsernameValid = validateUsername(username);
    const isEmailValid = validateEmail(email);
    const isPasswordValid = validatePassword(password);

    if (!isUsernameValid || !isEmailValid || !isPasswordValid) {
      return;
    }

    setIsSubmitting(true);

    try {
      await signup(username, email, password);
      // Redirection is handled in AuthContext
    } catch (err) {
      console.error("Signup error:", err);

      if (err instanceof AuthError) {
        switch (err.code) {
          case AuthErrorCode.USERNAME_EXISTS:
            setUsernameError("This username is already taken");
            break;
          case AuthErrorCode.EMAIL_EXISTS:
            setEmailError("This email is already registered");
            break;
          case AuthErrorCode.INVALID_INPUT:
            // Check if the error message contains specific details
            if (err.message.includes("username")) {
              setUsernameError(err.message);
            } else if (err.message.includes("email")) {
              setEmailError(err.message);
            } else if (err.message.includes("password")) {
              setPasswordError(err.message);
            } else {
              setError("Please check your input and try again");
            }
            break;
          case AuthErrorCode.SIGNUP_ERROR:
            setError(err.message);
            break;
          case AuthErrorCode.NETWORK_ERROR:
            setError("Network error, please try again");
            break;
          case AuthErrorCode.SERVER_ERROR:
            setError("Server error, please try again later");
            break;
          default:
            setError(err.message || "An unexpected error occurred");
        }
      } else if (err instanceof AxiosError) {
        setError("Server error. Please try again later.");
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
          htmlFor="username"
          className="block text-sm font-medium text-gray-300 mb-2"
        >
          Username
        </label>
        <input
          id="username"
          type="text"
          value={username}
          required
          minLength={3}
          maxLength={50}
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 ${
            usernameError ? "border-red-500" : "border-gray-600"
          }`}
        />
        {usernameError && (
          <p className="mt-1 text-sm text-red-500" role="alert">
            {usernameError}
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
          value={email}
          required
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 ${
            emailError ? "border-red-500" : "border-gray-600"
          }`}
        />
        {emailError && (
          <p className="mt-1 text-sm text-red-500" role="alert">
            {emailError}
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
          value={password}
          required
          minLength={8}
          disabled={isSubmitting}
          className={`mt-1 block w-full px-3 py-2 bg-deep-blue border rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 ${
            passwordError ? "border-red-500" : "border-gray-600"
          }`}
        />
        {passwordError && (
          <p className="mt-1 text-sm text-red-500" role="alert">
            {passwordError}
          </p>
        )}
        <p className="mt-1 text-sm text-gray-400">
          Must be at least 8 characters long
        </p>
      </div>

      {error && (
        <div
          className="p-3 bg-red-100 border border-red-400 text-red-700 rounded"
          role="alert"
        >
          {error}
        </div>
      )}

      <div>
        <button
          type="submit"
          disabled={
            isSubmitting || !!usernameError || !!emailError || !!passwordError
          }
          className={`w-full px-4 py-2 bg-gradient-blue text-white rounded-md ${
            isSubmitting || !!usernameError || !!emailError || !!passwordError
              ? "opacity-50 cursor-not-allowed"
              : "hover:bg-blue-700"
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
