import React from "react";
import { getAttackTypeForSpec, AttackType } from "@/utils/specmapping";
import { ArrowRight } from "lucide-react";

interface CharacterStatsProps {
  stats: any;
  specId: number;
}

interface StatDisplay {
  name: string;
  rating: number;
  value: number;
  color: string;
}

const CharacterStats: React.FC<CharacterStatsProps> = ({ stats, specId }) => {
  const getStatValues = (statType: "crit" | "haste") => {
    const attackType = getAttackTypeForSpec(specId);

    switch (attackType) {
      case AttackType.RANGED:
        return {
          rating:
            statType === "crit"
              ? stats.ranged_crit.rating
              : stats.ranged_haste.rating,
          value:
            statType === "crit"
              ? stats.ranged_crit.value
              : stats.ranged_haste.value,
        };
      case AttackType.SPELL:
        return {
          rating:
            statType === "crit"
              ? stats.spell_crit.rating
              : stats.spell_haste.rating,
          value:
            statType === "crit"
              ? stats.spell_crit.value
              : stats.spell_haste.value,
        };
      case AttackType.MELEE:
      default:
        return {
          rating:
            statType === "crit"
              ? stats.melee_crit.rating
              : stats.melee_haste.rating,
          value:
            statType === "crit"
              ? stats.melee_crit.value
              : stats.melee_haste.value,
        };
    }
  };

  const secondaryStats: StatDisplay[] = [
    {
      name: "Critical Strike",
      ...getStatValues("crit"),
      color: "rgb(255, 255, 255)", // White
    },
    {
      name: "Haste",
      ...getStatValues("haste"),
      color: "rgb(148, 130, 201)", // Purple
    },
    {
      name: "Mastery",
      rating: stats.mastery.rating,
      value: stats.mastery.value,
      color: "rgb(0, 112, 221)", // Blue
    },
    {
      name: "Versatility",
      rating: stats.versatility,
      value: stats.versatility_damage_done_bonus,
      color: "rgb(255, 140, 0)", // Orange
    },
  ].sort((a, b) => b.rating - a.rating);

  const minorStats: StatDisplay[] = [
    {
      name: "Speed",
      rating: stats.speed.rating || 0,
      value: stats.speed.rating_bonus || 0,
      color: "rgb(255, 140, 0)", // Orange
    },
    {
      name: "Avoidance",
      rating: stats.avoidance.rating,
      value: stats.avoidance.rating_bonus,
      color: "rgb(148, 130, 201)", // Purple
    },
    {
      name: "Leech",
      rating: stats.lifesteal.rating || 0,
      value: stats.lifesteal.value || 0,
      color: "rgb(0, 112, 221)", // Blue
    },
  ].sort((a, b) => b.rating - a.rating);

  const StatItem = ({ stat }: { stat: StatDisplay }) => (
    <div className="flex flex-col w-full">
      <div className="bg-black/40 rounded-lg px-4 py-3 flex flex-col">
        <div className="flex justify-between items-center">
          <span style={{ color: stat.color }} className="font-bold text-lg">
            {stat.name.toUpperCase()}
          </span>
          <span className="text-2xl font-bold" style={{ color: stat.color }}>
            {stat.value.toFixed(2)}%
          </span>
        </div>
        <div className="text-gray-400 text-sm mt-1 font-bold">
          {stat.rating.toLocaleString()}
        </div>
      </div>
    </div>
  );

  const StatPriorityChain = ({ stats }: { stats: StatDisplay[] }) => (
    <div className="flex flex-wrap items-center justify-center gap-2 mt-6 px-4 py-3 bg-black/20 rounded-lg">
      {stats.map((stat, index) => (
        <React.Fragment key={stat.name}>
          <span style={{ color: stat.color }} className="font-bold text-lg">
            {stat.name.split(" ")[0].toUpperCase()}
          </span>
          {index < stats.length - 1 && (
            <ArrowRight className="text-gray-600" size={20} />
          )}
        </React.Fragment>
      ))}
    </div>
  );

  return (
    <div className="w-full space-y-6 px-4 lg:px-6">
      <div>
        <h2 className="text-2xl font-bold text-white mb-4">Secondary Stats</h2>
        <div className="flex flex-col space-y-2 lg:grid lg:grid-cols-2 lg:gap-4 lg:space-y-0 xl:grid-cols-4 bg-deep-blue p-6 rounded-lg">
          {secondaryStats.map((stat) => (
            <StatItem key={stat.name} stat={stat} />
          ))}
        </div>
        <StatPriorityChain stats={secondaryStats} />
      </div>

      <div>
        <h2 className="text-2xl font-bold text-white mb-4">Minor Stats</h2>
        <div className="flex flex-col space-y-2 lg:grid lg:grid-cols-2 lg:gap-4 lg:space-y-0 xl:grid-cols-3 bg-deep-blue p-6 rounded-lg">
          {minorStats.map((stat) => (
            <StatItem key={stat.name} stat={stat} />
          ))}
        </div>
        <StatPriorityChain stats={minorStats} />
      </div>
    </div>
  );
};

export default CharacterStats;
