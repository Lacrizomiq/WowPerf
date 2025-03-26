// RaidbotTalentTree.tsx
import React from "react";

interface RaidbotsTreeProps {
  encodedString: string;
  width?: number;
  level?: number;
  locale?: string;
  hideHeader?: boolean;
  hideExport?: boolean;
  className?: string;
}

const RaidbotsTreeDisplay: React.FC<RaidbotsTreeProps> = ({
  encodedString,
  width = 1100,
  locale = "en_US",
  hideHeader = false,
  hideExport = false,
  className = "",
}) => {
  if (!encodedString) {
    return <div className="text-red-500">No talent data available</div>;
  }

  const params = new URLSearchParams({
    width: width.toString(),
    locale: locale,
  });

  if (hideHeader) params.append("hideHeader", "true");
  if (hideExport) params.append("hideExport", "true");

  const iframeUrl = `https://www.raidbots.com/simbot/render/talents/${encodeURIComponent(
    encodedString
  )}?${params}`;

  return (
    <div className={`relative ${className}`}>
      <div className="aspect-video w-full">
        <iframe
          src={iframeUrl}
          className="w-full h-full rounded-lg bg-transparent"
          allowFullScreen
        />
      </div>
    </div>
  );
};

export default RaidbotsTreeDisplay;
