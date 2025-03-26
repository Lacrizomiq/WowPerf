import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { Search, ChevronDown, ChevronUp, Loader2 } from "lucide-react";
import { useSearchBlizzardCharacter } from "@/hooks/useBlizzardApi";
import { eu, us, tw, kr } from "@/data/realms";
import { SidebarMenuButton, SidebarMenuSub } from "@/components/ui/sidebar";

interface Realm {
  id: number;
  name: string;
  slug: string;
}

interface SidebarSearchBarProps {
  isExpanded: boolean;
  searchOpen: boolean;
  setSearchOpen: (open: boolean) => void;
  onSuccessfulSearch?: () => void;
}

const MobileSidebarSearchBar: React.FC<SidebarSearchBarProps> = ({
  isExpanded,
  searchOpen,
  setSearchOpen,
  onSuccessfulSearch,
}) => {
  const router = useRouter();
  const [region, setRegion] = useState<string>("");
  const [realm, setRealm] = useState<string>("");
  const [character, setCharacter] = useState<string>("");
  const [realms, setRealms] = useState<Realm[]>([]);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const { refetch: checkCharacter } = useSearchBlizzardCharacter(
    region.toLowerCase(),
    realm.toLowerCase(),
    character.toLowerCase(),
    `profile-${region.toLowerCase()}`,
    "en_GB"
  );

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
    setRealms(selectedRealms.sort((a, b) => a.name.localeCompare(b.name)));
    setRealm("");
  }, [region]);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError(null);
    setIsSubmitting(true);

    if (!region || !realm || !character) {
      setError("Please fill in all fields");
      setIsSubmitting(false);
      return;
    }

    try {
      const result = await checkCharacter();
      if (result.error) {
        throw new Error("Character not found");
      }

      router.push(
        `/character/${region.toLowerCase()}/${realm.toLowerCase()}/${character.toLowerCase()}`
      );
      setSearchOpen(false);
      resetForm();
      if (onSuccessfulSearch) {
        onSuccessfulSearch();
      }
    } catch (err) {
      if (err instanceof Error) {
        setError("Character not found.");
      } else {
        setError("An error occurred while searching. Please try again.");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const resetForm = () => {
    setCharacter("");
    setRealm("");
    setRegion("");
    setError(null);
  };

  return (
    <>
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
              <option value="">Select Region</option>
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
              <option value="">Select Realm</option>
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

            {error && (
              <Alert variant="destructive" className="mb-2 text-red-500">
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <button
              type="submit"
              disabled={isSubmitting || !region || !realm || !character}
              className="w-full bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition duration-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
            >
              {isSubmitting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Searching...
                </>
              ) : (
                "Search"
              )}
            </button>
          </form>
        </SidebarMenuSub>
      )}
    </>
  );
};

export default MobileSidebarSearchBar;
