"use client";

import { useEffect, useRef } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import api from "@/libs/api";

export function BattleNetCallbackHandler() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const callbackProcessed = useRef(false); // Pour suivre si on a déjà traité le callback

  useEffect(() => {
    const handleCallback = async () => {
      // Si déjà traité, on sort
      if (callbackProcessed.current) {
        return;
      }

      callbackProcessed.current = true; // Marquer comme traité immédiatement

      try {
        const code = searchParams.get("code");
        const state = searchParams.get("state");

        console.log("Processing OAuth callback, first attempt:", {
          code,
          state,
        });

        if (!code || !state) {
          console.error("Missing OAuth parameters");
          router.push("/profile?error=missing_params");
          return;
        }

        // Faire la requête GET directement au callback
        const response = await api.get(`/auth/battle-net/callback`, {
          params: { code, state },
        });

        console.log("Backend response:", response.data);

        if (response.data.linked) {
          console.log("Successfully linked account, redirecting...");
          router.push("/profile?success=true");
        } else {
          throw new Error("Failed to link account");
        }
      } catch (error) {
        console.error("Detailed callback error:", error);
        router.push("/profile?error=unknown");
      }
    };

    if (searchParams.get("code")) {
      handleCallback();
    }
  }, [router, searchParams]);

  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-dark">
      <div className="text-center">
        <h2 className="text-xl font-bold mb-4">
          Processing Battle.net authentication...
        </h2>
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
      </div>
    </div>
  );
}
