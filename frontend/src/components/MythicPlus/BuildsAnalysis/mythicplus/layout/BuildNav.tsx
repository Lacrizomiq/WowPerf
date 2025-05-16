// BuildNav.tsx - Version mise à jour avec correction de la surbrillance
import { TabsContent, TabsList } from "@/components/ui/tabs";
import { ReactNode } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";

interface BuildNavProps {
  defaultTab?: string;
  children: ReactNode;
  className: WowClassParam;
  spec: WowSpecParam;
}

export default function BuildNav({
  defaultTab = "builds",
  children,
  className,
  spec,
}: BuildNavProps) {
  const pathname = usePathname();

  // Determine which tab is active based on the pathname
  const isActive = (tab: string) => {
    if (tab === "builds") {
      // La page builds est active si le chemin ne contient pas talents, gear ou enchants-gems
      return !(
        pathname.includes("/talents") ||
        pathname.includes("/gear") ||
        pathname.includes("/enchants-gems")
      );
    }

    // Pour les autres onglets, vérifier si le chemin se termine par cet onglet
    return pathname.endsWith(`/${tab}`) || pathname.includes(`/${tab}/`);
  };

  // CSS classes for the tabs
  const tabClass = "py-3 px-6 rounded-none transition-colors duration-200";
  const activeTabClass = "bg-purple-600 text-white"; // Couleur violette pour l'onglet actif
  const inactiveTabClass = "hover:bg-slate-700 text-slate-300"; // Hover plus clair, texte plus visible

  return (
    <div className="mt-6">
      <div className="bg-slate-800/30 p-0 mb-6 border-b border-slate-700 w-full flex justify-start rounded-none">
        <Link href={`/builds/mythic-plus/${className}/${spec}`}>
          <div
            className={`${tabClass} ${
              isActive("builds") ? activeTabClass : inactiveTabClass
            }`}
          >
            Builds
          </div>
        </Link>

        <Link href={`/builds/mythic-plus/${className}/${spec}/talents`}>
          <div
            className={`${tabClass} ${
              isActive("talents") ? activeTabClass : inactiveTabClass
            }`}
          >
            Talents
          </div>
        </Link>

        <Link href={`/builds/mythic-plus/${className}/${spec}/gear`}>
          <div
            className={`${tabClass} ${
              isActive("gear") ? activeTabClass : inactiveTabClass
            }`}
          >
            Gear
          </div>
        </Link>

        <Link href={`/builds/mythic-plus/${className}/${spec}/enchants-gems`}>
          <div
            className={`${tabClass} ${
              isActive("enchants-gems") ? activeTabClass : inactiveTabClass
            }`}
          >
            Enchants & Gems
          </div>
        </Link>
      </div>

      {children}
    </div>
  );
}
