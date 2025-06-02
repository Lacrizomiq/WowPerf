"use client";

import { useRouter, usePathname, useSearchParams } from "next/navigation";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";

type ContentType = "mythic-plus" | "raids" | "pvp";

export default function ContentTypeTabs({
  className,
  spec,
  activeTab = "mythic-plus",
}: {
  className: WowClassParam;
  spec: WowSpecParam;
  activeTab: ContentType;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  // Determine the current page/tab within the content type
  let currentPage = "";
  if (pathname.includes("/talents")) {
    currentPage = "/talents";
  } else if (pathname.includes("/gear")) {
    currentPage = "/gear";
  } else if (pathname.includes("/enchants-gems")) {
    currentPage = "/enchants-gems";
  }

  // Preserve current search params
  const currentSearchParams = searchParams.toString();
  const queryString = currentSearchParams ? `?${currentSearchParams}` : "";

  // Handle tab change
  const handleTabChange = (value: string) => {
    let newPath = "";

    if (value === "mythic-plus") {
      newPath = `/builds/mythic-plus/${className}/${spec}${currentPage}${queryString}`;
    } else if (value === "raids") {
      newPath = `/builds/raids/${className}/${spec}${currentPage}${queryString}`;
    } else if (value === "pvp") {
      newPath = `/builds/pvp/${className}/${spec}${currentPage}${queryString}`;
    }

    if (newPath) {
      router.push(newPath);
    }
  };

  return (
    <Tabs value={activeTab} onValueChange={handleTabChange} className="w-full">
      <TabsList className="grid w-full grid-cols-3 bg-slate-800/50 mb-6 ">
        <TabsTrigger
          value="mythic-plus"
          className="data-[state=active]:bg-purple-600"
        >
          Mythic+
        </TabsTrigger>
        <TabsTrigger
          value="raids"
          className="data-[state=active]:bg-purple-600 relative"
          disabled
        >
          Raids
          <Badge className="absolute -top-2 -right-2 bg-purple-600 text-[10px]">
            Soon
          </Badge>
        </TabsTrigger>
        <TabsTrigger
          value="pvp"
          className="data-[state=active]:bg-purple-600 relative"
          disabled
        >
          PvP
          <Badge className="absolute -top-2 -right-2 bg-purple-600 text-[10px]">
            Soon
          </Badge>
        </TabsTrigger>
      </TabsList>
    </Tabs>
  );
}
