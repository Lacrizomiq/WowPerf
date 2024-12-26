import React, { useEffect, useState } from "react";
import {
  useGetBlizzardCharacterSpecializations,
  useGetBlizzardCharacterProfile,
} from "@/hooks/useBlizzardApi";
import { Copy, Check, SquareArrowOutUpRight, LayoutGrid } from "lucide-react";
import TalentTree from "@/components/TalentTree/TalentTree";
import HeroTalentTree from "@/components/TalentTree/HeroTalentTree";
import { useWowheadTooltips } from "@/hooks/useWowheadTooltips";
import Image from "next/image";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  TalentNode,
  CharacterTalentProps,
  TalentLoadout,
  HeroTalent,
} from "@/types/talents";
import { toast } from "react-hot-toast";
import RaidbotsTreeDisplay from "@/components/TalentTree/RaidbotTalentTree";

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
        await navigator.clipboard.writeText(
          encodeURIComponent(talentLoadout.encoded_loadout_text)
        );
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
      <Card className="bg-gray-900/50 border-gray-800 shadow-xl">
        <CardHeader className="pb-2 border-b border-gray-800">
          <CardTitle className="flex items-center gap-3 text-xl text-white justify-center">
            <Image
              src={icon}
              alt={title}
              width={32}
              height={32}
              className="rounded"
              unoptimized
            />
            <span>{title}</span>
          </CardTitle>
        </CardHeader>
        <CardContent className="pt-4">
          <div className="grid grid-cols-4 xs:grid-cols-5 sm:grid-cols-6 md:grid-cols-7 lg:grid-cols-8 xl:grid-cols-9 gap-2">
            {selectedTalents.map((talent) => (
              <TalentIcon key={talent.id} talent={talent} />
            ))}
          </div>
        </CardContent>
      </Card>
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
      <Card className="bg-gray-900/50 border-gray-800 shadow-xl">
        <CardHeader className="pb-2 border-b border-gray-800 flex justify-center">
          <CardTitle className="flex items-center gap-3 text-xl text-white justify-center">
            <Image
              src={iconUrl}
              alt="Hero Talents"
              width={40}
              height={40}
              className="rounded"
              unoptimized
            />
            <span>{subTreeName} Hero Talents</span>
          </CardTitle>
        </CardHeader>
        <CardContent className="pt-4">
          <div className="grid grid-cols-4 xs:grid-cols-5 sm:grid-cols-6 md:grid-cols-7 lg:grid-cols-8 xl:grid-cols-9 gap-2">
            {heroTalents.map((talent) => (
              <HeroTalentIcon key={talent.id} talent={talent} />
            ))}
          </div>
        </CardContent>
      </Card>
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
    <div className="p-6 bg-gray-950">
      <style jsx global>{`
        .wowhead-tooltip {
          scale: 1.2;
          transform-origin: top left;
          max-width: 300px;
          font-size: 14px;
        }
      `}</style>

      {/* Header with Actions */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
        <h2 className="text-2xl font-bold text-white">Talent Build Summary</h2>
        <div className="flex flex-wrap gap-4">
          <Button
            variant="secondary"
            className="bg-blue-600 hover:bg-blue-700 text-white shadow-lg"
            onClick={toggleDisplayMode}
          >
            <LayoutGrid className="w-4 h-4 mr-2" />
            {displayMode === "list" ? "Show Full Tree" : "Show Talent List"}
          </Button>

          <Button
            variant="secondary"
            className="bg-blue-600 hover:bg-blue-700 text-white shadow-lg"
            onClick={copyLoadoutText}
          >
            {isCopied ? (
              <>
                <Check className="w-4 h-4 mr-2" />
                Copied
              </>
            ) : (
              <>
                <Copy className="w-4 h-4 mr-2" />
                Copy Talents
              </>
            )}
          </Button>

          {talentCalculatorUrl && (
            <Button
              variant="secondary"
              className="bg-blue-600 hover:bg-blue-700 text-white shadow-lg"
              asChild
            >
              <a
                href={talentCalculatorUrl}
                target="_blank"
                rel="noopener noreferrer"
              >
                Talent Calculator
                <SquareArrowOutUpRight className="w-4 h-4 ml-2" />
              </a>
            </Button>
          )}
        </div>
      </div>

      {displayMode === "list" ? (
        <div className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Class Talents */}
            <div>
              {renderTalentGroup(
                talentLoadout.class_talents,
                `${characterClass} Talents`,
                talentLoadout.class_icon
              )}
            </div>

            {/* Spec Talents */}
            <div>
              {renderTalentGroup(
                talentLoadout.spec_talents,
                `${activeSpecName} Talents`,
                talentLoadout.spec_icon
              )}
            </div>
          </div>

          {/* Hero Talents */}
          {talentLoadout.hero_talents.length > 0 && (
            <div className="md:w-1/2">
              {renderHeroTalentsGroup(talentLoadout.hero_talents)}
            </div>
          )}
        </div>
      ) : (
        <div className="space-y-6 p-4">
          <RaidbotsTreeDisplay
            encodedString={talentLoadout.encoded_loadout_text}
            width={1000}
            locale={locale}
            hideHeader={false}
            hideExport={false}
          />
        </div>
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
