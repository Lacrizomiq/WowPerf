"use client";

import React from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/providers/AuthContext";
import SignupForm from "@/components/auth/SignupForm";

export default function SignupPage() {
  const router = useRouter();
  const { signup } = useAuth();

  const handleSignup = async (
    username: string,
    email: string,
    password: string
  ): Promise<string | undefined> => {
    try {
      await signup(username, email, password);
      console.log("Signup successful");
      router.push("/login");
      return undefined;
    } catch (error) {
      console.error("Signup failed:", error);
      return "Invalid credentials. Please try again.";
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
        <SignupForm />
      </div>
    </div>
  );
}
