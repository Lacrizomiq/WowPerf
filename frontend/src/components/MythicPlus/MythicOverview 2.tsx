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
  MythicPlusSeasonInfo,
  Season,
} from "@/types/mythicPlusRuns";

const MythicDungeonOverview: React.FC<MythicDungeonProps> = ({
  characterName,
  realmSlug,
  region,
  namespace,
  locale,
  seasonSlug,
}) => {
  const [selectedSeason, setSelectedSeason] = useState<Season>(
    seasons.find((s) => s.slug === seasonSlug) || seasons[0]
  );

  const [selectedDungeon, setSelectedDungeon] = useState<Dungeon | null>(null);

  const {
    data: dungeonData,
    isLoading: isDungeonLoading,
    error: dungeonError,
  } = useGetBlizzardMythicDungeonPerSeason(selectedSeason.slug);

  const {
    data: mythicPlusSeasonInfo,
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

  const selectedRun = mythicPlusSeasonInfo?.BestRuns.find(
    (run) => run.Dungeon.ID === selectedDungeon?.ID
  );

  return (
    <div className="p-4 bg-gradient-dark shadow-lg rounded-lg glow-effect m-12 max-w-6xl mx-auto">
      <SeasonsSelector
        seasons={seasons}
        onSeasonChange={handleSeasonChange}
        selectedSeason={selectedSeason}
      />
      {!mythicPlusSeasonInfo ? (
        <div className="mt-4">
          <h2 className="text-2xl font-bold mb-4">{selectedSeason.name}</h2>
          <p>This character has no Mythic+ data for this season.</p>
        </div>
      ) : (
        <div className="mb-4 pt-2">
          <p className="text-xl">
            Season Mythic Rating :{" "}
            <span
              style={{ color: mythicPlusSeasonInfo.OverallMythicRatingHex }}
            >
              {mythicPlusSeasonInfo.OverallMythicRating.toFixed(2)}
            </span>
          </p>
        </div>
      )}
      <StaticDungeonList
        dungeons={dungeonData.dungeons}
        mythicPlusRuns={mythicPlusSeasonInfo?.BestRuns || []}
        onDungeonClick={handleDungeonClick}
      />
      {selectedRun && <DungeonDetails run={selectedRun} region={region} />}
    </div>
  );
};

export default MythicDungeonOverview;
