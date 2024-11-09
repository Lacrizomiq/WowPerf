// Basic encounter type
interface Encounter {
  id: number;
  name: string;
}

// All Stars ranking information
interface AllStars {
  partition: number;
  spec: string;
  points: number;
  possiblePoints: number;
  rank: number;
  regionRank: number;
  serverRank: number;
  rankPercent: number;
  total: number;
}

// Individual encounter ranking
interface EncounterRanking {
  encounter: Encounter;
  rankPercent: number;
  medianPercent: number;
  lockedIn: boolean;
  totalKills: number;
  fastestKill: number;
  allStars: AllStars;
  spec: string;
  bestSpec: string;
  bestAmount: number;
}

// Zone rankings structure
interface ZoneRankings {
  bestPerformanceAverage: number;
  medianPerformanceAverage: number;
  difficulty: number;
  metric: string;
  partition: number;
  zone: number;
  allStars: AllStars[];
  rankings: EncounterRanking[];
}

// Main response structure
interface RaidRankingsResponse {
  name: string;
  classID: number;
  id: number;
  zoneRankings: ZoneRankings;
}

// Export all interfaces
export type {
  Encounter,
  AllStars,
  EncounterRanking,
  ZoneRankings,
  RaidRankingsResponse,
};
