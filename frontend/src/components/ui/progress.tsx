// src/components/ui/progress.tsx
"use client";

import * as React from "react";
import { cn } from "@/lib/utils";

interface ProgressProps extends React.HTMLAttributes<HTMLDivElement> {
  value?: number;
  max?: number;
  className?: string;
}

const Progress = React.forwardRef<HTMLDivElement, ProgressProps>(
  ({ className, value = 0, max = 100, ...props }, ref) => {
    // Normaliser la valeur entre 0 et 100
    const normalizedValue = Math.min(Math.max(value, 0), max);
    const percentage = (normalizedValue / max) * 100;

    return (
      <div
        ref={ref}
        className={cn(
          "relative h-2 w-full overflow-hidden rounded-full bg-slate-700/50",
          className
        )}
        role="progressbar"
        aria-valuemin={0}
        aria-valuemax={max}
        aria-valuenow={normalizedValue}
        {...props}
      >
        <div
          className="h-full bg-gradient-to-r from-purple-500 to-purple-600 transition-all duration-500 ease-out"
          style={{
            width: `${percentage}%`,
            transform: `translateX(0%)`, // Assure que la barre commence à gauche
          }}
        />
      </div>
    );
  }
);

Progress.displayName = "Progress";

export { Progress };
export type { ProgressProps };
