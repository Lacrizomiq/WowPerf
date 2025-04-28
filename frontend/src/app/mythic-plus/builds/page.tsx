// # Home page with the selection of classes and specs
import ClassesGrid from "@/components/MythicPlus/BuildsAnalysis/home/ClassesGrid";
import { Metadata } from "next";

// Metadata for the page (for SEO)
export const metadata: Metadata = {
  title: "WoW Classes and Specs - Mythic+ Build Analyzer",
  description:
    "Discover the best combinations of talents, stats, and equipment for each class and specialization in Mythic+ in World of Warcraft.",
};

export default function BuildsPage() {
  return <ClassesGrid />;
}
