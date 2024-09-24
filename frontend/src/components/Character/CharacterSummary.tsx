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

  if (isLoading)
    return <div className="text-center p-4">Loading character data...</div>;
  if (error)
    return (
      <div className="text-center p-4 text-red-500">
        Error loading character data: {error.message}
      </div>
    );
  if (!character)
    return <div className="text-center p-4">No character data found</div>;

  const backgroundStyle = {
    backgroundColor: "rgba(0, 36, 64, 0.9)",
  };

  return (
    <div
      className="rounded-xl p-5 flex items-center space-x-5 shadow-2xl"
      style={backgroundStyle}
    >
      <div className="relative">
        {character.avatar_url && (
          <Image
            src={character.avatar_url}
            alt={character.name}
            width={76}
            height={76}
            className={`rounded-full border-2 border-class-color--${character.tree_id}`}
          />
        )}
      </div>
      <div>
        <h1 className={`text-2xl font-bold class-color--${character.tree_id}`}>
          {character.name}
        </h1>
        <p className="text-gray-400">
          {region.toUpperCase()} - {character.realm}
        </p>
        <p className="text-gray-400">
          {character.race} {character.active_spec_name} {character.class}
        </p>
      </div>
    </div>
  );
}
