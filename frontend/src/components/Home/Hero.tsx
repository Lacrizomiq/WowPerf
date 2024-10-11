"use client";

import Image from "next/image";
import { useState, useEffect } from "react";
import SearchBar from "@/components/Home/SearchBar";

export default function Hero() {
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  return (
    <div className="relative min-h-screen w-full overflow-hidden">
      {isMounted && (
        <Image
          src="/homepage.avif"
          alt="World of Warcraft Castle"
          layout="fill"
          objectFit="cover"
          quality={100}
          priority
          className="filter brightness-50"
        />
      )}
      <div className="absolute inset-0 flex flex-col items-center justify-center p-4">
        <div className="text-center text-white z-10 max-w-4xl bg-black/30 p-4 rounded-xl backdrop-blur-sm">
          <h1 className="text-5xl md:text-7xl font-bold mb-4 text-left">
            Elevate your WoW experience
          </h1>
          <p className="text-lg md:text-xl mb-8 text-gray-300 text-left">
            Explore characters equipments, talents, mythic + and raids
            progression to stay at the state of the art in World of Warcraft.
          </p>
          <div className="w-full max-w-xl">
            <SearchBar />
          </div>
        </div>
      </div>
    </div>
  );
}
