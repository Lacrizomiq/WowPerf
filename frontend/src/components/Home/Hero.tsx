"use client";

import Image from "next/image";
import { useState, useEffect } from "react";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function Hero() {
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  return (
    <div className="relative w-full overflow-hidden bg-[#1A1D21]">
      {/* Hauteur réduite à 65vh */}
      <div className="h-[45vh] relative">
        {/* Gradient overlay at the bottom to blend image with content */}
        <div className="absolute inset-0 bg-gradient-to-b from-transparent to-[#1A1D21] z-10"></div>

        {/* Left side gradient */}
        <div className="absolute inset-0 bg-gradient-to-r from-[#1A1D21] via-transparent to-transparent z-10"></div>

        {/* Right side gradient */}
        <div className="absolute inset-0 bg-gradient-to-l from-[#1A1D21] via-transparent to-transparent z-10"></div>

        {/* Top gradient */}
        <div className="absolute inset-0 bg-gradient-to-b from-[#1A1D21] via-transparent to-transparent z-10 opacity-50"></div>

        {isMounted && (
          <Image
            src="/alleria.jpg"
            alt="World of Warcraft"
            layout="fill"
            objectFit="cover"
            quality={100}
            priority
            className="opacity-80 z-0"
          />
        )}

        <div className="absolute inset-0 flex flex-col items-center justify-center p-4 z-20">
          <div className="text-center text-white max-w-4xl p-4">
            <h1 className="text-4xl md:text-6xl font-bold mb-3">
              WoW Perf:{" "}
              <span className="text-purple-300">
                Insight. Optimize. Conquer.
              </span>
            </h1>
            <p className="text-base md:text-xl mb-6 text-gray-200">
              Your ultimate companion for analytics and character improvement.
            </p>

            <div className="flex flex-col sm:flex-row gap-3 justify-center">
              <Button
                size="default"
                className="bg-purple-600 hover:bg-purple-700 text-white px-6"
              >
                Get Started
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
              <Button
                size="default"
                variant="outline"
                className="border-purple-600 text-purple-400 hover:bg-purple-900/30 hover:text-purple-300"
              >
                Learn More
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
