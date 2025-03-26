// contexts/StatsContext.tsx
"use client";

import React, { createContext, useContext, useState, ReactNode } from "react";

interface StatsContextType {
  season: string;
  region: string;
  dungeon: string;
  setRegion: (region: string) => void;
  setDungeon: (dungeon: string) => void;
}

const StatsContext = createContext<StatsContextType | undefined>(undefined);

export function StatsProvider({ children }: { children: ReactNode }) {
  const [season] = useState("season-tww-1");
  const [region, setRegion] = useState("world");
  const [dungeon, setDungeon] = useState("all");

  return (
    <StatsContext.Provider
      value={{
        season,
        region,
        dungeon,
        setRegion,
        setDungeon,
      }}
    >
      {children}
    </StatsContext.Provider>
  );
}

export function useStats() {
  const context = useContext(StatsContext);
  if (context === undefined) {
    throw new Error("useStats must be used within a StatsProvider");
  }
  return context;
}
