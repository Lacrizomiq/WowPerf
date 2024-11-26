// components/Sidebar/MobileSidebar.tsx
import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Sheet, SheetContent } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import {
  Menu,
  X,
  Home,
  Search,
  Sword,
  Hourglass,
  ChevronDown,
  ChevronUp,
  BicepsFlexed,
  ChartColumnDecreasing,
  Rows4,
} from "lucide-react";
import MobileUserMenu from "./MobileUserMenu";
import { eu, us, tw, kr } from "@/data/realms";

interface MobileSidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

const MobileSidebar: React.FC<MobileSidebarProps> = ({ isOpen, onClose }) => {
  const router = useRouter();
  const [searchOpen, setSearchOpen] = useState(false);
  const [mythicPlusExpanded, setMythicPlusExpanded] = useState(false);
  const [region, setRegion] = useState("");
  const [realm, setRealm] = useState("");
  const [character, setCharacter] = useState("");
  const [realms, setRealms] = useState<
    { id: number; name: string; slug: string }[]
  >([]);

  useEffect(() => {
    let selectedRealms: { id: number; name: string; slug: string }[] = [];
    switch (region) {
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
    setRealms(sortedRealms);
    setRealm("");
  }, [region]);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (region && realm && character) {
      router.push(
        `/character/${region.toLowerCase()}/${realm.toLowerCase()}/${character.toLowerCase()}`
      );
      onClose();
      setSearchOpen(false);
    }
  };

  return (
    <>
      <Sheet open={isOpen} onOpenChange={onClose}>
        <SheetContent
          side="left"
          className="p-0 w-80 bg-sidebar border-r border-sidebar-border overflow-y-auto"
        >
          <div className="h-full flex flex-col">
            {/* Header */}
            <div className="flex justify-between items-center p-4 border-b border-sidebar-border">
              <a href="/" className="flex items-center space-x-2">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                  <Home className="size-4" />
                </div>
                <span className="font-semibold">WoW Perf</span>
              </a>
            </div>

            {/* Content */}
            <div className="flex-1 px-4 py-2 ">
              {/* Search Section */}
              <div className="py-2">
                <button
                  onClick={() => setSearchOpen(!searchOpen)}
                  className="w-full flex items-center justify-between p-2 hover:bg-slate-800 rounded-md"
                >
                  <div className="flex items-center">
                    <Search className="mr-2 h-5 w-5" />
                    <span>Search</span>
                  </div>
                  {searchOpen ? <ChevronUp /> : <ChevronDown />}
                </button>

                {searchOpen && (
                  <form onSubmit={handleSubmit} className="mt-2 space-y-2 p-2">
                    <select
                      value={region}
                      onChange={(e) => setRegion(e.target.value)}
                      className="w-full px-3 py-2 bg-sidebar-secondary rounded-md border border-sidebar-border text-black"
                    >
                      <option value="" disabled>
                        Select Region
                      </option>
                      <option value="eu">EU</option>
                      <option value="us">US</option>
                      <option value="kr">KR</option>
                      <option value="tw">TW</option>
                    </select>

                    <select
                      value={realm}
                      onChange={(e) => setRealm(e.target.value)}
                      className="w-full px-3 py-2 bg-sidebar-secondary rounded-md border border-sidebar-border text-black"
                      disabled={!region}
                    >
                      <option value="" disabled>
                        Select Realm
                      </option>
                      {realms.map((realm) => (
                        <option key={realm.id} value={realm.slug}>
                          {realm.name}
                        </option>
                      ))}
                    </select>

                    <input
                      type="text"
                      placeholder="Character Name"
                      value={character}
                      onChange={(e) => setCharacter(e.target.value)}
                      className="w-full px-3 py-2 bg-sidebar-secondary rounded-md border border-sidebar-border text-black"
                    />

                    <Button
                      type="submit"
                      className="w-full bg-blue-600 hover:bg-blue-700 text-white"
                      disabled={!region || !realm || !character}
                    >
                      Search
                    </Button>
                  </form>
                )}
              </div>

              {/* Mythic+ Section */}
              <div className="py-2">
                <button
                  onClick={() => setMythicPlusExpanded(!mythicPlusExpanded)}
                  className="w-full flex items-center justify-between p-2 hover:bg-slate-800 rounded-md"
                >
                  <div className="flex items-center">
                    <Hourglass className="mr-2 h-5 w-5" />
                    <span>Mythic+</span>
                  </div>
                  {mythicPlusExpanded ? <ChevronUp /> : <ChevronDown />}
                </button>

                {mythicPlusExpanded && (
                  <div className="ml-4 mt-2 space-y-2">
                    <button
                      onClick={() => {
                        router.push("/mythic-plus/leaderboard");
                        onClose();
                      }}
                      className="w-full flex items-center p-2 hover:bg-slate-800 rounded-md"
                    >
                      <Rows4 className="mr-2 h-4 w-4" />
                      <span>Leaderboard</span>
                    </button>

                    <button
                      onClick={() => {
                        router.push("/mythic-plus/best-runs");
                        onClose();
                      }}
                      className="w-full flex items-center p-2 hover:bg-slate-800 rounded-md"
                    >
                      <BicepsFlexed className="mr-2 h-4 w-4" />
                      <span>Best Runs</span>
                    </button>

                    <button
                      onClick={() => {
                        router.push("/mythic-plus/statistics");
                        onClose();
                      }}
                      className="w-full flex items-center p-2 hover:bg-slate-800 rounded-md"
                    >
                      <ChartColumnDecreasing className="mr-2 h-4 w-4" />
                      <span>Statistics</span>
                    </button>
                  </div>
                )}
              </div>

              {/* Raids Section */}
              <div className="py-2">
                <button
                  onClick={() => {
                    router.push("/raids");
                    onClose();
                  }}
                  className="w-full flex items-center p-2 hover:bg-slate-800 rounded-md"
                >
                  <Sword className="mr-2 h-5 w-5" />
                  <span>Raids</span>
                </button>
              </div>
            </div>

            {/* Footer with User Menu */}
            <div className="mt-auto">
              <MobileUserMenu />
            </div>
          </div>
        </SheetContent>
      </Sheet>
    </>
  );
};

export default MobileSidebar;
