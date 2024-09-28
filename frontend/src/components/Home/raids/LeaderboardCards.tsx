import React from "react";
import Image from "next/image";
import { useGetRaiderioRaidLeaderboard } from "@/hooks/useRaiderioApi";
import { ProgressionItem, GuildProgression } from "@/types/raidLeaderboard";

interface LeaderBoardCardsProps {
  raid: string;
  difficulty: string;
  region: string;
}

const LeaderBoardCards: React.FC<LeaderBoardCardsProps> = ({
  raid,
  difficulty,
  region,
}) => {
  const { data, isLoading, error } = useGetRaiderioRaidLeaderboard(
    raid,
    difficulty,
    region
  );

  if (isLoading)
    return <div className="text-white">Loading leaderboard data...</div>;
  if (error)
    return <div className="text-red-500">Error loading leaderboard data.</div>;
  if (!data || !data.progression)
    return <div className="text-white">No leaderboard data available.</div>;

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {data.progression.map((item: ProgressionItem, index: number) => (
        <div key={index} className="bg-gray-800 rounded-lg shadow-lg p-4">
          <h3 className="text-xl font-bold mb-2 text-white">
            Boss {item.progress}
          </h3>
          <p className="text-sm mb-4 text-gray-300">
            Total Guilds: {item.totalGuilds}
          </p>
          {item.guilds.map(
            (guildProgression: GuildProgression, guildIndex: number) => (
              <div key={guildIndex} className="mb-4 last:mb-0">
                <div className="flex items-center mb-2">
                  {guildProgression.guild.logo && (
                    <Image
                      src={guildProgression.guild.logo}
                      alt={`${guildProgression.guild.name} logo`}
                      width={40}
                      height={40}
                      className="rounded-full mr-2"
                    />
                  )}
                  <div>
                    <h4 className="font-semibold text-white">
                      {guildProgression.guild.name}
                    </h4>
                    <p className="text-sm text-gray-300">
                      {guildProgression.guild.realm.name}
                    </p>
                  </div>
                </div>
                <p className="text-sm text-gray-300">
                  Defeated at:{" "}
                  {new Date(guildProgression.defeatedAt).toLocaleString()}
                </p>
                {guildProgression.streamers.count > 0 && (
                  <p className="text-sm text-blue-400">
                    Live Streamers: {guildProgression.streamers.count}
                  </p>
                )}
              </div>
            )
          )}
        </div>
      ))}
    </div>
  );
};

export default LeaderBoardCards;
