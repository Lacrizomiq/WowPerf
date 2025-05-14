// components/Header/Mobile/MobileSidebar.tsx
"use client";

import React from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAuth } from "@/providers/AuthContext";
import { useUserProfile } from "@/hooks/useUserProfile";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Sheet, SheetContent } from "@/components/ui/sheet";
import {
  Home,
  BarChart3,
  PieChart,
  Trophy,
  ListChecks,
  LayoutDashboard,
  TrendingUp,
  Search as SearchIcon,
  X,
  LogIn,
  LogOut,
  UserPlus,
  Settings,
  ChevronRight,
} from "lucide-react";
import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import { eu, us, tw, kr } from "@/data/realms";
import { useState, useEffect } from "react";
import { useSearchBlizzardCharacter } from "@/hooks/useBlizzardApi";

interface MobileSidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

const MobileSidebar: React.FC<MobileSidebarProps> = ({ isOpen, onClose }) => {
  const router = useRouter();
  const pathname = usePathname();
  const { isAuthenticated, logout } = useAuth();
  const { profile } = useUserProfile();

  const [searchRegion, setSearchRegion] = useState<string>("");
  const [searchRealm, setSearchRealm] = useState<string>("");
  const [searchCharacter, setSearchCharacter] = useState<string>("");
  const [availableRealms, setAvailableRealms] = useState<any[]>([]);
  const [filteredRealms, setFilteredRealms] = useState<any[]>([]);
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
    let selectedRealms: any[] = [];
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
      onClose(); // Close sidebar after navigation
    } catch (err) {
      setSearchError(
        err instanceof Error ? err.message : "An unexpected error occurred."
      );
    } finally {
      setIsSearching(false);
    }
  };

  const handleNavigation = (path: string) => {
    router.push(path);
    onClose();
  };

  const isActive = (path: string) => pathname === path;

  const getMenuItemClasses = (path: string) => {
    const active = isActive(path);
    return `w-full flex items-center justify-between py-3 px-4 rounded-md text-base
    ${
      active
        ? "bg-purple-900/30 text-white font-medium"
        : "text-slate-200 hover:bg-slate-800 hover:text-white"
    }`;
  };

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent
        side="left"
        className="w-full max-w-[300px] p-0 bg-[#1A1D21] border-r border-slate-800 text-slate-200"
        onCloseAutoFocus={(e) => e.preventDefault()}
        onOpenAutoFocus={(e) => e.preventDefault()}
      >
        <div className="flex items-center justify-between p-4 border-b border-slate-800">
          <div className="flex items-center gap-2">
            <div className="bg-purple-600 text-white p-1 rounded">
              <BarChart3 className="h-5 w-5" />
            </div>
            <div className="font-bold text-lg text-white">WoW Perf</div>
          </div>
        </div>

        <div className="overflow-y-auto h-full">
          {/* User Profile */}
          <div className="p-4 border-b border-slate-800">
            <div className="flex items-center">
              <Avatar className="h-10 w-10 bg-slate-700">
                <AvatarImage src={""} alt={profile?.username || "User"} />
                <AvatarFallback className="bg-slate-700 text-white">
                  {profile?.username
                    ? profile.username.charAt(0).toUpperCase()
                    : "U"}
                </AvatarFallback>
              </Avatar>
              <div className="ml-3">
                <p className="text-base font-medium">
                  {isAuthenticated ? profile?.username : "Guest"}
                </p>
                <p className="text-sm text-slate-400">
                  {isAuthenticated ? profile?.email : "Not logged in"}
                </p>
              </div>
            </div>

            <div className="mt-3 space-y-2">
              {isAuthenticated ? (
                <>
                  <Button
                    variant="outline"
                    className="w-full justify-start border-slate-700 hover:bg-slate-800"
                    onClick={() => handleNavigation("/profile")}
                  >
                    <Settings className="h-4 w-4 mr-2" />
                    Account
                  </Button>
                  <Button
                    variant="outline"
                    className="w-full justify-start border-slate-700 text-red-400 hover:text-red-300 hover:bg-slate-800"
                    onClick={() => {
                      logout();
                      onClose();
                    }}
                  >
                    <LogOut className="h-4 w-4 mr-2" />
                    Log out
                  </Button>
                </>
              ) : (
                <>
                  <Button
                    variant="outline"
                    className="w-full justify-start border-slate-700 hover:bg-slate-800"
                    onClick={() => handleNavigation("/login")}
                  >
                    <LogIn className="h-4 w-4 mr-2" />
                    Login
                  </Button>
                  <Button
                    variant="outline"
                    className="w-full justify-start border-slate-700 hover:bg-slate-800"
                    onClick={() => handleNavigation("/signup")}
                  >
                    <UserPlus className="h-4 w-4 mr-2" />
                    Register
                  </Button>
                </>
              )}
            </div>
          </div>

          {/* Search Character */}
          <div className="p-4 border-b border-slate-800">
            <h3 className="text-base font-medium mb-3">Search Character</h3>
            <form onSubmit={handleSearchSubmit} className="space-y-3">
              <Input
                placeholder="Character Name"
                value={searchCharacter}
                onChange={(e) => setSearchCharacter(e.target.value)}
                className="bg-slate-800 border-slate-700 text-white hover:border-purple-500 focus:border-purple-600 focus:ring-1 focus:ring-purple-500 transition-colors"
              />

              <div className="grid grid-cols-2 gap-2">
                <Select value={searchRegion} onValueChange={setSearchRegion}>
                  <SelectTrigger className="bg-slate-800 border-slate-700 text-white hover:border-purple-500 focus:border-purple-600 focus:ring-1 focus:ring-purple-500 transition-colors">
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
                  <SelectTrigger className="bg-slate-800 border-slate-700 text-white hover:border-purple-500 focus:border-purple-600 focus:ring-1 focus:ring-purple-500 transition-colors">
                    <SelectValue placeholder="Realm" />
                  </SelectTrigger>
                  <SelectContent
                    className="bg-[#1e2025] border-slate-700 text-white p-0"
                    position="popper"
                  >
                    <Input
                      placeholder="Search realms..."
                      className="bg-slate-800 border-slate-700 text-white mb-1 rounded-none focus:ring-0 focus:border-purple-500"
                      value={realmSearchValue}
                      onChange={(e) => {
                        const searchValue = e.target.value.toLowerCase();
                        setRealmSearchValue(e.target.value);
                        const filtered = availableRealms.filter(
                          (r) =>
                            r.name.toLowerCase().includes(searchValue) ||
                            r.slug.toLowerCase().includes(searchValue)
                        );
                        setFilteredRealms(filtered);
                      }}
                    />
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
                <div className="p-2 text-sm bg-red-900/50 border border-red-800 text-red-200 rounded-md">
                  {searchError}
                </div>
              )}

              <Button
                type="submit"
                className="w-full bg-purple-600 hover:bg-purple-700"
                disabled={
                  isSearching ||
                  !searchRegion ||
                  !searchRealm ||
                  !searchCharacter
                }
              >
                {isSearching ? (
                  <span className="flex items-center">
                    <svg
                      className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="4"
                      ></circle>
                      <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                      ></path>
                    </svg>
                    Searching...
                  </span>
                ) : (
                  <span className="flex items-center">
                    <SearchIcon className="mr-2 h-4 w-4" />
                    Search
                  </span>
                )}
              </Button>
            </form>
          </div>

          {/* Navigation */}
          <div className="p-4">
            <h3 className="text-sm font-medium uppercase tracking-wider text-slate-400 mb-2 px-1">
              Navigation
            </h3>
            <nav className="space-y-1">
              <button
                onClick={() => handleNavigation("/")}
                className={getMenuItemClasses("/")}
              >
                <span className="flex items-center">
                  <Home className="h-5 w-5 mr-3" />
                  Home
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/performance-analysis")}
                className={getMenuItemClasses("/performance-analysis")}
              >
                <span className="flex items-center">
                  <BarChart3 className="h-5 w-5 mr-3" />
                  Performance Analysis
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/statistics")}
                className={getMenuItemClasses("/statistics")}
              >
                <span className="flex items-center">
                  <PieChart className="h-5 w-5 mr-3" />
                  Statistics
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/best-runs")}
                className={getMenuItemClasses("/best-runs")}
              >
                <span className="flex items-center">
                  <Trophy className="h-5 w-5 mr-3" />
                  Best Runs
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/builds")}
                className={getMenuItemClasses("/builds")}
              >
                <span className="flex items-center">
                  <ListChecks className="h-5 w-5 mr-3" />
                  Builds
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/character-progress")}
                className={getMenuItemClasses("/character-progress")}
              >
                <span className="flex items-center">
                  <TrendingUp className="h-5 w-5 mr-3" />
                  Character Progress
                  <Badge
                    variant="outline"
                    className="ml-2 text-xs border-purple-600 text-purple-400"
                  >
                    Soon
                  </Badge>
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>

              <button
                onClick={() => handleNavigation("/dashboard")}
                className={getMenuItemClasses("/dashboard")}
              >
                <span className="flex items-center">
                  <LayoutDashboard className="h-5 w-5 mr-3" />
                  Dashboard
                  <Badge
                    variant="outline"
                    className="ml-2 text-xs border-purple-600 text-purple-400"
                  >
                    Soon
                  </Badge>
                </span>
                <ChevronRight className="h-4 w-4 text-slate-400" />
              </button>
            </nav>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
};

export default MobileSidebar;
