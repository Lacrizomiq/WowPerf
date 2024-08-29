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
}

const TalentTree: React.FC<TalentTreeProps> = ({
  talentTreeId,
  specId,
  region,
  namespace,
  locale,
  className,
  specName,
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

  return (
    <div className="talent-tree p-4 bg-gradient-dark shadow-lg rounded-lg overflow-auto">
      <h2 className="text-2xl font-bold mb-4 text-center">
        Talent Build Summary
      </h2>
      <div className="flex flex-col space-y-8">
        <ClassTalents talents={classTalents} className={className} />
        <SpecTalents talents={specTalents} specName={specName} />
      </div>
    </div>
  );
};

export default TalentTree;
