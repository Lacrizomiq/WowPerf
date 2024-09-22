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
    <div className="border-2 border-[#001830] shadow-2xl w-full mb-6 rounded-lg">
      <h3 className="text-xl p-4 font-bold text-center text-white bg-deep-blue flex items-center justify-center">
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
      <div className="p-12">
        <TalentGrid talents={talents} selectedTalents={selectedTalents} />
      </div>
    </div>
  );
};

export default ClassTalents;
