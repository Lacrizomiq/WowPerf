import React from "react";
import Image from "next/image";
import { useRaiderIoCharacterProfile } from "@/hooks/useRaiderioApi";

interface CharacterSummaryProps {
  region: string;
  realm: string;
  name: string;
}

export default function CharacterSummary({
  region,
  realm,
  name,
}: CharacterSummaryProps) {
  const {
    data: character,
    isLoading,
    error,
  } = useRaiderIoCharacterProfile(region, realm, name);

  if (isLoading) return <div>Loading character data...</div>;
  if (error) return <div>Error loading character data: {error.message}</div>;
  if (!character) return <div>No character data found</div>;

  return (
    <div>
      <div className="flex p-4 bg-gradient-dark">
        <div className="mr-4 w-24 h-24 bg-deep-blue bg-opacity-50 rounded-lg overflow-hidden shadow-lg glow-effect">
          {character.thumbnail_url && (
            <Image
              src={character.thumbnail_url}
              alt="World of Warcraft Logo"
              width={96}
              height={96}
              objectFit="cover"
            />
          )}
        </div>
        <div>
          <h2 className="text-3xl font-bold  text-gradient-glow">
            {character.name}
          </h2>
          <p className="text-blue-100">
            {character.region.toUpperCase()} - {character.realm}
          </p>
          <div className="flex gap-2">
            <p className="text-blue-400">{character.race}</p>
            <p>
              {character.active_spec_name} {character.class}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
