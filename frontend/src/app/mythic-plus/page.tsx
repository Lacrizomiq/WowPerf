"use client";

import React from "react";
import MythicPlusBestRuns from "@/components/Home/mythicplus/MythicBestRuns";
import Sidebar from "@/components/Header/Sidebar";
import { useState } from "react";

const MythicPlusPage = () => {
  const [mainMargin, setMainMargin] = useState(64);
  return (
    <main className="bg-[#002440]">
      <Sidebar setMainMargin={setMainMargin} />
      <div
        className="flex-1 transition-all duration-300"
        style={{ marginLeft: `${mainMargin}px` }}
      >
        <MythicPlusBestRuns />
        {/* <FeaturedContent /> */}
      </div>
    </main>
  );
};

export default MythicPlusPage;
