// components/Statistics/mythicplus/OverallStatsSection.tsx

import React from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useOverallStats } from "@/hooks/useMythicPlusRunsAnalysis";
import LoadingSpinner from "../shared/LoadingSpinner";
import ErrorDisplay from "../shared/ErrorDisplay";
import InfoTooltip from "@/components/Shared/InfoTooltip";

/**
 * Section qui affiche les statistiques générales du dataset Mythic+
 * Utilise le hook useOverallStats pour récupérer les données
 */
const OverallStatsSection: React.FC = () => {
  const { data: overallStats, isLoading, error, isError } = useOverallStats();

  if (isLoading) {
    return (
      <section>
        <h2 className="text-2xl font-bold mb-4">Dataset Overview</h2>
        <LoadingSpinner />
      </section>
    );
  }

  if (isError || !overallStats) {
    return (
      <section>
        <h2 className="text-2xl font-bold mb-4">Dataset Overview</h2>
        <ErrorDisplay
          error={error}
          message="Unable to load overall statistics"
        />
      </section>
    );
  }

  // Calcul de la période d'analyse
  const getAnalysisPeriod = () => {
    if (!overallStats.oldest_run || !overallStats.newest_run) {
      return "Data available";
    }

    const oldestDate = new Date(overallStats.oldest_run);
    const newestDate = new Date(overallStats.newest_run);

    const diffTime = Math.abs(newestDate.getTime() - oldestDate.getTime());
    const diffMonths = Math.ceil(diffTime / (1000 * 60 * 60 * 24 * 30));

    return `${diffMonths} months`;
  };

  const getAnalysisPeriodDescription = () => {
    if (!overallStats.oldest_run || !overallStats.newest_run) {
      return "Analysis period";
    }

    const oldestDate = new Date(overallStats.oldest_run);
    const newestDate = new Date(overallStats.newest_run);

    const formatDate = (date: Date) => {
      return date.toLocaleDateString("fr-FR", {
        month: "short",
        year: "numeric",
      });
    };

    return `${formatDate(oldestDate)} - ${formatDate(newestDate)}`;
  };

  return (
    <section>
      <h2 className="text-2xl font-bold mb-4 flex items-center">
        Dataset Overview
        <InfoTooltip
          content="This section provides an overview of the dataset used to generate the statistics."
          className="ml-2"
          size="lg"
        />
      </h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Total des Runs */}
        <StatCard
          title="Total Runs"
          value={overallStats.total_runs.toLocaleString("fr-FR")}
          description="Runs analyzed"
          tooltip="Total number of Mythic+ dungeon runs recorded in our dataset"
        />

        {/* Runs avec Score 
        <StatCard
          title="Runs with Score"
          value={overallStats.runs_with_score.toLocaleString("fr-FR")}
          description={`${(
            (overallStats.runs_with_score / overallStats.total_runs) *
            100
          ).toFixed(1)}% of total`}
        />
        */}

        {/* Compositions Uniques */}
        <StatCard
          title="Unique Compositions"
          value={overallStats.unique_compositions.toLocaleString("fr-FR")}
          description="Team combinations"
          tooltip="Number of distinct team compositions (tank, healer, 3 DPS specializations) found across all runs."
        />

        {/* Score Moyen 
        <StatCard
          title="Average Score"
          value={overallStats.avg_score.toFixed(1)}
          description="Mythic+ Score"
        />
        */}

        {/* Niveau de Clé Moyen */}
        <StatCard
          title="Average Key Level"
          value={`+${overallStats.avg_key_level.toFixed(1)}`}
          description="Average difficulty"
          tooltip="Average keystone level across all recorded runs. Higher levels indicate more challenging content with better rewards."
        />

        {/* Donjons Analysés 
        <StatCard
          title="Analyzed Dungeons"
          value={overallStats.unique_dungeons.toString()}
          description="Different dungeons"
        />
        */}

        {/* Régions Couvertes */}
        <StatCard
          title="Covered Regions"
          value={overallStats.unique_regions.toString()}
          description="Global regions"
          tooltip="Number of different regions covered by the dataset. This includes EU, US, KR and TW. CN runs are not included."
        />

        {/* Période d'Analyse 
        <StatCard
          title="Analysis Period"
          value={getAnalysisPeriod()}
          description={getAnalysisPeriodDescription()}
        />
        */}
      </div>
    </section>
  );
};

/**
 * Composant StatCard réutilisable pour afficher une statistique
 */
interface StatCardProps {
  title: string;
  value: string;
  description?: string;
  tooltip?: string;
  trend?: {
    value: number;
    isPositive: boolean;
  };
}

const StatCard: React.FC<StatCardProps> = ({
  title,
  value,
  description,
  tooltip,
  trend,
}) => {
  return (
    <Card className="bg-slate-800/30 border-slate-700 hover:bg-slate-800/40 transition-colors duration-200">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-slate-300">
          {title}
          {tooltip && (
            <InfoTooltip
              content={tooltip}
              side="right"
              delayDuration={100}
              size="sm"
              className="ml-1"
            />
          )}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between">
          <div className="text-2xl font-bold text-white">{value}</div>
          {trend && (
            <div
              className={`text-sm flex items-center ${
                trend.isPositive ? "text-green-500" : "text-red-500"
              }`}
            >
              <span className="mr-1">{trend.isPositive ? "↗" : "↘"}</span>
              {Math.abs(trend.value)}%
            </div>
          )}
        </div>
        {description && (
          <CardDescription className="mt-1 text-slate-400">
            {description}
          </CardDescription>
        )}
      </CardContent>
    </Card>
  );
};

export default OverallStatsSection;
