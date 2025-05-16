// RunsDetails.tsx - Version réajustée avec fond plus sombre
import React, { useEffect } from "react";
import { MythicPlusRun, Roster } from "@/types/runsDetails";
import { useGetRaiderioMythicPlusRunDetails } from "@/hooks/useRaiderioApi";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import RunsDetailsGear from "./RunsDetailsGear";
import { SquareArrowOutUpRight } from "lucide-react";
import Link from "next/link";

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
    return (
      <div className="flex justify-center items-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-600"></div>
      </div>
    );
  if (error)
    return (
      <div className="bg-slate-800/30 border border-red-500 rounded-md p-4 my-4">
        <h3 className="text-red-500 text-lg font-medium">
          Error: {(error as Error).message}
        </h3>
      </div>
    );
  if (!runDetails)
    return (
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 text-center">
        <p className="text-slate-400">No details available</p>
      </div>
    );

  return (
    <div className="px-4 pb-4 bg-slate-800/30">
      <h3 className="text-2xl font-bold mb-4 text-white">Run Details</h3>
      <p className="text-slate-400 text-sm mb-4">
        Here you can see the details of the mythic plus run, including the gear,
        talents and score of each player during that key.
      </p>
      <div className="flex flex-col gap-4">
        {runDetails.roster.map((member: Roster) => (
          <div
            key={member.character.id}
            className="bg-slate-800/40 rounded-lg p-4 border border-slate-700 shadow-lg"
          >
            <div className="flex items-center mb-2 justify-between">
              <div>
                <Link
                  href={`/character/${member.character.region.slug}/${
                    member.character.realm.slug
                  }/${member.character.name.toLowerCase()}`}
                >
                  <p
                    className={`font-bold flex items-center ${
                      member.character.class.slug
                        ? `class-color--${member.character.class.slug} hover:underline hover:decoration-current`
                        : ""
                    }`}
                  >
                    {member.character.name}
                    <SquareArrowOutUpRight className="w-4 h-4 ml-2" />
                  </p>
                </Link>
                <p className="text-white text-sm">
                  {member.character.spec.name} {member.character.class.name}
                </p>
              </div>
              <div className="flex flex-col gap-2">
                <p className="text-white">
                  Item Level: {member.items.item_level_equipped.toFixed(1)}
                </p>
                <p className="text-white">Score: {member.ranks.score}</p>
              </div>
            </div>
            <details className="mt-2">
              <summary className="cursor-pointer text-white hover:text-purple-400 transition-colors">
                Equipment Details
              </summary>
              <RunsDetailsGear items={member.items.items} />
            </details>

            <details className="mt-4">
              <summary className="cursor-pointer text-white hover:text-purple-400 transition-colors">
                Talents Details
              </summary>
              <div className="mt-2 bg-[#1a1c25] rounded-lg p-2 border border-slate-800 shadow-lg">
                <iframe
                  src={`https://www.raidbots.com/simbot/render/talents/${member.character.talentLoadout.loadoutText}?width=900&level=80`}
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
