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

  const {
    data: mythicPlusPlayerRankings,
    isLoading: isLoadingMythicPlusPlayerRankings,
    error: mythicPlusPlayerRankingsError,
  } = useGetPlayerMythicPlusRankings(characterName, realmSlug, region, 39);

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

  if (isDungeonLoading || isRunsLoading)
    return <div className="text-white text-center p-4">Loading data...</div>;
  if (dungeonError || runsError)
    return (
      <div className="text-red-500 text-center p-4">Error loading data</div>
    );
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
          {" "}
          {/* Cette div prend l'espace disponible */}
          {mythicPlusSeasonInfo ? (
            <p className="text-xl text-white">
              Season Mythic Rating:{" "}
              <span
                style={{ color: mythicPlusSeasonInfo.OverallMythicRatingHex }}
              >
                {mythicPlusSeasonInfo.OverallMythicRating.toFixed(2)}
              </span>
            </p>
          ) : (
            <p className="text-xl text-white">No score for this season</p>
          )}
        </div>

        <div>
          {" "}
          {/* Cette div ne prend que l'espace nécessaire */}
          <SeasonsSelector
            seasons={seasons}
            onSeasonChange={handleSeasonChange}
            selectedSeason={selectedSeason}
          />
        </div>
      </div>

      <StaticDungeonList
        dungeons={dungeonData.dungeons}
        mythicPlusRuns={mythicPlusSeasonInfo?.BestRuns || []}
        onDungeonClick={handleDungeonClick}
        selectedDungeon={selectedDungeon}
      />

      {selectedRun && <DungeonDetails run={selectedRun} region={region} />}
    </div>
  );
};

export default MythicDungeonOverview;
