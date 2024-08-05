"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export default function SearchBar() {
  const [region, setRegion] = useState("");
  const [realm, setRealm] = useState("");
  const [character, setCharacter] = useState("");
  const router = useRouter();

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (region && realm && character) {
      router.push(`/character/${region}/${realm}/${character}`);
    }
    console.log("Searching for:", { region, realm, character });
  };

  return (
    <div className="py-8 bg-gradient-dark">
      <div className="container mx-auto px-4">
        <form
          onSubmit={handleSubmit}
          className="flex flex-col md:flex-row justify-center items-center space-y-4 md:space-y-0 md:space-x-4 "
        >
          <select
            value={region}
            onChange={(e) => setRegion(e.target.value)}
            className="w-full md:w-1/6 px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400 bg-deep-blue text-gray-600 appearance-none cursor-pointer"
          >
            <option value="" disabled selected>
              Select Region
            </option>
            <option value="us">US</option>
            <option value="eu">EU</option>
            <option value="kr">KR</option>
            <option value="tw">TW</option>
          </select>
          <input
            type="text"
            placeholder="Realm"
            value={realm}
            onChange={(e) => setRealm(e.target.value)}
            className="w-full md:w-1/4 px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400 bg-deep-blue text-gray-600"
          />
          <input
            type="text"
            placeholder="Character Name"
            value={character}
            onChange={(e) => setCharacter(e.target.value)}
            className="w-full md:w-1/3 px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400 bg-deep-blue text-gray-600"
          />
          <button
            type="submit"
            className="w-full md:w-auto bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition duration-300 glow-effect"
          >
            Search
          </button>
        </form>
      </div>
    </div>
  );
}
