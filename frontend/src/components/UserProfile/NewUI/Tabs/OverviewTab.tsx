// OverviewTab.tsx
// Overview tab showing user's personal information, top character and account connections
import React from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Image from "next/image";
import { WoWProfile } from "../AccountProfile";
import { UserProfile } from "@/libs/userService";
import FavoriteCharacterSection from "../FavoriteCharacterSection";
interface OverviewTabProps {
  profile: UserProfile;
  linkStatus: { linked: boolean; battleTag?: string } | null;
  isLinkLoading: boolean;
  isUnlinking: boolean;
  onBattleNetLink: () => Promise<void>;
  onBattleNetUnlink: () => Promise<void>;
  onNavigate: (tab: string) => void;
}

const OverviewTab: React.FC<OverviewTabProps> = ({
  profile,
  linkStatus,
  isLinkLoading,
  isUnlinking,
  onBattleNetLink,
  onBattleNetUnlink,
  onNavigate,
}) => {
  return (
    <div className="space-y-6">
      <Card className="bg-[#131e33] border-gray-800 p-6">
        <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="text-blue-500"
          >
            <circle cx="12" cy="8" r="5" />
            <path d="M20 21a8 8 0 0 0-16 0" />
          </svg>
          Personal Information
        </h2>
        <div className="grid gap-4">
          <div>
            <p className="text-gray-400 font-bold">Username</p>
            <p className="font-medium">{profile.username}</p>
          </div>
          <div>
            <p className="text-gray-400 font-bold">Email</p>
            <p className="font-medium">{profile.email}</p>
          </div>
          {profile.battle_tag && (
            <div>
              <p className="text-gray-400 font-bold">Battle Tag</p>
              <p className="font-medium text-blue-400">{profile.battle_tag}</p>
            </div>
          )}
          <div>
            <p className="text-gray-400 font-bold">Member Since</p>
            <p className="font-medium">February 28, 2025</p>{" "}
            {/* Todo : Implement created_at from the api to render this dyanmicly */}
          </div>
        </div>
      </Card>

      <Card className="bg-[#131e33] border-gray-800 p-6">
        <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="text-blue-500"
          >
            <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
            <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
          </svg>
          Connected Accounts
        </h2>

        <div className="flex justify-between items-center">
          <div className="flex items-center gap-4">
            <Image
              src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
              alt="Battle.net"
              width={40}
              height={40}
            />
            <div>
              <h3 className="font-semibold">Battle.net</h3>
              {linkStatus?.linked ? (
                <p>Connected as: {profile.battle_tag}</p>
              ) : (
                <p className="text-gray-400">Not connected</p>
              )}
            </div>
          </div>

          {linkStatus?.linked ? (
            <Button
              variant="destructive"
              onClick={onBattleNetUnlink}
              disabled={isUnlinking}
            >
              {isUnlinking ? "Unlinking..." : "Disconnect"}
            </Button>
          ) : (
            <Button
              onClick={onBattleNetLink}
              disabled={isLinkLoading}
              className="flex items-center gap-2"
            >
              <Image
                src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
                alt="Battle.net"
                width={20}
                height={20}
              />
              {isLinkLoading ? "Connecting..." : "Connect"}
            </Button>
          )}
        </div>
      </Card>
    </div>
  );
};

export default OverviewTab;
