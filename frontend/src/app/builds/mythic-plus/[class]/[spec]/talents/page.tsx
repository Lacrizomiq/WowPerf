import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import TalentsContent from "@/components/BuildsAnalysis/mythicplus/talents/TalentsContent";

export default async function TalentsPage({
  params,
}: {
  params: Promise<{ class: string; spec: string }>;
}) {
  const resolvedParams = await params;
  const className = resolvedParams.class as WowClassParam;
  const spec = resolvedParams.spec as WowSpecParam;

  return <TalentsContent className={className} spec={spec} />;
}
