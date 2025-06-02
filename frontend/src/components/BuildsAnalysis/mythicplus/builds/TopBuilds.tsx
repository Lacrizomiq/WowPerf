import { TopTalentBuild } from "@/types/warcraftlogs/builds/buildsAnalysis";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import BuildCard from "./BuildCard";

interface TopBuildsProps {
  builds: TopTalentBuild[];
  className: WowClassParam;
  spec: WowSpecParam;
}

export default function TopBuilds({ builds, className, spec }: TopBuildsProps) {
  // We only display the top build (first in the array)
  const topBuild = builds[0];

  if (!topBuild) {
    return (
      <div className="bg-slate-800 rounded-lg p-5 text-center">
        <p className="text-slate-400">No builds available yet.</p>
      </div>
    );
  }

  return <BuildCard build={topBuild} className={className} spec={spec} />;
}
