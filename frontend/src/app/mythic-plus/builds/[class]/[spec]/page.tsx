import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import BuildsContent from "@/components/MythicPlus/BuildsAnalysis/builds/BuildsContent";

export default function BuildPage({
  params,
}: {
  params: { class: string; spec: string };
}) {
  const className = params.class as WowClassParam;
  const spec = params.spec as WowSpecParam;

  return <BuildsContent className={className} spec={spec} />;
}
