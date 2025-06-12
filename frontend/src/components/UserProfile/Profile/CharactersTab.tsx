// CharactersTab.tsx - Version simplifiée avec auto-relink
import React from "react";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { AlertCircle, Clock, UserPlus, RefreshCw, Users } from "lucide-react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useCharacters } from "@/hooks/useCharacters";
import { useBattleNetLink } from "@/hooks/useBattleNetLink";
import { EnrichedUserCharacter } from "@/types/character/character";
import { getClassIcon } from "@/utils/classandspecicons";

interface CharactersTabProps {
  isActive: boolean;
}

const CharactersTab: React.FC<CharactersTabProps> = ({ isActive }) => {
  const router = useRouter();
  const { linkStatus } = useBattleNetLink();
  const {
    characters: charactersData,
    hasCharacters,
    isLoadingCharacters,
    actions,
    isLoading,
    rateLimitState,
    ui,
    region,
  } = useCharacters();

  // Type assertion pour être sûr du type
  const characters: EnrichedUserCharacter[] = charactersData || [];

  if (!isActive) return null;

  const handleCharacterClick = (character: EnrichedUserCharacter) => {
    const characterName = character.name.toLowerCase();
    const realmSlug = character.realm;
    const safeRegion = region || "eu";
    router.push(`/character/${safeRegion}/${realmSlug}/${characterName}`);
  };

  const getClassColor = (className: string) => {
    const classColors: Record<string, string> = {
      "death knight": "text-red-400",
      "demon hunter": "text-purple-400",
      druid: "text-orange-400",
      evoker: "text-teal-400",
      hunter: "text-green-400",
      mage: "text-cyan-400",
      monk: "text-green-300",
      paladin: "text-pink-400",
      priest: "text-white",
      rogue: "text-yellow-400",
      shaman: "text-blue-400",
      warlock: "text-indigo-400",
      warrior: "text-amber-600",
    };
    return classColors[className.toLowerCase()] || "text-slate-300";
  };

  // Loading state
  if (isLoadingCharacters) {
    return (
      <div className="space-y-6">
        <div className="flex justify-center items-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-500" />
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header avec compteur et boutons */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-xl font-bold flex items-center gap-2">
            <Users className="w-5 h-5" />
            My Characters
            {hasCharacters && (
              <span className="text-sm text-slate-400 font-normal">
                ({characters.length} character
                {characters.length !== 1 ? "s" : ""})
              </span>
            )}
          </h2>
          <p className="text-sm text-slate-400 mt-1">
            {hasCharacters
              ? "Characters synchronized from your Battle.net account"
              : "Sync your characters to get started"}
          </p>
        </div>

        {/* Boutons d'action */}
        <div className="flex items-center gap-3">
          {/* Rate limit indicator */}
          {ui.showRateLimit && (
            <div className="text-sm text-orange-400 flex items-center gap-2">
              <Clock className="w-4 h-4" />
              {rateLimitState.formattedTime}
            </div>
          )}

          {/* Sync Button - TOUJOURS DISPONIBLE (auto-relink si nécessaire) */}
          <Button
            onClick={actions.syncAndEnrich}
            disabled={ui.isDisabled.sync}
            className="bg-purple-600 hover:bg-purple-700"
          >
            <UserPlus className="mr-2 h-4 w-4" />
            {isLoading.sync
              ? "Syncing..."
              : hasCharacters
              ? "Sync Characters"
              : "Import Characters"}
          </Button>

          {/* Refresh Button - seulement si connecté ET a des personnages */}
          {hasCharacters && linkStatus?.linked && (
            <Button
              onClick={actions.refreshAndEnrich}
              disabled={ui.isDisabled.refresh}
              variant="outline"
              className="border-blue-600 bg-blue-400 text-white hover:bg-blue-600"
            >
              <Image
                src="https://cdn.raiderio.net/assets/img/battlenet-icon-e75d33039b37cf7cd82eff67d292f478.png"
                alt="Battle.net"
                width={20}
                height={20}
                className="mr-2"
              />
              {isLoading.refresh
                ? "Refreshing..."
                : "Refresh Blizzard Account Data "}
            </Button>
          )}
        </div>
      </div>

      {/* Rate limit message */}
      {ui.showRateLimit && (
        <div className="bg-orange-900/30 border border-orange-500/50 rounded-lg p-4">
          <div className="flex items-start gap-3">
            <Clock className="h-4 w-4 mt-0.5 text-orange-400" />
            <div>
              <h4 className="font-semibold text-orange-400 mb-1">
                Please wait
              </h4>
              <p className="text-sm text-orange-300">
                {rateLimitState.message}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* AFFICHAGE CONDITIONNEL SIMPLIFIÉ */}
      {!hasCharacters ? (
        /* Aucun personnage - Invitation à la sync */
        <div className="text-center py-12">
          <div className="mb-6">
            <svg
              className="w-16 h-16 text-gray-500 mx-auto mb-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>

            <h3 className="text-xl font-semibold mb-2">No Characters Yet</h3>

            {isLoading.sync ? (
              /* Loading state */
              <>
                <p className="text-gray-400 mb-4">
                  Importing your characters from Battle.net...
                </p>
                <div className="text-sm text-gray-500 space-y-2">
                  <div className="flex items-center justify-center gap-2">
                    <div className="animate-pulse h-2 w-2 bg-purple-500 rounded-full"></div>
                    <span>Connecting to Battle.net</span>
                  </div>
                  <div className="flex items-center justify-center gap-2">
                    <div className="animate-pulse h-2 w-2 bg-purple-500 rounded-full"></div>
                    <span>Fetching character data</span>
                  </div>
                  <div className="flex items-center justify-center gap-2">
                    <div className="animate-pulse h-2 w-2 bg-purple-500 rounded-full"></div>
                    <span>Enriching character information</span>
                  </div>
                </div>
              </>
            ) : (
              /* Normal state */
              <>
                <p className="text-gray-400 mb-6">
                  Import your World of Warcraft characters from Battle.net to
                  get started.
                  {!linkStatus?.linked &&
                    " If your Battle.net account isn't linked, we'll guide you through the process."}
                </p>

                <Button
                  onClick={actions.syncAndEnrich}
                  disabled={ui.isDisabled.sync}
                  className="bg-purple-600 hover:bg-purple-700 min-w-[200px]"
                  size="lg"
                >
                  <UserPlus className="mr-2 h-5 w-5" />
                  Import My Characters
                </Button>
              </>
            )}
          </div>
        </div>
      ) : (
        /* PERSONNAGES TOUJOURS AFFICHÉS */
        <div className="bg-slate-800/30 border border-slate-700 rounded-lg">
          <Table>
            <TableHeader>
              <TableRow className="border-slate-700">
                <TableHead className="w-12"></TableHead>
                <TableHead>Name</TableHead>
                <TableHead className="hidden md:table-cell">Class</TableHead>
                <TableHead className="hidden md:table-cell">Realm</TableHead>
                <TableHead className="hidden md:table-cell">Level</TableHead>
                <TableHead className="hidden md:table-cell">
                  Last Update
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {characters.map((character: EnrichedUserCharacter) => {
                const normalizedClass = character.class.replace(/\s+/g, "");
                const classIcon = getClassIcon(normalizedClass);
                const avatarUrl = character.avatar_url || classIcon;

                return (
                  <TableRow
                    key={character.id}
                    className="border-slate-700 cursor-pointer hover:bg-slate-800/50"
                    onClick={() => handleCharacterClick(character)}
                  >
                    <TableCell>
                      <Image
                        src={avatarUrl}
                        alt={character.name}
                        width={32}
                        height={32}
                        className="rounded-full"
                        onError={(e) => {
                          e.currentTarget.src = classIcon;
                        }}
                      />
                    </TableCell>

                    <TableCell className="font-medium">
                      <div>
                        <div>{character.name}</div>
                        <div className="md:hidden text-sm text-slate-400">
                          {character.active_spec_name
                            ? `${character.active_spec_name} ${character.class}`
                            : character.class}
                        </div>
                      </div>
                    </TableCell>

                    <TableCell
                      className={`hidden md:table-cell ${getClassColor(
                        character.class
                      )}`}
                    >
                      {character.active_spec_name
                        ? `${character.active_spec_name} ${character.class}`
                        : character.class}
                    </TableCell>

                    <TableCell className="hidden md:table-cell">
                      {character.realm.charAt(0).toUpperCase() +
                        character.realm.slice(1)}{" "}
                      - {character.region.toUpperCase()}
                    </TableCell>

                    <TableCell className="hidden md:table-cell">
                      {character.level}
                    </TableCell>

                    <TableCell className="hidden md:table-cell">
                      {character.last_api_update
                        ? new Date(
                            character.last_api_update
                          ).toLocaleDateString()
                        : "Never"}
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </div>
      )}

      {/* Tips informatifs */}
      {!linkStatus?.linked && hasCharacters && (
        <div className="bg-amber-900/20 border border-amber-500/30 rounded-lg p-4">
          <div className="flex items-start gap-3">
            <AlertCircle className="h-4 w-4 mt-0.5 text-amber-400" />
            <div>
              <h4 className="font-semibold text-amber-400 mb-1">
                Battle.net Connection
              </h4>
              <p className="text-sm text-amber-300">
                Your Battle.net connection appears to be expired or not linked.
                You can still view your saved characters, but clicking
                &quot;Sync Characters&quot; will reconnect your account
                automatically.
              </p>
            </div>
          </div>
        </div>
      )}

      {!hasCharacters && (
        <div className="bg-slate-800/50 border border-slate-700 rounded-lg p-4">
          <div className="flex items-start gap-3">
            <AlertCircle className="h-4 w-4 mt-0.5 text-blue-400" />
            <div>
              <h4 className="font-semibold text-blue-400 mb-1">
                Getting Started
              </h4>
              <p className="text-sm text-slate-400">
                Import your characters to access detailed statistics, track your
                progress, and get personalized recommendations for your
                gameplay.
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default CharactersTab;
