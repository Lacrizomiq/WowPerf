"use client";

import React from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";
import { LoginForm } from "@/components/auth/LoginForm";

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();

  const handleLogin = async (username: string, password: string) => {
    try {
      await login(username, password);
      router.push("/dashboard");
    } catch (error) {
      console.error("Login failed:", error);
      // Handle error (e.g., show error message to user)
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-black py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Sign in to your account
          </h2>
        </div>
        <LoginForm onLogin={handleLogin} />
      </div>
    </div>
  );
}
