"use client";

import Header from "@/components/Header/Header";
import Hero from "@/components/Home/Hero";
import FeaturedContent from "@/components/Home/FeaturedContent";
import { useState } from "react";
import Sidebar from "@/components/Header/Sidebar";

export default function Home() {
  const [mainMargin, setMainMargin] = useState(64);
  return (
    <main className="bg-[#002440]">
      <Sidebar setMainMargin={setMainMargin} />
      <div
        className="flex-1 transition-all duration-300"
        style={{ marginLeft: `${mainMargin}px` }}
      >
        <Hero />
        <FeaturedContent />
      </div>
    </main>
  );
}
