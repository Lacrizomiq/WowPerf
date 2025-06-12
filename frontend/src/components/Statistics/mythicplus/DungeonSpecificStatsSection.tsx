// components/Statistics/mythicplus/DungeonSpecificStatsSection.tsx

import React, { useState, useRef, useMemo } from "react";
import Image from "next/image";
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
import {
  useSpecsByDungeonAndRole,
  useTopTeamCompositionsByDungeon,
} from "@/hooks/useMythicPlusRunsAnalysis";
import { useGetDungeonMedia } from "@/hooks/useWarcraftLogsApi";
import { ClassColoredText, getRoleStyle } from "../shared/ClassColorUtils";
import LoadingSpinner from "../shared/LoadingSpinner";
import ErrorDisplay from "../shared/ErrorDisplay";
import { Role } from "@/types/raiderio/mythicplus_runs/mythicPlusRuns";
import InfoTooltip from "@/components/Shared/InfoTooltip";

/**
 * Section qui affiche les statistiques spécifiques à un donjon
 * Combine spécialisations par rôle et compositions populaires
 */
const DungeonSpecificStatsSection: React.FC = () => {
  const [selectedDungeon, setSelectedDungeon] = useState("cinderbrew-meadery");
  const [selectedRole, setSelectedRole] = useState<Role>(Role.DPS);

  // Load compositions to extract available dungeons
  const { data: allCompositions } = useTopTeamCompositionsByDungeon({
    top_n: 0, // Get all to extract dungeon list
    min_usage: 3,
  });

  // Load dungeon media for icons
  const { data: dungeonMedia } = useGetDungeonMedia();

  // Extract unique dungeons from compositions data and merge with media
  const availableDungeons = useMemo(() => {
    if (!allCompositions) return [];

    const uniqueDungeons = Array.from(
      new Map(
        allCompositions.map((comp) => [
          comp.dungeon_slug,
          {
            slug: comp.dungeon_slug,
            name: comp.dungeon_name,
          },
        ])
      ).values()
    );

    // Merge with media data to get icons
    const dungeonsWithMedia = uniqueDungeons.map((dungeon) => {
      const media = dungeonMedia?.find((m) => m.dungeon_slug === dungeon.slug);
      return {
        ...dungeon,
        icon: media?.icon,
        media_url: media?.media_url,
      };
    });

    return dungeonsWithMedia.sort((a, b) => a.name.localeCompare(b.name));
  }, [allCompositions, dungeonMedia]);

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

  const selectedDungeonData = availableDungeons.find(
    (d) => d.slug === selectedDungeon
  );
  const selectedDungeonName = selectedDungeonData?.name || "Unknown Dungeon";

  return (
    <section>
      <h2 className="text-2xl font-bold mb-4 flex items-center">
        Dungeon-Specific Statistics
        <InfoTooltip
          content="This section shows the statistics for a specific dungeon. It includes the top specializations and compositions for each role."
          className="ml-2"
          size="lg"
        />
      </h2>

      {/* Dungeon Selector */}
      <div className="mb-6">
        <div className="flex items-center gap-4">
          <label className="text-sm font-medium text-slate-300">Dungeon:</label>
          <Select value={selectedDungeon} onValueChange={setSelectedDungeon}>
            <SelectTrigger className="w-80 bg-slate-800/50 border-slate-700 text-slate-200">
              <div className="flex items-center gap-3">
                {selectedDungeonData?.icon && (
                  <Image
                    src={`https://wow.zamimg.com/images/wow/icons/large/${selectedDungeonData.icon}.jpg`}
                    alt={selectedDungeonName}
                    width={24}
                    height={24}
                    unoptimized
                    className="rounded-sm"
                  />
                )}
                <span>{selectedDungeonName}</span>
              </div>
            </SelectTrigger>
            <SelectContent className="bg-slate-800 border-slate-700">
              {availableDungeons.map((dungeon) => (
                <SelectItem
                  key={dungeon.slug}
                  value={dungeon.slug}
                  className="text-slate-200 focus:bg-slate-700 focus:text-white"
                >
                  <div className="flex items-center gap-3">
                    {dungeon.icon && (
                      <Image
                        src={`https://wow.zamimg.com/images/wow/icons/large/${dungeon.icon}.jpg`}
                        alt={dungeon.name}
                        width={24}
                        height={24}
                        unoptimized
                        className="rounded-sm"
                      />
                    )}
                    <span>{dungeon.name}</span>
                  </div>
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Two Column Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Left Column: Specializations by Role */}
        <div>
          <h3 className="text-xl font-semibold mb-4">
            Top Specializations - {selectedDungeonName}
          </h3>

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
              <DungeonSpecTable
                dungeonSlug={selectedDungeon}
                role={Role.TANK}
              />
            </TabsContent>

            <TabsContent value={Role.HEALER}>
              <DungeonSpecTable
                dungeonSlug={selectedDungeon}
                role={Role.HEALER}
              />
            </TabsContent>

            <TabsContent value={Role.DPS}>
              <DungeonSpecTable dungeonSlug={selectedDungeon} role={Role.DPS} />
            </TabsContent>
          </Tabs>
        </div>

        {/* Right Column: Top Compositions */}
        <div>
          <h3 className="text-xl font-semibold mb-4">
            Top Compositions - {selectedDungeonName}
          </h3>
          <DungeonCompositionsCard dungeonSlug={selectedDungeon} />
        </div>
      </div>
    </section>
  );
};

/**
 * Table component for dungeon-specific specializations
 */
interface DungeonSpecTableProps {
  dungeonSlug: string;
  role: Role;
}

const DungeonSpecTable: React.FC<DungeonSpecTableProps> = ({
  dungeonSlug,
  role,
}) => {
  const [showAll, setShowAll] = useState(false);
  const sectionRef = useRef<HTMLDivElement>(null);

  // Determine top_n parameter based on showAll state
  const topN = showAll ? 0 : 5;

  const {
    data: specs,
    isLoading,
    error,
    isError,
  } = useSpecsByDungeonAndRole(dungeonSlug, role, { top_n: topN });

  // Toggle function for View More/Less with auto-scroll
  const toggleViewAll = () => {
    const newShowAll = !showAll;
    setShowAll(newShowAll);

    if (!newShowAll && sectionRef.current) {
      setTimeout(() => {
        sectionRef.current?.scrollIntoView({
          behavior: "smooth",
          block: "start",
        });
      }, 100);
    }
  };

  if (isLoading) {
    return (
      <Card className="bg-slate-800/30 border-slate-700">
        <CardContent className="py-8">
          <LoadingSpinner />
        </CardContent>
      </Card>
    );
  }

  if (isError || !specs) {
    return (
      <Card className="bg-slate-800/30 border-slate-700">
        <CardContent className="py-8">
          <ErrorDisplay
            error={error}
            message={`Unable to load ${role} specializations for this dungeon`}
          />
        </CardContent>
      </Card>
    );
  }

  if (specs.length === 0) {
    return (
      <Card className="bg-slate-800/30 border-slate-700">
        <CardContent className="py-8">
          <div className="text-center text-slate-400">
            No {role.toLowerCase()} data available for this dungeon
          </div>
        </CardContent>
      </Card>
    );
  }

  const maxPercentage = Math.max(...specs.map((spec) => spec.percentage));

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
          <CardTitle className="text-lg">
            {getRoleDisplayName(role)} Specializations
          </CardTitle>
        </div>

        <CardDescription>
          Usage statistics for {role.toLowerCase()} in this dungeon
          {!showAll && specs.length >= 5 && (
            <span className="text-slate-300">
              {" "}
              • Showing top 5 specializations
            </span>
          )}
          {showAll && (
            <span className="text-slate-300">
              {" "}
              • Showing all {specs.length} specializations
            </span>
          )}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow className="border-slate-700">
                <TableHead>Rank</TableHead>
                <TableHead>Specialization</TableHead>
                <TableHead className="text-right">%</TableHead>
                <TableHead className="w-24">Popularity</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {specs.map((spec) => (
                <TableRow
                  key={`${spec.class}-${spec.spec}`}
                  className="border-slate-700 hover:bg-slate-800/30 transition-colors duration-150"
                >
                  {/* Rank */}
                  <TableCell className="font-medium">
                    <Badge
                      variant="outline"
                      className="bg-slate-700/50 text-slate-300 border-slate-600"
                    >
                      #{spec.rank_in_dungeon}
                    </Badge>
                  </TableCell>

                  {/* Specialization */}
                  <TableCell>
                    <div className="space-y-1">
                      <ClassColoredText
                        className={spec.class}
                        additionalClasses="font-semibold"
                      >
                        {spec.spec}
                      </ClassColoredText>
                      <div className="text-xs text-slate-400">{spec.class}</div>
                    </div>
                  </TableCell>

                  {/* Percentage */}
                  <TableCell className="text-right">
                    <div className="text-base font-bold text-white">
                      {spec.percentage.toFixed(1)}%
                    </div>
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

        {/* View More/Less Button */}
        <div className="mt-4 flex justify-center">
          <Button
            variant="outline"
            size="sm"
            onClick={toggleViewAll}
            className="flex items-center gap-2 bg-slate-700/50 border-slate-600 hover:bg-slate-600/50 text-slate-200 hover:text-white transition-all duration-200 px-4 py-1"
          >
            {showAll ? (
              <>
                <ChevronUp className="h-3 w-3" />
                Top 5
              </>
            ) : (
              <>
                <ChevronDown className="h-3 w-3" />
                View All
              </>
            )}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

/**
 * Card component for dungeon-specific compositions
 */
interface DungeonCompositionsCardProps {
  dungeonSlug: string;
}

const DungeonCompositionsCard: React.FC<DungeonCompositionsCardProps> = ({
  dungeonSlug,
}) => {
  const [showAll, setShowAll] = useState(false);
  const sectionRef = useRef<HTMLDivElement>(null);

  // Determine limit based on showAll state
  const limit = showAll ? 15 : 5;

  const {
    data: allCompositions,
    isLoading,
    error,
    isError,
  } = useTopTeamCompositionsByDungeon({
    top_n: limit,
    min_usage: 3,
  });

  // Filter compositions for selected dungeon
  const compositions =
    allCompositions?.filter((comp) => comp.dungeon_slug === dungeonSlug) || [];

  // Toggle function for View More/Less with auto-scroll
  const toggleViewAll = () => {
    const newShowAll = !showAll;
    setShowAll(newShowAll);

    if (!newShowAll && sectionRef.current) {
      setTimeout(() => {
        sectionRef.current?.scrollIntoView({
          behavior: "smooth",
          block: "start",
        });
      }, 100);
    }
  };

  // Helper functions for extracting class/spec names
  const extractClassName = (display: string): string => {
    return display.split(" - ")[0] || display;
  };

  const extractSpecName = (display: string): string => {
    return display.split(" - ")[1] || display;
  };

  if (isLoading) {
    return (
      <Card className="bg-slate-800/30 border-slate-700">
        <CardContent className="py-8">
          <LoadingSpinner />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card className="bg-slate-800/30 border-slate-700">
        <CardContent className="py-8">
          <ErrorDisplay
            error={error}
            message="Unable to load compositions for this dungeon"
          />
        </CardContent>
      </Card>
    );
  }

  if (compositions.length === 0) {
    return (
      <Card className="bg-slate-800/30 border-slate-700">
        <CardContent className="py-8">
          <div className="text-center text-slate-400">
            No composition data available for this dungeon
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card ref={sectionRef} className="bg-slate-800/30 border-slate-700">
      <CardHeader>
        <CardTitle className="text-lg">Top Compositions</CardTitle>
        <CardDescription>
          Most popular team compositions in this dungeon
          {!showAll && compositions.length >= 5 && (
            <span className="text-slate-300">
              {" "}
              • Showing top 5 compositions
            </span>
          )}
          {showAll && (
            <span className="text-slate-300">
              {" "}
              • Showing {compositions.length} compositions
            </span>
          )}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {compositions.map((comp) => (
            <div
              key={`${comp.rank_in_dungeon}-${comp.tank}-${comp.healer}`}
              className="p-4 bg-slate-800/50 rounded-lg border border-slate-700 hover:bg-slate-800/70 transition-colors duration-150"
            >
              <div className="flex justify-between items-center mb-3">
                <Badge className="bg-purple-600 text-white">
                  #{comp.rank_in_dungeon}
                </Badge>
                <div className="text-sm text-slate-400">
                  {comp.percentage.toFixed(1)}% • Score:{" "}
                  {comp.avg_score.toFixed(0)}
                </div>
              </div>

              <div className="grid grid-cols-1 gap-2 text-sm">
                {/* Tank */}
                <div className="flex items-center gap-2">
                  <Badge
                    variant="outline"
                    className="bg-blue-500/20 text-blue-400 border-blue-500/50 text-xs"
                  >
                    Tank
                  </Badge>
                  <ClassColoredText
                    className={extractClassName(comp.tank)}
                    additionalClasses="font-medium"
                  >
                    {extractSpecName(comp.tank)}
                  </ClassColoredText>
                  <span className="text-slate-400">
                    ({extractClassName(comp.tank)})
                  </span>
                </div>

                {/* Healer */}
                <div className="flex items-center gap-2">
                  <Badge
                    variant="outline"
                    className="bg-green-500/20 text-green-400 border-green-500/50 text-xs"
                  >
                    Heal
                  </Badge>
                  <ClassColoredText
                    className={extractClassName(comp.healer)}
                    additionalClasses="font-medium"
                  >
                    {extractSpecName(comp.healer)}
                  </ClassColoredText>
                  <span className="text-slate-400">
                    ({extractClassName(comp.healer)})
                  </span>
                </div>

                {/* DPS */}
                <div className="flex items-center gap-2">
                  <Badge
                    variant="outline"
                    className="bg-red-500/20 text-red-400 border-red-500/50 text-xs"
                  >
                    DPS
                  </Badge>
                  <div className="flex flex-wrap gap-1">
                    {[comp.dps1, comp.dps2, comp.dps3].map((dps, index) => (
                      <span key={index} className="flex items-center">
                        <ClassColoredText
                          className={extractClassName(dps)}
                          additionalClasses="font-medium"
                        >
                          {extractSpecName(dps)}
                        </ClassColoredText>
                        {index < 2 && (
                          <span className="text-slate-500 mx-1">•</span>
                        )}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* View More/Less Button */}
        <div className="mt-4 flex justify-center">
          <Button
            variant="outline"
            size="sm"
            onClick={toggleViewAll}
            className="flex items-center gap-2 bg-slate-700/50 border-slate-600 hover:bg-slate-600/50 text-slate-200 hover:text-white transition-all duration-200 px-4 py-1"
          >
            {showAll ? (
              <>
                <ChevronUp className="h-3 w-3" />
                Top 5
              </>
            ) : (
              <>
                <ChevronDown className="h-3 w-3" />
                View More
              </>
            )}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

export default DungeonSpecificStatsSection;
