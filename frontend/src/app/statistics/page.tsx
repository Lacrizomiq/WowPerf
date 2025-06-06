// app/statistics/page.tsx
import React from "react";
import StatisticsLayout from "@/components/Statistics/Layout/StatisticsLayout";

/**
 * Page principale des statistiques - Mythic+ par défaut
 * Route: /statistics
 */
export default function StatisticsPage() {
  return <StatisticsLayout activeTab="mythic" />;
}
