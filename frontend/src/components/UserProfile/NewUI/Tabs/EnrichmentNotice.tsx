import React, { useState } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface EnrichmentNoticeProps {
  hasUnenrichedCharacters: boolean;
  onRefresh: () => void;
  isRefreshing: boolean;
}

export const EnrichmentNotice: React.FC<EnrichmentNoticeProps> = ({
  hasUnenrichedCharacters,
  onRefresh,
  isRefreshing,
}) => {
  const [isDismissed, setIsDismissed] = useState(false);

  if (!hasUnenrichedCharacters || isDismissed) {
    return null;
  }

  return (
    <Card className="bg-yellow-900/30 border-yellow-500/50 p-4 mb-4">
      <div className="flex items-start gap-3">
        {/* Warning icon */}
        <div className="flex-shrink-0 mt-0.5">
          <svg
            className="w-5 h-5 text-yellow-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        </div>

        <div className="flex-1">
          <h4 className="font-semibold text-yellow-400 mb-1">
            Limited Character Data
          </h4>
          <p className="text-sm text-yellow-200 mb-3">
            Some characters show basic data only. This can happen when:
          </p>
          <ul className="text-xs text-yellow-100 space-y-1 mb-3">
            <li>
              • Character is not fully synced or an error occurred during sync
            </li>
            <li>• Character names contain special characters</li>
            <li>• Characters were recently transferred or renamed</li>
            <li>• Battle.net API is temporarily unavailable</li>
            <li>• Character privacy settings restrict access</li>
          </ul>

          <div className="flex gap-2">
            <Button
              size="sm"
              onClick={onRefresh}
              disabled={isRefreshing}
              className="bg-yellow-600 hover:bg-yellow-700 text-white"
            >
              {isRefreshing ? "Refreshing..." : "Try Refresh"}
            </Button>

            <Button
              size="sm"
              variant="outline"
              onClick={() => setIsDismissed(true)}
              className="border-yellow-500 text-yellow-400 hover:bg-yellow-500/10"
            >
              Dismiss
            </Button>
          </div>
        </div>
      </div>
    </Card>
  );
};
