// components/Statistics/shared/LoadingSpinner.tsx
import React from "react";

interface LoadingSpinnerProps {
  size?: "sm" | "md" | "lg";
  message?: string;
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = "md",
  message,
}) => {
  const sizeClasses = {
    sm: "h-6 w-6",
    md: "h-12 w-12",
    lg: "h-16 w-16",
  };

  return (
    <div className="flex flex-col items-center justify-center py-12 space-y-4">
      <div
        className={`animate-spin rounded-full border-t-2 border-b-2 border-purple-600 ${sizeClasses[size]}`}
      />
      {message && <p className="text-slate-400 text-sm">{message}</p>}
    </div>
  );
};

export default LoadingSpinner;
