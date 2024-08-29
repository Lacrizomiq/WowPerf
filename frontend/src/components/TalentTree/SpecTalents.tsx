import React from "react";
import { TalentNode } from "@/types/talents";
import TalentGrid from "@/components/TalentTree/TalentGrid";

interface SpecTalentsProps {
  talents: TalentNode[];
  specName: string;
  selectedTalents: TalentNode[];
}

const SpecTalents: React.FC<SpecTalentsProps> = ({
  talents,
  specName,
  selectedTalents,
}) => {
  return (
    <div className="mb-8">
      <h3 className="text-xl font-bold mb-2 text-center">{specName} Talents</h3>
      <TalentGrid talents={talents} selectedTalents={selectedTalents} />
    </div>
  );
};

export default SpecTalents;
