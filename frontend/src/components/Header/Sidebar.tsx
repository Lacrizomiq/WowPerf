// components/Header/Sidebar.tsx
"use client";

import React, { useState, useEffect } from "react";
import { useRouter, usePathname } from "next/navigation";
import { SidebarMenu, SidebarMenuItem } from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import {
  Home,
  BarChart3,
  PieChart,
  Trophy,
  ListChecks,
  LayoutDashboard,
  TrendingUp,
  Search as SearchIcon,
  ShieldQuestion,
  Loader2,
  AlertCircle as AlertCircleIcon,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import Link from "next/link";
import UserMenuOverlay from "./UserMenuOverlay";
import { Badge } from "@/components/ui/badge";
import { useSearchBlizzardCharacter } from "@/hooks/useBlizzardApi";
import { eu, us, tw, kr } from "@/data/realms";

interface Realm {
  id: number;
  name: string;
  slug: string;
}

interface AppSidebarProps {
  isExpanded: boolean;
  onToggleSidebar: () => void;
}

const AppSidebar: React.FC<AppSidebarProps> = ({
  isExpanded,
  onToggleSidebar,
}) => {
  const router = useRouter();
  const pathname = usePathname();

  const [searchRegion, setSearchRegion] = useState<string>("");
  const [searchRealm, setSearchRealm] = useState<string>("");
  const [searchCharacter, setSearchCharacter] = useState<string>("");
  const [availableRealms, setAvailableRealms] = useState<Realm[]>([]);
  const [filteredRealms, setFilteredRealms] = useState<Realm[]>([]);
  const [realmSearchValue, setRealmSearchValue] = useState<string>("");
  const [isSearching, setIsSearching] = useState<boolean>(false);
  const [searchError, setSearchError] = useState<string | null>(null);

  const { refetch: checkCharacter } = useSearchBlizzardCharacter(
    searchRegion.toLowerCase(),
    searchRealm.toLowerCase(),
    searchCharacter.toLowerCase(),
    `profile-${searchRegion.toLowerCase()}`,
    "en_GB"
  );

  useEffect(() => {
    let selectedRealms: Realm[] = [];
    switch (searchRegion) {
      case "eu":
        selectedRealms = eu.realms;
        break;
      case "us":
        selectedRealms = us.realms;
        break;
      case "tw":
        selectedRealms = tw.realms;
        break;
      case "kr":
        selectedRealms = kr.realms;
        break;
      default:
        selectedRealms = [];
    }
    const sortedRealms = selectedRealms.sort((a, b) =>
      a.name.localeCompare(b.name)
    );
    setAvailableRealms(sortedRealms);
    setFilteredRealms(sortedRealms);
    setSearchRealm("");
    setRealmSearchValue("");
  }, [searchRegion]);

  const handleSearchSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSearchError(null);
    if (!searchRegion || !searchRealm || !searchCharacter) {
      setSearchError("Please fill in all fields.");
      return;
    }
    setIsSearching(true);
    try {
      const result = await checkCharacter();
      if (result.isError || result.error) {
        throw new Error(
          (result.error as any)?.message || "Character not found or API error."
        );
      }
      router.push(
        `/character/${searchRegion.toLowerCase()}/${searchRealm.toLowerCase()}/${searchCharacter.toLowerCase()}`
      );
    } catch (err) {
      setSearchError(
        err instanceof Error ? err.message : "An unexpected error occurred."
      );
    } finally {
      setIsSearching(false);
    }
  };

  const isActive = (path: string) => pathname === path;

  // Styles améliorés pour les éléments de menu
  const getMenuItemClasses = (path: string) => {
    const active = isActive(path);
    return `w-full flex items-center py-2 px-3 rounded-md text-sm transition-colors duration-200 
    ${
      active
        ? "bg-purple-900/30 text-white"
        : "text-slate-200 hover:bg-slate-800/80 hover:text-white"
    }
    ${!isExpanded ? "justify-center" : ""}`;
  };

  // Badge "Soon" pour les fonctionnalités à venir
  const soonBadge = (
    <Badge
      variant="outline"
      className="ml-2 text-xs py-0 h-5 border-purple-600 text-purple-400"
    >
      Soon
    </Badge>
  );

  return (
    <div className="flex flex-col h-full bg-[#1A1D21] text-slate-200 border-r border-slate-800">
      {/* Header - Complètement différent selon l'état expanded/collapsed */}
      {isExpanded ? (
        <div className="p-4 flex items-center justify-between h-16 border-b border-slate-800">
          <Link href="/" className="flex items-center gap-2">
            <div className="bg-purple-600 text-white p-1 rounded flex-shrink-0">
              <BarChart3 className="h-5 w-5" />
            </div>
            <span className="font-bold text-lg">WoW Perf</span>
          </Link>
          <Button
            variant="ghost"
            size="icon"
            onClick={onToggleSidebar}
            className="bg-slate-800/80 text-white hover:bg-purple-700/80"
            aria-label="Collapse sidebar"
          >
            <ChevronLeft className="h-5 w-5" />
          </Button>
        </div>
      ) : (
        <div className="flex flex-col items-center py-3 border-b border-slate-800">
          <Button
            variant="ghost"
            size="icon"
            onClick={onToggleSidebar}
            className="bg-purple-600 text-white hover:bg-purple-700 mb-2"
            aria-label="Expand sidebar"
          >
            <ChevronRight className="h-5 w-5" />
          </Button>
          <Link href="/" className="bg-purple-600 text-white p-1 rounded">
            <BarChart3 className="h-5 w-5" />
          </Link>
        </div>
      )}

      {/* Content - Scrollable */}
      <div className="flex-1 overflow-y-auto">
        {/* User Profile Section */}
        <div className="px-4 py-2">
          <UserMenuOverlay isExpanded={isExpanded} />
        </div>

        {/* Search Section - Only show expanded form when sidebar is expanded */}
        <div className="px-4 py-2">
          {!isExpanded ? (
            <Button
              size="icon"
              className="w-full h-8 bg-purple-600 hover:bg-purple-700"
              onClick={onToggleSidebar}
            >
              <SearchIcon className="h-4 w-4" />
            </Button>
          ) : (
            <>
              <div className="text-sm font-medium mb-1">Search Character</div>
              <form onSubmit={handleSearchSubmit} className="space-y-1.5">
                <Input
                  placeholder="Character Name"
                  value={searchCharacter}
                  onChange={(e) => setSearchCharacter(e.target.value)}
                  className="h-8 bg-slate-800/50 border-slate-700 text-white hover:border-purple-500 focus:border-purple-600 focus:ring-1 focus:ring-purple-500 transition-colors"
                />

                <div className="grid grid-cols-2 gap-1.5">
                  <Select value={searchRegion} onValueChange={setSearchRegion}>
                    <SelectTrigger className="h-8 bg-slate-800/50 border-slate-700 text-white hover:border-purple-500 focus:border-purple-600 focus:ring-1 focus:ring-purple-500 transition-colors">
                      <SelectValue placeholder="Region" />
                    </SelectTrigger>
                    <SelectContent
                      className="bg-[#1e2025] border-slate-700 text-white p-0"
                      position="popper"
                    >
                      <div className="rounded-sm overflow-hidden">
                        <SelectItem
                          value="eu"
                          className="py-3 hover:bg-slate-700 focus:bg-slate-700"
                        >
                          EU
                        </SelectItem>
                        <SelectItem
                          value="us"
                          className="py-3 hover:bg-slate-700 focus:bg-slate-700"
                        >
                          US
                        </SelectItem>
                        <SelectItem
                          value="kr"
                          className="py-3 hover:bg-slate-700 focus:bg-slate-700"
                        >
                          KR
                        </SelectItem>
                        <SelectItem
                          value="tw"
                          className="py-3 hover:bg-slate-700 focus:bg-slate-700"
                        >
                          TW
                        </SelectItem>
                      </div>
                    </SelectContent>
                  </Select>

                  <Select
                    value={searchRealm}
                    onValueChange={setSearchRealm}
                    disabled={!searchRegion}
                  >
                    <SelectTrigger className="h-8 bg-slate-800/50 border-slate-700 text-white hover:border-purple-500 focus:border-purple-600 focus:ring-1 focus:ring-purple-500 transition-colors">
                      <SelectValue placeholder="Realm" />
                    </SelectTrigger>
                    <SelectContent
                      className="bg-[#1e2025] border-slate-700 text-white p-0"
                      position="popper"
                    >
                      <div
                        onClick={(e) => e.stopPropagation()}
                        onKeyDown={(e) => e.stopPropagation()}
                      >
                        <Input
                          placeholder="Search realms..."
                          className="bg-slate-800 border-slate-700 text-white mb-1 rounded-none focus:ring-0 focus:border-purple-500"
                          value={realmSearchValue}
                          onClick={(e) => e.stopPropagation()}
                          onKeyDown={(e) => e.stopPropagation()}
                          onChange={(e) => {
                            const searchValue = e.target.value.toLowerCase();
                            setRealmSearchValue(e.target.value);
                            setFilteredRealms(
                              availableRealms.filter(
                                (r) =>
                                  r.name.toLowerCase().includes(searchValue) ||
                                  r.slug.toLowerCase().includes(searchValue)
                              )
                            );
                          }}
                        />
                      </div>
                      <div className="max-h-[200px] overflow-y-auto">
                        {filteredRealms.map((r) => (
                          <SelectItem
                            key={r.id}
                            value={r.slug}
                            className="py-3 hover:bg-slate-700 focus:bg-slate-700"
                          >
                            {r.name}
                          </SelectItem>
                        ))}
                      </div>
                    </SelectContent>
                  </Select>
                </div>

                {searchError && (
                  <Alert
                    variant="destructive"
                    className="py-2 text-sm bg-red-900/50 border-red-800 text-red-200"
                  >
                    <AlertCircleIcon className="h-4 w-4" />
                    <AlertDescription className="ml-2">
                      {searchError}
                    </AlertDescription>
                  </Alert>
                )}

                <Button
                  type="submit"
                  className="w-full h-8 bg-purple-600 hover:bg-purple-700"
                  disabled={
                    isSearching ||
                    !searchRegion ||
                    !searchRealm ||
                    !searchCharacter
                  }
                >
                  {isSearching ? (
                    <Loader2 className="h-4 w-4 animate-spin mr-2" />
                  ) : (
                    <SearchIcon className="h-4 w-4 mr-2" />
                  )}
                  Search
                </Button>
              </form>
            </>
          )}
        </div>

        {/* Navigation Section */}
        <div className="px-4 py-3">
          {isExpanded && (
            <div className="text-xs font-medium uppercase tracking-wider text-slate-400 mb-2 px-3">
              Navigation
            </div>
          )}
          <SidebarMenu>
            <SidebarMenuItem>
              <Link href="/" className={getMenuItemClasses("/")}>
                {isExpanded ? (
                  <Home className="h-4 w-4 mr-3" />
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <Home className="w-5 h-5" />
                  </div>
                )}
                {isExpanded && <span>Home</span>}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link
                href="/performance-analysis"
                className={getMenuItemClasses("/performance-analysis")}
              >
                {isExpanded ? (
                  <BarChart3 className="h-4 w-4 mr-3" />
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <BarChart3 className="w-5 h-5" />
                  </div>
                )}
                {isExpanded && <span>Performance Analysis</span>}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link
                href="/statistics"
                className={getMenuItemClasses("/statistics")}
              >
                {isExpanded ? (
                  <PieChart className="h-4 w-4 mr-3" />
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <PieChart className="w-5 h-5" />
                  </div>
                )}
                {isExpanded && <span>Statistics</span>}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link
                href="/best-runs"
                className={getMenuItemClasses("/best-runs")}
              >
                {isExpanded ? (
                  <Trophy className="h-4 w-4 mr-3" />
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <Trophy className="w-5 h-5" />
                  </div>
                )}
                {isExpanded && <span>Best Runs</span>}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link href="/builds" className={getMenuItemClasses("/builds")}>
                {isExpanded ? (
                  <ListChecks className="h-4 w-4 mr-3" />
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <ListChecks className="w-5 h-5" />
                  </div>
                )}
                {isExpanded && <span>Builds</span>}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link
                href="/character-progress"
                className={getMenuItemClasses("/character-progress")}
              >
                {isExpanded ? (
                  <>
                    <TrendingUp className="h-4 w-4 mr-3" />
                    <div className="flex-1 flex justify-between items-center">
                      <span>Character Progress</span>
                      {soonBadge}
                    </div>
                  </>
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <TrendingUp className="w-5 h-5" />
                  </div>
                )}
              </Link>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <Link
                href="/dashboard"
                className={getMenuItemClasses("/dashboard")}
              >
                {isExpanded ? (
                  <>
                    <LayoutDashboard className="h-4 w-4 mr-3" />
                    <div className="flex-1 flex justify-between items-center">
                      <span>Dashboard</span>
                      {soonBadge}
                    </div>
                  </>
                ) : (
                  <div className="w-6 h-6 flex items-center justify-center">
                    <LayoutDashboard className="w-5 h-5" />
                  </div>
                )}
              </Link>
            </SidebarMenuItem>
          </SidebarMenu>
        </div>
      </div>
    </div>
  );
};

export default AppSidebar;
