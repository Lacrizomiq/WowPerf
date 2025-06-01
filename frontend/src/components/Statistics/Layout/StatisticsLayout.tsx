// components/Statistics/Layout/StatisticsLayout.tsx
"use client";

import React from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { useRouter } from "next/navigation";

// Import Mythic+ sections
import OverallStatsSection from "../mythicplus/OverallStatsSection";
import SpecByRoleSection from "../mythicplus/SpecByRoleSection";
import KeyLevelDistributionSection from "../mythicplus/KeyLevelDistributionSection";
import TopCompositionsSection from "../mythicplus/TopCompositionsSection";

// Import Coming Soon components
import ComingSoon from "../shared/ComingSoon";

interface StatisticsLayoutProps {
  activeTab: string;
}

export default function StatisticsLayout({ activeTab }: StatisticsLayoutProps) {
  const router = useRouter();

  const handleTabChange = (value: string) => {
    let newPath = "";

    if (value === "mythic") {
      newPath = `/statistics`;
    } else if (value === "raids") {
      newPath = `/statistics/raids`;
    } else if (value === "pvp") {
      newPath = `/statistics/pvp`;
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
          <h1 className="text-3xl md:text-4xl font-bold mb-2">
            Trends & Statistics
          </h1>
          <p className="text-muted-foreground text-base md:text-lg">
            Explore trends, meta, and detailed statistics for World of Warcraft
            content.
          </p>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 container mx-auto px-4 md:px-8 py-6">
        {/* Content Type Tabs */}
        <Tabs
          value={activeTab}
          onValueChange={handleTabChange}
          className="w-full"
        >
          <TabsList className="grid w-full grid-cols-3 bg-slate-800/50 mb-6">
            <TabsTrigger
              value="mythic"
              className="data-[state=active]:bg-purple-600 hover:bg-slate-700 transition-colors duration-200"
            >
              Mythic+
            </TabsTrigger>
            <TabsTrigger
              value="raids"
              disabled
              className="data-[state=active]:bg-purple-600 hover:bg-slate-700 transition-colors duration-200 relative"
            >
              Raids
              <Badge className="absolute -top-2 -right-2 bg-purple-600 text-[10px]">
                Soon
              </Badge>
            </TabsTrigger>
            <TabsTrigger
              value="pvp"
              disabled
              className="data-[state=active]:bg-purple-600 hover:bg-slate-700 transition-colors duration-200 relative"
            >
              PvP
              <Badge className="absolute -top-2 -right-2 bg-purple-600 text-[10px]">
                Soon
              </Badge>
            </TabsTrigger>
          </TabsList>

          {/* Mythic+ Content */}
          <TabsContent value="mythic" className="space-y-8">
            {/* Section 1: Overall Overview */}
            <OverallStatsSection />

            {/* Section 2: Key Level Distribution */}
            <KeyLevelDistributionSection />

            {/* Section 3: Specialization Usage by Role */}
            <SpecByRoleSection />

            {/* Section 4: Top Team Compositions */}
            <TopCompositionsSection />
          </TabsContent>

          {/* Raids Content - Coming Soon */}
          <TabsContent value="raids" className="space-y-6">
            <ComingSoon
              title="Raid Analytics"
              description="Detailed raid statistics will be available soon."
            />
          </TabsContent>

          {/* PvP Content - Coming Soon */}
          <TabsContent value="pvp" className="space-y-6">
            <ComingSoon
              title="PvP Analytics"
              description="Detailed PvP statistics will be available soon."
            />
          </TabsContent>
        </Tabs>
      </main>
    </div>
  );
}
