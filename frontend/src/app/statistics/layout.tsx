// app/statistics/layout.tsx (optionnel - layout spécifique aux statistics)
import React from "react";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Trends & Statistics - WoWPerf",
  description:
    "Explore trends, meta, and detailed statistics for World of Warcraft content",
  keywords: [
    "World of Warcraft",
    "Mythic+",
    "Raids",
    "PvP",
    "Statistics",
    "Meta",
  ],
};

interface StatisticsLayoutProps {
  children: React.ReactNode;
}

/**
 * Layout partagé pour toutes les pages /statistics/*
 * Ajoute les métadonnées SEO et peut inclure des composants communs
 */
export default function StatisticsRootLayout({
  children,
}: StatisticsLayoutProps) {
  return <div className="min-h-screen bg-[#1A1D21]">{children}</div>;
}
