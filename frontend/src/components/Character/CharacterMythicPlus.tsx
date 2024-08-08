import React from "react";
import { useRaiderIoCharacterProfile } from "@/hooks/useRaiderioApi";

interface CharacterMythicPlusProps {
  region: string;
  realm: string;
  name: string;
}

export default function CharacterMythicPlus({
  region,
  realm,
  name,
}: CharacterMythicPlusProps) {
  const {
    data: character,
    isLoading,
    error,
  } = useRaiderIoCharacterProfile(region, realm, name, [
    "mythic_plus_scores_by_season:current",
    "mythic_plus_highest_level_runs",
  ]);

  if (isLoading) return <div>Loading character data...</div>;
  if (error) return <div>Error loading character data: {error.message}</div>;
  if (!character || !character.mythic_plus_scores_by_season)
    return <div>No Mythic+ data found</div>;

  const season = character.mythic_plus_scores_by_season[0]; // Assuming we're only getting the current season

  const relevantScores = ["dps", "healer", "tank"];
  const sortedScores = relevantScores
    .filter((key) => season.scores[key] > 0)
    .sort((a, b) => season.scores[b] - season.scores[a])
    .map((key) => ({
      type: key,
      score: season.scores[key],
      color: season.segments[key].color,
    }));

  const highestRun =
    character.mythic_plus_highest_level_runs &&
    character.mythic_plus_highest_level_runs.length > 0
      ? character.mythic_plus_highest_level_runs[0]
      : null;

  return (
    <div className="p-4 bg-gradient-dark shadow-lg glow-effect m-12">
      <div className="flex p-4">
        <h3 className="text-xl font-semibold mb-2 text-gradient-glow">
          Mythic+ Progression
        </h3>
      </div>
      <div className="flex px-8 shadow-xl glow-effect rounded-lg">
        <div className="flex justify-between items-center p-2 py-4">
          {sortedScores.map(({ type, score, color }, index) => (
            <div
              key={index}
              className="flex flex-col justify-between items-center p-2 py-4"
            >
              <span
                className="font-bold text-xl  p-2 rounded-sm shadow-xl glow-effect"
                style={{ color }}
              >
                {score}
              </span>
              <span className="font-semibold text-xl uppercase">
                {type.charAt(0).toUpperCase() + type.slice(1)}
              </span>
              <span>Mythic+ score</span>
            </div>
          ))}
          {highestRun && (
            <div className="flex flex-col justify-between items-center p-2 py-4 px-12">
              <span className="font-bold text-xl  p-2 rounded-sm shadow-xl glow-effect">
                {character.mythic_plus_highest_level_runs[0].mythic_level}
              </span>
              <span className="font-semibold text-xl">Highest key level </span>
              <span>This season</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
