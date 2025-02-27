"use client";

import React, { useMemo } from "react";
import {
  useGetMaxKeyLevelPerSpecAndDungeon,
  useGetDungeonMedia,
} from "@/hooks/useWarcraftLogsApi";

const DungeonPerformance: React.FC<{
  className: string;
  specName: string;
}> = ({ className, specName }) => {
  const {
    data: dungeonData,
    isLoading,
    error,
  } = useGetMaxKeyLevelPerSpecAndDungeon();

  const { data: dungeonMedia, isLoading: isLoadingDungeonMedia } =
    useGetDungeonMedia();

  // Map the media by slug
  const mediaBySlug = useMemo(() => {
    if (!dungeonMedia) return {};

    return dungeonMedia.reduce((acc, media) => {
      acc[media.dungeon_slug] = media;
      return acc;
    }, {} as Record<string, (typeof dungeonMedia)[0]>);
  }, [dungeonMedia]);

  // Filtrer class and spec data
  const specDungeonData = useMemo(() => {
    if (!dungeonData) return [];
    return dungeonData.filter(
      (dungeon) =>
        dungeon.class.toLowerCase() === className.toLowerCase() &&
        dungeon.spec.toLowerCase() === specName.toLowerCase()
    );
  }, [dungeonData, className, specName]);

  if (isLoading || isLoadingDungeonMedia) {
    return (
      <div className="flex items-center justify-center min-h-[200px]">
        <div className="text-white text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-white mx-auto"></div>
          <p className="mt-4">Loading dungeon performance...</p>
        </div>
      </div>
    );
  }

  if (error || !specDungeonData.length) {
    return (
      <div className="flex items-center justify-center min-h-[200px]">
        <p className="text-white text-center">
          No dungeon performance data available.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {specDungeonData.map((dungeon, index) => (
        <div
          key={index}
          className="card p-4 relative overflow-hidden rounded-lg"
          style={{
            backgroundColor: "#112240",
            backgroundImage: mediaBySlug[dungeon.dungeon_slug]?.media_url
              ? `url('${mediaBySlug[dungeon.dungeon_slug].media_url}')`
              : undefined,
            backgroundSize: "cover",
            backgroundPosition: "center",
            backgroundRepeat: "no-repeat",
          }}
        >
          <div className="absolute inset-0 bg-black/50" /> {/* Dark overlay */}
          <div className="relative z-10">
            {" "}
            {/* Content wrapper */}
            <div className="flex justify-between items-center mb-2">
              <h3 className="font-bold text-white">{dungeon.dungeon_name}</h3>
              <span className="text-lg font-bold text-white">
                +{dungeon.max_key_level}
              </span>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default DungeonPerformance;
