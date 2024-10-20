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
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { Separator } from "@/components/ui/separator";
import {
  Home,
  Search,
  Sword,
  Hourglass,
  LogIn,
  LogOut,
  UserPlus,
  ChevronDown,
  ChevronUp,
  ChartColumnDecreasing,
  BicepsFlexed,
  User,
  LayoutDashboard,
} from "lucide-react";
import { eu, us, tw, kr } from "@/data/realms";

interface Realm {
  id: number;
  name: string;
  slug: string;
}

const AppSidebar: React.FC = () => {
  const [searchOpen, setSearchOpen] = useState(false);
  const [region, setRegion] = useState("");
  const [realm, setRealm] = useState("");
  const [character, setCharacter] = useState("");
  const [realms, setRealms] = useState<Realm[]>([]);
  const [mythicPlusExpanded, setMythicPlusExpanded] = useState(false);
  const [userDropdownExpanded, setUserDropdownExpanded] = useState(false);

  const router = useRouter();
  const { isAuthenticated, logout } = useAuth();

  useEffect(() => {
    let selectedRealms: Realm[] = [];
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
      router.push(`/character/${region}/${realm}/${character.toLowerCase()}`);
      setSearchOpen(false);
    }
  };

  const handleLogout = async () => {
    await logout();
    router.push("/");
  };

  return (
    <Sidebar>
      <SidebarHeader className="p-4">
        <SidebarMenuButton asChild>
          <a href="/" className="flex items-center space-x-2">
            <Home />
            <span className="font-bold">WoWPerf</span>
          </a>
        </SidebarMenuButton>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton onClick={() => setSearchOpen(!searchOpen)}>
                  <Search className="mr-2" />
                  <span>Search</span>
                  {searchOpen ? (
                    <ChevronUp className="ml-auto" />
                  ) : (
                    <ChevronDown className="ml-auto" />
                  )}
                </SidebarMenuButton>
                {searchOpen && (
                  <SidebarMenuSub>
                    <form onSubmit={handleSubmit} className="px-4 py-2">
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
              <SidebarMenuItem>
                <SidebarMenuButton
                  onClick={() => setMythicPlusExpanded(!mythicPlusExpanded)}
                >
                  <Hourglass className="mr-2" />
                  <span>Mythic +</span>
                  {mythicPlusExpanded ? (
                    <ChevronUp className="ml-auto" />
                  ) : (
                    <ChevronDown className="ml-auto" />
                  )}
                </SidebarMenuButton>
                {mythicPlusExpanded && (
                  <SidebarMenuSub>
                    <SidebarMenuSubItem>
                      <SidebarMenuSubButton
                        onClick={() => router.push("/mythic-plus/best-runs")}
                      >
                        <BicepsFlexed className="mr-2" />
                        <span>Best Runs</span>
                      </SidebarMenuSubButton>
                    </SidebarMenuSubItem>
                    <SidebarMenuSubItem>
                      <SidebarMenuSubButton
                        onClick={() => router.push("/mythic-plus/statistics")}
                      >
                        <ChartColumnDecreasing className="mr-2" />
                        <span>Statistics</span>
                      </SidebarMenuSubButton>
                    </SidebarMenuSubItem>
                  </SidebarMenuSub>
                )}
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton onClick={() => router.push("/raids")}>
                  <Sword className="mr-2" />
                  <span>Raids</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          {isAuthenticated ? (
            <SidebarMenuItem>
              <SidebarMenuButton
                onClick={() => setUserDropdownExpanded(!userDropdownExpanded)}
              >
                <User className="mr-2" />
                <span>User</span>
                {userDropdownExpanded ? (
                  <ChevronUp className="ml-auto" />
                ) : (
                  <ChevronDown className="ml-auto" />
                )}
              </SidebarMenuButton>
              {userDropdownExpanded && (
                <SidebarMenuSub>
                  <SidebarMenuSubItem>
                    <SidebarMenuSubButton
                      onClick={() => router.push("/dashboard")}
                    >
                      <LayoutDashboard className="mr-2" />
                      <span>Dashboard</span>
                    </SidebarMenuSubButton>
                  </SidebarMenuSubItem>
                  <SidebarMenuSubItem>
                    <SidebarMenuSubButton
                      onClick={() => router.push("/profile")}
                    >
                      <User className="mr-2" />
                      <span>Profile</span>
                    </SidebarMenuSubButton>
                  </SidebarMenuSubItem>
                  <SidebarMenuSubItem>
                    <SidebarMenuSubButton onClick={handleLogout}>
                      <LogOut className="mr-2" />
                      <span>Logout</span>
                    </SidebarMenuSubButton>
                  </SidebarMenuSubItem>
                </SidebarMenuSub>
              )}
            </SidebarMenuItem>
          ) : (
            <>
              <SidebarMenuItem>
                <SidebarMenuButton onClick={() => router.push("/login")}>
                  <LogIn className="mr-2" />
                  <span>Login</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton onClick={() => router.push("/signup")}>
                  <UserPlus className="mr-2" />
                  <span>Register</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </>
          )}
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
};

export default AppSidebar;
