import React from "react";
import Image from "next/image";
import { useGetBlizzardCharacterProfile } from "@/hooks/useBlizzardApi";

interface CharacterSummaryProps {
  region: string;
  realm: string;
  name: string;
  namespace: string;
  locale: string;
}

export default function CharacterSummary({
  region,
  realm,
  name,
  namespace,
  locale,
}: CharacterSummaryProps) {
  const {
    data: character,
    isLoading,
    error,
  } = useGetBlizzardCharacterProfile(region, realm, name, namespace, locale);

  if (isLoading) return <div>Loading character data...</div>;
  if (error) return <div>Error loading character data: {error.message}</div>;
  if (!character) return <div>No character data found</div>;

  return (
    <div>
      <div className="flex p-4 bg-gradient-dark">
        <div className="mr-4  bg-deep-blue bg-opacity-50 rounded-lg overflow-hidden shadow-lg glow-effect">
          {character.inset_avatar_url && (
            <Image
              src={character.inset_avatar_url}
              alt="World of Warcraft Logo"
              width={196}
              height={196}
              objectFit="cover"
            />
          )}
        </div>
        <div>
          <h2
            className={`text-4xl font-bold class-color--${character.tree_id}`}
          >
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
