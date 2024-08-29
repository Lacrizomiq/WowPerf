import React from "react";
import { TalentNode } from "@/types/talents";
import TalentGrid from "@/components/TalentTree/TalentGrid";

interface SpecTalentsProps {
  talents: TalentNode[];
  specName: string;
}

const SpecTalents: React.FC<SpecTalentsProps> = ({ talents, specName }) => {
  return (
    <div className="mb-8">
      <h3 className="text-xl font-bold mb-4 text-center">{specName} Talents</h3>
      <TalentGrid talents={talents} />
    </div>
  );
};

export default SpecTalents;
