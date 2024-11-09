import React, { useState } from "react";
import Link from "next/link";
import { useAuth } from "@/providers/AuthContext";
import { AuthError, AuthErrorCode } from "@/libs/authService";

const LoginForm: React.FC = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (isSubmitting) {
      return;
    }

    setIsSubmitting(true);
    setError("");

    try {
      await login(username, password);
    } catch (err) {
      if (err instanceof AuthError) {
        switch (err.code) {
          case AuthErrorCode.INVALID_CREDENTIALS:
            setError("Invalid username or password");
            break;
          case AuthErrorCode.NETWORK_ERROR:
            setError("Network error, please try again");
            break;
          default:
            setError(err.message);
        }
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
          onChange={(e) => setUsername(e.target.value)}
          required
          disabled={isSubmitting}
          className="mt-1 block w-full px-3 py-2 bg-deep-blue border border-gray-600 rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
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
          onChange={(e) => setPassword(e.target.value)}
          required
          disabled={isSubmitting}
          className="mt-1 block w-full px-3 py-2 bg-deep-blue border border-gray-600 rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      {error && (
        <div className="p-3 bg-red-100 border border-red-400 text-red-700 rounded">
          {error}
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
          {isSubmitting ? "Logging in..." : "Log In"}
        </button>
      </div>
      <div className="flex items-center justify-between">
        <Link
          href="/forgot-password"
          className="text-sm text-blue-400 hover:text-blue-300"
        >
          Forgot password?
        </Link>
        <Link
          href="/signup"
          className="text-sm text-blue-400 hover:text-blue-300"
        >
          Not a user? Sign up
        </Link>
      </div>
    </form>
  );
};

export default LoginForm;
