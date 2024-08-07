import React from "react";
import Image from "next/image";
import { useGetRaiderIoCharacterTalents } from "@/hooks/useRaiderioApi";

interface ClassIconsProps {
  region: string;
  realm: string;
  name: string;
}

const ClassIcons = ({ region, realm, name }: ClassIconsProps) => {
  const {
    data: characterData,
    isLoading,
    error,
  } = useGetRaiderIoCharacterTalents(region, realm, name);

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading class icon</div>;
  if (!characterData) return null;

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
