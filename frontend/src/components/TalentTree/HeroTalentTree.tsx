import React from "react";
import { HeroTalent } from "@/types/talents";
import HeroTalentGrid from "@/components/TalentTree/HeroTalentGrid";
import { useGetBlizzardTalentTree } from "@/hooks/useBlizzardApi";
import Image from "next/image";

interface TalentTreeProps {
  talentTreeId: number;
  specId: number;
  region: string;
  namespace: string;
  locale: string;
  className: string;
  specName: string;
  selectedHeroTalentTree: HeroTalent[];
}

const HeroTalentTree: React.FC<TalentTreeProps> = ({
  talentTreeId,
  specId,
  region,
  namespace,
  locale,
  className,
  specName,
  selectedHeroTalentTree,
}) => {
  const {
    data: talentData,
    isLoading,
    error,
  } = useGetBlizzardTalentTree(talentTreeId, specId, region, namespace, locale);

  if (isLoading) return <div>Loading talents...</div>;
  if (error)
    return <div>Error loading talents: {(error as Error).message}</div>;

  const heroTalents = talentData?.heroNodes || [];
  const heroTalentsName = talentData.subTreeNodes[0]?.entries[0]?.name || "";
  const heroTalentsIcon =
    talentData.subTreeNodes[0]?.entries[0]?.atlasMemberName || "";
  const iconUrl = heroTalentsIcon
    ? `https://wow.zamimg.com/images/wow/TextureAtlas/live/${heroTalentsIcon}.webp`
    : "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg";

  return (
    <div className="p-4 shadow-lg rounded-lg overflow-auto">
      <div className="flex flex-col border-2 border-black shadow-2xl rounded-lg overflow-hidden">
        <h3 className="text-lg font-semibold text-white bg-black bg-opacity-70 p-4 items-center flex justify-center">
          <Image
            src={iconUrl}
            alt="Hero Talents"
            width={40}
            height={40}
            className="mr-2"
            unoptimized
          />
          <span>{heroTalentsName} Hero Talents</span>
        </h3>
        <div className="px-40 py-20">
          <HeroTalentGrid
            selectedHeroTalentTree={selectedHeroTalentTree.filter(
              (r) => r.rank > 0
            )}
          />
        </div>
      </div>
    </div>
  );
};

export default HeroTalentTree;
