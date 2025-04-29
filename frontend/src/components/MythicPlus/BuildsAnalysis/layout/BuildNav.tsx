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
      // The Builds tab is active only on the base route (without subpath)
      return pathname === `/mythic-plus/builds/${className}/${spec}`;
    }

    // For other tabs, check if the pathname ends with this tab
    return pathname.endsWith(`/${tab}`) || pathname.includes(`/${tab}/`);
  };

  // CSS classes for the tabs
  const tabClass = "py-3 px-6 rounded-none";
  const activeTabClass = "bg-indigo-600 text-white";
  const inactiveTabClass = "hover:bg-slate-700";

  return (
    <div className="mt-6">
      <div className="bg-slate-800 p-0 mb-6 border-b border-slate-700 w-full flex justify-start rounded-none">
        <Link href={`/mythic-plus/builds/${className}/${spec}`}>
          <div
            className={`${tabClass} ${
              isActive("builds") ? activeTabClass : inactiveTabClass
            }`}
          >
            Builds
          </div>
        </Link>

        <Link href={`/mythic-plus/builds/${className}/${spec}/talents`}>
          <div
            className={`${tabClass} ${
              isActive("talents") ? activeTabClass : inactiveTabClass
            }`}
          >
            Talents
          </div>
        </Link>

        <Link href={`/mythic-plus/builds/${className}/${spec}/gear`}>
          <div
            className={`${tabClass} ${
              isActive("gear") ? activeTabClass : inactiveTabClass
            }`}
          >
            Gear
          </div>
        </Link>

        <Link href={`/mythic-plus/builds/${className}/${spec}/enchants-gems`}>
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
