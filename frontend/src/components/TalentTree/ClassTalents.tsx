import React from "react";
import { TalentNode } from "@/types/talents";
import TalentGrid from "@/components/TalentTree/TalentGrid";

interface ClassTalentsProps {
  talents: TalentNode[];
  className: string;
}

const ClassTalents: React.FC<ClassTalentsProps> = ({ talents, className }) => {
  return (
    <div className="mb-8">
      <h3 className="text-xl font-bold mb-4 text-center">
        {className} Talents
      </h3>
      <TalentGrid talents={talents} />
    </div>
  );
};

export default ClassTalents;
