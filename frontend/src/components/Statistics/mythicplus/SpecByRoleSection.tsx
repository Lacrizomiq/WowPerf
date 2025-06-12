// components/Statistics/mythicplus/SpecByRoleSection.tsx

import React, { useState, useRef } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
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
import { useSpecsByRole } from "@/hooks/useMythicPlusRunsAnalysis";
import { ClassColoredText, getRoleStyle } from "../shared/ClassColorUtils";
import LoadingSpinner from "../shared/LoadingSpinner";
import ErrorDisplay from "../shared/ErrorDisplay";
import { Role } from "@/types/raiderio/mythicplus_runs/mythicPlusRuns";
import InfoTooltip from "@/components/Shared/InfoTooltip";

/**
 * Section that displays specialization usage by role
 * Uses tabs to organize Tank/Healer/DPS
 */
const SpecByRoleSection: React.FC = () => {
  const [selectedRole, setSelectedRole] = useState<Role>(Role.DPS);

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

  return (
    <section>
      <h2 className="text-2xl font-bold mb-4 flex items-center">
        Specialization Usage by Role
        <InfoTooltip
          content="This section shows the usage of specializations by role. It represents the global usage of specializations across all runs and dungeons."
          className="ml-2"
          size="lg"
        />
      </h2>

      <Tabs
        value={selectedRole}
        onValueChange={handleTabChange}
        className="w-full"
      >
        <TabsList className="grid w-full grid-cols-3 bg-slate-800/50 mb-4">
          <TabsTrigger
            value={Role.TANK}
            className="data-[state=active]:bg-indigo-600 hover:bg-slate-700 transition-colors duration-200"
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
          <RoleSpecTable role={Role.TANK} />
        </TabsContent>

        <TabsContent value={Role.HEALER}>
          <RoleSpecTable role={Role.HEALER} />
        </TabsContent>

        <TabsContent value={Role.DPS}>
          <RoleSpecTable role={Role.DPS} />
        </TabsContent>
      </Tabs>
    </section>
  );
};

export default SpecByRoleSection;

/**
 * Table component to display specializations for a role
 */
interface RoleSpecTableProps {
  role: Role;
}

const RoleSpecTable: React.FC<RoleSpecTableProps> = ({ role }) => {
  const [showAll, setShowAll] = useState(false);
  const sectionRef = useRef<HTMLDivElement>(null);

  // Determine top_n parameter based on showAll state
  const topN = showAll ? 0 : 5; // 0 = all specs, 5 = top 5 only

  const {
    data: specs,
    isLoading,
    error,
    isError,
  } = useSpecsByRole(role, {
    top_n: topN,
  });

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (isError || !specs) {
    return (
      <ErrorDisplay
        error={error}
        message={`Unable to load ${role} specializations`}
      />
    );
  }

  // Calculate maximum value for progress bars
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
          <CardTitle>
            Top Specializations - {getRoleDisplayName(role)}
          </CardTitle>
        </div>

        <CardDescription>
          Ranking based on usage in high-level Mythic+ runs
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
              <TableRow className="border-slate-700 hover:bg-slate-800/50">
                <TableHead className="w-16">Rank</TableHead>
                <TableHead>Specialization</TableHead>
                <TableHead className="text-right">Usage</TableHead>
                <TableHead className="text-right">Percentage</TableHead>
                <TableHead className="w-32">Popularity</TableHead>
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
                {specs.length}
              </div>
              <div className="text-sm text-slate-400">
                {showAll
                  ? `total ${role.toLowerCase()} specs`
                  : `top ${role.toLowerCase()} specs`}
              </div>
            </div>
            <div>
              <div className="text-2xl font-bold text-white">
                {specs
                  .reduce((sum, spec) => sum + spec.usage_count, 0)
                  .toLocaleString("en-US")}
              </div>
              <div className="text-sm text-slate-400">
                {showAll ? "Total usage" : "Top 5 usage"}
              </div>
            </div>
            <div>
              <div className="text-2xl font-bold text-white">
                {specs[0]?.percentage.toFixed(1)}%
              </div>
              <div className="text-sm text-slate-400">Most popular spec</div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
