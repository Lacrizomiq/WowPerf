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
    <div className="mb-8">
      <h3 className="text-xl font-bold mb-8 text-center text-gradient-glow flex items-center justify-center">
        <Image
          src={specIcon}
          alt={specName}
          width={32}
          height={32}
          className="mr-2"
          unoptimized
        />
        <span>{specName} Talents</span>
      </h3>
      <div className="mb-6 shadow-2xl border-4 p-12" style={{ width: "100%" }}>
        <TalentGrid talents={talents} selectedTalents={selectedTalents} />
      </div>
    </div>
  );
};

export default SpecTalents;
