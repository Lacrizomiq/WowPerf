import React from "react";
import Image from "next/image";
import { useGetBlizzardCharacterProfile } from "@/hooks/useBlizzardApi";
import { useGetPlayerRankings } from "@/hooks/useWarcraftLogsApi";
import MythicPlusRanking from "@/components/MythicPlus/CharacterPersonalRanking/Summary/MythicPlusSummary";
import RaidRanking from "@/components/MythicPlus/CharacterPersonalRanking/Summary/RaidSummary";

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

  const {
    data: mythicPlusPlayerRankings,
    isLoading: isLoadingMythicPlusPlayerRankings,
    error: mythicPlusPlayerRankingsError,
  } = useGetPlayerRankings(name, realm, region, 39);

  const {
    data: raidPlayerRankings,
    isLoading: isLoadingRaidPlayerRankings,
    error: raidPlayerRankingsError,
  } = useGetPlayerRankings(name, realm, region, 38);

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
    backgroundSize: "cover",
    backgroundPosition: "top",
  };

  const defaultBackgroundClass = "bg-deep-blue";

  const fallbackMythicPlusImg =
    "https://wow.zamimg.com/images/wow/icons/large/ability_racial_chillofnight.jpg";

  return (
    <div
      className={`p-5 flex items-center bg-deep-blue space-x-5 shadow-2xl rounded-2xl ${
        character?.spec_id
          ? `bg-spec-${character.spec_id}`
          : defaultBackgroundClass
      }`}
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
      <div>
        <MythicPlusRanking
          seasonName="Mythic+ Season 1"
          rank={mythicPlusPlayerRankings?.zoneRankings.allStars[0].rank}
          classId={mythicPlusPlayerRankings?.classID || 0}
          spec={mythicPlusPlayerRankings?.zoneRankings.allStars[0].spec || ""}
          fallbackImageUrl={fallbackMythicPlusImg}
          isLoading={isLoadingMythicPlusPlayerRankings}
        />
      </div>
      <div>
        <RaidRanking
          raidName="Nerubar Palace"
          rank={raidPlayerRankings?.zoneRankings.allStars[0].rank}
          classId={raidPlayerRankings?.classID || 0}
          spec={raidPlayerRankings?.zoneRankings.allStars[0].spec || ""}
          isLoading={isLoadingRaidPlayerRankings}
        />
      </div>
    </div>
  );
}
