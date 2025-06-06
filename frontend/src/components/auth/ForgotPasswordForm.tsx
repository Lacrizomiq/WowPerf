"use client";

import React from "react";
import { useRouter } from "next/navigation";
import api from "@/libs/api";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { CheckCircle } from "lucide-react";

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
      <div className="space-y-6">
        {/* Success message avec le thème harmonisé */}
        <div className="bg-green-500/10 border border-green-500/20 rounded-lg p-6">
          <div className="flex items-start space-x-3">
            <div className="flex-shrink-0">
              <CheckCircle className="h-6 w-6 text-green-500" />
            </div>
            <div className="flex-1">
              <h3 className="text-lg font-medium text-green-400 mb-2">
                Password Reset Email Sent
              </h3>
              <p className="text-sm text-green-300/80 mb-4">
                If an account exists with the email{" "}
                <span className="font-medium">{submittedEmail}</span>, you will
                receive password reset instructions shortly.
              </p>
              <Button
                onClick={() => router.push("/login")}
                variant="outline"
                className="border-green-500/30 text-green-400 hover:bg-green-500/10"
              >
                Return to login
              </Button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="space-y-2">
        <label
          htmlFor="email"
          className="block text-sm font-medium text-foreground"
        >
          Email Address
        </label>
        <Input
          id="email"
          type="email"
          disabled={isSubmitting}
          placeholder="Enter your email address"
          className={errors.email || errors.root ? "border-destructive" : ""}
          {...register("email", {
            required: "Email is required",
            pattern: {
              value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
              message: "Please enter a valid email address",
            },
          })}
        />
        {errors.email && (
          <p className="text-sm text-destructive" role="alert">
            {errors.email.message}
          </p>
        )}
      </div>

      {errors.root && (
        <Alert variant="destructive">
          <AlertDescription>{errors.root.message}</AlertDescription>
        </Alert>
      )}

      <div>
        <Button type="submit" disabled={isSubmitting} className="w-full">
          {isSubmitting ? "Sending..." : "Send Reset Instructions"}
        </Button>
      </div>

      <div className="text-center">
        <Button
          type="button"
          variant="ghost"
          onClick={() => router.push("/login")}
          className="text-white bg-slate-800/40 border border-lg hover:text-purple-600"
        >
          Back to login
        </Button>
      </div>
    </form>
  );
};

export default ForgotPasswordForm;
