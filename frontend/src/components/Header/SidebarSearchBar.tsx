// components/SidebarSearchBar.tsx
"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Search as SearchIcon,
  Loader2,
  AlertCircle as AlertCircleIcon,
} from "lucide-react";
import { useSearchBlizzardCharacter } from "@/hooks/useBlizzardApi";
import { eu, us, tw, kr } from "@/data/realms";

interface Realm {
  id: number;
  name: string;
  slug: string;
}

interface SidebarSearchBarProps {
  isExpanded: boolean;
  onToggleSidebar: () => void;
}

const SidebarSearchBar: React.FC<SidebarSearchBarProps> = ({
  isExpanded,
  onToggleSidebar,
}) => {
  const router = useRouter();
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

  // Render differently based on expanded state
  if (!isExpanded) {
    return (
      <div className="px-4 py-2">
        <Button
          size="icon"
          className="w-full h-8 bg-purple-600 hover:bg-purple-700"
          onClick={onToggleSidebar}
        >
          <SearchIcon className="h-4 w-4" />
        </Button>
      </div>
    );
  }

  return (
    <div className="px-4 py-2">
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
            <AlertDescription className="ml-2">{searchError}</AlertDescription>
          </Alert>
        )}

        <Button
          type="submit"
          className="w-full h-8 bg-purple-600 hover:bg-purple-700"
          disabled={
            isSearching || !searchRegion || !searchRealm || !searchCharacter
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
    </div>
  );
};

export default SidebarSearchBar;
