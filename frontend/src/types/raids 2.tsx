export interface StaticRaid {
  ID: number;
  Slug: string;
  Name: string;
  ShortName: string;
  Expansion: string;
  MediaURL: string;
  Icon: string;
  Starts: Record<string, string>;
  Ends: Record<string, string>;
  Encounters: Array<{ id: number; slug: string; name: string }>;
}

export interface RaidProgressionData {
  expansions: Expansion[];
}

export interface Expansion {
  id: number;
  name: string;
  raids: Raid[];
}

export interface Raid {
  id: number;
  name: string;
  modes: RaidMode[];
}

export interface RaidMode {
  difficulty: string;
  progress: {
    completed_count: number;
    total_count: number;
    encounters: Array<{
      id: number;
      name: string;
      completed_count: number;
      last_kill_timestamp: number;
    }>;
  };
  status: "COMPLETE" | "IN_PROGRESS";
}

export interface RaidProgress {
  completed_count: number;
  total_count: number;
  encounters: Encounter[];
}

export interface Encounter {
  id: number;
  name: string;
  completed_count: number;
  last_kill_timestamp: number;
}

export type CombinedRaid = StaticRaid & Partial<Raid>;
