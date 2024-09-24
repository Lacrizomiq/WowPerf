import React, { useState, useEffect } from "react";
import { Home, Search, Book, Settings, LogIn, UserPlus } from "lucide-react";
import { useRouter } from "next/navigation";
import { eu, us, tw, kr } from "@/data/realms";

const SidebarItem = ({ icon: Icon, label, isExpanded, onClick }) => (
  <div
    className={`flex items-center p-4 mt-1 hover:bg-blue-700 transition-all duration-300 cursor-pointer ${
      isExpanded ? "justify-start" : "justify-center"
    }`}
    onClick={onClick}
  >
    <Icon size={24} />
    {isExpanded && <span className="ml-4">{label}</span>}
  </div>
);

const Sidebar = ({ setMainMargin }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [searchOpen, setSearchOpen] = useState(false);
  const [region, setRegion] = useState("");
  const [realm, setRealm] = useState("");
  const [character, setCharacter] = useState("");
  const [realms, setRealms] = useState([]);
  const router = useRouter();

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

  const handleSubmit = (e) => {
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
      <SidebarItem
        icon={Home}
        label="Home"
        isExpanded={isExpanded}
        onClick={toggleSidebar}
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
            className="w-full px-2 py-1 mb-2 rounded text-black"
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
            className="w-full px-2 py-1 mb-2 rounded text-black"
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
            className="w-full px-2 py-1 mb-2 rounded text-black"
          />
          <button
            type="submit"
            className="w-full bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition duration-300"
          >
            Search
          </button>
        </form>
      )}
      <SidebarItem
        icon={Book}
        label="Guides"
        isExpanded={isExpanded}
        onClick={toggleSidebar}
      />
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
  );
};

export default Sidebar;
