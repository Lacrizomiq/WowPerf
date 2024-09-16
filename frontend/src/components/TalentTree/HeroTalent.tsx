import React from "react";
import Image from "next/image";
import { HeroTalent } from "@/types/talents";
import HeroTalentGrid from "@/components/TalentTree/HeroTalentGrid";

interface HeroTalentProps {
  talents: HeroTalent[];
  selectedHeroTalentTree: HeroTalent[];
  HeroTalentName: string;
  heroTalentIcon: string;
}

const HeroSpecTalent: React.FC<HeroTalentProps> = ({
  talents,
  selectedHeroTalentTree,
  HeroTalentName,
  heroTalentIcon,
}) => {
  return (
    <div className="mb-6">
      <h3 className="text-lg font-semibold text-gradient-glow mb-4 items-center flex justify-center">
        <Image
          src={heroTalentIcon}
          alt="Hero Talents"
          width={40}
          height={40}
          className="mr-2"
          unoptimized
        />
        <span>{HeroTalentName} Hero Talents</span>
      </h3>
      <div className="mb-6 shadow-2xl border-4 p-12" style={{ width: "100%" }}>
        <HeroTalentGrid
          talents={talents}
          selectedHeroTalentTree={selectedHeroTalentTree}
        />
      </div>
    </div>
  );
};

export default HeroSpecTalent;
