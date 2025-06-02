// components/Statistics/mythicplus/KeyLevelDistributionSection.tsx
import React from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Progress } from "@/components/ui/progress";
import { useKeyLevelDistribution } from "@/hooks/useMythicPlusRunsAnalysis";
import LoadingSpinner from "../shared/LoadingSpinner";
import ErrorDisplay from "../shared/ErrorDisplay";
import InfoTooltip from "@/components/Shared/InfoTooltip";

/**
 * Section qui affiche la distribution des runs par niveau de clé Mythic+
 */
const KeyLevelDistributionSection: React.FC = () => {
  const {
    data: distribution,
    isLoading,
    error,
    isError,
  } = useKeyLevelDistribution();

  if (isLoading) {
    return (
      <section>
        <h2 className="text-2xl font-bold mb-4">Key Level Distribution</h2>
        <LoadingSpinner />
      </section>
    );
  }

  if (isError || !distribution) {
    return (
      <section>
        <h2 className="text-2xl font-bold mb-4">Key Level Distribution</h2>
        <ErrorDisplay
          error={error}
          message="Unable to load key level distribution"
        />
      </section>
    );
  }

  // Calculer la valeur maximale pour les progress bars
  const maxCount = Math.max(...distribution.map((level) => level.count));

  return (
    <section>
      <h2 className="text-2xl font-bold mb-4 flex items-center">
        Key Level Distribution
        <InfoTooltip
          content="This section shows the distribution of runs by key level difficulty. It includes the percentage of runs by key level and the average score for each key level."
          className="ml-2"
        />
      </h2>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Graphique de répartition */}
        <Card className="bg-slate-800/30 border-slate-700">
          <CardHeader>
            <CardTitle>Distribution by Key Level</CardTitle>
            <CardDescription>
              Percentage of runs by key level difficulty
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3 max-h-96 overflow-y-auto">
              {distribution.slice(0, 15).map((level) => (
                <div
                  key={level.mythic_level}
                  className="flex items-center gap-3"
                >
                  <div className="w-12 text-sm font-medium text-white">
                    +{level.mythic_level}
                  </div>
                  <div className="flex-1">
                    <Progress
                      value={(level.count / maxCount) * 100}
                      className="h-3"
                    />
                  </div>
                  <div className="w-16 text-sm text-slate-400 text-right">
                    {level.percentage.toFixed(1)}%
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Tableau détaillé */}
        <Card className="bg-slate-800/30 border-slate-700">
          <CardHeader>
            <CardTitle>Detailed Data</CardTitle>
            <CardDescription>
              Top key levels with complete statistics
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto max-h-96 overflow-y-auto">
              <Table>
                <TableHeader>
                  <TableRow className="border-slate-700">
                    <TableHead>Key Level</TableHead>
                    <TableHead className="text-right">Runs</TableHead>
                    <TableHead className="text-right">%</TableHead>
                    <TableHead className="text-right">Average Score</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {distribution.slice(0, 10).map((level) => (
                    <TableRow
                      key={level.mythic_level}
                      className="border-slate-700 hover:bg-slate-800/30 transition-colors duration-150"
                    >
                      <TableCell className="font-medium">
                        <span className="text-white font-bold">
                          +{level.mythic_level}
                        </span>
                      </TableCell>
                      <TableCell className="text-right font-mono text-white">
                        {level.count.toLocaleString("en-US")}
                      </TableCell>
                      <TableCell className="text-right font-semibold text-white">
                        {level.percentage.toFixed(1)}%
                      </TableCell>
                      <TableCell className="text-right font-mono text-white">
                        {level.avg_score.toFixed(0)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      </div>
    </section>
  );
};

export default KeyLevelDistributionSection;
