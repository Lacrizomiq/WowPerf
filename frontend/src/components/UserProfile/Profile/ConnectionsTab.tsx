// ConnectionsTab.tsx - Tab connections adapté au style V0
import React from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { Check, LinkIcon, Unlink } from "lucide-react";
import Image from "next/image";
import { UserProfile } from "@/libs/userService";

interface ConnectionsTabProps {
  profile: UserProfile;
  linkStatus: { linked: boolean; battleTag?: string } | null;
  isLinkLoading: boolean;
  isUnlinking: boolean;
  onBattleNetLink: () => Promise<void>;
  onBattleNetUnlink: () => Promise<void>;
  isActive: boolean;
}

const ConnectionsTab: React.FC<ConnectionsTabProps> = ({
  profile,
  linkStatus,
  isLinkLoading,
  isUnlinking,
  onBattleNetLink,
  onBattleNetUnlink,
  isActive,
}) => {
  if (!isActive) return null;
  return (
    <div className="space-y-6">
      {/* Battle.net Connection */}
      <Card className="bg-slate-800/30 border-slate-700">
        <CardHeader>
          <CardTitle>Battle.net Connection</CardTitle>
          <CardDescription>
            Link your Battle.net account to synchronize your characters and
            access exclusive features on WoW Perf.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center p-4 rounded-md bg-slate-800/50 border border-slate-700">
            <div className="mr-4">
              <div className="w-12 h-12 rounded-md bg-blue-900/30 flex items-center justify-center">
                <Image
                  src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
                  alt="Battle.net"
                  width={24}
                  height={24}
                />
              </div>
            </div>
            <div className="flex-1">
              <h3 className="font-medium">Battle.net</h3>
              <p className="text-sm text-slate-400">
                {linkStatus?.linked
                  ? `Connected as ${profile.battle_tag}`
                  : "Not connected"}
              </p>
            </div>
            <Button
              onClick={linkStatus?.linked ? onBattleNetUnlink : onBattleNetLink}
              variant="outline"
              disabled={isLinkLoading || isUnlinking}
              className={
                linkStatus?.linked
                  ? "border-red-700 text-red-400 hover:bg-red-900/30"
                  : "border-blue-700 text-blue-400 hover:bg-blue-900/30"
              }
            >
              {linkStatus?.linked ? (
                <>
                  <Unlink className="mr-2 h-4 w-4" />
                  {isUnlinking ? "Unlinking..." : "Unlink Account"}
                </>
              ) : (
                <>
                  <LinkIcon className="mr-2 h-4 w-4" />
                  {isLinkLoading ? "Connecting..." : "Link Account"}
                </>
              )}
            </Button>
          </div>

          <Separator className="my-4 bg-slate-700" />

          <div className="space-y-4">
            <h3 className="font-medium">Benefits of Battle.net connection</h3>
            <ul className="space-y-2 text-sm text-slate-400">
              <li className="flex items-start">
                <Check className="h-4 w-4 mr-2 text-green-400 mt-0.5" />
                Synchronization of all your characters on demand
              </li>
              <li className="flex items-start">
                <Check className="h-4 w-4 mr-2 text-green-400 mt-0.5" />
                Access to detailed Mythic+ run statistics in Dashboard page
              </li>
              <li className="flex items-start">
                <Check className="h-4 w-4 mr-2 text-green-400 mt-0.5" />
                Tracking of your progression in Dashboard page
              </li>
              <li className="flex items-start">
                <Check className="h-4 w-4 mr-2 text-green-400 mt-0.5" />
                Personalized recommendations based on your playstyle in
                Dashboard page
              </li>
            </ul>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default ConnectionsTab;
