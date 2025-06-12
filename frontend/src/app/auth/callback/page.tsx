"use client";

import React from "react";
import { useGoogleCallback } from "@/hooks/useGoogleAuth";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { CheckCircle, AlertTriangle } from "lucide-react";

const CallbackPage = () => {
  const { isProcessing, error, isSuccess, errorDisplay } = useGoogleCallback();

  // 🔄 PROCESSING - Pendant que le hook traite
  if (isProcessing) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="max-w-md w-full space-y-8 p-8">
          <div className="text-center">
            {/* Logo */}
            <div className="bg-primary text-primary-foreground p-2 rounded-lg mb-6 w-fit mx-auto">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                className="h-6 w-6"
              >
                <rect width="18" height="18" x="3" y="3" rx="2" />
                <path d="M3 9h18" />
                <path d="M9 21V9" />
              </svg>
            </div>

            {/* Spinner animé avec les couleurs du thème */}
            <div className="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-muted border-t-primary mb-6"></div>

            <h2 className="text-3xl font-extrabold text-foreground mb-4">
              Signing you in...
            </h2>

            <p className="text-sm text-muted-foreground">
              Please wait while we complete your Google sign-in
            </p>
          </div>
        </div>
      </div>
    );
  }

  // ✅ SUCCESS - Login réussi (avant redirection automatique)
  if (isSuccess) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="max-w-md w-full space-y-8 p-8">
          <div className="text-center">
            {/* Logo */}
            <div className="bg-primary text-primary-foreground p-2 rounded-lg mb-6 w-fit mx-auto">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                className="h-6 w-6"
              >
                <rect width="18" height="18" x="3" y="3" rx="2" />
                <path d="M3 9h18" />
                <path d="M9 21V9" />
              </svg>
            </div>

            {/* Success icon avec les couleurs du thème */}
            <div className="mx-auto h-12 w-12 bg-green-500/10 rounded-full flex items-center justify-center mb-6">
              <CheckCircle className="h-6 w-6 text-green-500" />
            </div>

            <h2 className="text-3xl font-extrabold text-foreground mb-4">
              Success!
            </h2>

            <p className="text-sm text-muted-foreground">
              Redirecting to your dashboard...
            </p>
          </div>
        </div>
      </div>
    );
  }

  // ❌ ERROR - Gestion des erreurs OAuth
  if (error && errorDisplay) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="max-w-md w-full space-y-8 p-8">
          <div className="text-center">
            {/* Logo */}
            <div className="bg-primary text-primary-foreground p-2 rounded-lg mb-6 w-fit mx-auto">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                className="h-6 w-6"
              >
                <rect width="18" height="18" x="3" y="3" rx="2" />
                <path d="M3 9h18" />
                <path d="M9 21V9" />
              </svg>
            </div>

            {/* Error icon avec les couleurs du thème */}
            <div className="mx-auto h-12 w-12 bg-destructive/10 rounded-full flex items-center justify-center mb-6">
              <AlertTriangle className="h-6 w-6 text-destructive" />
            </div>

            <h2 className="text-3xl font-extrabold text-foreground mb-4">
              {errorDisplay.title}
            </h2>

            <Alert variant="destructive" className="mb-6">
              <AlertDescription>{errorDisplay.message}</AlertDescription>
            </Alert>

            {/* Actions selon l'erreur */}
            <div className="space-y-3">
              {errorDisplay.actions?.map((action, index) => (
                <div key={index}>
                  {action.href ? (
                    <Button
                      asChild
                      className="w-full"
                      variant={index === 0 ? "default" : "outline"}
                    >
                      <a href={action.href}>{action.label}</a>
                    </Button>
                  ) : (
                    <Button
                      onClick={action.onClick}
                      className="w-full"
                      variant={index === 0 ? "default" : "outline"}
                    >
                      {action.label}
                    </Button>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    );
  }

  // 🔄 Fallback - Ne devrait jamais arriver
  return (
    <div className="min-h-screen bg-background flex items-center justify-center">
      <div className="text-center">
        <p className="text-foreground">Processing authentication...</p>
      </div>
    </div>
  );
};

export default CallbackPage;
