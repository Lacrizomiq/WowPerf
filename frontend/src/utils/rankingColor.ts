/**
 * Returns the appropriate color code based on a performance percentage
 * @param value - Number between 0 and 100
 * @returns Hex color code as a string
 */
export const getPerformanceColor = (value: number): string => {
  if (value === 100) return "#e5cc80"; // Gold
  if (value >= 99) return "#e268a8"; // Purple
  if (value >= 95) return "#ff8000"; // Orange
  if (value >= 75) return "#a335ee"; // Purple
  if (value >= 50) return "#0070ff"; // Blue
  if (value >= 25) return "#1eff00"; // Green
  return "#666666"; // Gray
};

/**
 * Returns a tailwind text color class based on a performance percentage
 * @param value - Number between 0 and 100
 * @returns Tailwind text color class as a string
 */
export const getPerformanceColorClass = (value: number): string => {
  if (value === 100) return "text-[#e5cc80]";
  if (value >= 99) return "text-[#e268a8]";
  if (value >= 95) return "text-[#ff8000]";
  if (value >= 75) return "text-[#a335ee]";
  if (value >= 50) return "text-[#0070ff]";
  if (value >= 25) return "text-[#1eff00]";
  return "text-[#666666]";
};
