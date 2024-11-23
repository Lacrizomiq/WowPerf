"use client";

import { useEffect, useRef } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import axios from "axios";
import api from "@/libs/api";

export function BattleNetCallbackHandler() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const callbackProcessed = useRef(false);
  const queryClient = useQueryClient();

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
          ]);
          toast.success(`Successfully linked to ${response.data.battleTag}`, {
            id: toastId,
          });

          // Wait a bit before redirecting to let the user see the message
          setTimeout(() => {
            router.push("/profile?success=link_successful");
          }, 1000);
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

  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-dark">
      <div className="text-center max-w-md mx-auto px-4">
        <h2 className="text-xl font-bold mb-4">
          Linking your Battle.net account...
        </h2>
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mb-4 mx-auto"></div>
        <p className="text-sm text-gray-400 mb-2">
          This may take a few seconds, please don&apos;t close this window.
        </p>
        <p className="text-xs text-gray-500">
          We&apos;re connecting to Battle.net to verify your account...
        </p>
      </div>
    </div>
  );
}
