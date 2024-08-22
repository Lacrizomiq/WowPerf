import React, { useEffect } from "react";
import {
  useGetBlizzardCharacterSpecializations,
  useGetBlizzardCharacterProfile,
} from "@/hooks/useBlizzardApi";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import { SquareArrowOutUpRight } from "lucide-react";
import ClassIcons from "@/components/ui/ClassIcons";
import SpecIcons from "@/components/ui/SpecIcon";
import Image from "next/image";

interface CharacterTalentProps {
  region: string;
  realm: string;
  name: string;
  namespace: string;
  locale: string;
}

interface TalentNode {
  node: {
    id: number;
    entries: Array<{
      spell: {
        id: number;
        name: string;
        icon_url: string;
      };
    }>;
  };
  entryIndex: number;
  rank: number;
}

interface SpecializationsData {
  encoded_loadout: string;
  talent_loadout: {
    class_talents: TalentNode[];
    spec_talents: TalentNode[];
    encoded_loadout_text: string;
    loadout_spec_id: number;
    loadout_text: string;
  };
}

interface ProfileData {
  class: string;
  active_spec_name: string;
}

export default function CharacterTalent({
  region,
  realm,
  name,
  namespace,
  locale,
}: CharacterTalentProps) {
  const {
    data: specializationsData,
    isLoading: isLoadingSpecializations,
    error: specializationsError,
  } = useGetBlizzardCharacterSpecializations(
    region,
    realm,
    name,
    namespace,
    locale
  );

  const {
    data: profileData,
    isLoading: isLoadingProfile,
    error: profileError,
  } = useGetBlizzardCharacterProfile(region, realm, name, namespace, locale);

  useWowheadTooltips();

  useEffect(() => {
    if (specializationsData && window.$WowheadPower) {
      window.$WowheadPower.refreshLinks();
    }
  }, [specializationsData]);

  if (isLoadingSpecializations || isLoadingProfile)
    return <div className="text-white">Loading talent data...</div>;

  if (specializationsError || profileError) {
    console.error("Specializations Error:", specializationsError);
    console.error("Profile Error:", profileError);
    return (
      <div className="text-red-500">
        Error loading talent data:{" "}
        {((specializationsError || profileError) as Error)?.message ||
          "Unknown error"}
      </div>
    );
  }

  console.log("Full Specializations Data:", specializationsData);
  console.log("Class Talents:", specializationsData?.class_talents);
  console.log("Spec Talents:", specializationsData?.spec_talents);
  console.log("Full Profile Data:", profileData);

  if (!specializationsData && !profileData) {
    return (
      <div className="text-yellow-500">No talent or profile data found</div>
    );
  }

  if (!specializationsData) {
    return <div className="text-yellow-500">No talent data found</div>;
  }

  if (!profileData) {
    return <div className="text-yellow-500">No profile data found</div>;
  }

  const classTalents = specializationsData.talent_loadout?.class_talents || [];
  const specTalents = specializationsData.talent_loadout?.spec_talents || [];
  const loadoutText = specializationsData.talent_loadout?.loadout_text || "";

  const characterClass = profileData.class || "Unknown Class";
  const activeSpecName = profileData.active_spec_name || "Unknown Spec";

  const renderTalentGroup = (
    talents: TalentNode[],
    title: string,
    isClassTalents: boolean
  ) => {
    if (!talents || talents.length === 0) {
      return <div className="text-yellow-500">No {title} found</div>;
    }

    return (
      <div className="mb-6 shadow-xl glow-effect p-4">
        <h3 className="text-lg font-semibold text-gradient-glow mb-4 items-center flex justify-center">
          {isClassTalents ? (
            <ClassIcons
              region={region}
              realm={realm}
              name={name}
              namespace={namespace}
              locale={locale}
            />
          ) : (
            <SpecIcons
              region={region}
              realm={realm}
              name={name}
              namespace={namespace}
              locale={locale}
            />
          )}
          <span className="ml-2">{title}</span>
        </h3>
        <div className="grid grid-cols-7 gap-2 mb-4">
          {talents.map((talent) => {
            const spellEntry = talent.node.entries[talent.entryIndex];
            if (!spellEntry) return null;
            return (
              <div key={talent.node.id} className="relative">
                <a
                  href={`https://www.wowhead.com/spell=${spellEntry.spell.id}`}
                  data-wowhead={`spell=${spellEntry.spell.id}`}
                  className="block cursor-pointer talent active relative"
                  data-wh-icon-size="medium"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <div className="relative w-10 h-10">
                    <Image
                      src={spellEntry.spell.icon_url}
                      alt={spellEntry.spell.name}
                      width={40}
                      height={40}
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
  };

  const getTalentCalculatorUrl = () => {
    if (!characterClass || !activeSpecName || !loadoutText) {
      return "";
    }

    const classSlug = characterClass.toLowerCase().replace(" ", "");
    const specSlug = activeSpecName.toLowerCase().replace(" ", "");
    const encodedLoadout = encodeURIComponent(loadoutText);

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
              className="font-bold flex items-center gap-2 align-center mb-4 hover:text-blue-300"
            >
              Talent Calculator <SquareArrowOutUpRight className="ml-2" />
            </a>
          )}
        </div>
      </div>
      <div className="flex flex-col md:flex-row gap-4">
        <div className="flex-1">
          {renderTalentGroup(classTalents, `${characterClass} Talents`, true)}
        </div>
        <div className="flex-1">
          {renderTalentGroup(specTalents, `${activeSpecName} Talents`, false)}
        </div>
      </div>
    </div>
  );
}
