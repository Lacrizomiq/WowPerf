import React, { useState } from "react";
import { seasons } from "@/data/seasons";
import SeasonsSelector from "@/components/MythicPlus/SeasonsSelector";
import StaticDungeonList from "./StaticDungeonList";
import DungeonDetails from "./DungeonDetails";
import {
  useGetBlizzardMythicDungeonPerSeason,
  useGetBlizzardCharacterMythicPlusBestRuns,
} from "@/hooks/useBlizzardApi";
import {
  MythicDungeonProps,
  Dungeon,
  MythicPlusRuns,
} from "@/types/mythicPlusRuns";

const MythicDungeonOverview: React.FC<MythicDungeonProps> = ({
  characterName,
  realmSlug,
  region,
  namespace,
  locale,
  seasonSlug,
}) => {
  const [selectedSeason, setSelectedSeason] = useState<{
    slug: string;
    id: number;
  }>(seasons.find((s) => s.slug === seasonSlug) || seasons[0]);

  const [selectedDungeon, setSelectedDungeon] = useState<Dungeon | null>(null);

  const {
    data: dungeonData,
    isLoading: isDungeonLoading,
    error: dungeonError,
  } = useGetBlizzardMythicDungeonPerSeason(selectedSeason.slug);

  const {
    data: mythicPlusRuns,
    isLoading: isRunsLoading,
    error: runsError,
  } = useGetBlizzardCharacterMythicPlusBestRuns(
    region,
    realmSlug,
    characterName,
    namespace,
    locale,
    selectedSeason.id.toString()
  );

  console.log("Mythic+ Runs Data:", mythicPlusRuns);

  const handleSeasonChange = (seasonSlug: string) => {
    const newSeason = seasons.find((s) => s.slug === seasonSlug);
    if (newSeason) {
      setSelectedSeason(newSeason);
      setSelectedDungeon(null);
    }
  };

  const handleDungeonClick = (dungeon: Dungeon) => {
    setSelectedDungeon(dungeon);
  };

  if (isDungeonLoading || isRunsLoading) return <div>Loading data...</div>;
  if (dungeonError || runsError) return <div>Error loading data</div>;
  if (!dungeonData) return <div>No dungeon data found</div>;

  const selectedRun = mythicPlusRuns?.find(
    (run) => run.Dungeon.ID === selectedDungeon?.ID
  );

  return (
    <div className="p-4 bg-gradient-dark shadow-lg rounded-lg glow-effect m-12 max-w-6xl mx-auto">
      <SeasonsSelector seasons={seasons} onSeasonChange={handleSeasonChange} />
      <StaticDungeonList
        dungeons={dungeonData.dungeons}
        mythicPlusRuns={mythicPlusRuns || []}
        onDungeonClick={handleDungeonClick}
      />
      {selectedRun && <DungeonDetails run={selectedRun} />}
    </div>
  );
};

export default MythicDungeonOverview;
