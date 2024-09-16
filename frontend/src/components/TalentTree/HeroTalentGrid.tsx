import React from "react";
import Image from "next/image";
import { HeroTalent } from "@/types/talents";

interface TalentGridProps {
  talents: HeroTalent[];
  selectedHeroTalentTree: HeroTalent[];
}

const HeroTalentGrid: React.FC<TalentGridProps> = ({
  talents,
  selectedHeroTalentTree,
}) => {
  const cellSize = 5; // Taille en pourcentage
  const minX = Math.min(...talents.map((t) => t.posX));
  const minY = Math.min(...talents.map((t) => t.posY));
  const maxX = Math.max(...talents.map((t) => t.posX));
  const maxY = Math.max(...talents.map((t) => t.posY));

  const aspectRatio = (maxY - minY) / (maxX - minX);

  const gridStyle: React.CSSProperties = {
    position: "relative",
    width: "100%",
    paddingTop: `${aspectRatio * 100}%`,
    borderRadius: "8px",
    margin: "0 auto",
    overflow: "visible",
  };

  return (
    <div style={gridStyle} className="talent-grid">
      {talents.map((talent) => (
        <TalentIcon
          key={talent.id}
          talent={talent}
          cellSize={cellSize}
          minX={minX}
          minY={minY}
          maxX={maxX}
          maxY={maxY}
          isSelected={selectedHeroTalentTree.some((t) => t.id === talent.id)}
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

  const talentEntry = talent.entries[0];

  return (
    <div
      className={`talent-icon ${isSelected ? "selected" : "unselected"}`}
      style={iconStyle}
    >
      <div className="relative" style={{ width: "100%", height: "100%" }}>
        <Image
          src={
            imageError
              ? "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg"
              : `https://wow.zamimg.com/images/wow/icons/large/${talentEntry.icon}.jpg`
          }
          alt={talent.name}
          layout="fill"
          objectFit="contain"
          className={`rounded-full border-2 ${
            isSelected
              ? "border-yellow-400 glow-effect"
              : "border-gray-700 opacity-50"
          }`}
          onError={() => setImageError(true)}
        />
        {isSelected && (
          <div className="absolute bottom-0 right-0 bg-black bg-opacity-70 text-white text-[8px] font-bold px-1 rounded-full">
            {talent.rank}/{talentEntry.maxRanks}
          </div>
        )}
      </div>
    </div>
  );
};

export default HeroTalentGrid;
