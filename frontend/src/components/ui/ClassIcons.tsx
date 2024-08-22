import React from "react";
import Image from "next/image";
import { useGetBlizzardCharacterSpecializations } from "@/hooks/useBlizzardApi";

interface ClassIconsProps {
  region: string;
  realm: string;
  name: string;
  namespace: string;
  locale: string;
}

const ClassIcons = ({
  region,
  realm,
  name,
  namespace,
  locale,
}: ClassIconsProps) => {
  const {
    data: characterData,
    isLoading,
    error,
  } = useGetBlizzardCharacterSpecializations(
    region,
    realm,
    name,
    namespace,
    locale
  );

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading class icon</div>;
  if (!characterData || !characterData.class) return null;

  const classNameForUrl = characterData.class.replace(/\s+/g, "");

  return (
    <Image
      width={26}
      height={26}
      src={`https://assets.rpglogs.com/img/warcraft/icons/${classNameForUrl}.jpg`}
      alt={characterData.class}
    />
  );
};

export default ClassIcons;
