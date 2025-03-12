import React, { useState } from "react";
import { seasons } from "@/data/seasons";
import SeasonsSelector from "@/components/MythicPlus/SeasonsSelector";
import StaticDungeonList from "./StaticDungeonList";
import DungeonDetails from "./DungeonDetails";
import {
  useGetBlizzardMythicDungeonPerSeason,
  useGetBlizzardCharacterMythicPlusBestRuns,
} from "@/hooks/useBlizzardApi";
import { MythicDungeonProps, Dungeon, Season } from "@/types/mythicPlusRuns";
import { useGetPlayerMythicPlusRankings } from "@/hooks/useWarcraftLogsApi";
import MythicPlusPlayerPerformance from "./CharacterPersonalRanking/MythicPlusPlayerPerformance";

const MythicDungeonOverview: React.FC<MythicDungeonProps> = ({
  characterName,
  realmSlug,
  region,
  namespace,
  locale,
  seasonSlug,
}) => {
  // Get the selected season
  const [selectedSeason, setSelectedSeason] = useState<Season>(
    seasons.find((s) => s.slug === seasonSlug) || seasons[0]
  );

  // Get the selected dungeon
  const [selectedDungeon, setSelectedDungeon] = useState<Dungeon | null>(null);

  // Get the dungeon data via a custom Hook
  const {
    data: dungeonData,
    isLoading: isDungeonLoading,
    error: dungeonError,
  } = useGetBlizzardMythicDungeonPerSeason(selectedSeason.slug);

  // Get the season Mythic Rating via a custom Hook
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

  // Get the player's personal performance for the selected season via a custom Hook
  const {
    data: mythicPlusPlayerRankings,
    isLoading: isLoadingMythicPlusPlayerRankings,
    error: mythicPlusPlayerRankingsError,
  } = useGetPlayerMythicPlusRankings(characterName, realmSlug, region, 43);

  // Handle season change
  const handleSeasonChange = (seasonSlug: string) => {
    const newSeason = seasons.find((s) => s.slug === seasonSlug);
    if (newSeason) {
      setSelectedSeason(newSeason);
      setSelectedDungeon(null);
    }
  };

  // Handle dungeon click
  const handleDungeonClick = (dungeon: Dungeon) => {
    setSelectedDungeon(dungeon);
  };

  // Display loading message
  if (isDungeonLoading || isRunsLoading)
    return <div className="text-white text-center p-4">Loading data...</div>;

  // Display error message
  if (dungeonError || runsError)
    return (
      <div className="text-red-500 text-center p-4">Error loading data</div>
    );

  // Display no dungeon data message
  if (!dungeonData)
    return (
      <div className="text-yellow-500 text-center p-4">
        No dungeon data found
      </div>
    );

  const selectedRun = mythicPlusSeasonInfo?.BestRuns.find(
    (run) => run.Dungeon.ID === selectedDungeon?.ID
  );

  return (
    <div className="p-6 rounded-xl shadow-lg m-4">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold text-white">Mythic+ Dungeons</h2>
      </div>

      {/* Display the player's personal performance for the selected season */}
      {mythicPlusPlayerRankings && (
        <div className="mb-6">
          <h3 className="text-xl text-white mb-2">
            Personal Performance for {selectedSeason.name}
          </h3>
          <MythicPlusPlayerPerformance
            playerData={mythicPlusPlayerRankings}
            dungeonData={dungeonData}
          />
        </div>
      )}

      <div className="flex justify-between items-center mb-6">
        <div className="flex-1">
          {/* Display the season Mythic Rating */}
          {mythicPlusSeasonInfo ? (
            <p className="text-xl text-white">
              Season Mythic Rating:{" "}
              <span
                style={{ color: mythicPlusSeasonInfo.OverallMythicRatingHex }}
                className="font-bold"
              >
                {mythicPlusSeasonInfo.OverallMythicRating.toFixed(0)}
              </span>
            </p>
          ) : (
            <p className="text-xl text-white">No score for this season</p>
          )}
        </div>

        <div>
          {" "}
          {/* Season selector */}
          <SeasonsSelector
            seasons={seasons}
            onSeasonChange={handleSeasonChange}
            selectedSeason={selectedSeason}
          />
        </div>
      </div>

      {/* Display the dungeon list */}
      <StaticDungeonList
        dungeons={dungeonData.dungeons}
        mythicPlusRuns={mythicPlusSeasonInfo?.BestRuns || []}
        onDungeonClick={handleDungeonClick}
        selectedDungeon={selectedDungeon}
      />

      {/* Display the dungeon details with the team composition */}
      {selectedRun && <DungeonDetails run={selectedRun} region={region} />}
    </div>
  );
};

export default MythicDungeonOverview;
