import React, { useEffect } from "react";
import { useGetRaiderIoCharacterTalents } from "@/hooks/useRaiderioApi";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import { SquareArrowOutUpRight } from "lucide-react";
import ClassIcons from "@/components/ui/ClassIcons";
import SpecIcons from "@/components/ui/SpecIcon";

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

  const renderTalentGroup = (
    talents: any[],
    title: string,
    isClassTalents: boolean
  ) => (
    <div className="mb-6 shadow-xl glow-effect p-4">
      <h3 className="text-lg font-semibold text-gradient-glow mb-4 items-center flex justify-center">
        {isClassTalents ? (
          <ClassIcons region={region} realm={realm} name={name} />
        ) : (
          <SpecIcons region={region} realm={realm} name={name} />
        )}
        <span className="ml-2">{title}</span>
      </h3>
      <div className="grid grid-cols-7 gap-2 mb-4">
        {talents.map((talent) => {
          const spellEntry = talent.node.entries[talent.entryIndex];
          const iconUrl = `https://wow.zamimg.com/images/wow/icons/large/${spellEntry.spell.icon}.jpg`;
          return (
            <div key={talent.node.id} className="relative">
              <a
                href={`https://www.wowhead.com/spell=${spellEntry.spell.id}`}
                data-wowhead={`spell=${spellEntry.spell.id}`}
                className="block cursor-pointer talent active relative"
                data-wh-icon-size="medium"
                target="_blank"
              >
                <div className="relative w-10 h-10">
                  <img
                    src={iconUrl}
                    alt={spellEntry.spell.name}
                    className="w-full h-full rounded-md border-2 border-gray-700"
                    onError={(e) => {
                      e.currentTarget.src =
                        "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg";
                    }}
                  />
                  {talent.rank > 1 && (
                    <div className="absolute bottom-0 right-0 bg-black bg-opacity-70 text-white text-xs font-bold px-1 rounded">
                      {talent.rank}/2
                    </div>
                  )}
                </div>
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
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold text-gradient-glow flex justify-between mb-6">
          Talent Build Summary
        </h2>
        <div>
          {talentCalculatorUrl && (
            <a
              href={talentCalculatorUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="font-bold flex items-center gap-2 align-center mb-4  hover:text-blue-300"
            >
              Talent Calculator <SquareArrowOutUpRight className="ml-2" />
            </a>
          )}
        </div>
      </div>
      <div className="flex flex-col md:flex-row gap-4">
        <div className="flex-1">
          {renderTalentGroup(
            classTalents,
            `${characterData.class} Talents`,
            true
          )}
        </div>
        <div className="flex-1">
          {renderTalentGroup(
            specTalents,
            `${characterData.active_spec_name} Talents`,
            false
          )}
        </div>
      </div>
    </div>
  );
}
