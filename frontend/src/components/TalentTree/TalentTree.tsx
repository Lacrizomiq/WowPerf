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

interface TalentEntry {
  spellId: number;
  // Ajoutez ici d'autres propriétés si nécessaire
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

  const findTreeTalent = (selectedTalent: TalentNode) => {
    return [...classTalents, ...specTalents].find(
      (t) =>
        t.id === selectedTalent.id ||
        (t.entries &&
          t.entries.some(
            (entry: TalentEntry) => entry.spellId === selectedTalent.id
          ))
    );
  };

  const combineTalents = (
    treeTalents: TalentNode[],
    selectedTalents: TalentNode[]
  ) => {
    const combined = [...treeTalents];
    selectedTalents.forEach((selectedTalent) => {
      const treeTalent = findTreeTalent(selectedTalent);
      if (treeTalent) {
        const index = combined.findIndex((t) => t.id === treeTalent.id);
        if (index !== -1) {
          combined[index] = { ...treeTalent, ...selectedTalent };
        }
      } else {
        combined.push(selectedTalent);
      }
    });
    return combined;
  };

  const combinedClassTalents = combineTalents(
    classTalents,
    selectedTalents.filter((t) => t.nodeType === "class")
  );
  const combinedSpecTalents = combineTalents(
    specTalents,
    selectedTalents.filter((t) => t.nodeType === "spec")
  );

  return (
    <div className="p-4 shadow-lg rounded-lg overflow-auto">
      <div className="flex flex-col space-y-2">
        <ClassTalents
          talents={combinedClassTalents}
          className={className}
          classIcon={talentData?.classIcon || ""}
          selectedTalents={selectedTalents.filter(
            (t) => t.nodeType === "class"
          )}
        />
        <SpecTalents
          talents={combinedSpecTalents}
          specName={specName}
          specIcon={talentData?.specIcon || ""}
          selectedTalents={selectedTalents.filter((t) => t.nodeType === "spec")}
        />
      </div>
    </div>
  );
};

export default TalentTree;
