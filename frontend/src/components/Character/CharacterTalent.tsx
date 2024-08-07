import React, { useEffect } from "react";
import Image from "next/image";
import Link from "next/link";
import { useGetRaiderIoCharacterTalents } from "@/hooks/useRaiderioApi";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import { SquareArrowOutUpRight } from "lucide-react";

interface CharacterTalentProps {
  region: string;
  realm: string;
  name: string;
}

export default function CharacterTalent({
  region,
  realm,
  name,
}: CharacterTalentProps) {
  const {
    data: characterData,
    isLoading,
    error,
  } = useGetRaiderIoCharacterTalents(region, realm, name);
  useWowheadTooltips();

  useEffect(() => {
    if (characterData && window.$WowheadPower) {
      window.$WowheadPower.refreshLinks();
    }
  }, [characterData]);

  if (isLoading)
    return <div className="text-white">Loading talent data...</div>;
  if (error)
    return (
      <div className="text-red-500">
        Error loading talent data:{" "}
        {error instanceof Error ? error.message : "Unknown error"}
      </div>
    );

  if (
    !characterData ||
    (!characterData.talentLoadout && !characterData.talents)
  )
    return <div className="text-yellow-500">No talent data found</div>;

  const classTalents = characterData.talentLoadout?.class_talents || [];
  const specTalents = characterData.talentLoadout?.spec_talents || [];

  const renderTalentGroup = (talents: any[], title: string) => (
    <div className="mb-6 shadow-xl glow-effect p-4">
      <h3 className="text-lg font-semibold text-gradient-glow mb-4 items-center flex justify-center">
        {title}
      </h3>
      <div className="grid grid-cols-7 gap-2 mb-4">
        {talents.map((talent) => {
          const spellEntry = talent.node.entries[talent.entryIndex];
          return (
            <div key={talent.node.id} className="relative">
              <a
                href={`https://www.wowhead.com/spell=${spellEntry.spell.id}`}
                data-wowhead={`spell=${spellEntry.spell.id}`}
                className="block cursor-pointer"
                data-wh-icon-size="medium"
              >
                <Image
                  src={`https://wow.zamimg.com/images/wow/icons/large/${spellEntry.spell.icon}.jpg`}
                  alt={spellEntry.spell.name}
                  width={40}
                  height={40}
                  className="rounded-md border-2 border-gray-700"
                />
              </a>
            </div>
          );
        })}
      </div>
    </div>
  );

  const getTalentCalculatorUrl = () => {
    if (
      !characterData.class ||
      !characterData.active_spec_name ||
      !characterData.talentLoadout?.loadout_text
    ) {
      return "";
    }

    const classSlug = characterData.class.toLowerCase().replace(" ", "");
    const specSlug = characterData.active_spec_name
      .toLowerCase()
      .replace(" ", "");
    const encodedLoadout = encodeURIComponent(
      characterData.talentLoadout.loadout_text
    );

    return `https://www.wowhead.com/talent-calc/${classSlug}/${specSlug}/${encodedLoadout}`;
  };

  const talentCalculatorUrl = getTalentCalculatorUrl();

  return (
    <div className="p-6 bg-gradient-dark shadow-lg rounded-lg glow-effect m-12">
      <style jsx global>{`
        .wowhead-tooltip {
          scale: 1.2;
          transform-origin: top left;
          max-width: 300px;
          font-size: 14px;
        }
      `}</style>
      <h2 className="text-2xl font-bold text-gradient-glow flex justify-center mb-6">
        Talent Build Summary
      </h2>
      <div className="flex flex-col md:flex-row gap-4">
        <div className="flex-1">
          {renderTalentGroup(classTalents, "CLASS TALENTS")}
        </div>
        <div className="flex-1">
          {renderTalentGroup(specTalents, "SPEC TALENTS")}
        </div>
      </div>
      <div className="flex justify-center">
        {talentCalculatorUrl && (
          <a
            href={talentCalculatorUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="font-bold flex items-center gap-2 align-center mb-4 text-blue-400 hover:text-blue-300"
          >
            Talent Calculator <SquareArrowOutUpRight className="ml-2" />
          </a>
        )}
      </div>
    </div>
  );
}
