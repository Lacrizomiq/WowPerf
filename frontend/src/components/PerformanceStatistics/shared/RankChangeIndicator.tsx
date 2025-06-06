// components/performance/shared/RankChangeIndicator.tsx
import { ArrowUp, ArrowDown } from "lucide-react";

interface RankChangeIndicatorProps {
  change: number;
}

export default function RankChangeIndicator({
  change,
}: RankChangeIndicatorProps) {
  if (change > 0) {
    return (
      <span className="ml-2 text-green-400 text-sm flex items-center">
        <ArrowUp className="h-3 w-3 mr-0.5" />
        {change}
      </span>
    );
  } else if (change < 0) {
    return (
      <span className="ml-2 text-red-400 text-sm flex items-center">
        <ArrowDown className="h-3 w-3 mr-0.5" />
        {Math.abs(change)}
      </span>
    );
  } else {
    return <span className="ml-2 text-slate-400 text-sm">=</span>;
  }
}
