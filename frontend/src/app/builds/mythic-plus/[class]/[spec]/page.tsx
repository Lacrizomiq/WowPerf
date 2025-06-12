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

export const generateMetadata = ({
  params,
}: {
  params: { class: string; spec: string };
}): Metadata => {
  const className = params.class as WowClassParam;
  const spec = params.spec as WowSpecParam;

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

export default function BuildPage({
  params,
}: {
  params: { class: string; spec: string };
}) {
  const className = params.class as WowClassParam;
  const spec = params.spec as WowSpecParam;

  return <BuildsContent className={className} spec={spec} />;
}
