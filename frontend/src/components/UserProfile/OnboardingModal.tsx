"use client";

import React, { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useCharacters } from "@/hooks/useCharacters";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import Image from "next/image";

interface OnboardingModalProps {
  isOpen: boolean;
}

export const OnboardingModal: React.FC<OnboardingModalProps> = ({ isOpen }) => {
  const router = useRouter();
  const { actions, isLoading, rateLimitState, ui, characters } =
    useCharacters();
  const [syncCompleted, setSyncCompleted] = useState(false);

  // ✅ Internal navigation handlers (MUST be declared FIRST)
  const handleComplete = useCallback(() => {
    router.push("/profile?success=sync_complete");
  }, [router]);

  const handleSkip = useCallback(() => {
    router.push("/profile?info=sync_skipped");
  }, [router]);

  const handleSync = useCallback(async () => {
    try {
      await actions.syncAndEnrich();
      setSyncCompleted(true);

      // Auto-close after 2 seconds to show success
      setTimeout(() => {
        handleComplete();
      }, 2000);
    } catch (error) {
      // Error handling is done in the hook with toast
      console.error("Sync failed:", error);
    }
  }, [actions, handleComplete]);

  // ✅ Debug logs (AFTER handlers are declared)
  console.log("OnboardingModal:", {
    isOpen,
    characters,
    charactersLength: characters?.length,
    hasCharacters: characters && characters.length > 0,
  });

  // ✅ If user already has characters, auto-skip modal (AFTER handlers are declared)
  React.useEffect(() => {
    if (isOpen && characters && characters.length > 0) {
      console.log("User already has characters, auto-completing...");
      handleComplete();
    }
  }, [isOpen, characters, handleComplete]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4">
      <Card className="bg-[#131e33] border-gray-800 max-w-md w-full">
        <CardHeader>
          <CardTitle className="flex items-center gap-3">
            <Image
              src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
              alt="Battle.net"
              width={32}
              height={32}
            />
            <span>Sync Your WoW Characters</span>
          </CardTitle>
        </CardHeader>

        <CardContent className="space-y-6">
          {/* Success state */}
          {syncCompleted ? (
            <div className="text-center py-4">
              <div className="mb-4">
                <svg
                  className="w-16 h-16 text-green-500 mx-auto"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </div>
              <h3 className="text-xl font-semibold text-green-500 mb-2">
                Sync Complete!
              </h3>
              <p className="text-gray-400">
                Your characters have been synchronized and enriched. Redirecting
                to your profile...
              </p>
            </div>
          ) : (
            <>
              {/* Loading state */}
              {isLoading.sync ? (
                <div className="text-center py-4">
                  <div className="mb-4">
                    <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-500 mx-auto"></div>
                  </div>
                  <h3 className="text-xl font-semibold mb-2">
                    Sync in Progress
                  </h3>
                  <p className="text-gray-400">
                    We&apos;re fetching and enriching your character data from
                    Battle.net. This may take a few moments...
                  </p>
                  <div className="mt-4 space-y-2 text-sm text-gray-500">
                    <div className="flex items-center justify-center gap-2">
                      <div className="animate-pulse h-2 w-2 bg-blue-500 rounded-full"></div>
                      <span>Connecting to Battle.net</span>
                    </div>
                    <div className="flex items-center justify-center gap-2">
                      <div className="animate-pulse h-2 w-2 bg-blue-500 rounded-full"></div>
                      <span>Fetching character data</span>
                    </div>
                    <div className="flex items-center justify-center gap-2">
                      <div className="animate-pulse h-2 w-2 bg-blue-500 rounded-full"></div>
                      <span>Enriching character information</span>
                    </div>
                  </div>
                </div>
              ) : (
                <>
                  {/* Rate limit state */}
                  {ui.showRateLimit ? (
                    <div className="text-center py-4">
                      <div className="mb-4">
                        <svg
                          className="w-16 h-16 text-orange-500 mx-auto"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                      </div>
                      <h3 className="text-xl font-semibold text-orange-500 mb-2">
                        Please Wait
                      </h3>
                      <p className="text-gray-400 mb-4">
                        {rateLimitState.message}
                      </p>
                      <div className="text-2xl font-mono text-blue-400">
                        {rateLimitState.formattedTime}
                      </div>
                    </div>
                  ) : (
                    /* Normal state */
                    <div className="text-center">
                      <div className="mb-4">
                        <svg
                          className="w-16 h-16 text-blue-500 mx-auto"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                          />
                        </svg>
                      </div>
                      <h3 className="text-xl font-semibold mb-2">
                        Welcome to WoW Perf!
                      </h3>
                      <p className="text-gray-400 mb-6">
                        We&apos;ll fetch and enrich your WoW character data to
                        give you the best experience. This includes specs,
                        achievements, gear information, and more!
                      </p>

                      <div className="bg-blue-900/30 border border-blue-500/50 rounded-lg p-4 mb-6">
                        <h4 className="font-semibold text-blue-400 mb-2">
                          What we&apos;ll sync:
                        </h4>
                        <ul className="text-sm text-gray-300 space-y-1">
                          <li>• Level 80+ characters only</li>
                          <li>• Character specs and roles</li>
                          <li>• Achievement points</li>
                          <li>• Character avatars</li>
                          <li>• Basic equipment information</li>
                        </ul>
                      </div>
                    </div>
                  )}
                </>
              )}
            </>
          )}

          {/* Action buttons */}
          {!syncCompleted && !isLoading.sync && (
            <div className="flex gap-3">
              <Button
                onClick={handleSync}
                disabled={ui.isDisabled.sync}
                className="flex-1 bg-blue-500 hover:bg-blue-600"
              >
                {ui.showRateLimit ? "Wait..." : "Sync Characters"}
              </Button>

              <Button
                onClick={handleSkip}
                variant="outline"
                className="flex-1"
                disabled={isLoading.sync}
              >
                Skip for Now
              </Button>
            </div>
          )}

          {!syncCompleted && !isLoading.sync && (
            <p className="text-xs text-gray-500 text-center">
              You can always sync your characters later from your profile.
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  );
};
