import React from "react";
import Image from "next/image";
import { useGetRaiderIoCharacterTalents } from "@/hooks/useRaiderioApi";

interface SpecIconsProps {
  region: string;
  realm: string;
  name: string;
}

const SpecIcons = ({ region, realm, name }: SpecIconsProps) => {
  const {
    data: characterData,
    isLoading,
    error,
  } = useGetRaiderIoCharacterTalents(region, realm, name);
  return (
    <Image
      width="26"
      height="26"
      src={`https://assets.rpglogs.com/img/warcraft/icons/${characterData.class}-${characterData.active_spec_name}.jpg`}
      alt=""
    />
  );
};

export default SpecIcons;
