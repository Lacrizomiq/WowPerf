import React from "react";
import Image from "next/image";
import { TalentNode } from "@/types/talents";
import TalentGrid from "./TalentGrid";

interface SpecTalentsProps {
  talents: TalentNode[];
  specName: string;
  selectedTalents: TalentNode[];
  specIcon: string;
}

const SpecTalents: React.FC<SpecTalentsProps> = ({
  talents,
  specName,
  selectedTalents,
  specIcon,
}) => {
  return (
    <div className="border-2 border-[#001830] shadow-2xl w-full mb-6 rounded-lg">
      <h3 className="text-xl p-4 font-bold text-center text-white bg-deep-blue flex items-center justify-center">
        <Image
          src={specIcon}
          alt={specName}
          width={32}
          height={32}
          className="mr-2 rounded-full"
          unoptimized
        />
        <span>{specName} Talents</span>
      </h3>
      <div className="p-12">
        <TalentGrid talents={talents} selectedTalents={selectedTalents} />
      </div>
    </div>
  );
};

export default SpecTalents;
