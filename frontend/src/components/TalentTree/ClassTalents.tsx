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
    <div className="">
      <h3 className="text-xl font-bold mb-8 text-center text-gradient-glow flex items-center justify-center">
        <Image
          src={classIcon}
          alt={className}
          width={32}
          height={32}
          className="mr-2 rounded-full"
          unoptimized
        />
        <span>{className} Talents</span>
      </h3>
      <div className="mb-6 shadow-2xl border-4 p-12" style={{ width: "100%" }}>
        <TalentGrid talents={talents} selectedTalents={selectedTalents} />
      </div>
    </div>
  );
};

export default ClassTalents;
