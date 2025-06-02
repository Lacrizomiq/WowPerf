// LeaderBoardCards.tsx - Refactoré avec tableau shadcn/ui

"use client";

import React, { useState, useEffect } from "react";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useGetRaiderioRaidLeaderboard } from "@/hooks/useRaiderioApi";
import {
  RaidRankings,
  RaidRanking,
  EncounterDefeated,
  EncounterPulled,
} from "@/types/raidLeaderboard";
import Pagination from "../../../Shared/Pagination";

interface LeaderBoardCardsProps {
  raid: string;
  difficulty: string;
  region: string;
  limit: number;
  page: number;
  onPageChange?: (page: number) => void;
}

const LeaderBoardCards: React.FC<LeaderBoardCardsProps> = ({
  raid,
  difficulty,
  region,
  limit,
  page,
  onPageChange,
}) => {
  const [totalPages, setTotalPages] = useState(5);

  const { data, isLoading, error } = useGetRaiderioRaidLeaderboard(
    raid,
    difficulty,
    region,
    limit,
    page
  );

  // Mise à jour du nombre total de pages basé sur les données
  useEffect(() => {
    if (data && data.totalCount) {
      // Si l'API fournit un compte total de guildes
      setTotalPages(Math.ceil(data.totalCount / limit));
    } else if (
      data &&
      data.raidRankings &&
      data.raidRankings.length === limit
    ) {
      // Si nous avons une page complète, supposons qu'il y a au moins une page de plus
      setTotalPages(Math.max(totalPages, page + 2));
    } else if (data && data.raidRankings) {
      // Si nous n'avons pas une page complète, c'est probablement la dernière
      setTotalPages(page + 1);
    }
  }, [data, limit, page, totalPages]);

  // Gérer le changement de page
  const handlePageChange = (newPage: number) => {
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
  if (!data || !data.raidRankings)
    return (
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 text-center">
        <p className="text-slate-400">No data available</p>
      </div>
    );

  return (
    <div>
      <div className="rounded-md border border-slate-700 overflow-hidden">
        <Table>
          <TableHeader className="bg-slate-800/50">
            <TableRow className="hover:bg-slate-800/70 border-slate-700">
              <TableHead className="text-white font-medium w-16 text-center">
                Rank
              </TableHead>
              <TableHead className="text-white font-medium">Guild</TableHead>
              <TableHead className="text-white font-medium">Realm</TableHead>
              <TableHead className="text-white font-medium">Region</TableHead>
              <TableHead className="text-white font-medium text-center">
                Progress
              </TableHead>
              <TableHead className="text-white font-medium text-right pr-6">
                Region Rank
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.raidRankings.map((ranking: RaidRanking, index: number) => (
              <TableRow
                key={index}
                className="hover:bg-slate-800/70 border-slate-700"
              >
                <TableCell className="font-medium text-white text-center">
                  {ranking.rank}
                </TableCell>
                <TableCell className="font-bold text-white">
                  {ranking.guild.name}
                </TableCell>
                <TableCell className="text-slate-300">
                  {ranking.guild.realm.name}
                </TableCell>
                <TableCell className="text-slate-300">
                  {ranking.guild.region.name}
                </TableCell>
                <TableCell className="text-center">
                  <span className="bg-slate-800/70 px-3 py-1 rounded-full text-white font-medium">
                    {ranking.encountersDefeated.length}/8M
                  </span>
                </TableCell>
                <TableCell className="text-right pr-6 font-medium text-slate-300">
                  {ranking.regionRank}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
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

export default LeaderBoardCards;
