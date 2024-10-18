"use client";

import React from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/providers/AuthContext";
import LoginForm from "@/components/auth/LoginForm";

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();

  const handleLogin = async (
    username: string,
    password: string
  ): Promise<string | undefined> => {
    try {
      await login(username, password);
      console.log("Login successful");
      router.push("/");
      return undefined;
    } catch (error) {
      console.error("Login failed:", error);
      if (error instanceof Error) {
        return error.message;
      }
      return "An unknown error occurred";
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-black py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-white">
            Sign in to your account
          </h2>
        </div>
        <LoginForm onLogin={handleLogin} />
      </div>
    </div>
  );
}
