// RunsCard.tsx - Avec pagination
import React, { useMemo, useState, useRef, useEffect } from "react";
import { useGetRaiderioMythicPlusBestRuns } from "@/hooks/useRaiderioApi";
import { useGetBlizzardMythicDungeonPerSeason } from "@/hooks/useBlizzardApi";
import { useGetRaiderioMythicPlusRunDetails } from "@/hooks/useRaiderioApi";
import { Dungeon } from "@/types/mythicPlusRuns";
import Image from "next/image";
import { Star, ChevronDown, ChevronUp } from "lucide-react";
import Link from "next/link";
import RunsDetails from "./runsDetails/RunsDetails";
import Pagination from "@/components/Shared/Pagination";

interface RunsCardProps {
  season: string;
  region: string;
  dungeon: string;
  page: number;
  onPageChange?: (page: number) => void; // Nouveau prop pour le changement de page
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
  onPageChange,
}) => {
  const [selectedRunId, setSelectedRunId] = useState<number | null>(null);
  const cardRefs = useRef<{ [key: number]: HTMLDivElement | null }>({});
  // Nombre total de pages estimé (vous devrez l'ajuster selon votre API)
  const [totalPages, setTotalPages] = useState(5);

  const { data: dungeonData } =
    useGetBlizzardMythicDungeonPerSeason("season-tww-2");
  const {
    data: mythicPlusData,
    isLoading,
    error,
  } = useGetRaiderioMythicPlusBestRuns(season, region, dungeon, page);
  const { data: runDetails } = useGetRaiderioMythicPlusRunDetails(
    season,
    mythicPlusData?.rankings[0]?.run.keystone_run_id
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

  // Mise à jour du nombre total de pages basé sur les données
  useEffect(() => {
    console.log("mythicPlusData", mythicPlusData);
    if (mythicPlusData && mythicPlusData.totalCount) {
      // Si l'API fournit un compte total
      const pageSize = 20; // Nombre d'entrées par page, ajustez selon votre API
      setTotalPages(Math.ceil(mythicPlusData.totalCount / pageSize));
    } else if (
      mythicPlusData &&
      mythicPlusData.rankings &&
      mythicPlusData.rankings.length === 20
    ) {
      // Si nous avons une page complète, supposons qu'il y a au moins une page de plus
      setTotalPages(Math.max(totalPages, page + 2));
    } else if (mythicPlusData && mythicPlusData.rankings) {
      // Si nous n'avons pas une page complète, c'est probablement la dernière
      setTotalPages(page + 1);
    }
  }, [mythicPlusData, page, totalPages]);

  useEffect(() => {
    if (selectedRunId && cardRefs.current[selectedRunId]) {
      cardRefs.current[selectedRunId]?.scrollIntoView({
        behavior: "smooth",
        block: "start",
      });
    }
  }, [selectedRunId]);

  const sortRoster = (roster: RosterMember[]) => {
    const roleOrder: { [key: string]: number } = { tank: 1, healer: 2, dps: 3 };
    return roster.sort((a, b) => roleOrder[a.role] - roleOrder[b.role]);
  };

  // Gérer le changement de page
  const handlePageChange = (newPage: number) => {
    setSelectedRunId(null); // Réinitialiser le détail développé lors du changement de page
    if (onPageChange) {
      onPageChange(newPage);
    }
  };

  if (isLoading)
    return (
      <div className="flex justify-center items-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-600"></div>
      </div>
    );
  if (error)
    return (
      <div className="bg-red-900/20 border border-red-500 rounded-md p-4 my-4">
        <h3 className="text-red-500 text-lg font-medium">
          Error: {(error as Error).message}
        </h3>
      </div>
    );
  if (!mythicPlusData || !mythicPlusData.rankings)
    return (
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 text-center">
        <p className="text-slate-400">No data available</p>
      </div>
    );

  return (
    <div>
      <div className="space-y-6">
        {mythicPlusData.rankings.map((ranking: any) => {
          const dungeonSlug = ranking.run.dungeon.slug.toLowerCase();
          const dungeonInfo = dungeonMap[dungeonSlug];
          const sortedRoster = sortRoster(ranking.run.roster);
          const isSelected = selectedRunId === ranking.run.keystone_run_id;

          return (
            <div
              key={ranking.run.keystone_run_id}
              ref={(el) => {
                cardRefs.current[ranking.run.keystone_run_id] = el;
              }}
              className="bg-slate-800/30 rounded-lg border border-slate-700 overflow-hidden"
            >
              <div className="flex flex-col md:flex-row">
                <div className="w-full md:w-1/3 relative">
                  <div className="aspect-w-16 aspect-h-9 md:h-full">
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
                </div>
                <div className="w-full md:w-2/3 p-4 flex flex-col justify-between">
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
                        className="text-purple-500 hover:text-purple-400 font-semibold flex items-center"
                      >
                        {isSelected ? (
                          <>
                            Hide Details <ChevronUp className="ml-1 w-4 h-4" />
                          </>
                        ) : (
                          <>
                            View Details{" "}
                            <ChevronDown className="ml-1 w-4 h-4" />
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
                      {new Date(ranking.run.completed_at).toLocaleString()}
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
                    <div className="flex flex-wrap justify-between">
                      {sortedRoster.map((member: RosterMember) => (
                        <div
                          key={member.character.id}
                          className="text-center mb-2 px-2"
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
                <div className="border-t border-slate-700 pt-4">
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

      {/* Composant de pagination */}
      <Pagination
        currentPage={page}
        totalPages={totalPages}
        onPageChange={handlePageChange}
      />
    </div>
  );
};

export default RunsCard;
