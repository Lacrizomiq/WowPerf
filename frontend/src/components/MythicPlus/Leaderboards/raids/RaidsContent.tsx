// components/leaderboards/raids/RaidsContent.tsx
import { Badge } from "@/components/ui/badge";

export default function RaidsContent() {
  return (
    <div className="flex flex-col items-center justify-center py-20 text-center">
      <Badge className="mb-4 bg-purple-600 text-white px-3 py-1">
        Coming Soon
      </Badge>
      <h3 className="text-2xl font-bold mb-2">Raid Leaderboards</h3>
      <p className="text-muted-foreground max-w-md">
        We&apos;re working hard to bring you detailed Raid leaderboards. Check
        back soon!
      </p>
    </div>
  );
}
