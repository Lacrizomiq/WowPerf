// components/Sidebar/MobileSidebar.tsx
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { Sheet, SheetContent } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import {
  Menu,
  X,
  Home,
  Sword,
  Hourglass,
  ChevronDown,
  ChevronUp,
  BicepsFlexed,
  ChartColumnDecreasing,
  Rows4,
} from "lucide-react";
import MobileUserMenu from "./MobileUserMenu";
import Link from "next/link";
import MobileSidebarSearchBar from "./MobileSidebarSearchBar";

interface MobileSidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

const MobileSidebar: React.FC<MobileSidebarProps> = ({ isOpen, onClose }) => {
  const router = useRouter();
  const [searchOpen, setSearchOpen] = useState(false);
  const [mythicPlusExpanded, setMythicPlusExpanded] = useState(false);

  const handleSuccessfulSearch = () => {
    onClose();
    setSearchOpen(false);
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
              <Link href="/" className="flex items-center space-x-2">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                  <Home className="size-4" />
                </div>
                <span className="font-semibold">WoW Perf</span>
              </Link>
            </div>

            {/* Content */}
            <div className="flex-1 px-4 py-2">
              {/* Search Section */}
              <div className="py-2">
                <MobileSidebarSearchBar
                  isExpanded={true}
                  searchOpen={searchOpen}
                  setSearchOpen={setSearchOpen}
                  onSuccessfulSearch={handleSuccessfulSearch}
                />
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
