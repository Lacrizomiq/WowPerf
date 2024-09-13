import React from "react";
import Image from "next/image";
import { TalentNode } from "@/types/talents";

interface TalentGridProps {
  talents: TalentNode[];
  selectedTalents: TalentNode[];
}
const TalentGrid: React.FC<TalentGridProps> = ({
  talents,
  selectedTalents,
}) => {
  const cellSize = 40;
  const padding = 10;
  const scaleFactor = 0.1; // Adjust this value to scale the grid

  const minX = Math.min(...talents.map((t) => t.posX));
  const minY = Math.min(...talents.map((t) => t.posY));
  const maxX = Math.max(...talents.map((t) => t.posX));
  const maxY = Math.max(...talents.map((t) => t.posY));

  const gridWidth = (maxX - minX) * scaleFactor + cellSize + padding * 2;
  const gridHeight = (maxY - minY) * scaleFactor + cellSize + padding * 2;

  const gridStyle: React.CSSProperties = {
    position: "relative",
    width: `${gridWidth}px`,
    height: `${gridHeight}px`,
    backgroundColor: "rgba(0, 0, 0, 0.5)",
    borderRadius: "8px",
    padding: `${padding}px`,
    margin: "0 auto",
  };

  return (
    <div style={gridStyle} className="p-20 talent-grid">
      {talents.map((talent) => (
        <TalentIcon
          key={talent.id}
          talent={talent}
          cellSize={cellSize}
          minX={minX}
          minY={minY}
          scaleFactor={scaleFactor}
          padding={padding}
          isSelected={selectedTalents.some((t) => t.id === talent.id)}
        />
      ))}
    </div>
  );
};
interface TalentIconProps {
  talent: TalentNode;
  cellSize: number;
  minX: number;
  minY: number;
  scaleFactor: number;
  padding: number;
  isSelected: boolean;
}

const TalentIcon: React.FC<TalentIconProps> = ({
  talent,
  cellSize,
  minX,
  minY,
  scaleFactor,
  padding,
  isSelected,
}) => {
  const [imageError, setImageError] = React.useState(false);

  const iconStyle: React.CSSProperties = {
    position: "absolute",
    left: `${(talent.posX - minX) * scaleFactor + padding}px`,
    top: `${(talent.posY - minY) * scaleFactor + padding}px`,
    width: `${cellSize}px`,
    height: `${cellSize}px`,
  };

  return (
    <div
      className={`talent-icon ${isSelected ? "selected" : "unselected"}`}
      style={iconStyle}
    >
      <div className={`relative ${isSelected ? "" : ""}`}>
        <Image
          src={
            imageError
              ? "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg"
              : `https://wow.zamimg.com/images/wow/icons/large/${talent.entries[0].icon}.jpg`
          }
          alt={talent.name}
          width={cellSize}
          height={cellSize}
          className={`rounded-full border-2  ${
            isSelected
              ? "border-yellow-400 glow-effect"
              : "border-gray-700 opacity-50"
          }`}
          onError={() => setImageError(true)}
        />
        {isSelected && (
          <div className="absolute bottom-0 right-0 bg-black bg-opacity-70 text-white text-[8px] font-bold px-1 rounded-full">
            {talent.rank}/{talent.maxRanks}
          </div>
        )}
      </div>
    </div>
  );
};

export default TalentGrid;
