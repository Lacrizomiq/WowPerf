/* Not used anymore */

"use client";

import React, { useState, useEffect, useMemo } from "react";
import { useRouter } from "next/navigation";
import { eu, us, tw, kr } from "@/data/realms";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { Loader2, CircleX } from "lucide-react";
import { useSearchBlizzardCharacter } from "@/hooks/useBlizzardApi";

interface Realm {
  id: number;
  name: string;
  slug: string;
}

export default function SearchBar() {
  const [region, setRegion] = useState("");
  const [realm, setRealm] = useState("");
  const [character, setCharacter] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const { refetch: checkCharacter } = useSearchBlizzardCharacter(
    region.toLowerCase(),
    realm.toLowerCase(),
    character.toLowerCase(),
    `profile-${region.toLowerCase()}`,
    "en_GB"
  );

  const sortedRealms = useMemo(() => {
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
    return selectedRealms.sort((a, b) => a.name.localeCompare(b.name));
  }, [region]);

  useEffect(() => {
    setRealm("");
    setError(null);
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

  return (
    <div className="w-full max-w-3xl mx-auto space-y-4">
      {error && (
        <Alert
          variant="destructive"
          className="mb-2 text-red-500 flex items-center"
        >
          <AlertTitle>
            <CircleX className="size-4 text-red-500 mr-2" />
          </AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="flex space-x-4">
          <select
            value={region}
            onChange={(e) => setRegion(e.target.value)}
            className="w-1/2 px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400 bg-deep-blue text-white appearance-none cursor-pointer"
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
            className="w-1/2 px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400 bg-deep-blue text-white appearance-none cursor-pointer"
            disabled={!region}
          >
            <option value="">Select Realm</option>
            {sortedRealms.map((realm) => (
              <option key={realm.id} value={realm.slug}>
                {realm.name}
              </option>
            ))}
          </select>
        </div>
        <div className="flex space-x-4">
          <input
            type="text"
            placeholder="Character Name"
            value={character}
            onChange={(e) => setCharacter(e.target.value)}
            className="w-3/4 px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400 bg-deep-blue text-white"
          />
          <button
            type="submit"
            disabled={isSubmitting || !region || !realm || !character}
            className="w-1/4 bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition duration-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center glow-effect"
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
        </div>
      </form>
    </div>
  );
}
