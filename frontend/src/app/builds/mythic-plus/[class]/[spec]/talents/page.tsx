import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import TalentsContent from "@/components/BuildsAnalysis/mythicplus/talents/TalentsContent";

export default function TalentsPage({
  params,
}: {
  params: { class: string; spec: string };
}) {
  const className = params.class as WowClassParam;
  const spec = params.spec as WowSpecParam;

  return <TalentsContent className={className} spec={spec} />;
}
