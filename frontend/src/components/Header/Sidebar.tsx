"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/providers/AuthContext";
import {
  Sidebar,
  SidebarHeader,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarMenuSub,
  SidebarMenuSubItem,
  SidebarMenuSubButton,
} from "@/components/ui/sidebar";
import {
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
import { eu, us, tw, kr } from "@/data/realms";

import UserMenuOverlay from "./UserMenuOverlay";

interface AppSidebarProps {
  isExpanded: boolean;
  setIsExpanded: (expanded: boolean) => void;
}

const AppSidebar: React.FC<AppSidebarProps> = ({
  isExpanded,
  setIsExpanded,
}) => {
  const router = useRouter();
  const { isAuthenticated, logout } = useAuth();
  const [mythicPlusExpanded, setMythicPlusExpanded] = useState(false);
  const [searchOpen, setSearchOpen] = useState(false);
  const [region, setRegion] = useState("");
  const [realm, setRealm] = useState("");
  const [character, setCharacter] = useState("");
  const [realms, setRealms] = useState<
    { id: number; name: string; slug: string }[]
  >([]);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (region && realm && character) {
      router.push(
        `/character/${region.toLowerCase()}/${realm.toLowerCase()}/${character.toLowerCase()}`
      );
      setSearchOpen(false);
    }
  };

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

  return (
    <Sidebar
      className={`
      transition-all duration-300 h-full 
      bg-[color:var(--sidebar-background)] 
      text-[color:var(--sidebar-foreground)]
      shadow-lg
      ${isExpanded ? "w-64" : "w-16"}
    `}
    >
      {/* Header */}
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <a href="/">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                  <Home className="size-4" />
                </div>
                {isExpanded && (
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-semibold">WoW Perf</span>
                  </div>
                )}
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      {/* Content */}

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem
                className={!isExpanded ? "flex justify-center mt-4" : "mt-4"}
              >
                <SidebarMenuButton onClick={() => setSearchOpen(!searchOpen)}>
                  <div
                    className={`flex items-center w-full mb-2 ${
                      isExpanded ? "w-full" : "justify-center"
                    }`}
                  >
                    <Search className={isExpanded ? "mr-4" : ""} />
                    {isExpanded && (
                      <>
                        <span>Search</span>
                        {searchOpen ? (
                          <ChevronUp className="ml-auto" />
                        ) : (
                          <ChevronDown className="ml-auto" />
                        )}
                      </>
                    )}
                  </div>
                </SidebarMenuButton>
                {isExpanded && searchOpen && (
                  <SidebarMenuSub>
                    <form onSubmit={handleSubmit} className="px-4 py-4">
                      <select
                        value={region}
                        onChange={(e) => setRegion(e.target.value)}
                        className="w-full px-2 py-2 mb-2 text-white border-2 rounded-md bg-deep-blue"
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
                        className="w-full px-2 py-2 mb-2 text-white border-2 rounded-md bg-deep-blue"
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
                        className="w-full px-2 py-2 mb-2 text-white border-2 rounded-md bg-deep-blue"
                      />
                      <button
                        type="submit"
                        className="w-full bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition duration-300"
                      >
                        Search
                      </button>
                    </form>
                  </SidebarMenuSub>
                )}
              </SidebarMenuItem>
              <SidebarMenuItem
                className={!isExpanded ? "flex justify-center mt-4" : "mt-4"}
              >
                <SidebarMenuButton
                  onClick={() => setMythicPlusExpanded(!mythicPlusExpanded)}
                >
                  <div
                    className={`flex items-center w-full mb-2 mt-2 ${
                      isExpanded ? "w-full" : "justify-center"
                    }`}
                  >
                    <Hourglass className={isExpanded ? "mr-4" : ""} />
                    {isExpanded && (
                      <>
                        <span>Mythic +</span>
                        {mythicPlusExpanded ? (
                          <ChevronUp className="ml-auto" />
                        ) : (
                          <ChevronDown className="ml-auto" />
                        )}
                      </>
                    )}
                  </div>
                </SidebarMenuButton>
                {isExpanded && mythicPlusExpanded && (
                  <SidebarMenuSub>
                    <SidebarMenuSubItem className="pt-4">
                      <SidebarMenuSubButton
                        onClick={() => router.push("/mythic-plus/leaderboard")}
                      >
                        <Rows4 />
                        <span className="cursor-pointer">Leaderboard</span>
                      </SidebarMenuSubButton>
                    </SidebarMenuSubItem>
                    <SidebarMenuSubItem className="py-4 ">
                      <SidebarMenuSubButton
                        onClick={() => router.push("/mythic-plus/best-runs")}
                      >
                        <BicepsFlexed />
                        <span className="cursor-pointer">Best Runs</span>
                      </SidebarMenuSubButton>
                    </SidebarMenuSubItem>
                    <SidebarMenuSubItem>
                      <SidebarMenuSubButton
                        onClick={() => router.push("/mythic-plus/statistics")}
                      >
                        <ChartColumnDecreasing />
                        <span className="cursor-pointer">Statistics</span>
                      </SidebarMenuSubButton>
                    </SidebarMenuSubItem>
                  </SidebarMenuSub>
                )}
              </SidebarMenuItem>
              <SidebarMenuItem
                className={!isExpanded ? "flex justify-center mt-4" : "mt-4"}
              >
                <SidebarMenuButton onClick={() => router.push("/raids")}>
                  <div
                    className={`flex items-center w-full mt-2 mb-2 ${
                      isExpanded ? "w-full" : "justify-center"
                    }`}
                  >
                    <Sword className={isExpanded ? "mr-4" : ""} />
                    {isExpanded && <span>Raids</span>}
                  </div>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      {/* Footer */}
      <SidebarFooter>
        <div className="h-16 bg-[color:var(--sidebar-background)] text-[color:var(--sidebar-foreground)] rounded-lg hover:bg-slate-800">
          <UserMenuOverlay isExpanded={isExpanded} />
        </div>
      </SidebarFooter>
    </Sidebar>
  );
};

export default AppSidebar;
