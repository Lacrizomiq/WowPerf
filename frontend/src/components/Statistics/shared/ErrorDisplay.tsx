// components/Statistics/shared/ErrorDisplay.tsx
import React from "react";
import { AlertTriangle, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";

interface ErrorDisplayProps {
  error?: Error | null;
  message?: string;
  onRetry?: () => void;
  showRetry?: boolean;
}

const ErrorDisplay: React.FC<ErrorDisplayProps> = ({
  error,
  message = "Une erreur est survenue",
  onRetry,
  showRetry = false,
}) => {
  return (
    <div className="bg-red-900/20 border border-red-500/50 rounded-lg p-6 my-4">
      <div className="flex items-center gap-3 mb-3">
        <AlertTriangle className="h-5 w-5 text-red-500" />
        <h3 className="text-red-400 text-lg font-medium">{message}</h3>
      </div>

      {error?.message && (
        <p className="text-red-300 text-sm mb-4">{error.message}</p>
      )}

      {showRetry && onRetry && (
        <Button
          onClick={onRetry}
          variant="outline"
          size="sm"
          className="border-red-500/50 text-red-400 hover:bg-red-500/10"
        >
          <RefreshCw className="w-4 h-4 mr-2" />
          Réessayer
        </Button>
      )}
    </div>
  );
};

export default ErrorDisplay;
