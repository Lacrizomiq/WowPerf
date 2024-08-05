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
      <div>
        {character.thumbnail_url && (
          <Image
            src={character.thumbnail_url}
            alt="World of Warcraft Logo"
            width={100}
            height={100}
          />
        )}
      </div>
      <div>
        <h2 className="text-3xl font-bold  text-gradient-glow">
          {character.name}
        </h2>

        <p className="text-blue-200">
          {character.region.toUpperCase()} - {character.realm}
        </p>
      </div>
    </div>
  );
}
