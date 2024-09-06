import React from "react";
import Image from "next/image";
import { Raid } from "@/types/raids";

interface StaticRaidsListProps {
  raids: Raid[];
}

const StaticRaidsList: React.FC<StaticRaidsListProps> = ({ raids }) => {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mt-4">
      {raids.map((raid) => (
        <div
          key={raid.ID}
          className="bg-deep-blue p-4 rounded-lg cursor-pointer hover:bg-blue-700 transition-colors duration-200 flex flex-col items-center"
        >
          <Image
            src={raid.MediaURL}
            alt={raid.Name}
            width={200}
            height={200}
            className="rounded-md mb-2"
          />
          <h3 className="font-bold text-lg mb-2 text-center">{raid.Name}</h3>
        </div>
      ))}
    </div>
  );
};

export default StaticRaidsList;
