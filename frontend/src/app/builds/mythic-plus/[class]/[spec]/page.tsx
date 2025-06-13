import { Metadata } from "next";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import BuildsContent from "@/components/BuildsAnalysis/mythicplus/builds/BuildsContent";
import {
  formatDisplayClassName,
  formatDisplaySpecName,
} from "@/utils/classandspecicons";

export const generateMetadata = async ({
  params,
}: {
  params: Promise<{ class: string; spec: string }>;
}): Promise<Metadata> => {
  const resolvedParams = await params;
  const className = resolvedParams.class as WowClassParam;
  const spec = resolvedParams.spec as WowSpecParam;

  const displayClassName = formatDisplayClassName(className);
  const displaySpecName = formatDisplaySpecName(spec);

  return {
    title: `${displaySpecName} ${displayClassName} Builds - Mythic+`,
    description: `Discover the best talents, gear, enchants, and stat priorities for ${displaySpecName} ${displayClassName} in Mythic+.`,
    keywords: [
      `wow builds`,
      `${displaySpecName}`,
      `${displayClassName}`,
      `mythic+`,
      `wow talents`,
      `wow gear`,
    ],
    openGraph: {
      title: `${displaySpecName} ${displayClassName} Builds - Mythic+`,
      description: `Discover the best talents, gear, enchants, and stat priorities for ${displaySpecName} ${displayClassName} in Mythic+.`,
      type: "website",
    },
  };
};

export default async function BuildPage({
  params,
}: {
  params: Promise<{ class: string; spec: string }>;
}) {
  const resolvedParams = await params;
  const className = resolvedParams.class as WowClassParam;
  const spec = resolvedParams.spec as WowSpecParam;

  return <BuildsContent className={className} spec={spec} />;
}
