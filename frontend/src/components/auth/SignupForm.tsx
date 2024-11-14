import React, { useState } from "react";
import { useAuth } from "@/providers/AuthContext";
import Link from "next/link";
import { AuthError, AuthErrorCode } from "@/libs/authService";

const SignupForm: React.FC = () => {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { signup } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (isSubmitting) {
      return;
    }

    setIsSubmitting(true);
    setError("");

    try {
      await signup(username, email, password);
      // La redirection est gérée dans le AuthContext
    } catch (err) {
      if (err instanceof AuthError) {
        switch (err.code) {
          case AuthErrorCode.USERNAME_EXISTS:
            setError("This username is already taken");
            break;
          case AuthErrorCode.EMAIL_EXISTS:
            setError("This email is already registered");
            break;
          case AuthErrorCode.INVALID_INPUT:
            setError("Please check your input and try again");
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
          minLength={3}
          maxLength={50}
          disabled={isSubmitting}
          className="mt-1 block w-full px-3 py-2 bg-deep-blue border border-gray-600 rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <p className="mt-1 text-sm text-gray-400">3 characters min</p>
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
          onChange={(e) => setEmail(e.target.value)}
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
          minLength={8}
          disabled={isSubmitting}
          className="mt-1 block w-full px-3 py-2 bg-deep-blue border border-gray-600 rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <p className="mt-1 text-sm text-gray-400">
          Must be at least 8 characters long and contain at least one special
          character (!@#$%^&*()_+).
        </p>
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
          {isSubmitting ? "Creating account..." : "Sign Up"}
        </button>
      </div>
      <div className="flex items-center justify-center">
        <Link
          href="/login"
          className="text-sm text-blue-400 hover:text-blue-300"
        >
          Already a user? Sign in
        </Link>
      </div>
    </form>
  );
};

export default SignupForm;
