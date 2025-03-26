"use client";

import React, { useState } from "react";
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
  Sword,
  Hourglass,
  ChevronDown,
  ChevronUp,
  BicepsFlexed,
  ChartColumnDecreasing,
  Rows4,
} from "lucide-react";
import Link from "next/link";
import UserMenuOverlay from "./UserMenuOverlay";
import SidebarSearchBar from "./SidebarSearchBar";

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
              <Link href="/">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                  <Home className="size-4" />
                </div>
                {isExpanded && (
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-semibold">WoW Perf</span>
                  </div>
                )}
              </Link>
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
                <SidebarSearchBar
                  isExpanded={isExpanded}
                  searchOpen={searchOpen}
                  setSearchOpen={setSearchOpen}
                />
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
                        onClick={() => router.push("/mythic-plus/analysis")}
                      >
                        <Rows4 />
                        <span className="cursor-pointer">
                          Performance Analysis
                        </span>
                      </SidebarMenuSubButton>
                    </SidebarMenuSubItem>
                    <SidebarMenuSubItem className="py-4">
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
