"use client";

import React, { useMemo, useState } from "react";
import { useGetRaiderioMythicPlusBestRuns } from "@/hooks/useRaiderioApi";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import { useGetRaiderioMythicPlusRunDetails } from "@/hooks/useRaiderioApi";
import { Dungeon } from "@/types/mythicPlusRuns";
import Image from "next/image";
import { Star, ChevronDown, ChevronUp } from "lucide-react";
import Link from "next/link";
import RunsDetails from "./runsDetails/RunsDetails";

interface RunsCardProps {
  season: string;
  region: string;
  dungeon: string;
  page: number;
}

interface RosterMember {
  role: "tank" | "healer" | "dps";
  character: {
    id: number;
    name: string;
    class: {
      id: number;
      name: string;
      slug: string;
    };
    realm: {
      altSlug: string;
      slug: string;
    };
    region: {
      slug: string;
    };
    spec: {
      id: number;
      name: string;
      slug: string;
    };
  };
}

const roleIcons = {
  tank: "https://cdn.raiderio.net/assets/img/role_tank-6cee7610058306ba277e82c392987134.png",
  healer:
    "https://cdn.raiderio.net/assets/img/role_healer-984e5e9867d6508a714a9c878d87441b.png",
  dps: "https://cdn.raiderio.net/assets/img/role_dps-eb25989187d4d3ac866d609dc009f090.png",
};

const RunsCard: React.FC<RunsCardProps> = ({
  season,
  region,
  dungeon,
  page,
}) => {
  const [selectedRunId, setSelectedRunId] = useState<number | null>(null);

  const { data: dungeonData } =
    useGetBlizzardMythicDungeonPerSeason("season-tww-1");
  const {
    data: mythicPlusData,
    isLoading,
    error,
  } = useGetRaiderioMythicPlusBestRuns(season, region, dungeon, page);
  const { data: runDetails } = useGetRaiderioMythicPlusRunDetails(
    season,
    mythicPlusData?.rankings[0].run.keystone_run_id
  );

  const dungeonMap = useMemo(() => {
    if (dungeonData?.dungeons) {
      return dungeonData.dungeons.reduce(
        (acc: Record<string, Dungeon>, dungeon: Dungeon) => {
          acc[dungeon.Slug.toLowerCase()] = dungeon;
          return acc;
        },
        {}
      );
    }
    return {};
  }, [dungeonData]);

  const sortRoster = (roster: RosterMember[]) => {
    const roleOrder: { [key: string]: number } = { tank: 1, healer: 2, dps: 3 };
    return roster.sort((a, b) => roleOrder[a.role] - roleOrder[b.role]);
  };

  if (isLoading)
    return <div className="text-white text-center p-4">Loading...</div>;
  if (error)
    return (
      <div className="text-red-500 text-center p-4">
        Error: {(error as Error).message}
      </div>
    );
  if (!mythicPlusData || !mythicPlusData.rankings)
    return (
      <div className="text-yellow-500 text-center p-4">No data available</div>
    );

  return (
    <div className="space-y-6">
      {mythicPlusData.rankings.map((ranking: any) => {
        const dungeonSlug = ranking.run.dungeon.slug.toLowerCase();
        const dungeonInfo = dungeonMap[dungeonSlug];
        const sortedRoster = sortRoster(ranking.run.roster);
        const isSelected = selectedRunId === ranking.run.keystone_run_id;

        return (
          <div
            key={ranking.run.keystone_run_id}
            className="bg-deep-blue bg-opacity-80 rounded-2xl overflow-hidden shadow-2xl glow-effect"
          >
            <div className="flex">
              <div className="w-1/3 relative">
                <Image
                  src={dungeonInfo?.MediaURL || "/placeholder.jpg"}
                  alt={ranking.run.dungeon.name}
                  layout="fill"
                  objectFit="cover"
                />
                <div className="absolute inset-0 bg-black bg-opacity-50 flex flex-col items-center justify-center">
                  <h2 className="text-white text-xl font-bold mb-2">
                    {ranking.run.dungeon.name}
                  </h2>
                  <div className="flex items-center">
                    <span className="text-white text-4xl font-bold">
                      +{ranking.run.mythic_level}
                    </span>
                    <div className="ml-2">
                      {[...Array(ranking.run.num_chests)].map((_, i) => (
                        <Star
                          key={i}
                          className="inline-block text-yellow-400 w-6 h-6"
                        />
                      ))}
                    </div>
                  </div>
                </div>
              </div>
              <div className="w-2/3 p-4 flex flex-col justify-between">
                <div>
                  <div className="flex justify-between items-center mb-4">
                    <p className="text-white text-2xl font-bold">
                      Rank: {ranking.rank}
                    </p>
                    <p className="text-white text-xl">
                      Score: {ranking.score.toFixed(1)}
                    </p>
                    <button
                      onClick={() =>
                        setSelectedRunId(
                          isSelected ? null : ranking.run.keystone_run_id
                        )
                      }
                      className="text-blue-500 hover:text-blue-400 font-semibold flex items-center"
                    >
                      {isSelected ? (
                        <>
                          Hide Details <ChevronUp className="ml-1 w-4 h-4" />
                        </>
                      ) : (
                        <>
                          View Details <ChevronDown className="ml-1 w-4 h-4" />
                        </>
                      )}
                    </button>
                  </div>
                  <p className="text-white text-sm mb-2">
                    Time: {(ranking.run.clear_time_ms / 1000 / 60).toFixed(2)}{" "}
                    min
                  </p>
                  <p className="text-white text-sm mb-2">
                    Completed at:{" "}
                    {new Date(ranking.run.completed_at).toUTCString()}
                  </p>
                  <p className="text-white mb-4">
                    Affixes:{" "}
                    {ranking.run.weekly_modifiers
                      .map((mod: any) => mod.name)
                      .join(", ")}
                  </p>
                </div>
                <div>
                  <h3 className="text-white font-semibold mb-2">
                    Team Composition:
                  </h3>
                  <div className="flex justify-between">
                    {sortedRoster.map((member: RosterMember) => (
                      <div
                        key={member.character.id}
                        className="text-center mb-2"
                      >
                        <div
                          className={`w-12 h-12 mx-auto rounded-full bg-${member.character.class.slug} flex items-center justify-center relative`}
                        >
                          <Image
                            src={roleIcons[member.role]}
                            alt={member.role}
                            width={24}
                            height={24}
                          />
                        </div>
                        <Link
                          href={`/character/${member.character.region.slug}/${
                            member.character.realm.slug
                          }/${member.character.name.toLowerCase()}`}
                        >
                          <p
                            className={`font-bold ${
                              member.character.class.slug
                                ? `class-color--${member.character.class.slug} hover:underline hover:decoration-current`
                                : ""
                            }`}
                          >
                            {member.character.name}
                          </p>
                        </Link>
                        <p className="text-white text-xs">
                          {member.character.spec.name}
                        </p>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
            {isSelected && (
              <div className="border-t border-gray-600  pt-4">
                <RunsDetails
                  season={season}
                  runId={ranking.run.keystone_run_id}
                />
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
};

export default RunsCard;
