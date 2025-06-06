import SignupForm from "@/components/auth/SignupForm";

export default function SignupPage() {
  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 relative z-20">
      <div className="max-w-md w-full space-y-8">
        <SignupForm />
      </div>
    </div>
  );
}
