"use client";

import React from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";
import { SignupForm } from "@/components/auth/SignupForm";

export default function SignupPage() {
  const router = useRouter();
  const { signup } = useAuth();

  const handleSignup = async (
    username: string,
    email: string,
    password: string
  ) => {
    try {
      await signup(username, email, password);
      router.push("/login");
    } catch (error) {
      console.error("Signup failed:", error);
      // Handle error (e.g., show error message to user)
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-black py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Create your account
          </h2>
        </div>
        <SignupForm onSignup={handleSignup} />
      </div>
    </div>
  );
}
