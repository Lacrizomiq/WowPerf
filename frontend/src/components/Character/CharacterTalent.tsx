import React, { useEffect, useState } from "react";
import {
  useGetBlizzardCharacterSpecializations,
  useGetBlizzardCharacterProfile,
} from "@/hooks/useBlizzardApi";
import TalentTree from "@/components/TalentTree/TalentTree";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import { SquareArrowOutUpRight } from "lucide-react";
import Image from "next/image";
import { TalentNode } from "@/types/talents";

interface CharacterTalentProps {
  region: string;
  realm: string;
  name: string;
  namespace: string;
  locale: string;
}

interface TalentLoadout {
  loadout_spec_id: number;
  tree_id: number;
  loadout_text: string;
  encoded_loadout_text: string;
  class_icon: string;
  spec_icon: string;
  class_talents: TalentNode[];
  spec_talents: TalentNode[];
  sub_tree_nodes: SubTreeNode[];
  hero_talents: HeroTalent[];
}

interface SubTreeNode {
  id: number;
  name: string;
  type: string;
  entries: SubTreeEntry[];
}

interface SubTreeEntry {
  id: number;
  type: string;
  name: string;
  traitSubTreeId: number;
  atlasMemberName: string;
  nodes: number[];
}

interface HeroTalent {
  id: number;
  type: string;
  name: string;
  traitSubTreeId: number;
  posX: number;
  posY: number;
  nodes: number[];
  rank: number;
  entries: HeroEntry[];
}

interface HeroEntry {
  id: number;
  name: string;
  type: string;
  maxRanks: number;
  entryNode: boolean;
  subTreeId: number;
  freeNode: boolean;
  spellId: number;
  icon: string;
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

  const [displayMode, setDisplayMode] = useState<"list" | "tree">("list");

  const toggleDisplayMode = () => {
    setDisplayMode((prevMode) => (prevMode === "list" ? "tree" : "list"));
  };

  useWowheadTooltips();

  useEffect(() => {
    if (specializationsData) {
      console.log("Specializations Data:", specializationsData);
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

  if (!specializationsData?.talent_loadout || !profileData) {
    console.log("No talent or profile data found");
    return (
      <div className="text-yellow-500">No talent or profile data found</div>
    );
  }

  const talentLoadout: TalentLoadout = specializationsData.talent_loadout;
  const characterClass = profileData.class || "Unknown Class";
  const activeSpecName = profileData.active_spec_name || "Unknown Spec";

  console.log("Talent Loadout:", talentLoadout);
  console.log("Character Class:", characterClass);
  console.log("Active Spec Name:", activeSpecName);

  const renderTalentGroup = (
    talents: TalentNode[],
    title: string,
    icon: string
  ) => {
    const selectedTalents = talents.filter((talent) => talent.rank > 0);

    if (selectedTalents.length === 0) {
      return <div className="text-yellow-500">No {title} found</div>;
    }

    return (
      <div className="mb-6 shadow-xl glow-effect p-4">
        <h3 className="text-lg font-semibold text-gradient-glow mb-4 items-center flex justify-center">
          <Image
            src={icon}
            alt={title}
            width={32}
            height={32}
            className="mr-2"
            unoptimized
          />
          <span>{title}</span>
        </h3>

        <div className="grid grid-cols-7 gap-2 mb-4">
          {selectedTalents.map((talent) => (
            <TalentIcon key={talent.id} talent={talent} />
          ))}
        </div>
      </div>
    );
  };

  const renderHeroTalentsGroup = (heroTalents: HeroTalent[]) => {
    const subTreeName = talentLoadout.sub_tree_nodes[0]?.name;
    const subtreeIcon =
      talentLoadout.sub_tree_nodes[0]?.entries[0]?.atlasMemberName;
    const iconUrl = subtreeIcon
      ? `https://wow.zamimg.com/images/wow/TextureAtlas/live/${subtreeIcon}.webp`
      : "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg";

    return (
      <div className="mb-6">
        <h3 className="text-lg font-semibold text-gradient-glow mb-4 items-center flex justify-center">
          <Image
            src={iconUrl}
            alt="Hero Talents"
            width={40}
            height={40}
            className="mr-2"
            unoptimized
          />
          <span>{subTreeName} Hero Talents</span>
        </h3>

        <div className="grid grid-cols-7 gap-2 mb-4">
          {heroTalents.map((talent) => (
            <HeroTalentIcon key={talent.id} talent={talent} />
          ))}
        </div>
      </div>
    );
  };

  const getTalentCalculatorUrl = () => {
    if (
      !characterClass ||
      !activeSpecName ||
      !talentLoadout.encoded_loadout_text
    ) {
      return "";
    }

    const classSlug = characterClass.toLowerCase().replace(" ", "");
    const specSlug = activeSpecName.toLowerCase().replace(" ", "");
    const encodedLoadout = encodeURIComponent(
      talentLoadout.encoded_loadout_text
    );

    return `https://www.wowhead.com/talent-calc/${classSlug}/${specSlug}/${encodedLoadout}`;
  };

  const talentCalculatorUrl = getTalentCalculatorUrl();

  return (
    <div className="p-6 bg-gradient-dark shadow-lg rounded-lg glow-effect m-12 max-w-6xl mx-auto">
      <style jsx global>{`
        .wowhead-tooltip {
          scale: 1.2;
          transform-origin: top left;
          max-width: 300px;
          font-size: 14px;
        }
      `}</style>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold text-gradient-glow">
          Talent Build Summary
        </h2>
        <div className="flex gap-4">
          <button
            onClick={toggleDisplayMode}
            className="bg-purple-600 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded"
          >
            {displayMode === "list" ? "Show Full Tree" : "Show Talent List"}
          </button>
          {talentCalculatorUrl && (
            <a
              href={talentCalculatorUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="font-bold flex items-center gap-2 align-center hover:text-blue-300"
            >
              Talent Calculator <SquareArrowOutUpRight className="ml-2" />
            </a>
          )}
        </div>
      </div>
      {displayMode === "list" ? (
        <div className="flex flex-col md:flex-col gap-4">
          <div className="flex-1">
            {renderTalentGroup(
              talentLoadout.class_talents,
              `${characterClass} Talents`,
              talentLoadout.class_icon
            )}
          </div>
          <div className="flex-1">
            {renderTalentGroup(
              talentLoadout.spec_talents,
              `${activeSpecName} Talents`,
              talentLoadout.spec_icon
            )}
          </div>
          {talentLoadout.hero_talents.length > 0 && (
            <div className="flex-1">
              {renderHeroTalentsGroup(talentLoadout.hero_talents)}
            </div>
          )}
        </div>
      ) : (
        <TalentTree
          talentTreeId={talentLoadout.tree_id}
          specId={profileData.spec_id}
          region={region}
          namespace={namespace}
          locale={locale}
          className={characterClass}
          specName={activeSpecName}
          selectedTalents={[
            ...talentLoadout.class_talents,
            ...talentLoadout.spec_talents,
          ].filter((t) => t.rank > 0)}
        />
      )}
    </div>
  );
}

interface TalentIconProps {
  talent: TalentNode;
}

const TalentIcon: React.FC<TalentIconProps> = ({ talent }) => {
  const [imageError, setImageError] = useState(false);

  return (
    <div className="relative">
      <a
        href={`https://www.wowhead.com/spell=${talent.entries[0].spellId}`}
        data-wowhead={`spell=${talent.entries[0].spellId}`}
        className="block cursor-pointer talent active relative"
        data-wh-icon-size="medium"
        target="_blank"
        rel="noopener noreferrer"
      >
        <div className="relative w-10 h-10">
          <Image
            src={
              imageError
                ? "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg"
                : `https://wow.zamimg.com/images/wow/icons/large/${talent.entries[0].icon}.jpg`
            }
            alt={talent.name}
            width={40}
            height={40}
            className="w-full h-full rounded-md border-2 border-gray-700"
            onError={() => setImageError(true)}
            unoptimized
          />
          {talent.rank > 0 && (
            <div className="absolute bottom-0 right-0 bg-black bg-opacity-70 text-white text-xs font-bold px-1 rounded">
              {talent.rank}/{talent.maxRanks}
            </div>
          )}
        </div>
      </a>
    </div>
  );
};

interface HeroTalentIconProps {
  talent: HeroTalent;
}

const HeroTalentIcon: React.FC<HeroTalentIconProps> = ({ talent }) => {
  const [imageError, setImageError] = useState(false);

  return (
    <div className="relative">
      <a
        href={`https://www.wowhead.com/spell=${talent.entries[0].spellId}`}
        data-wowhead={`spell=${talent.entries[0].spellId}`}
        className="block cursor-pointer talent active relative"
        data-wh-icon-size="medium"
        target="_blank"
        rel="noopener noreferrer"
      >
        <div className="relative w-10 h-10">
          <Image
            src={
              imageError
                ? "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg"
                : `https://wow.zamimg.com/images/wow/icons/large/${talent.entries[0].icon}.jpg`
            }
            alt={talent.name}
            width={40}
            height={40}
            className="w-full h-full rounded-md border-2 border-gray-700"
            onError={() => setImageError(true)}
            unoptimized
          />
          {talent.rank > 0 && (
            <div className="absolute bottom-0 right-0 bg-black bg-opacity-70 text-white text-xs font-bold px-1 rounded">
              {talent.rank}/{talent.entries[0].maxRanks}
            </div>
          )}
        </div>
      </a>
    </div>
  );
};
