import React, { useEffect } from "react";
import { MythicPlusRun, Roster } from "@/types/runsDetails";
import { useGetRaiderioMythicPlusRunDetails } from "@/hooks/useRaiderioApi";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import RunsDetailsGear from "./RunsDetailsGear";

interface RunsDetailsProps {
  season: string;
  runId: number;
}

const RunsDetails: React.FC<RunsDetailsProps> = ({ season, runId }) => {
  const {
    data: runDetails,
    isLoading,
    error,
  } = useGetRaiderioMythicPlusRunDetails(season, runId);

  useWowheadTooltips();

  useEffect(() => {
    if (runDetails && window.$WowheadPower) {
      window.$WowheadPower.refreshLinks();
    }
  }, [runDetails]);

  if (isLoading)
    return <div className="text-white text-center p-4">Loading details...</div>;
  if (error)
    return (
      <div className="text-red-500 text-center p-4">
        Error: {(error as Error).message}
      </div>
    );
  if (!runDetails)
    return (
      <div className="text-yellow-500 text-center p-4">
        No details available
      </div>
    );

  return (
    <div className="p-4 bg-deep-blue-light rounded-lg">
      <h3 className="text-2xl font-bold mb-4 text-white">Run Details</h3>

      <div className="flex flex-col gap-4">
        {runDetails.roster.map((member: Roster) => (
          <div
            key={member.character.id}
            className="bg-deep-blue-lighter rounded-lg p-4 bg-black bg-opacity-30"
          >
            <div className="flex items-center mb-2 justify-between">
              <div>
                <h5
                  className={`font-bold class-color--${member.character.class.slug}`}
                >
                  {member.character.name}
                </h5>
                <p className="text-white text-sm">
                  {member.character.spec.name} {member.character.class.name}
                </p>
              </div>
              <div className="flex flex-col gap-2">
                <p className="text-white">
                  Item Level: {member.items.item_level_equipped}
                </p>
                <p className="text-white">Score: {member.ranks.score}</p>
              </div>
            </div>
            <details className="mt-2">
              <summary className="cursor-pointer text-white">
                Equipment Details
              </summary>
              <RunsDetailsGear items={member.items.items} />
            </details>
            {/* Vous pouvez ajouter ici un autre details pour les talents */}
            <details className="mt-4">
              <summary className="cursor-pointer text-white">
                Talents Details
              </summary>
              <div className="mt-2">
                <iframe
                  src={`https://www.raidbots.com/simbot/render/talents/${member.character.talentLoadout.loadoutText}?width=900&level=80&hideHeader=1`}
                  width="100%"
                  height="600px"
                ></iframe>
              </div>
            </details>
          </div>
        ))}
      </div>
    </div>
  );
};

export default RunsDetails;
