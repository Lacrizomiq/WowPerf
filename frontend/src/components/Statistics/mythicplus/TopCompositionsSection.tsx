// components/Statistics/mythicplus/TopCompositionsSection.tsx
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
import { Badge } from "@/components/ui/badge";
import { useTopTeamCompositionsGlobal } from "@/hooks/useMythicPlusRunsAnalysis";
import { ClassColoredText } from "../shared/ClassColorUtils";
import LoadingSpinner from "../shared/LoadingSpinner";
import ErrorDisplay from "../shared/ErrorDisplay";

/**
 * Section qui affiche les compositions d'équipes les plus populaires
 */
const TopCompositionsSection: React.FC = () => {
  const {
    data: compositions,
    isLoading,
    error,
    isError,
  } = useTopTeamCompositionsGlobal({
    limit: 10,
    min_usage: 10,
  });

  if (isLoading) {
    return (
      <section>
        <h2 className="text-2xl font-bold mb-4">Top Team Compositions</h2>
        <LoadingSpinner />
      </section>
    );
  }

  if (isError || !compositions) {
    return (
      <section>
        <h2 className="text-2xl font-bold mb-4">Top Team Compositions</h2>
        <ErrorDisplay
          error={error}
          message="Unable to load top team compositions"
        />
      </section>
    );
  }

  // Fonction pour extraire le nom de classe d'un display string comme "Warrior - Protection"
  const extractClassName = (display: string): string => {
    return display.split(" - ")[0] || display;
  };

  // Fonction pour extraire le nom de spec d'un display string
  const extractSpecName = (display: string): string => {
    return display.split(" - ")[1] || display;
  };

  return (
    <section>
      <h2 className="text-2xl font-bold mb-4">Top Team Compositions</h2>

      <Card className="bg-slate-800/30 border-slate-700">
        <CardHeader>
          <CardTitle>Most Popular Compositions</CardTitle>
          <CardDescription>
            Based on usage in Mythic+ high-level runs (minimum 10 usages)
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-slate-700">
                  <TableHead className="w-16">Rank</TableHead>
                  <TableHead>Tank</TableHead>
                  <TableHead>Healer</TableHead>
                  <TableHead>DPS 1</TableHead>
                  <TableHead>DPS 2</TableHead>
                  <TableHead>DPS 3</TableHead>
                  <TableHead className="text-right">Usage</TableHead>
                  <TableHead className="text-right">%</TableHead>
                  <TableHead className="text-right">Average Score</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {compositions.map((comp) => (
                  <TableRow
                    key={comp.rank}
                    className="border-slate-700 hover:bg-slate-800/30 transition-colors duration-150"
                  >
                    {/* Rang */}
                    <TableCell className="font-medium">
                      <Badge
                        variant="outline"
                        className="bg-slate-700/50 text-slate-300 border-slate-600"
                      >
                        #{comp.rank}
                      </Badge>
                    </TableCell>

                    {/* Tank */}
                    <TableCell className="text-sm">
                      <div className="space-y-1">
                        <ClassColoredText
                          className={extractClassName(comp.tank)}
                          additionalClasses="font-semibold"
                        >
                          {extractSpecName(comp.tank)}
                        </ClassColoredText>
                        <div className="text-xs text-slate-400">
                          {extractClassName(comp.tank)}
                        </div>
                      </div>
                    </TableCell>

                    {/* Healer */}
                    <TableCell className="text-sm">
                      <div className="space-y-1">
                        <ClassColoredText
                          className={extractClassName(comp.healer)}
                          additionalClasses="font-semibold"
                        >
                          {extractSpecName(comp.healer)}
                        </ClassColoredText>
                        <div className="text-xs text-slate-400">
                          {extractClassName(comp.healer)}
                        </div>
                      </div>
                    </TableCell>

                    {/* DPS 1 */}
                    <TableCell className="text-sm">
                      <div className="space-y-1">
                        <ClassColoredText
                          className={extractClassName(comp.dps1)}
                          additionalClasses="font-semibold"
                        >
                          {extractSpecName(comp.dps1)}
                        </ClassColoredText>
                        <div className="text-xs text-slate-400">
                          {extractClassName(comp.dps1)}
                        </div>
                      </div>
                    </TableCell>

                    {/* DPS 2 */}
                    <TableCell className="text-sm">
                      <div className="space-y-1">
                        <ClassColoredText
                          className={extractClassName(comp.dps2)}
                          additionalClasses="font-semibold"
                        >
                          {extractSpecName(comp.dps2)}
                        </ClassColoredText>
                        <div className="text-xs text-slate-400">
                          {extractClassName(comp.dps2)}
                        </div>
                      </div>
                    </TableCell>

                    {/* DPS 3 */}
                    <TableCell className="text-sm">
                      <div className="space-y-1">
                        <ClassColoredText
                          className={extractClassName(comp.dps3)}
                          additionalClasses="font-semibold"
                        >
                          {extractSpecName(comp.dps3)}
                        </ClassColoredText>
                        <div className="text-xs text-slate-400">
                          {extractClassName(comp.dps3)}
                        </div>
                      </div>
                    </TableCell>

                    {/* Utilisation */}
                    <TableCell className="text-right font-mono">
                      <div className="text-base font-semibold text-white">
                        {comp.usage_count.toLocaleString("fr-FR")}
                      </div>
                    </TableCell>

                    {/* Pourcentage */}
                    <TableCell className="text-right">
                      <div className="text-base font-bold text-white">
                        {comp.percentage.toFixed(1)}%
                      </div>
                    </TableCell>

                    {/* Score moyen */}
                    <TableCell className="text-right font-mono">
                      <div className="text-base font-semibold text-white">
                        {comp.avg_score.toFixed(0)}
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {/* Statistiques résumées */}
          <div className="mt-6 p-4 bg-slate-800/50 rounded-lg border border-slate-700">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-center">
              <div>
                <div className="text-2xl font-bold text-white">
                  {compositions.length}
                </div>
                <div className="text-sm text-slate-400">
                  Popular compositions
                </div>
              </div>
              <div>
                <div className="text-2xl font-bold text-white">
                  {compositions
                    .reduce((sum, comp) => sum + comp.usage_count, 0)
                    .toLocaleString("fr-FR")}
                </div>
                <div className="text-sm text-slate-400">Total usages</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-white">
                  {compositions[0]?.avg_score.toFixed(0) || "N/A"}
                </div>
                <div className="text-sm text-slate-400">
                  Average score top 1
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </section>
  );
};

export default TopCompositionsSection;
