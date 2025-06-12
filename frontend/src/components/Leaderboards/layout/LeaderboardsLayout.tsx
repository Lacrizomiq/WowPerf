// components/leaderboards/layout/LeaderboardsLayout.tsx
"use client";

import React from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import MythicPlusBestRuns from "@/components/Leaderboards/mythicplus/MythicBestRuns";
import RaidLeaderboard from "@/components/Leaderboards/raids/RaidLeaderboard";
import PvPContent from "@/components/Leaderboards/pvp/PvPContent";
import { useRouter } from "next/navigation";

interface LeaderboardsLayoutProps {
  activeTab: string;
}

export default function LeaderboardsLayout({
  activeTab,
}: LeaderboardsLayoutProps) {
  const router = useRouter();

  const handleTabChange = (value: string) => {
    let newPath = "";

    if (value === "mythic-plus") {
      newPath = `/leaderboards`;
    } else if (value === "raids") {
      newPath = `/leaderboards/raids`;
    }

    if (newPath) {
      router.push(newPath);
    }
  };

  return (
    <div className="flex flex-col min-h-screen bg-[#1A1D21] text-[#EAEAEA]">
      {/* Page Header */}
      <header className="pt-8 pb-6 px-4 md:px-8 border-b border-slate-800">
        <div className="container mx-auto">
          <h1 className="text-3xl md:text-4xl font-bold mb-2">Leaderboards</h1>
          <p className="text-muted-foreground text-base md:text-lg">
            Explore detailed teams, players and guilds rankings across Mythic +,
            Raids and PvP game content.
          </p>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 container mx-auto px-4 md:px-8 py-6">
        {/* Content Type Filter */}
        <Tabs
          value={activeTab}
          onValueChange={handleTabChange}
          className="w-full"
        >
          <TabsList className="grid w-full grid-cols-3 bg-slate-800/50 mb-6">
            <TabsTrigger
              value="mythic-plus"
              className="data-[state=active]:bg-purple-600 hover:bg-slate-700 transition-colors duration-200"
            >
              Mythic+
            </TabsTrigger>
            <TabsTrigger
              value="raids"
              className="data-[state=active]:bg-purple-600 hover:bg-slate-700 transition-colors duration-200 relative"
            >
              Raids
            </TabsTrigger>
            <TabsTrigger
              value="pvp"
              className="data-[state=active]:bg-purple-600 hover:bg-slate-700 transition-colors duration-200 relative"
              disabled
            >
              PvP
              <Badge className="absolute -top-2 -right-2 bg-purple-600 text-[10px]">
                Soon
              </Badge>
            </TabsTrigger>
          </TabsList>

          {/* Tab Contents */}
          <TabsContent value="mythic-plus" className="space-y-6">
            <MythicPlusBestRuns />
          </TabsContent>

          <TabsContent value="raids" className="space-y-6">
            <RaidLeaderboard />
          </TabsContent>

          <TabsContent value="pvp" className="space-y-6">
            <PvPContent />
          </TabsContent>
        </Tabs>
      </main>
    </div>
  );
}
