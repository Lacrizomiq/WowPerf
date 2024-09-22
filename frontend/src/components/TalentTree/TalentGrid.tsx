import React from "react";
import Image from "next/image";

interface TalentNode {
  id: number;
  name: string;
  posX: number;
  posY: number;
  entries: { icon: string; spellId: number }[];
  rank: number;
  maxRanks: number;
  next?: number[];
}

interface TalentGridProps {
  talents: TalentNode[];
  selectedTalents: TalentNode[];
}

const TalentGrid: React.FC<TalentGridProps> = ({
  talents,
  selectedTalents,
}) => {
  const cellSize = 6; // Taille en pourcentage
  const iconRadiusPercent = cellSize / 2;

  const minX = Math.min(...talents.map((t) => t.posX));
  const minY = Math.min(...talents.map((t) => t.posY));
  const maxX = Math.max(...talents.map((t) => t.posX));
  const maxY = Math.max(...talents.map((t) => t.posY));

  const gridStyle: React.CSSProperties = {
    position: "relative",
    width: "100%",
    paddingTop: `${((maxY - minY) / (maxX - minX)) * 100}%`,
    borderRadius: "8px",
    margin: "0 auto",
    overflow: "visible",
  };

  const connections = talents.flatMap(
    (talent) =>
      talent.next?.map((nextId) => ({
        from: talent.id,
        to: nextId,
      })) || []
  );

  return (
    <div style={gridStyle} className="talent-grid">
      <svg
        viewBox="0 0 100 100"
        preserveAspectRatio="none"
        style={{
          position: "absolute",
          top: 0,
          left: 0,
          width: "100%",
          height: "100%",
        }}
      >
        {connections.map((conn, index) => {
          const fromTalent = talents.find((t) => t.id === conn.from);
          const toTalent = talents.find((t) => t.id === conn.to);

          if (!fromTalent || !toTalent) return null;

          const fromX = ((fromTalent.posX - minX) / (maxX - minX)) * 100;
          const fromY = ((fromTalent.posY - minY) / (maxY - minY)) * 100;
          const toX = ((toTalent.posX - minX) / (maxX - minX)) * 100;
          const toY = ((toTalent.posY - minY) / (maxY - minY)) * 100;

          const deltaX = toX - fromX;
          const deltaY = toY - fromY;
          const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);

          const unitX = deltaX / distance;
          const unitY = deltaY / distance;

          const adjustedFromX = fromX + unitX * iconRadiusPercent;
          const adjustedFromY = fromY + unitY * iconRadiusPercent;
          const adjustedToX = toX - unitX * iconRadiusPercent;
          const adjustedToY = toY - unitY * iconRadiusPercent;

          return (
            <line
              key={index}
              x1={`${adjustedFromX}`}
              y1={`${adjustedFromY}`}
              x2={`${adjustedToX}`}
              y2={`${adjustedToY}`}
              stroke="white"
              strokeWidth="0.2"
              vectorEffect="non-scaling-stroke"
            />
          );
        })}
      </svg>
      {talents.map((talent) => {
        const selectedTalent = selectedTalents.find((t) => t.id === talent.id);
        return (
          <TalentIcon
            key={talent.id}
            talent={talent}
            cellSize={cellSize}
            minX={minX}
            minY={minY}
            maxX={maxX}
            maxY={maxY}
            isSelected={!!selectedTalent}
            selectedRank={selectedTalent?.rank || 0}
          />
        );
      })}
    </div>
  );
};

interface TalentIconProps {
  talent: TalentNode;
  cellSize: number;
  minX: number;
  minY: number;
  maxX: number;
  maxY: number;
  isSelected: boolean;
  selectedRank: number;
}

const TalentIcon: React.FC<TalentIconProps> = ({
  talent,
  cellSize,
  minX,
  minY,
  maxX,
  maxY,
  isSelected,
  selectedRank,
}) => {
  const [imageError, setImageError] = React.useState(false);

  const selectedEntry =
    isSelected && talent.entries.length > 1
      ? talent.entries[selectedRank - 1]
      : talent.entries[0];

  const normalizedPosX = (talent.posX - minX) / (maxX - minX);
  const normalizedPosY = (talent.posY - minY) / (maxY - minY);

  const iconStyle: React.CSSProperties = {
    position: "absolute",
    left: `calc(${normalizedPosX * 100}% - ${cellSize / 2}%)`,
    top: `calc(${normalizedPosY * 100}% - ${cellSize / 2}%)`,
    width: `${cellSize}%`,
    height: `${cellSize}%`,
  };

  return (
    <div
      className={`talent-icon ${isSelected ? "selected" : "unselected"}`}
      style={iconStyle}
    >
      <div className="relative" style={{ width: "90%", height: "90%" }}>
        <Image
          src={
            imageError
              ? "https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg"
              : `https://wow.zamimg.com/images/wow/icons/large/${selectedEntry.icon}.jpg`
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
            {selectedRank}/{talent.maxRanks}
          </div>
        )}
      </div>
    </div>
  );
};

export default TalentGrid;
