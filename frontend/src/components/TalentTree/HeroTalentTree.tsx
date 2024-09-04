import React from "react";
import { HeroTalent } from "@/types/talents";
import HeroSpecTalents from "./HeroTalent";
import { useGetBlizzardTalentTree } from "@/hooks/useBlizzardApi";

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

  console.log("Talent Data:", talentData);

  if (isLoading) return <div>Loading talents...</div>;
  if (error)
    return <div>Error loading talents: {(error as Error).message}</div>;

  const heroTalents = talentData?.heroNodes || [];
  const heroTalentsName = talentData.subTreeNodes[0]?.entries[0]?.name;
  const heroTalentsIcon =
    talentData.subTreeNodes[0]?.entries[0]?.atlasMemberName;
  const iconUrl = heroTalentsIcon
    ? `https://wow.zamimg.com/images/wow/TextureAtlas/live/${heroTalentsIcon}.webp`
    : "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg";

  return (
    <div className="talent-tree p-4 bg-gradient-dark shadow-lg rounded-lg overflow-auto">
      <div className="flex flex-col space-y-2">
        <HeroSpecTalents
          talents={heroTalents}
          selectedHeroTalentTree={selectedHeroTalentTree.filter(
            (r) => r.rank > 0
          )}
          HeroTalentName={heroTalentsName || ""}
          heroTalentIcon={iconUrl || ""}
        />
      </div>
    </div>
  );
};

export default HeroTalentTree;
