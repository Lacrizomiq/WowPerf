"use client";

import React, { useState, useEffect, useMemo } from "react";
import { useRouter } from "next/navigation";
import { eu, us, tw, kr } from "@/data/realms";

interface Realm {
  id: number;
  name: string;
  slug: string;
}

export default function SearchBar() {
  const [region, setRegion] = useState("");
  const [realm, setRealm] = useState("");
  const [character, setCharacter] = useState("");
  const router = useRouter();

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
    setRealm(""); // Reset realm when region changes
  }, [region]);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (region && realm && character) {
      const lowerCaseCharacter = character.toLowerCase();
      router.push(`/character/${region}/${realm}/${lowerCaseCharacter}`);
    }
  };

  return (
    <div className="w-full max-w-3xl mx-auto">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="flex space-x-4">
          <select
            value={region}
            onChange={(e) => setRegion(e.target.value)}
            className="w-1/2 px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400 bg-deep-blue text-white appearance-none cursor-pointer"
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
            className="w-1/2 px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400 bg-deep-blue text-white appearance-none cursor-pointer"
            disabled={!region}
          >
            <option value="" disabled>
              Select Realm
            </option>
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
            className="w-1/4 bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition duration-300 glow-effect"
          >
            Search
          </button>
        </div>
      </form>
    </div>
  );
}
