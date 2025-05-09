// components/performance/layout/PerformanceLayout.tsx
import React from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import MythicPlusContent from "@/components/MythicPlus/PerformanceStatistics/mythicplus/MythicPlusContent";
import RaidsContent from "@/components/MythicPlus/PerformanceStatistics/raids/RaidsContent";
import PvPContent from "@/components/MythicPlus/PerformanceStatistics/pvp/PvPContent";

export default function PerformanceLayout() {
  return (
    <div className="flex flex-col min-h-screen bg-[#1A1D21] text-[#EAEAEA]">
      {/* Page Header */}
      <header className="pt-8 pb-6 px-4 md:px-8 border-b border-slate-800">
        <div className="container mx-auto">
          <h1 className="text-3xl md:text-4xl font-bold mb-2">
            Performance Analysis
          </h1>
          <p className="text-muted-foreground text-base md:text-lg">
            Explore detailed class and specialization rankings across various
            game content.
          </p>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 container mx-auto px-4 md:px-8 py-6">
        {/* Content Type Filter */}
        <Tabs defaultValue="mythic" className="w-full">
          <TabsList className="grid w-full grid-cols-3 bg-slate-800/50 mb-6">
            <TabsTrigger
              value="mythic"
              className="data-[state=active]:bg-purple-600"
            >
              Mythic+
            </TabsTrigger>
            <TabsTrigger
              value="raids"
              className="data-[state=active]:bg-purple-600 relative"
              disabled
            >
              Raids
              <Badge className="absolute -top-2 -right-2 bg-purple-600 text-[10px]">
                Soon
              </Badge>
            </TabsTrigger>
            <TabsTrigger
              value="pvp"
              className="data-[state=active]:bg-purple-600 relative"
              disabled
            >
              PvP
              <Badge className="absolute -top-2 -right-2 bg-purple-600 text-[10px]">
                Soon
              </Badge>
            </TabsTrigger>
          </TabsList>

          {/* Tab Contents */}
          <TabsContent value="mythic" className="space-y-6">
            <MythicPlusContent />
          </TabsContent>

          <TabsContent value="raids" className="space-y-6">
            <RaidsContent />
          </TabsContent>

          <TabsContent value="pvp" className="space-y-6">
            <PvPContent />
          </TabsContent>
        </Tabs>
      </main>
    </div>
  );
}
