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
    <div className="border-4 shadow-2xl w-full">
      <h3 className="text-xl p-4 font-bold text-center text-white bg-black bg-opacity-70 flex items-center justify-center border-b-4">
        <Image
          src={heroTalentIcon}
          alt={HeroTalentName}
          width={40}
          height={40}
          className="mr-2 rounded-full"
          unoptimized
        />
        <span className="text-xl">{HeroTalentName} Hero Talents</span>
      </h3>
      <div className="p-12">
        <HeroTalentGrid selectedHeroTalentTree={selectedHeroTalentTree} />
      </div>
    </div>
  );
};

export default HeroSpecTalent;
