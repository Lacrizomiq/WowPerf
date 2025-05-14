// components/SidebarSearchBar.tsx
"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert"; // shadcn/ui Alert
import { Input } from "@/components/ui/input"; // shadcn/ui Input
import { Button } from "@/components/ui/button"; // shadcn/ui Button
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"; // shadcn/ui Select
import {
  Search,
  ChevronDown,
  ChevronUp,
  Loader2,
  AlertCircle,
} from "lucide-react"; // AlertCircle for the error icon
import { useSearchBlizzardCharacter } from "@/hooks/useBlizzardApi";
import { eu, us, tw, kr } from "@/data/realms";
import {
  SidebarMenuItem, // Imported to use as a wrapper if needed, or removed if not
  SidebarMenuButton,
  SidebarMenuSub,
} from "@/components/ui/sidebar"; // Your existing sidebar components

interface Realm {
  id: number;
  name: string;
  slug: string;
}

interface SidebarSearchBarProps {
  isExpanded: boolean;
}

// Note : searchOpen et setSearchOpen are removed from the props here,
// because AppSidebar will handle the opening of this panel directly.
const SidebarSearchBar: React.FC<SidebarSearchBarProps> = ({ isExpanded }) => {
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

    if (!region || !realm || !character) {
      setError("Please fill in all fields.");
      return;
    }
    setIsSubmitting(true);
    try {
      const result = await checkCharacter();
      // The error checking of react-query is often via result.isError or result.error
      if (result.isError || result.error) {
        // Adapt this to the return structure of your hook
        throw new Error(
          result.error?.message || "Character not found or API error."
        );
      }
      // If checkCharacter() throws an exception in case of error, the catch will handle it.
      // Otherwise, explicitly check result.data or a success indicator.

      router.push(
        `/character/${region.toLowerCase()}/${realm.toLowerCase()}/${character.toLowerCase()}`
      );
      // The closing of the panel (formerly setSearchOpen(false)) will be handled by AppSidebar
      resetForm();
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message || "Character not found.");
      } else {
        setError("An unexpected error occurred. Please try again.");
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

  // The button to open/close the search will be in AppSidebar.
  // This component will only render the form itself.
  // It will be displayed conditionally by AppSidebar.

  return (
    <div className="p-3 space-y-3">
      {" "}
      {/* Padding and spacing for the form */}
      <form onSubmit={handleSubmit} className="space-y-3">
        <Select value={region} onValueChange={(value) => setRegion(value)}>
          <SelectTrigger className="w-full bg-[color:var(--input)]">
            <SelectValue placeholder="Select Region" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="eu">EU</SelectItem>
            <SelectItem value="us">US</SelectItem>
            <SelectItem value="kr">KR</SelectItem>
            <SelectItem value="tw">TW</SelectItem>
          </SelectContent>
        </Select>

        <Select
          value={realm}
          onValueChange={(value) => setRealm(value)}
          disabled={!region}
        >
          <SelectTrigger className="w-full bg-[color:var(--input)]">
            <SelectValue placeholder="Select Realm" />
          </SelectTrigger>
          <SelectContent>
            {realms.map((r) => (
              <SelectItem key={r.id} value={r.slug}>
                {r.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Input
          type="text"
          placeholder="Character Name"
          value={character}
          onChange={(e) => setCharacter(e.target.value)}
          className="bg-[color:var(--input)]"
        />

        {error && (
          <Alert
            variant="destructive"
            className="bg-[color:var(--destructive)] text-[color:var(--destructive-foreground)]"
          >
            <AlertCircle className="h-4 w-4" />{" "}
            {/* Icon for the destructive alert of shadcn/ui */}
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <Button
          type="submit"
          disabled={isSubmitting || !region || !realm || !character}
          className="w-full bg-primary text-primary-foreground hover:bg-primary/90" // Main button style
        >
          {isSubmitting ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Searching...
            </>
          ) : (
            <>
              <Search className="mr-2 h-4 w-4" />
              Search
            </>
          )}
        </Button>
      </form>
    </div>
  );
};

export default SidebarSearchBar;
