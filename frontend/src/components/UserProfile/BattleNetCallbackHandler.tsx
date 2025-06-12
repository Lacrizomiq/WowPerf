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

  // State for the flow
  const [linkSuccess, setLinkSuccess] = useState(false);
  const [autoSyncTriggered, setAutoSyncTriggered] = useState(false);
  const [showOnboardingModal, setShowOnboardingModal] = useState(false);

  // Hook pour la sync
  const { characters, actions, isLoading } = useCharacters();

  // 🔥 Détection si on vient d'un auto-relink (URL param OU réponse backend)
  const [isAutoRelink, setIsAutoRelink] = useState(
    searchParams.get("auto_relink") === "true"
  );

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
          params: { code, state },
          withCredentials: true,
        });

        if (response.data.linked) {
          // Détecter auto_relink depuis la réponse backend aussi
          const autoRelinkFromResponse = response.data.auto_relink === true;
          if (autoRelinkFromResponse) {
            setIsAutoRelink(true);
          }

          // Invalider les queries
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

          // 🔥 LOGIQUE SIMPLIFIÉE
          if (isAutoRelink || autoRelinkFromResponse) {
            // Auto-relink depuis un click sync → lancer la sync automatiquement
            toast.success("Now syncing your characters...", { duration: 2000 });

            // Attendre un peu que les queries se mettent à jour
            setTimeout(() => {
              actions.syncAndEnrich();
              setAutoSyncTriggered(true);

              // Rediriger après la sync
              setTimeout(() => {
                router.push("/profile?tab=characters&success=auto_sync");
              }, 2000);
            }, 1000);
          } else {
            // Link normal → montrer le modal d'onboarding pour les nouveaux utilisateurs
            setTimeout(() => {
              setShowOnboardingModal(true);
            }, 1500);
          }
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
          toast.error("An unexpected error occurred");
          router.push("/profile?error=unknown");
        }
      }
    };

    if (searchParams.get("code")) {
      handleCallback();
    }
  }, [router, searchParams, queryClient, isAutoRelink, actions]);

  // Affichage conditionnel selon le type de flow
  const shouldShowModal = showOnboardingModal && linkSuccess && !isAutoRelink;

  return (
    <>
      {/* Loading screen adaptatif */}
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

              {isAutoRelink && autoSyncTriggered ? (
                /* Auto-sync en cours */
                <>
                  <div className="mb-4">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-500 mx-auto"></div>
                  </div>
                  <p className="text-sm text-gray-400">
                    {isLoading.sync
                      ? "Syncing your characters..."
                      : "Redirecting to your profile..."}
                  </p>
                </>
              ) : (
                /* Préparation onboarding */
                <p className="text-sm text-gray-400">
                  {isAutoRelink
                    ? "Starting character sync..."
                    : "Preparing your character sync..."}
                </p>
              )}
            </>
          ) : (
            /* État de linking */
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

      {/* Onboarding Modal - seulement pour les nouveaux utilisateurs */}
      <OnboardingModal isOpen={shouldShowModal} />
    </>
  );
}
