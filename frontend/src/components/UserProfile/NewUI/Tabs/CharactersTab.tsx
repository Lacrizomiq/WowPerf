// CharactersTab.tsx
// Tab displaying all user characters
import React from "react";
import { Card } from "@/components/ui/card";
import { WoWProfile } from "../AccountProfile";

const CharactersTab: React.FC = () => {
  return (
    <Card className="bg-[#131e33] border-gray-800 p-6">
      <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          className="text-blue-500"
        >
          <path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l8.57-8.57A4 4 0 1 1 18 8.84l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48" />
        </svg>
        Your Characters
      </h2>

      <p className="text-gray-400 mb-6">
        Only characters from the same region as your Battle.net account will be
        displayed, and only level 80 characters will be shown.
      </p>

      <WoWProfile showTooltips={true} />
    </Card>
  );
};

export default CharactersTab;
