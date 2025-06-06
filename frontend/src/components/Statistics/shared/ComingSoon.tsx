// components/Statistics/shared/ComingSoon.tsx
import React from "react";
import { Badge } from "@/components/ui/badge";
import { Clock, Sparkles } from "lucide-react";

interface ComingSoonProps {
  title: string;
  description: string;
  features?: string[];
}

const ComingSoon: React.FC<ComingSoonProps> = ({
  title,
  description,
  features,
}) => {
  return (
    <div className="flex flex-col items-center justify-center py-20 text-center">
      <div className="mb-6">
        <Sparkles className="h-16 w-16 text-purple-500 mx-auto mb-4" />
        <Badge className="mb-4 bg-purple-600 text-white px-4 py-2">
          <Clock className="w-4 h-4 mr-2" />
          Coming Soon
        </Badge>
      </div>

      <h3 className="text-3xl font-bold mb-4 text-white">{title}</h3>

      <p className="text-muted-foreground max-w-md mb-6 text-lg">
        {description}
      </p>

      {features && features.length > 0 && (
        <div className="bg-slate-800/30 rounded-lg p-6 border border-slate-700 max-w-lg">
          <h4 className="text-lg font-semibold mb-4 text-white">
            Fonctionnalités prévues :
          </h4>
          <ul className="space-y-2 text-slate-300">
            {features.map((feature, index) => (
              <li key={index} className="flex items-center gap-2">
                <div className="w-1.5 h-1.5 bg-purple-500 rounded-full" />
                {feature}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default ComingSoon;
