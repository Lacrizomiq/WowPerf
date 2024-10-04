import React, { useState, useEffect } from "react";
import {
  Home,
  Search,
  Sword,
  Hourglass,
  LogIn,
  UserPlus,
  ChevronUp,
  ChevronDown,
} from "lucide-react";
import { useRouter } from "next/navigation";
import { eu, us, tw, kr } from "@/data/realms";

interface Realm {
  id: number;
  name: string;
  slug: string;
}

interface SidebarItemProps {
  icon: React.ElementType;
  label: string;
  isExpanded: boolean;
  onClick: () => void;
  route?: string;
}

const SidebarItem: React.FC<SidebarItemProps> = ({
  icon: Icon,
  label,
  isExpanded,
  onClick,
  route,
}) => {
  const router = useRouter();

  const handleClick = () => {
    if (route) {
      router.push(route);
    }
    onClick();
  };

  return (
    <div
      className={`flex items-center p-4 mt-1 hover:bg-blue-700 transition-all duration-300 cursor-pointer ${
        isExpanded ? "justify-start" : "justify-center"
      }`}
      onClick={handleClick}
    >
      <Icon size={24} />
      {isExpanded && <span className="ml-4">{label}</span>}
    </div>
  );
};

interface SidebarProps {
  setMainMargin: (margin: number) => void;
}

const Sidebar: React.FC<SidebarProps> = ({ setMainMargin }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [searchOpen, setSearchOpen] = useState(false);
  const [region, setRegion] = useState("");
  const [realm, setRealm] = useState("");
  const [character, setCharacter] = useState("");
  const [realms, setRealms] = useState<Realm[]>([]);
  const router = useRouter();
  const [mythicPlusExpanded, setMythicPlusExpanded] = useState(false);

  useEffect(() => {
    switch (region) {
      case "eu":
        setRealms(eu.realms);
        break;
      case "us":
        setRealms(us.realms);
        break;
      case "tw":
        setRealms(tw.realms);
        break;
      case "kr":
        setRealms(kr.realms);
        break;
      default:
        setRealms([]);
    }
    setRealm("");
  }, [region]);

  const toggleSidebar = () => {
    setIsExpanded(!isExpanded);
    setMainMargin(isExpanded ? 64 : 240);
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (region && realm && character) {
      router.push(`/character/${region}/${realm}/${character.toLowerCase()}`);
      setSearchOpen(false);
    }
  };

  return (
    <div
      className={`fixed left-0 top-0 h-full pt-2 bg-deep-blue text-white transition-all duration-300 ${
        isExpanded ? "w-60" : "w-16"
      }`}
      onMouseEnter={() => !isExpanded && toggleSidebar()}
      onMouseLeave={() => isExpanded && toggleSidebar()}
    >
      <div className="flex flex-col justify-between h-full">
        <div>
          <SidebarItem
            icon={Home}
            label="Home"
            isExpanded={isExpanded}
            onClick={toggleSidebar}
            route="/"
          />
          <SidebarItem
            icon={Search}
            label="Search"
            isExpanded={isExpanded}
            onClick={() => setSearchOpen(!searchOpen)}
          />
          {isExpanded && searchOpen && (
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
          )}
          <div className="flex flex-col justify-between">
            <SidebarItem
              icon={Hourglass}
              label="Mythic +"
              isExpanded={isExpanded}
              onClick={() => setMythicPlusExpanded(!mythicPlusExpanded)}
            />
            {isExpanded && (
              <span className="flex items-center justify-end">
                {mythicPlusExpanded ? (
                  <ChevronUp size={16} />
                ) : (
                  <ChevronDown size={16} />
                )}
              </span>
            )}
            {isExpanded && mythicPlusExpanded && (
              <div className="pl-8">
                <SidebarItem
                  icon={() => <></>}
                  label="Best Runs"
                  isExpanded={isExpanded}
                  onClick={() => router.push("/mythic-plus/best-runs")}
                />
                <SidebarItem
                  icon={() => <></>}
                  label="Statistics"
                  isExpanded={isExpanded}
                  onClick={() => router.push("/mythic-plus/statistics")}
                />
              </div>
            )}
          </div>
          <SidebarItem
            icon={Sword}
            label="Raids"
            isExpanded={isExpanded}
            onClick={toggleSidebar}
            route="/raids"
          />
        </div>
        <div>
          <SidebarItem
            icon={LogIn}
            label="Login"
            isExpanded={isExpanded}
            onClick={toggleSidebar}
          />
          <SidebarItem
            icon={UserPlus}
            label="Register"
            isExpanded={isExpanded}
            onClick={toggleSidebar}
          />
        </div>
      </div>
    </div>
  );
};

export default Sidebar;
