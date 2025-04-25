import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ReactNode } from "react";

interface BuildNavProps {
  defaultTab?: string;
  children: ReactNode;
}

export default function BuildNav({
  defaultTab = "builds",
  children,
}: BuildNavProps) {
  return (
    <Tabs defaultValue={defaultTab} className="mt-6">
      <TabsList className="bg-slate-800 p-0 mb-6 border-b border-slate-700 w-full flex justify-start rounded-none">
        <TabsTrigger
          value="builds"
          className="py-3 px-6 rounded-none data-[state=active]:bg-indigo-600 data-[state=active]:text-white"
        >
          Builds
        </TabsTrigger>
        <TabsTrigger
          value="talents"
          className="py-3 px-6 rounded-none data-[state=active]:bg-indigo-600 data-[state=active]:text-white"
        >
          Talents
        </TabsTrigger>
        <TabsTrigger
          value="gear"
          className="py-3 px-6 rounded-none data-[state=active]:bg-indigo-600 data-[state=active]:text-white"
        >
          Gear
        </TabsTrigger>
        <TabsTrigger
          value="enchants"
          className="py-3 px-6 rounded-none data-[state=active]:bg-indigo-600 data-[state=active]:text-white"
        >
          Enchants & Gems
        </TabsTrigger>
      </TabsList>

      {children}
    </Tabs>
  );
}
