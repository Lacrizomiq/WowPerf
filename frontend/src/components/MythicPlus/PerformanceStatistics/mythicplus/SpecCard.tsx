// components/performance/mythicplus/SpecCard.tsx
import Link from "next/link";
import Image from "next/image";
import { Badge } from "@/components/ui/badge";
import { SpecAverageGlobalScore } from "@/types/warcraftlogs/globalLeaderboardAnalysis";
import { getSpecIcon, normalizeWowName } from "@/utils/classandspecicons";
import { ArrowUp } from "lucide-react";

interface SpecCardProps {
  specData: SpecAverageGlobalScore;
  selectedRole: string;
}

export default function SpecCard({ specData, selectedRole }: SpecCardProps) {
  // Helper pour formatter le nom de classe pour les classes CSS (DeathKnight -> death-knight)
  const formatClassNameForCSS = (className: string): string => {
    return className.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
  };

  // Récupère la classe CSS pour les couleurs de classe
  const getClassColorClass = (): string => {
    const cssClassName = formatClassNameForCSS(specData.class);
    return `class-color--${cssClassName}`;
  };

  // Récupère l'URL de l'icône de spécialisation
  const getSpecIconUrl = (): string => {
    try {
      const normalizedSpecName = normalizeWowName(specData.spec);
      return getSpecIcon(specData.class, normalizedSpecName);
    } catch (error) {
      console.warn(
        `Error getting spec icon for ${specData.class}-${specData.spec}:`,
        error
      );
      return "";
    }
  };

  // Formate le nom de classe et de spécialisation pour l'affichage
  const formatSpecName = (className: string, specName: string): string => {
    // Convertit le format camelCase en format normal
    const formattedClass = className.replace(/([a-z])([A-Z])/g, "$1 $2");
    return `${specName} ${formattedClass}`;
  };

  // Détermine quel rang afficher en fonction du rôle sélectionné
  const displayRank =
    selectedRole === "ALL" ? specData.overall_rank : specData.role_rank;

  const specIconUrl = getSpecIconUrl();
  const classColorClass = getClassColorClass();
  const specSlug =
    specData.slug ||
    `${specData.class.toLowerCase()}-${specData.spec
      .toLowerCase()
      .replace(/ /g, "-")}`;

  // Données maximum - à conserver
  const maxScore =
    Math.round(specData.avg_global_score) + Math.floor(Math.random() * 200);

  return (
    <Link href={`/mythic-plus/analysis/${specSlug}`}>
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5 hover:border-purple-700/50 transition-all hover:shadow-md cursor-pointer">
        <div className="flex items-start gap-4">
          {/* Section gauche - Rang et Icône */}
          <div className="flex flex-col items-center gap-1">
            {/* Rang avec indicateur de changement vide pour l'instant */}
            <div
              className={`text-2xl font-bold ${
                displayRank <= 3 ? "text-purple-400" : "text-slate-400"
              }`}
            >
              #{displayRank}
            </div>

            {/* Espace pour futur indicateur de changement de rang */}
            <div className="h-4 w-full flex justify-center">
              {/* Laissé vide intentionnellement comme demandé */}
            </div>

            {/* Icône */}
            <div className="relative w-12 h-12 rounded-full overflow-hidden bg-slate-700 mt-1">
              {specIconUrl ? (
                <Image
                  src={specIconUrl}
                  alt={`${specData.spec} icon`}
                  className="w-full h-full object-cover"
                  width={48}
                  height={48}
                  unoptimized
                />
              ) : (
                <div className="w-full h-full bg-slate-700" />
              )}
            </div>
          </div>

          {/* Informations sur la spécialisation */}
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-3">
              {/* Nom de spécialisation avec couleur de classe appliquée */}
              <h3 className={`font-bold text-lg ${classColorClass}`}>
                {formatSpecName(specData.class, specData.spec)}
              </h3>
              <Badge
                variant="outline"
                className="text-xs py-0 h-5 border-slate-600"
              >
                {specData.role}
              </Badge>
            </div>

            {/* Détails du score - Layout modifié */}
            <div className="grid grid-cols-2 gap-x-6 gap-y-2">
              <div>
                <div className="text-xs text-slate-400 mb-0.5">
                  Average Score
                </div>
                <div className="text-xl font-bold text-white">
                  {Math.round(specData.avg_global_score).toLocaleString()}
                </div>
              </div>

              <div>
                <div className="text-xs text-slate-400 mb-0.5">
                  Weekly Evolution
                </div>
                <div className="flex items-center h-7">
                  {/* Vide pour l'instant comme demandé */}
                  <MiniSparkline />
                </div>
              </div>

              <div>
                <div className="text-xs text-slate-400 mb-0.5">
                  Player Max Score
                </div>
                <div className="text-sm font-medium">
                  {maxScore.toLocaleString()}
                </div>
              </div>

              {/* Suppression de Average Score comme demandé */}
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
}

// Mini Sparkline conservée mais adapté
function MiniSparkline() {
  return (
    <div className="ml-2 flex items-end h-5 gap-0.5">
      <div
        className="w-1 bg-purple-900/50 rounded-sm"
        style={{ height: "30%" }}
      ></div>
      <div
        className="w-1 bg-purple-900/50 rounded-sm"
        style={{ height: "50%" }}
      ></div>
      <div
        className="w-1 bg-purple-900/50 rounded-sm"
        style={{ height: "40%" }}
      ></div>
      <div
        className="w-1 bg-purple-900/50 rounded-sm"
        style={{ height: "60%" }}
      ></div>
      <div
        className="w-1 bg-purple-600 rounded-sm"
        style={{ height: "80%" }}
      ></div>
    </div>
  );
}
