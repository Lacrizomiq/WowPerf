"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import axios from "axios";
import api from "@/libs/api";
import { OnboardingModal } from "@/components/UserProfile/OnboardingModal";
import { useCharacters } from "@/hooks/useCharacters";

export function BattleNetCallbackHandler() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const callbackProcessed = useRef(false);
  const queryClient = useQueryClient();

  // State for onboarding modal
  const [showOnboardingModal, setShowOnboardingModal] = useState(false);
  const [linkSuccess, setLinkSuccess] = useState(false);

  // Check if user has characters to determine if we should show modal
  const { characters, region, isAuthenticated } = useCharacters();

  // Debug logs pour comprendre le problème
  console.log("BattleNetCallback Debug:", {
    showOnboardingModal,
    linkSuccess,
    characters,
    charactersLength: characters?.length,
    isAuthenticated,
    region,
    shouldShow:
      showOnboardingModal &&
      linkSuccess &&
      (!characters || characters.length === 0),
  });

  useEffect(() => {
    const handleCallback = async () => {
      if (callbackProcessed.current) return;

      callbackProcessed.current = true;
      const toastId = toast.loading("Linking your Battle.net account...");

      try {
        const code = searchParams.get("code");
        const state = searchParams.get("state");

        if (!code || !state) {
          toast.error("Missing authentication parameters", { id: toastId });
          router.push("/profile?error=missing_params");
          return;
        }

        const response = await api.get("/auth/battle-net/callback", {
          params: {
            code,
            state,
          },
          withCredentials: true,
        });

        console.log("Battle.net callback response:", response.data);

        if (response.data.linked) {
          // Invalidate both queries
          await Promise.all([
            queryClient.invalidateQueries({
              queryKey: ["battleNetLinkStatus"],
            }),
            queryClient.invalidateQueries({ queryKey: ["userProfile"] }),
            queryClient.invalidateQueries({ queryKey: ["characters"] }),
          ]);

          toast.success(`Successfully linked to ${response.data.battleTag}`, {
            id: toastId,
          });

          setLinkSuccess(true);

          // Wait a bit for queries to invalidate, then check characters
          setTimeout(() => {
            console.log("Checking characters after delay...");
            setShowOnboardingModal(true);
          }, 2000); // Increased delay to 2 seconds
        } else {
          throw new Error(response.data.error || "Failed to link account");
        }
      } catch (error) {
        console.error("Battle.net callback error:", error);

        if (axios.isAxiosError(error)) {
          const errorMessage =
            error.response?.data?.error || "Failed to link Battle.net account";
          const errorCode = error.response?.data?.code || "unknown_error";

          toast.error(errorMessage, { id: toastId });
          router.push(`/profile?error=${errorCode}`);
        } else {
          toast.error("An unexpected error occurred", { id: toastId });
          router.push("/profile?error=unknown");
        }
      }
    };

    if (searchParams.get("code")) {
      handleCallback();
    }
  }, [router, searchParams, queryClient]);

  // Handle onboarding modal completion - no callbacks needed
  // Modal will handle its own navigation internally

  // Show onboarding modal when:
  // 1. Link was successful
  // 2. Modal should be shown (after timeout)
  // Note: We'll let the modal itself handle the character checking
  const shouldShowModal = showOnboardingModal && linkSuccess;

  console.log("shouldShowModal:", shouldShowModal, {
    showOnboardingModal,
    linkSuccess,
  });

  return (
    <>
      {/* Loading screen */}
      <div className="flex items-center justify-center min-h-screen bg-gradient-dark">
        <div className="text-center max-w-md mx-auto px-4">
          {linkSuccess ? (
            <>
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
              <h2 className="text-xl font-bold mb-4 text-green-500">
                Battle.net Account Linked!
              </h2>
              <p className="text-sm text-gray-400">
                Preparing your character sync...
              </p>
            </>
          ) : (
            <>
              <h2 className="text-xl font-bold mb-4">
                Linking your Battle.net account...
              </h2>
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mb-4 mx-auto"></div>
              <p className="text-sm text-gray-400 mb-2">
                This may take a few seconds, please don&apos;t close this
                window.
              </p>
              <p className="text-xs text-gray-500">
                We&apos;re connecting to Battle.net to verify your account...
              </p>
            </>
          )}
        </div>
      </div>

      {/* Onboarding Modal */}
      <OnboardingModal isOpen={shouldShowModal} />
    </>
  );
}
