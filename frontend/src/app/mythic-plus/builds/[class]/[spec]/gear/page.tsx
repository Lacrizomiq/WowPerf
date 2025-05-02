// app/mythic-plus/builds/[class]/[spec]/gear/page.tsx
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import GearContent from "@/components/MythicPlus/BuildsAnalysis/gear/GearContent";

export default function BuildGearPage({
  params,
}: {
  params: { class: string; spec: string };
}) {
  const className = params.class as WowClassParam;
  const spec = params.spec as WowSpecParam;

  return <GearContent className={className} spec={spec} />;
}
