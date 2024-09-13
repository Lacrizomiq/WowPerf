import React from "react";
import Image from "next/image";
import { TalentNode } from "@/types/talents";
import TalentGrid from "@/components/TalentTree/TalentGrid";

interface ClassTalentsProps {
  talents: TalentNode[];
  className: string;
  selectedTalents: TalentNode[];
  classIcon: string;
}

const ClassTalents: React.FC<ClassTalentsProps> = ({
  talents,
  className,
  selectedTalents,
  classIcon,
}) => {
  return (
    <div className="mb-8">
      <h3 className="text-xl font-bold mb-2 text-center text-gradient-glow flex items-center justify-center">
        <Image
          src={classIcon}
          alt={className}
          width={32}
          height={32}
          className="mr-2"
          unoptimized
        />
        <span>{className} Talents</span>
      </h3>
      <TalentGrid talents={talents} selectedTalents={selectedTalents} />
    </div>
  );
};

export default ClassTalents;
