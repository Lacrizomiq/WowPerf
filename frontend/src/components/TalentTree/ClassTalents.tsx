import React from "react";
import { TalentNode } from "@/types/talents";
import TalentGrid from "./TalentGrid";

interface ClassTalentsProps {
  talents: TalentNode[];
  className: string;
  selectedTalents: TalentNode[];
}

const ClassTalents: React.FC<ClassTalentsProps> = ({
  talents,
  className,
  selectedTalents,
}) => {
  return (
    <div className="mb-8">
      <h3 className="text-xl font-bold mb-2 text-center">
        {className} Talents
      </h3>
      <TalentGrid talents={talents} selectedTalents={selectedTalents} />
    </div>
  );
};

export default ClassTalents;
