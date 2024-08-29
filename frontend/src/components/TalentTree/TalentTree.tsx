import React from "react";
import { TalentNode } from "@/types/talents";
import ClassTalents from "./ClassTalents";
import SpecTalents from "./SpecTalents";
import { useGetBlizzardTalentTree } from "@/hooks/useBlizzardApi";

interface TalentTreeProps {
  talentTreeId: number;
  specId: number;
  region: string;
  namespace: string;
  locale: string;
  className: string;
  specName: string;
  selectedTalents: TalentNode[]; // Add this prop
}

const TalentTree: React.FC<TalentTreeProps> = ({
  talentTreeId,
  specId,
  region,
  namespace,
  locale,
  className,
  specName,
  selectedTalents, // Add this prop
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
    <div className="talent-tree p-4 bg-gradient-dark shadow-lg rounded-lg overflow-auto">
      <h2 className="text-2xl font-bold mb-4 text-center">
        {className} - {specName} Talents
      </h2>
      <div className="flex flex-col space-y-2">
        <ClassTalents
          talents={filteredClassTalents}
          className={className}
          selectedTalents={selectedTalents.filter(
            (t) => t.nodeType === "class"
          )}
        />
        <SpecTalents
          talents={filteredSpecTalents}
          specName={specName}
          selectedTalents={selectedTalents.filter((t) => t.nodeType === "spec")}
        />
      </div>
    </div>
  );
};

export default TalentTree;
