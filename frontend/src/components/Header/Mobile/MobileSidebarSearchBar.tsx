// components/Header/Mobile/MobileSidebarSearchBar.tsx
"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Search as SearchIcon, Loader2 } from "lucide-react";
import { useSearchBlizzardCharacter } from "@/hooks/useBlizzardApi";
import { eu, us, tw, kr } from "@/data/realms";

interface Realm {
  id: number;
  name: string;
  slug: string;
}

interface MobileSidebarSearchBarProps {
  onClose: () => void;
}

const MobileSidebarSearchBar: React.FC<MobileSidebarSearchBarProps> = ({
  onClose,
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
      onClose(); // Close sidebar after navigation
    } catch (err) {
      setSearchError(
        err instanceof Error ? err.message : "An unexpected error occurred."
      );
    } finally {
      setIsSearching(false);
    }
  };

  return (
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
            isSearching || !searchRegion || !searchRealm || !searchCharacter
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
  );
};

export default MobileSidebarSearchBar;
