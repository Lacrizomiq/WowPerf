import React, { useEffect, useState } from "react";
import {
  useGetBlizzardCharacterSpecializations,
  useGetBlizzardCharacterProfile,
} from "@/hooks/useBlizzardApi";
import { Copy, Check } from "lucide-react";
import TalentTree from "@/components/TalentTree/TalentTree";
import HeroTalentTree from "@/components/TalentTree/HeroTalentTree";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import { SquareArrowOutUpRight } from "lucide-react";
import Image from "next/image";
import {
  TalentNode,
  CharacterTalentProps,
  TalentLoadout,
  HeroTalent,
} from "@/types/talents";
import { toast } from "react-hot-toast";
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
  const [isCopied, setIsCopied] = useState(false);

  const toggleDisplayMode = () => {
    setDisplayMode((prevMode) => (prevMode === "list" ? "tree" : "list"));
  };

  // Copy the talent loadout text to the clipboard
  const copyLoadoutText = async () => {
    if (talentLoadout.encoded_loadout_text) {
      try {
        await navigator.clipboard.writeText(talentLoadout.encoded_loadout_text);
        toast.success("Talent loadout copied to clipboard");
        setIsCopied(true);
        setTimeout(() => {
          setIsCopied(false);
        }, 2000);
      } catch (error) {
        toast.error("Failed to copy talent loadout to clipboard");
        console.error("Failed to copy talent loadout to clipboard", error);
      }
    } else {
      toast.error("No talent loadout text found");
    }
  };

  useWowheadTooltips();

  useEffect(() => {
    if (window.$WowheadPower && window.$WowheadPower.refreshLinks) {
      window.$WowheadPower.refreshLinks();
    }
  }, []);

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
      <div className="mb-6 shadow-2xl border-2 border-[#001830] rounded-lg glow-effect">
        <h3 className="text-lg font-semibold text-white bg-deep-blue p-4 items-center flex justify-center">
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

        <div className="grid grid-cols-4 xs:grid-cols-5 sm:grid-cols-6 md:grid-cols-7 lg:grid-cols-8 xl:grid-cols-9 gap-2 p-4">
          {selectedTalents.map((talent) => (
            <TalentIcon key={talent.id} talent={talent} />
          ))}
        </div>
      </div>
    );
  };

  const renderHeroTalentsGroup = (heroTalents: HeroTalent[]) => {
    const subTreeName =
      talentLoadout.sub_tree_nodes?.[0]?.entries?.[0]?.name ?? "Unknown";
    const subtreeIcon =
      talentLoadout.sub_tree_nodes?.[0]?.entries?.[0]?.atlasMemberName;
    const iconUrl = subtreeIcon
      ? `https://wow.zamimg.com/images/wow/TextureAtlas/live/${subtreeIcon}.webp`
      : "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg";

    return (
      <div className="mb-6 shadow-2xl border-2 border-[#001830] rounded-lg md:w-1/2 glow-effect">
        <h3 className="text-lg font-semibold text-white bg-deep-blue p-4 items-center flex justify-center">
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

        <div className="grid grid-cols-4 xs:grid-cols-5 sm:grid-cols-6 md:grid-cols-7 lg:grid-cols-8 xl:grid-cols-9 gap-2 p-4">
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
    <div className="p-4 sm:p-6 shadow-lg m-2 sm:m-4">
      <style jsx global>{`
        .wowhead-tooltip {
          scale: 1.2;
          transform-origin: top left;
          max-width: 300px;
          font-size: 14px;
        }
      `}</style>
      <div className="flex pb-2 flex-col sm:flex-row justify-between items-center mb-4 space-y-4 sm:space-y-0">
        <h2 className="text-xl sm:text-2xl font-bold">Talent Build Summary</h2>
        <div className="flex flex-col sm:flex-row gap-2 sm:gap-4 w-full sm:w-auto">
          {/* Display Mode Talent Tree or List */}
          <button
            onClick={toggleDisplayMode}
            className="bg-gradient-blue hover:bg-purple-700 text-white font-bold py-2 px-4 rounded-lg w-full sm:w-auto flex items-center justify-center shadow-2xl"
          >
            {displayMode === "list" ? "Show Full Tree" : "Show Talent List"}
          </button>
          {/* Copy Loadout Text */}
          <button
            onClick={copyLoadoutText}
            className="bg-gradient-blue hover:bg-purple-700 text-white font-bold py-2 px-4 rounded-lg w-full sm:w-auto flex items-center justify-center shadow-2xl"
          >
            {isCopied ? (
              <>
                <Check className="mr-2 h-4 w-4" />
                Copied
              </>
            ) : (
              <>
                <Copy className="mr-2 h-4 w-4" />
                Copy Talents
              </>
            )}
          </button>
          {/* Talent Calculator to wowhead */}
          {talentCalculatorUrl && (
            <a
              href={talentCalculatorUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="bg-gradient-blue hover:bg-purple-700 text-white font-bold py-2 px-4 rounded-lg w-full sm:w-auto flex items-center justify-center shadow-2xl "
            >
              Talent Calculator <SquareArrowOutUpRight className="ml-2" />
            </a>
          )}
        </div>
      </div>
      {displayMode === "list" ? (
        <div className="flex flex-col gap-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              {renderTalentGroup(
                talentLoadout.class_talents,
                `${characterClass} Talents`,
                talentLoadout.class_icon
              )}
            </div>
            <div>
              {renderTalentGroup(
                talentLoadout.spec_talents,
                `${activeSpecName} Talents`,
                talentLoadout.spec_icon
              )}
            </div>
          </div>
          {talentLoadout.hero_talents.length > 0 && (
            <div className="">
              {renderHeroTalentsGroup(talentLoadout.hero_talents)}
            </div>
          )}
        </div>
      ) : (
        <>
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
          <HeroTalentTree
            talentTreeId={talentLoadout.tree_id}
            specId={profileData.spec_id}
            region={region}
            namespace={namespace}
            locale={locale}
            className={characterClass}
            specName={activeSpecName}
            selectedHeroTalentTree={[...talentLoadout.hero_talents].filter(
              (t) => t.rank > 0
            )}
          />
        </>
      )}
    </div>
  );
}

interface TalentIconProps {
  talent: TalentNode;
}

const TalentIcon: React.FC<TalentIconProps> = ({ talent }) => {
  const [imageError, setImageError] = useState(false);

  const hasEntries = talent.entries && talent.entries.length > 0;
  const spellId = hasEntries ? talent.entries[0].spellId : undefined;
  const icon = hasEntries ? talent.entries[0].icon : "inv_misc_questionmark";

  return (
    <div className="relative w-full pb-[100%]">
      {" "}
      {/* Aspect ratio 1:1 */}
      {spellId ? (
        <a
          href={`https://www.wowhead.com/spell=${spellId}`}
          data-wowhead={`spell=${spellId}`}
          className="absolute inset-0 block cursor-pointer talent active"
          data-wh-icon-size="medium"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Image
            src={`https://wow.zamimg.com/images/wow/icons/large/${icon}.jpg`}
            alt={talent.name}
            fill
            className="rounded-md border-2 border-gray-700 object-cover"
            onError={() => setImageError(true)}
            unoptimized
          />
          {talent.rank > 0 && (
            <div className="absolute bottom-0 right-0 bg-black bg-opacity-70 text-white text-xs font-bold px-1 rounded">
              {talent.rank}/{talent.maxRanks}
            </div>
          )}
        </a>
      ) : (
        <div className="relative w-8 h-8 sm:w-10 sm:h-10">
          <Image
            src="https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg"
            alt={talent.name}
            width={40}
            height={40}
            className="w-full h-full rounded-md border-2 border-gray-700"
            unoptimized
          />
          {talent.rank > 0 && (
            <div className="absolute bottom-0 right-0 bg-black bg-opacity-70 text-white text-xs font-bold px-1 rounded">
              {talent.rank}/{talent.maxRanks}
            </div>
          )}
        </div>
      )}
    </div>
  );
};

interface HeroTalentIconProps {
  talent: HeroTalent;
}

const HeroTalentIcon: React.FC<HeroTalentIconProps> = ({ talent }) => {
  const [imageError, setImageError] = useState(false);

  const selectedEntry =
    talent.entries.find((entry) => entry.id === talent.id) || talent.entries[0];

  return (
    <div className="relative w-full pb-[80%]">
      <a
        href={`https://www.wowhead.com/spell=${selectedEntry.spellId}`}
        data-wowhead={`spell=${selectedEntry.spellId}`}
        className="absolute inset-0 block cursor-pointer talent active"
        data-wh-icon-size="medium"
        target="_blank"
        rel="noopener noreferrer"
      >
        <Image
          src={
            imageError
              ? "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg"
              : `https://wow.zamimg.com/images/wow/icons/large/${selectedEntry.icon}.jpg`
          }
          alt={talent.name}
          fill
          className="rounded-md border-2 border-gray-700 object-cover"
          onError={() => setImageError(true)}
          unoptimized
        />
        {talent.rank > 0 && (
          <div className="absolute bottom-0 right-0 bg-black bg-opacity-70 text-white text-xs font-bold px-1 rounded">
            {talent.rank}/{selectedEntry.maxRanks}
          </div>
        )}
      </a>
    </div>
  );
};
