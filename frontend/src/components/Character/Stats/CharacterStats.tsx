import React from "react";
import {
  Sword,
  Zap,
  Crown,
  Activity,
  Forward,
  Gauge,
  Heart,
  ArrowRight,
  Ban,
} from "lucide-react";
import { getAttackTypeForSpec, AttackType } from "@/utils/specmapping";

interface CharacterStatsProps {
  stats: any;
  specId: number;
}

interface StatDisplay {
  name: string;
  value: number;
  rating: number;
  icon: React.ReactNode;
  color: string;
  bgColor: string;
  borderColor: string;
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
      name: "CRITICAL STRIKE",
      value: getStatValues("crit").value,
      rating: getStatValues("crit").rating,
      icon: <Sword className="w-8 h-8" />,
      color: "rgb(255, 45, 45)",
      bgColor: "rgba(255, 45, 45, 0.1)",
      borderColor: "rgba(255, 45, 45, 0.3)",
    },
    {
      name: "HASTE",
      value: getStatValues("haste").value,
      rating: getStatValues("haste").rating,
      icon: <Zap className="w-8 h-8" />,
      color: "rgb(0, 255, 200)",
      bgColor: "rgba(0, 255, 200, 0.1)",
      borderColor: "rgba(0, 255, 200, 0.3)",
    },
    {
      name: "MASTERY",
      value: stats.mastery.value,
      rating: stats.mastery.rating,
      icon: <Crown className="w-8 h-8" />,
      color: "rgb(153, 102, 255)",
      bgColor: "rgba(153, 102, 255, 0.1)",
      borderColor: "rgba(153, 102, 255, 0.3)",
    },
    {
      name: "VERSATILITY",
      value: stats.versatility_damage_done_bonus,
      rating: stats.versatility,
      icon: <Activity className="w-8 h-8" />,
      color: "rgb(128, 128, 128)",
      bgColor: "rgba(128, 128, 128, 0.1)",
      borderColor: "rgba(128, 128, 128, 0.3)",
    },
  ].sort((a, b) => b.rating - a.rating);

  const minorStats: StatDisplay[] = [
    {
      name: "SPEED",
      value: stats.speed.rating_bonus || 0,
      rating: stats.speed.rating || 0,
      icon: <Gauge className="w-8 h-8" />,
      color: "rgb(255, 140, 0)",
      bgColor: "rgba(255, 140, 0, 0.1)",
      borderColor: "rgba(255, 140, 0, 0.3)",
    },
    {
      name: "AVOIDANCE",
      value: stats.avoidance.rating_bonus,
      rating: stats.avoidance.rating,
      icon: <Ban className="w-8 h-8" />,
      color: "rgb(148, 130, 201)",
      bgColor: "rgba(148, 130, 201, 0.1)",
      borderColor: "rgba(148, 130, 201, 0.3)",
    },
    {
      name: "LEECH",
      value: stats.lifesteal.value || 0,
      rating: stats.lifesteal.rating || 0,
      icon: <Heart className="w-8 h-8" />,
      color: "rgb(0, 112, 221)",
      bgColor: "rgba(0, 112, 221, 0.1)",
      borderColor: "rgba(0, 112, 221, 0.3)",
    },
  ].sort((a, b) => b.rating - a.rating);

  const StatBox = ({ stat }: { stat: StatDisplay }) => (
    <div
      className="relative flex items-center p-3 sm:p-4 rounded-lg transition-all duration-300 hover:bg-opacity-20"
      style={{
        backgroundColor: stat.bgColor,
        border: `1px solid ${stat.borderColor}`,
      }}
    >
      <div className="mr-3 sm:mr-4 flex-shrink-0" style={{ color: stat.color }}>
        {React.cloneElement(stat.icon as React.ReactElement, {
          className: "w-6 h-6 sm:w-8 sm:h-8",
        })}
      </div>
      <div className="flex flex-col min-w-0">
        <span className="text-xs sm:text-sm text-gray-400 truncate">
          {stat.name}
        </span>
        <span
          className="text-lg sm:text-xl font-bold"
          style={{ color: stat.color }}
        >
          {stat.value.toFixed(1)}%
        </span>
        {stat.rating > 0 && (
          <span className="text-xs text-gray-500 hidden sm:block">
            Rating: {stat.rating.toLocaleString()}
          </span>
        )}
      </div>
    </div>
  );

  const StatPriorityChain = ({ stats }: { stats: StatDisplay[] }) => (
    <div className="flex flex-wrap items-center justify-center gap-1 sm:gap-2 mt-3 sm:mt-4 p-2 sm:px-4 sm:py-3 bg-black/20 rounded-lg text-center">
      {stats.map((stat, index) => (
        <React.Fragment key={stat.name}>
          <span
            style={{ color: stat.color }}
            className="font-bold text-sm sm:text-base md:text-lg"
          >
            {stat.name.split(" ")[0]}
          </span>
          {index < stats.length - 1 && (
            <ArrowRight className="text-gray-600 hidden sm:block" size={16} />
          )}
          {index < stats.length - 1 && (
            <span className="text-gray-600 sm:hidden">→</span>
          )}
        </React.Fragment>
      ))}
    </div>
  );

  return (
    <div className="w-full space-y-6 sm:space-y-8 p-2 sm:p-4">
      <div>
        <h2 className="text-lg sm:text-xl font-bold text-white mb-3 sm:mb-4 px-1">
          Secondary Stats
        </h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-2 sm:gap-4">
          {secondaryStats.map((stat) => (
            <StatBox key={stat.name} stat={stat} />
          ))}
        </div>
        <StatPriorityChain stats={secondaryStats} />
      </div>

      <div>
        <h2 className="text-lg sm:text-xl font-bold text-white mb-3 sm:mb-4 px-1">
          Minor Stats
        </h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2 sm:gap-4">
          {minorStats.map((stat) => (
            <StatBox key={stat.name} stat={stat} />
          ))}
        </div>
        <StatPriorityChain stats={minorStats} />
      </div>
    </div>
  );
};

export default CharacterStats;
