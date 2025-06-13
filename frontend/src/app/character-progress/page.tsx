import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Construction } from "lucide-react";
import Link from "next/link";

function CharacterProgressPage() {
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <h1 className="text-xl font-bold mb-4 text-center">
        <div className="flex items-center">
          <Construction className="mr-2" /> Character Progress Page is still
          under construction
          <Construction className="ml-2" />
        </div>
      </h1>
      <div className="flex flex-col items-center justify-center gap-4">
        <Badge
          variant="outline"
          className="border-purple-600 text-purple-400 text-center"
        >
          Coming Soon...
        </Badge>
      </div>
    </div>
  );
}

export default CharacterProgressPage;
