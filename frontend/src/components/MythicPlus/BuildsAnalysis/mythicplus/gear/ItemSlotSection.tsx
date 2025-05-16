// components/builds/gear/ItemSlotSection.tsx
import {
  GlobalPopularItem,
  PopularItem,
} from "@/types/warcraftlogs/builds/buildsAnalysis";
import ItemCard from "./ItemCard";

interface ItemSlotSectionProps {
  slotId: number;
  slotName: string;
  items: PopularItem[] | GlobalPopularItem[];
}

export default function ItemSlotSection({
  slotId,
  slotName,
  items,
}: ItemSlotSectionProps) {
  // Check if we have items to display
  if (!items || items.length === 0) {
    return null;
  }

  // Sort items by rank (most popular first)
  const sortedItems = [...items].sort((a, b) => a.rank - b.rank);

  return (
    <div className="bg-slate-900 rounded-lg border border-slate-800 p-4">
      <h3 className="text-xl font-bold text-white mb-4">{slotName}</h3>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Limit to 4 items maximum per slot */}
        {sortedItems.slice(0, 4).map((item) => (
          <ItemCard key={`${item.item_id}-${item.rank}`} item={item} />
        ))}
      </div>
    </div>
  );
}
