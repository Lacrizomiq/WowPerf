import { FC } from "react";
import Image from "next/image";
import { getSpecIconById } from "@/utils/classandspecicons";

interface RaidRankingProps {
  raidName: string;
  rank?: number;
  classId: number;
  spec: string;
  isLoading?: boolean;
  raidProgressionData?: any;
}

const RaidRanking: FC<RaidRankingProps> = ({
  raidName,
  rank,
  classId,
  spec,
  isLoading = false,
  raidProgressionData,
}) => {
  const imageUrl = getSpecIconById(classId, spec);

  const getHighestDifficulty = () => {
    if (!raidProgressionData?.expansions) return null;

    const currentExpansion = raidProgressionData.expansions.find(
      (exp: any) => exp.name === "Current Season"
    );
    if (!currentExpansion) return null;

    const currentRaid = currentExpansion.raids.find(
      (raid: any) => raid.name === "Nerub-ar Palace"
    );
    if (!currentRaid) return null;

    const difficultyOrder = ["Mythic", "Heroic", "Normal", "LFR"];

    for (const difficulty of difficultyOrder) {
      const mode = currentRaid.modes.find(
        (mode: any) =>
          mode.difficulty === difficulty && mode.progress.completed_count > 0
      );
      if (mode) {
        return mode;
      }
    }

    return null;
  };

  const highestMode = getHighestDifficulty();

  return (
    <div className="bg-slate-900/40 backdrop-blur-sm rounded-lg px-3 sm:px-4 py-2 sm:py-2.5 min-w-0 sm:min-w-[150px]">
      <div className="text-sm text-slate-300/90 mb-2 w-full text-center border-b border-slate-700/50 pb-2">
        {raidName} {highestMode && `${highestMode.difficulty}`}
      </div>
      <div className="flex items-center justify-between">
        <div className="flex items-center justify-center">
          {imageUrl && (
            <Image src={imageUrl} alt={spec} width={32} height={32} />
          )}
        </div>
        <div className="flex flex-col md:flex-row items-center gap-2 text-center justify-center md:px-2">
          <span className="text-slate-400 text-sm">Rank</span>
          {isLoading ? (
            <div className="animate-pulse w-6 h-6 bg-slate-700/50 rounded" />
          ) : (
            <span className="text-xl font-semibold text-white/90">{rank}</span>
          )}
          <span className="text-slate-400 text-sm">Progress</span>
          {isLoading ? (
            <div className="animate-pulse w-6 h-6 bg-slate-700/50 rounded" />
          ) : (
            <span className="text-xl font-semibold text-white/90">
              {highestMode
                ? `${highestMode.progress.completed_count}/${highestMode.progress.total_count}`
                : "0/0"}
            </span>
          )}
        </div>
      </div>
    </div>
  );
};

export default RaidRanking;
