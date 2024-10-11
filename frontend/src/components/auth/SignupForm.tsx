import React, { useState } from "react";
import { useAuth } from "@/providers/AuthContext";
import { useRouter } from "next/navigation";
import Link from "next/link";

const SignupForm: React.FC = () => {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const { signup } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await signup(username, email, password);
      router.push("/login");
    } catch (err) {
      setError("Failed to sign up. Please try again.");
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
          className="mt-1 block w-full px-3 py-2 bg-deep-blue border border-gray-600 rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <p className="mt-1 text-sm text-gray-400">
          Enter a valid email address.
        </p>
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
          pattern="(?=.*[!@#$%^&*()_+]).{8,}"
          className="mt-1 block w-full px-3 py-2 bg-deep-blue border border-gray-600 rounded-md text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <p className="mt-1 text-sm text-gray-400">
          Must be at least 8 characters long and contain at least one special
          character (!@#$%^&*()_+).
        </p>
      </div>
      {error && <p className="text-red-500 text-sm">{error}</p>}
      <div>
        <button
          type="submit"
          className="w-full px-4 py-2 bg-gradient-blue text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-800"
        >
          Sign Up
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
