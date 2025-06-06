import ForgotPasswordForm from "@/components/auth/ForgotPasswordForm";

export default function ForgotPasswordPage() {
  return (
    <div className="min-h-screen bg-background flex items-center justify-center">
      <div className="max-w-md w-full space-y-8 p-8">
        {/* Header harmonisé avec le thème */}
        <div className="flex flex-col items-center space-y-2 mb-8">
          <div className="bg-primary text-primary-foreground p-2 rounded-lg mb-2">
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
              <path d="M9 12l2 2 4-4" />
              <path d="M21 12c-1 0-3-1-3-3s2-3 3-3 3 1 3 3-2 3-3 3" />
              <path d="M3 12c1 0 3-1 3-3s-2-3-3-3-3 1-3 3 2 3 3 3" />
            </svg>
          </div>
          <h1 className="text-3xl font-bold text-foreground">
            Forgot Password
          </h1>
          <p className="text-muted-foreground text-sm text-center">
            Enter your email address and we&apos;ll send you instructions to
            reset your password.
          </p>
        </div>

        <ForgotPasswordForm />
      </div>
    </div>
  );
}
