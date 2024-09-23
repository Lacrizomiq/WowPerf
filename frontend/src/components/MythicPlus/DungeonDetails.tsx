import React, { useEffect } from "react";
import { MythicPlusRuns } from "@/types/mythicPlusRuns";
import Image from "next/image";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import { useGetBlizzardCharacterProfile } from "@/hooks/useBlizzardApi";
import Link from "next/link";

interface DungeonDetailsProps {
  run: MythicPlusRuns;
  region: string;
}

const DungeonDetails: React.FC<DungeonDetailsProps> = ({ run, region }) => {
  useWowheadTooltips();

  useEffect(() => {
    if (window.$WowheadPower) {
      window.$WowheadPower.refreshLinks();
    }
  }, [run]);

  return (
    <div className="mt-8 p-6 bg-deep-blue rounded-lg glow-effect">
      <h2 className="text-2xl font-bold mb-4">
        Detailed run for {run.Dungeon.Name}
      </h2>
      <div className="grid grid-cols-2 gap-4">
        <div>
          <p>Key level: {run.KeystoneLevel}</p>
          <p>Score: {run.MythicRating.toFixed(2)}</p>
          <p>Time: {(run.Duration / 1000 / 60).toFixed(2)} minutes</p>
          <p>Completed: {new Date(run.CompletedTimestamp).toDateString()}</p>
        </div>
        <div>
          <ul className="flex flex-wrap gap-1">
            {run.Affixes.map((affix) => (
              <li key={affix.ID} className="flex flex-col items-center w-24">
                <a
                  href={affix.WowheadURL}
                  data-wowhead={`affix=${affix.ID}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex flex-col items-center"
                >
                  <Image
                    src={`https://wow.zamimg.com/images/wow/icons/large/${affix.Icon}.jpg`}
                    alt={affix.Name}
                    width={44}
                    height={44}
                    unoptimized
                  />
                  <span className="text-xs text-center mt-1 break-words">
                    {affix.Name}
                  </span>
                </a>
              </li>
            ))}
          </ul>
        </div>
      </div>
      <div className="mt-4">
        <h3 className="text-xl font-bold mb-2">Party Members:</h3>
        <ul className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
          {run.Members.map((member) => (
            <PartyMember
              key={member.CharacterID}
              member={member}
              region={region}
            />
          ))}
        </ul>
      </div>
    </div>
  );
};

interface PartyMemberProps {
  member: MythicPlusRuns["Members"][0];
  region: string;
}

const PartyMember: React.FC<PartyMemberProps> = ({ member, region }) => {
  const {
    data: character,
    isLoading,
    error,
  } = useGetBlizzardCharacterProfile(
    region.toLowerCase(),
    member.RealmSlug.toLowerCase(),
    member.CharacterName.toLowerCase(),
    `profile-${region.toLowerCase()}`,
    "en_GB"
  );

  const characterUrl = `/character/${region.toLowerCase()}/${member.RealmSlug.toLowerCase()}/${member.CharacterName.toLowerCase()}`;

  return (
    <li className="bg-deep-blue p-2 rounded-lg">
      <Link
        href={characterUrl}
        className="flex items-center space-x-2 hover:bg-opacity-50 transition-colors duration-200"
      >
        <div className="w-12 h-12 bg-deep-blue bg-opacity-50 rounded-lg overflow-hidden shadow-lg glow-effect">
          {character && character.avatar_url ? (
            <Image
              src={character.avatar_url}
              alt={member.CharacterName}
              width={48}
              height={48}
              unoptimized
            />
          ) : (
            <div className="w-full h-full bg-gray-700" />
          )}
        </div>
        <div>
          <p
            className={`font-bold ${
              character ? `class-color--${character.tree_id}` : ""
            }`}
          >
            {member.CharacterName}
          </p>
          <p className="text-sm">
            {member.Specialization} {member.RaceName}
          </p>
        </div>
      </Link>
    </li>
  );
};

export default DungeonDetails;
