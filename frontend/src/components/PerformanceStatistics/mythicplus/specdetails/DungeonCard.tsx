import React from "react";

interface DungeonCardProps {
  name: string;
  keyLevel: number;
  score: number;
  maxScore: number;
  backgroundUrl?: string;
}

const DungeonCard: React.FC<DungeonCardProps> = ({
  name,
  keyLevel,
  score,
  maxScore,
  backgroundUrl,
}) => {
  return (
    <div
      className="relative bg-slate-800/30 rounded-lg border border-slate-700 p-4 hover:border-purple-700/50 transition-all overflow-hidden"
      style={{
        backgroundImage: backgroundUrl ? `url('${backgroundUrl}')` : undefined,
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
      }}
    >
      {/* Dark overlay for text readability */}
      <div className="absolute inset-0 bg-black/60" />

      {/* Content */}
      <div className="relative z-10">
        <h3 className="font-bold text-lg mb-3">{name}</h3>
        <div className="grid grid-cols-3 gap-2 text-center">
          <div>
            <div className="text-xs text-slate-400 mb-1">Highest Key</div>
            <div className="text-xl font-bold text-purple-400">+{keyLevel}</div>
          </div>
          <div>
            <div className="text-xs text-slate-400 mb-1">Max Score</div>
            <div className="text-lg font-medium">{maxScore}</div>
          </div>
          <div>
            <div className="text-xs text-slate-400 mb-1">Avg Score</div>
            <div className="text-lg font-medium">{score}</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DungeonCard;
