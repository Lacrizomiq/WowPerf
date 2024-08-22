import React from "react";
import Image from "next/image";
import { useGetBlizzardCharacterSpecializations } from "@/hooks/useBlizzardApi";

interface SpecIconsProps {
  region: string;
  realm: string;
  name: string;
  namespace: string;
  locale: string;
}

const SpecIcons = ({
  region,
  realm,
  name,
  namespace,
  locale,
}: SpecIconsProps) => {
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
  if (error) return <div>Error loading spec icon</div>;
  if (!characterData || !characterData.class || !characterData.active_spec_name)
    return null;

  const classNameForUrl = characterData.class.replace(/\s+/g, "");
  const specNameForUrl = characterData.active_spec_name.replace(/\s+/g, "");

  return (
    <Image
      width={26}
      height={26}
      src={`https://assets.rpglogs.com/img/warcraft/icons/${classNameForUrl}-${specNameForUrl}.jpg`}
      alt={`${characterData.class} - ${characterData.active_spec_name}`}
    />
  );
};

export default SpecIcons;
