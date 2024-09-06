import RaidOverview from "@/components/Raids/RaidOverview";

export default function RaidPage() {
  const initialExpansion = "DF";
  return <RaidOverview initialExpansion={initialExpansion} />;
}
