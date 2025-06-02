// components/shared/InfoTooltip.tsx

import React, { useState } from "react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { HelpCircle } from "lucide-react";

interface InfoTooltipProps {
  content: string;
  className?: string;
  side?: "top" | "bottom" | "left" | "right";
  size?: "sm" | "md" | "lg";
  delayDuration?: number;
}

/**
 * Composant InfoTooltip réutilisable avec support mobile
 * Affiche une icône "?" avec un tooltip explicatif au survol (desktop) ou au clic (mobile)
 *
 * @param content - Texte à afficher dans le tooltip
 * @param className - Classes CSS additionnelles
 * @param side - Position du tooltip par rapport à l'icône
 * @param size - Taille de l'icône
 * @param delayDuration - Délai avant affichage (ms), défaut: 300ms
 */
const InfoTooltip: React.FC<InfoTooltipProps> = ({
  content,
  className = "",
  side = "right",
  size = "sm",
  delayDuration = 100,
}) => {
  const [open, setOpen] = useState(false);

  const sizeClasses = {
    sm: "h-3 w-3",
    md: "h-4 w-4",
    lg: "h-5 w-5",
  };

  const handleClick = () => {
    // Sur mobile, toggle le tooltip au clic
    setOpen(!open);
  };

  const handleOpenChange = (newOpen: boolean) => {
    setOpen(newOpen);
  };

  return (
    <TooltipProvider delayDuration={delayDuration}>
      <Tooltip open={open} onOpenChange={handleOpenChange}>
        <TooltipTrigger asChild>
          <button
            type="button"
            onClick={handleClick}
            className={`inline-flex items-center justify-center ml-1 text-slate-400 hover:text-slate-200 transition-colors cursor-help touch-manipulation ${className}`}
            aria-label="More information"
            aria-describedby={open ? "tooltip-content" : undefined}
          >
            <HelpCircle className={sizeClasses[size]} />
          </button>
        </TooltipTrigger>
        <TooltipContent
          id="tooltip-content"
          side={side}
          className="max-w-xs bg-slate-900 border-slate-700 text-slate-200 text-sm p-3"
          sideOffset={4}
          onPointerDownOutside={() => setOpen(false)}
        >
          <p className="leading-relaxed">{content}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default InfoTooltip;
