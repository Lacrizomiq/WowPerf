import { FC } from "react";
import Image from "next/image";
import { getSpecIconById } from "@/utils/classandspecicons";

interface MythicPlusRankingProps {
  seasonName: string;
  rank?: number;
  classId: number;
  spec: string;
  points: number;
  isLoading?: boolean;
  fallbackImageUrl?: string;
}

const MythicPlusRanking: FC<MythicPlusRankingProps> = ({
  seasonName,
  rank,
  classId,
  spec,
  points,
  isLoading = false,
  fallbackImageUrl,
}) => {
  const imageUrl = getSpecIconById(classId, spec) || fallbackImageUrl;

  return (
    <div className="bg-slate-900/40 backdrop-blur-sm rounded-lg px-3 sm:px-4 py-2 sm:py-2.5 min-w-0 sm:min-w-[150px]">
      <div className="text-sm text-slate-300/90 mb-2 w-full text-center border-b border-slate-700/50 pb-2">
        {seasonName}
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
          <span className="text-slate-400 text-sm">Score</span>
          {isLoading ? (
            <div className="animate-pulse w-6 h-6 bg-slate-700/50 rounded" />
          ) : (
            <span className="text-xl font-semibold text-white/90">
              {points.toFixed(0)}
            </span>
          )}
        </div>
      </div>
    </div>
  );
};

export default MythicPlusRanking;
