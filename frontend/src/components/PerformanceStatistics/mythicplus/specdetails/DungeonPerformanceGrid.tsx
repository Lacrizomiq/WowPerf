import React, { useMemo } from "react";
import {
  useGetMaxKeyLevelPerSpecAndDungeon,
  useGetDungeonMedia,
  useGetSpecDungeonScoreAverages,
} from "@/hooks/useWarcraftLogsApi";
import DungeonCard from "./DungeonCard";

interface DungeonPerformanceGridProps {
  className: string;
  specName: string;
}

const DungeonPerformanceGrid: React.FC<DungeonPerformanceGridProps> = ({
  className,
  specName,
}) => {
  const { data: dungeonData, isLoading: isLoadingDungeons } =
    useGetMaxKeyLevelPerSpecAndDungeon();

  const { data: dungeonMedia, isLoading: isLoadingMedia } =
    useGetDungeonMedia();

  const { data: scoreAverages, isLoading: isLoadingScores } =
    useGetSpecDungeonScoreAverages(className, specName);

  // Map the media by slug
  const mediaBySlug = useMemo(() => {
    if (!dungeonMedia) return {};
    return dungeonMedia.reduce((acc, media) => {
      acc[media.dungeon_slug] = media;
      return acc;
    }, {} as Record<string, (typeof dungeonMedia)[0]>);
  }, [dungeonMedia]);

  // Map scores by encounter_id
  const scoresByEncounterId = useMemo(() => {
    if (!scoreAverages) return {};
    return scoreAverages.reduce((acc, score) => {
      acc[score.encounter_id] = score;
      return acc;
    }, {} as Record<number, (typeof scoreAverages)[0]>);
  }, [scoreAverages]);

  // Filter class and spec data
  const specDungeonData = useMemo(() => {
    if (!dungeonData) return [];
    return dungeonData.filter(
      (dungeon) =>
        dungeon.class.toLowerCase() === className.toLowerCase() &&
        dungeon.spec.toLowerCase() === specName.toLowerCase()
    );
  }, [dungeonData, className, specName]);

  const isLoading = isLoadingDungeons || isLoadingMedia || isLoadingScores;

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[200px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-600 mx-auto"></div>
          <p className="mt-4 text-slate-300">Loading dungeon performance...</p>
        </div>
      </div>
    );
  }

  if (!specDungeonData.length) {
    return (
      <div className="flex items-center justify-center min-h-[200px]">
        <p className="text-slate-300 text-center">
          No dungeon performance data available.
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {specDungeonData.map((dungeon, index) => (
        <DungeonCard
          key={index}
          name={dungeon.dungeon_name}
          keyLevel={dungeon.max_key_level}
          score={Math.round(
            scoresByEncounterId[dungeon.encounter_id]?.avg_dungeon_score || 0
          )}
          maxScore={Math.round(
            scoresByEncounterId[dungeon.encounter_id]?.max_score || 0
          )}
          backgroundUrl={mediaBySlug[dungeon.dungeon_slug]?.media_url}
        />
      ))}
    </div>
  );
};

export default DungeonPerformanceGrid;
