import { FC } from "react";
import Image from "next/image";
import { getSpecIconById } from "@/utils/classandspecicons";

interface RaidRankingProps {
  raidName: string;
  rank?: number;
  classId: number;
  spec: string;
  isLoading?: boolean;
}

const RaidRanking: FC<RaidRankingProps> = ({
  raidName,
  rank,
  classId,
  spec,
  isLoading = false,
}) => {
  const imageUrl = getSpecIconById(classId, spec);

  return (
    <div className="bg-slate-900/40 backdrop-blur-sm rounded-lg px-3 sm:px-4 py-2 sm:py-2.5 min-w-0 sm:min-w-[180px]">
      <div className="text-sm text-slate-300/90 mb-2 w-full text-center border-b border-slate-700/50 pb-2">
        {raidName}
      </div>
      <div className="flex items-center justify-between">
        <div className="flex items-center justify-center">
          {imageUrl && (
            <Image src={imageUrl} alt={spec} width={32} height={32} />
          )}
        </div>
        <div className="flex flex-col items-center gap-2 text-center">
          <span className="text-slate-400 text-sm">World Rank</span>
          {isLoading ? (
            <div className="animate-pulse w-6 h-6 bg-slate-700/50 rounded" />
          ) : (
            <span className="text-xl font-semibold text-white/90">{rank}</span>
          )}
        </div>
      </div>
    </div>
  );
};

export default RaidRanking;
