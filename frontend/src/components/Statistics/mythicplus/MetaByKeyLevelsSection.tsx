// components/Statistics/mythicplus/MetaByKeyLevelsSection.tsx

import React, { useState, useRef } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
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
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronUp } from "lucide-react";
import { useMetaByKeyLevels } from "@/hooks/useMythicPlusRunsAnalysis";
import { ClassColoredText, getRoleStyle } from "../shared/ClassColorUtils";
import LoadingSpinner from "../shared/LoadingSpinner";
import ErrorDisplay from "../shared/ErrorDisplay";
import { Role } from "@/types/raiderio/mythicplus_runs/mythicPlusRuns";
import InfoTooltip from "@/components/Shared/InfoTooltip";

/**
 * Section qui affiche les métadonnées par niveau de clé
 * Permet de voir comment la méta évolue selon la difficulté
 */
const MetaByKeyLevelsSection: React.FC = () => {
  const [selectedRole, setSelectedRole] = useState<Role>(Role.DPS);
  const [selectedBracket, setSelectedBracket] = useState("High Keys (18-19)");

  const handleTabChange = (value: string) => {
    // Safe conversion from string to Role enum
    const roleMapping: Record<string, Role> = {
      [Role.TANK]: Role.TANK,
      [Role.HEALER]: Role.HEALER,
      [Role.DPS]: Role.DPS,
    };

    if (roleMapping[value]) {
      setSelectedRole(roleMapping[value]);
    }
  };

  const keyLevelBrackets = [
    { value: "Very High Keys (20+)", label: "Very High Keys (20+)" },
    { value: "High Keys (18-19)", label: "High Keys (18-19)" },
    { value: "Mid Keys (16-17)", label: "Mid Keys (16-17)" },
  ];

  return (
    <section>
      <h2 className="text-2xl font-bold mb-4 flex items-center">
        Meta by Key Levels
        <InfoTooltip
          content="This section shows the meta analysis for each specialization in a specific key level bracket, filtered by the selected role."
          className="ml-2"
          size="lg"
        />
      </h2>

      {/* Key Level Bracket Selector */}
      <div className="mb-6">
        <div className="flex items-center gap-4">
          <label className="text-sm font-medium text-slate-300">
            Key Level Bracket:
          </label>
          <Select value={selectedBracket} onValueChange={setSelectedBracket}>
            <SelectTrigger className="w-64 bg-slate-800/50 border-slate-700 text-slate-200">
              <SelectValue />
            </SelectTrigger>
            <SelectContent className="bg-slate-800 border-slate-700">
              {keyLevelBrackets.map((bracket) => (
                <SelectItem
                  key={bracket.value}
                  value={bracket.value}
                  className="text-slate-200 focus:bg-slate-700 focus:text-white"
                >
                  {bracket.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      <Tabs
        value={selectedRole}
        onValueChange={handleTabChange}
        className="w-full"
      >
        <TabsList className="grid w-full grid-cols-3 bg-slate-800/50 mb-4">
          <TabsTrigger
            value={Role.TANK}
            className="data-[state=active]:bg-purple-600 hover:bg-slate-700 transition-colors duration-200"
          >
            <div className="flex items-center gap-2">
              <div
                className={`w-2 h-2 rounded-full ${getRoleStyle(
                  "tank",
                  "bgColor"
                )}`}
              />
              Tank
            </div>
          </TabsTrigger>
          <TabsTrigger
            value={Role.HEALER}
            className="data-[state=active]:bg-purple-600 hover:bg-slate-700 transition-colors duration-200"
          >
            <div className="flex items-center gap-2">
              <div
                className={`w-2 h-2 rounded-full ${getRoleStyle(
                  "healer",
                  "bgColor"
                )}`}
              />
              Healer
            </div>
          </TabsTrigger>
          <TabsTrigger
            value={Role.DPS}
            className="data-[state=active]:bg-purple-600 hover:bg-slate-700 transition-colors duration-200"
          >
            <div className="flex items-center gap-2">
              <div
                className={`w-2 h-2 rounded-full ${getRoleStyle(
                  "dps",
                  "bgColor"
                )}`}
              />
              DPS
            </div>
          </TabsTrigger>
        </TabsList>

        <TabsContent value={Role.TANK}>
          <KeyLevelMetaTable
            role={Role.TANK}
            selectedBracket={selectedBracket}
          />
        </TabsContent>

        <TabsContent value={Role.HEALER}>
          <KeyLevelMetaTable
            role={Role.HEALER}
            selectedBracket={selectedBracket}
          />
        </TabsContent>

        <TabsContent value={Role.DPS}>
          <KeyLevelMetaTable
            role={Role.DPS}
            selectedBracket={selectedBracket}
          />
        </TabsContent>
      </Tabs>
    </section>
  );
};

/**
 * Table component to display key level meta for a specific role
 */
interface KeyLevelMetaTableProps {
  role: Role;
  selectedBracket: string;
}

const KeyLevelMetaTable: React.FC<KeyLevelMetaTableProps> = ({
  role,
  selectedBracket,
}) => {
  const [showAll, setShowAll] = useState(false);
  const sectionRef = useRef<HTMLDivElement>(null);

  // Determine top_n parameter based on showAll state
  const topN = showAll ? 0 : 5; // 0 = all specs, 5 = top 5 only

  const {
    data: metaData,
    isLoading,
    error,
    isError,
  } = useMetaByKeyLevels({
    top_n: topN,
    min_usage: 5,
  });

  // Toggle function for View More/Less with auto-scroll
  const toggleViewAll = () => {
    const newShowAll = !showAll;
    setShowAll(newShowAll);

    // If switching back to "top 5", scroll to section top
    if (!newShowAll && sectionRef.current) {
      setTimeout(() => {
        sectionRef.current?.scrollIntoView({
          behavior: "smooth",
          block: "start",
        });
      }, 100); // Small delay to let the content update
    }
  };

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (isError || !metaData) {
    return (
      <ErrorDisplay
        error={error}
        message={`Unable to load ${role} meta by key levels`}
      />
    );
  }

  // Filter data by selected bracket and role
  // Convert role enum to match API data format
  const roleForAPI =
    role === "dps"
      ? "DPS"
      : role.charAt(0).toUpperCase() + role.slice(1).toLowerCase();

  const filteredData = metaData.filter(
    (item) =>
      item.key_level_bracket === selectedBracket && item.role === roleForAPI
  );

  if (filteredData.length === 0) {
    return (
      <Card className="bg-slate-800/30 border-slate-700">
        <CardContent className="py-8">
          <div className="text-center text-slate-400">
            No data available for {role} in {selectedBracket}
          </div>
        </CardContent>
      </Card>
    );
  }

  // Calculate maximum value for progress bars
  const maxPercentage = Math.max(
    ...filteredData.map((spec) => spec.percentage)
  );

  const getRoleDisplayName = (role: Role): string => {
    const roleNames: Record<Role, string> = {
      [Role.TANK]: "TANK",
      [Role.HEALER]: "HEALER",
      [Role.DPS]: "DPS",
    };
    return roleNames[role] || role.toString().toUpperCase();
  };

  const getRoleBadgeColor = (role: Role): string => {
    const colors: Record<Role, string> = {
      [Role.TANK]: "bg-blue-500/20 text-blue-400 border-blue-500/50",
      [Role.HEALER]: "bg-green-500/20 text-green-400 border-green-500/50",
      [Role.DPS]: "bg-red-500/20 text-red-400 border-red-500/50",
    };
    return colors[role] || "";
  };

  return (
    <Card ref={sectionRef} className="bg-slate-800/30 border-slate-700">
      <CardHeader>
        <div className="flex items-center gap-3">
          <Badge
            variant="outline"
            className={`${getRoleBadgeColor(role)} border`}
          >
            {getRoleDisplayName(role)}
          </Badge>
          <CardTitle>Top Specializations - {selectedBracket}</CardTitle>
        </div>

        <CardDescription>
          Meta analysis for {role.toLowerCase()} specializations in{" "}
          {selectedBracket}
          {!showAll && filteredData.length >= 5 && (
            <span className="text-slate-300">
              {" "}
              • Showing top 5 specializations
            </span>
          )}
          {showAll && (
            <span className="text-slate-300">
              {" "}
              • Showing all {filteredData.length} specializations
            </span>
          )}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow className="border-slate-700 hover:bg-slate-800/50">
                <TableHead className="w-16">Rank</TableHead>
                <TableHead>Specialization</TableHead>
                <TableHead className="text-right">Usage</TableHead>
                <TableHead className="text-right">Percentage</TableHead>
                <TableHead className="text-right">Avg Score per Key</TableHead>
                <TableHead className="w-32">Popularity</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredData.map((spec) => (
                <TableRow
                  key={`${spec.class}-${spec.spec}-${spec.key_level_bracket}`}
                  className="border-slate-700 hover:bg-slate-800/30 transition-colors duration-150"
                >
                  {/* Rank */}
                  <TableCell className="font-medium">
                    <Badge
                      variant="outline"
                      className="bg-slate-700/50 text-slate-300 border-slate-600"
                    >
                      #{spec.rank}
                    </Badge>
                  </TableCell>

                  {/* Specialization */}
                  <TableCell>
                    <div className="space-y-1">
                      <ClassColoredText
                        className={spec.class}
                        additionalClasses="font-semibold text-base"
                      >
                        {spec.spec}
                      </ClassColoredText>
                      <div className="text-sm text-slate-400">{spec.class}</div>
                    </div>
                  </TableCell>

                  {/* Usage */}
                  <TableCell className="text-right font-mono">
                    <div className="text-base font-semibold text-white">
                      {spec.usage_count.toLocaleString("en-US")}
                    </div>
                    <div className="text-xs text-slate-400">uses</div>
                  </TableCell>

                  {/* Percentage */}
                  <TableCell className="text-right">
                    <div className="text-base font-bold text-white">
                      {spec.percentage.toFixed(1)}%
                    </div>
                  </TableCell>

                  {/* Average Score */}
                  <TableCell className="text-right font-mono">
                    <div className="text-base font-semibold text-white">
                      {spec.avg_score.toFixed(0)}
                    </div>
                    <div className="text-xs text-slate-400">score</div>
                  </TableCell>

                  {/* Popularity bar */}
                  <TableCell>
                    <div className="space-y-1">
                      <Progress
                        value={(spec.percentage / maxPercentage) * 100}
                        className="w-full h-2"
                      />
                      <div className="text-xs text-slate-400 text-center">
                        {spec.percentage < 1
                          ? "< 1%"
                          : `${spec.percentage.toFixed(0)}%`}
                      </div>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>

        {/* View More/Less Button - Positioned after table */}
        <div className="mt-6 flex justify-center">
          <Button
            variant="outline"
            size="default"
            onClick={toggleViewAll}
            className="flex items-center gap-2 bg-slate-700/50 border-slate-600 hover:bg-slate-600/50 text-slate-200 hover:text-white transition-all duration-200 px-6 py-2"
          >
            {showAll ? (
              <>
                <ChevronUp className="h-4 w-4" />
                Show Top 5 Only
              </>
            ) : (
              <>
                <ChevronDown className="h-4 w-4" />
                View All Specializations
              </>
            )}
          </Button>
        </div>

        {/* Summary statistics */}
        <div className="mt-6 p-4 bg-slate-800/50 rounded-lg border border-slate-700">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-center">
            <div>
              <div className="text-2xl font-bold text-white">
                {filteredData.length}
              </div>
              <div className="text-sm text-slate-400">
                {showAll
                  ? `total ${role.toLowerCase()} specs`
                  : `top ${role.toLowerCase()} specs`}
              </div>
            </div>
            <div>
              <div className="text-2xl font-bold text-white">
                {filteredData
                  .reduce((sum, spec) => sum + spec.usage_count, 0)
                  .toLocaleString("en-US")}
              </div>
              <div className="text-sm text-slate-400">
                {showAll ? "Total usage" : "Top 5 usage"}
              </div>
            </div>
            <div>
              <div className="text-2xl font-bold text-white">
                {filteredData[0]?.avg_score.toFixed(0) || "N/A"}
              </div>
              <div className="text-sm text-slate-400">
                Best avg score per key
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default MetaByKeyLevelsSection;
