import ResetPasswordForm from "@/components/auth/ResetPasswordForm";

export default function ResetPasswordPage() {
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
              <rect width="18" height="18" x="3" y="3" rx="2" />
              <path d="M7 7v10" />
              <path d="M17 7v10" />
              <path d="M7 7l5-5 5 5" />
            </svg>
          </div>
          <h1 className="text-3xl font-bold text-foreground">Reset Password</h1>
          <p className="text-muted-foreground text-sm text-center">
            Please enter your new password below.
          </p>
        </div>

        <ResetPasswordForm />
      </div>
    </div>
  );
}
