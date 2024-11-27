import React from "react";
import Image from "next/image";
import { useGetBlizzardCharacterProfile } from "@/hooks/useBlizzardApi";
import {
  useGetPlayerRaidRankings,
  useGetPlayerMythicPlusRankings,
} from "@/hooks/useWarcraftLogsApi";
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
  } = useGetPlayerMythicPlusRankings(name, realm, region, 39);

  const {
    data: raidPlayerRankings,
    isLoading: isLoadingRaidPlayerRankings,
    error: raidPlayerRankingsError,
  } = useGetPlayerRaidRankings(name, realm, region, 38);

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

  // Check if allStars data is available
  const allStarsMythicPlusData =
    mythicPlusPlayerRankings?.zoneRankings?.allStars?.[0];
  const allStarsRaidData = raidPlayerRankings?.zoneRankings?.allStars?.[0];
  const isDataAvailable = !!allStarsMythicPlusData || !!allStarsRaidData;

  const backgroundStyle = {
    backgroundSize: "cover",
    backgroundPosition: "top",
  };

  const defaultBackgroundClass = "bg-deep-blue";

  const fallbackMythicPlusImg =
    "https://wow.zamimg.com/images/wow/icons/large/ability_racial_chillofnight.jpg";

  const is500Error = (error: any) => {
    return error?.response?.status === 500 || error?.status === 500;
  };

  const shouldShowMythicPlusRanking =
    !is500Error(mythicPlusPlayerRankingsError) &&
    !isLoadingMythicPlusPlayerRankings &&
    allStarsMythicPlusData?.rank !== undefined &&
    allStarsMythicPlusData?.rank !== null;

  const shouldShowRaidRanking =
    !is500Error(raidPlayerRankingsError) &&
    !isLoadingRaidPlayerRankings &&
    allStarsRaidData?.rank !== undefined &&
    allStarsRaidData?.rank !== null;

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
      {shouldShowMythicPlusRanking && allStarsMythicPlusData && (
        <div>
          <MythicPlusRanking
            seasonName="TWW M+ S1"
            rank={allStarsMythicPlusData?.rank}
            classId={mythicPlusPlayerRankings?.classID || 0}
            spec={allStarsMythicPlusData?.spec || ""}
            fallbackImageUrl={fallbackMythicPlusImg}
            isLoading={isLoadingMythicPlusPlayerRankings}
          />
        </div>
      )}
      {shouldShowRaidRanking && allStarsRaidData && (
        <div>
          <RaidRanking
            raidName="Nerubar Palace"
            rank={allStarsRaidData?.rank}
            classId={raidPlayerRankings?.classID || 0}
            spec={allStarsRaidData?.spec || ""}
            isLoading={isLoadingRaidPlayerRankings}
          />
        </div>
      )}
    </div>
  );
}
