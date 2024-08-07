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
  return (
    <Image
      width="26"
      height="26"
      src={`https://assets.rpglogs.com/img/warcraft/icons/${characterData.class}.jpg`}
      alt=""
    />
  );
};

export default ClassIcons;
