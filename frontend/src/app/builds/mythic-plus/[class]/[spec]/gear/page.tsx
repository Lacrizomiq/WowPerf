// app/mythic-plus/builds/[class]/[spec]/gear/page.tsx
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import GearContent from "@/components/BuildsAnalysis/mythicplus/gear/GearContent";

export default async function BuildGearPage({
  params,
}: {
  params: Promise<{ class: string; spec: string }>;
}) {
  const resolvedParams = await params;
  const className = resolvedParams.class as WowClassParam;
  const spec = resolvedParams.spec as WowSpecParam;

  return <GearContent className={className} spec={spec} />;
}
