import React from "react";
import Image from "next/image";
import { HeroTalent } from "@/types/talents";

interface TalentGridProps {
  selectedHeroTalentTree: HeroTalent[];
}

const HeroTalentGrid: React.FC<TalentGridProps> = ({
  selectedHeroTalentTree,
}) => {
  const cellSize = 6; // Taille en pourcentage
  const minX = Math.min(...selectedHeroTalentTree.map((t) => t.posX));
  const minY = Math.min(...selectedHeroTalentTree.map((t) => t.posY));
  const maxX = Math.max(...selectedHeroTalentTree.map((t) => t.posX));
  const maxY = Math.max(...selectedHeroTalentTree.map((t) => t.posY));

  const gridStyle: React.CSSProperties = {
    position: "relative",
    width: "80%",
    paddingTop: `${((maxY - minY) / (maxX - minX)) * 50}%`,
    borderRadius: "8px",
    margin: "0 auto",
    overflow: "visible",
  };

  return (
    <div style={gridStyle} className="talent-grid">
      {selectedHeroTalentTree.map((talent) => (
        <TalentIcon
          key={talent.id}
          talent={talent}
          cellSize={cellSize}
          minX={minX}
          minY={minY}
          maxX={maxX}
          maxY={maxY}
          isSelected={true}
        />
      ))}
    </div>
  );
};

interface TalentIconProps {
  talent: HeroTalent;
  cellSize: number;
  minX: number;
  minY: number;
  maxX: number;
  maxY: number;
  isSelected: boolean;
}

const TalentIcon: React.FC<TalentIconProps> = ({
  talent,
  cellSize,
  minX,
  minY,
  maxX,
  maxY,
  isSelected,
}) => {
  const [imageError, setImageError] = React.useState(false);

  const normalizedPosX = (talent.posX - minX) / (maxX - minX);
  const normalizedPosY = (talent.posY - minY) / (maxY - minY);

  const iconStyle: React.CSSProperties = {
    position: "absolute",
    left: `calc(${normalizedPosX * 100}% - ${cellSize / 2}%)`,
    top: `calc(${normalizedPosY * 100}% - ${cellSize / 2}%)`,
    width: `36px`,
    height: `36px`,
  };

  // Utiliser le talent sélectionné ou le premier si c'est un choix multiple
  const selectedEntry =
    talent.entries.find((entry) => entry.id === talent.id) || talent.entries[0];

  return (
    <div
      className={`talent-icon ${isSelected ? "selected" : "unselected"}`}
      style={iconStyle}
    >
      <div className="relative" style={{ width: "150%", height: "150%" }}>
        <Image
          src={
            imageError
              ? "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg"
              : `https://wow.zamimg.com/images/wow/icons/large/${selectedEntry.icon}.jpg`
          }
          alt={talent.name}
          layout="fill"
          objectFit="cover"
          className={`rounded-full border-2 ${
            isSelected
              ? "border-yellow-400 glow-effect"
              : "border-gray-700 opacity-50"
          }`}
          onError={() => setImageError(true)}
        />
        {isSelected && (
          <div className="absolute bottom-0 right-0 bg-black bg-opacity-70 text-white text-[8px] font-bold px-1 rounded-full">
            {talent.rank}/{selectedEntry.maxRanks}
          </div>
        )}
      </div>
    </div>
  );
};

export default HeroTalentGrid;
