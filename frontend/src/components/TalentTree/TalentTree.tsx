import React from "react";
import { TalentNode } from "@/types/talents";
import ClassTalents from "./ClassTalents";
import SpecTalents from "./SpecTalent";
import HeroSpecTalents from "./HeroTalent";
import { useGetBlizzardTalentTree } from "@/hooks/useBlizzardApi";

interface TalentTreeProps {
  talentTreeId: number;
  specId: number;
  region: string;
  namespace: string;
  locale: string;
  className: string;
  specName: string;
  selectedTalents: TalentNode[];
}

const TalentTree: React.FC<TalentTreeProps> = ({
  talentTreeId,
  specId,
  region,
  namespace,
  locale,
  className,
  specName,
  selectedTalents,
}) => {
  const {
    data: talentData,
    isLoading,
    error,
  } = useGetBlizzardTalentTree(talentTreeId, specId, region, namespace, locale);

  if (isLoading) return <div>Loading talents...</div>;
  if (error)
    return <div>Error loading talents: {(error as Error).message}</div>;

  const classTalents = talentData?.classNodes || [];
  const specTalents = talentData?.specNodes || [];

  // Filter out hero talents and ensure we only have class and spec talents
  const filteredClassTalents = classTalents.filter(
    (talent: TalentNode) => talent.nodeType === "class"
  );
  const filteredSpecTalents = specTalents.filter(
    (talent: TalentNode) => talent.nodeType === "spec"
  );

  console.log("Filtered Class Talents:", filteredClassTalents);
  console.log("Filtered Spec Talents:", filteredSpecTalents);

  return (
    <div className="p-4 shadow-lg rounded-lg overflow-auto">
      <div className="flex flex-col space-y-2">
        <ClassTalents
          talents={filteredClassTalents}
          className={className}
          selectedTalents={selectedTalents.filter(
            (t) => t.nodeType === "class"
          )}
          classIcon={talentData?.classIcon || ""}
        />
        <SpecTalents
          talents={filteredSpecTalents}
          specName={specName}
          selectedTalents={selectedTalents.filter((t) => t.nodeType === "spec")}
          specIcon={talentData?.specIcon || ""}
        />
      </div>
    </div>
  );
};

export default TalentTree;
