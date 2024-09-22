"use client";

import React, { useState, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";
import { Search } from "lucide-react";
import { eu, us, tw, kr } from "@/data/realms";

interface Realm {
  id: number;
  name: string;
  slug: string;
}

const MiniSearchBar = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [region, setRegion] = useState("");
  const [realm, setRealm] = useState("");
  const [character, setCharacter] = useState("");
  const [realms, setRealms] = useState<Realm[]>([]);
  const router = useRouter();
  const dropdownRef = useRef<HTMLDivElement>(null);

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

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (region && realm && character) {
      const lowerCaseCharacter = character.toLowerCase();
      router.push(`/character/${region}/${realm}/${lowerCaseCharacter}`);
      setIsOpen(false);
    }
  };

  return (
    <div
      className="relative border-2 border-blue-400 rounded-xl p-2"
      ref={dropdownRef}
    >
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center text-white hover:text-blue-300 transition-colors"
      >
        <Search size={20} />
        <span className="ml-2">Search</span>
      </button>
      {isOpen && (
        <div className="absolute right-0 mt-2 w-64 bg-deep-blue rounded-lg shadow-lg p-4">
          <form onSubmit={handleSubmit} className="flex flex-col space-y-2">
            <select
              value={region}
              onChange={(e) => setRegion(e.target.value)}
              className="w-full px-2 py-1 rounded text-black "
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
              className="w-full px-2 py-1 rounded text-black"
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
              className="w-full px-2 py-1 rounded text-black"
            />
            <button
              type="submit"
              className="w-full bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition duration-300"
            >
              Search
            </button>
          </form>
        </div>
      )}
    </div>
  );
};

export default MiniSearchBar;
