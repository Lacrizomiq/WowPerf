// components/Home/Feature.tsx
"use client";

import { ArrowRight, ArrowLeft, ArrowRightCircle } from "lucide-react";
import Image from "next/image";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import Link from "next/link";
import { useState } from "react";

export default function Feature() {
  const [activeSpotlight, setActiveSpotlight] = useState(0);
  const spotlightItems = [
    { name: "Balance Druid", badge: "Spotlight" },
    { name: "Unholy Death Knight", badge: "Spotlight" },
    { name: "Arcane Mage", badge: "Spotlight" },
  ];

  const nextSpotlight = () => {
    setActiveSpotlight((prev) => (prev + 1) % spotlightItems.length);
  };

  const prevSpotlight = () => {
    setActiveSpotlight(
      (prev) => (prev - 1 + spotlightItems.length) % spotlightItems.length
    );
  };

  return (
    <section className="pt-6 pb-16 bg-[#1A1D21]">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold mb-4 text-white">
            Get a Glimpse of the Power
          </h2>
          <p className="text-slate-400 text-base max-w-2xl mx-auto">
            Experience a preview of our analytical tools designed to enhance
            your World of Warcraft gameplay.
          </p>
        </div>

        <Tabs defaultValue="mythic" className="w-full max-w-4xl mx-auto">
          <TabsList className="grid w-full grid-cols-3 bg-slate-800/50">
            <TabsTrigger
              value="mythic"
              className="data-[state=active]:bg-purple-600"
            >
              Mythic+
            </TabsTrigger>
            <TabsTrigger
              value="raids"
              disabled
              className="relative text-slate-500"
            >
              Raids
              <Badge
                variant="outline"
                className="ml-2 text-xs py-0 h-5 border-purple-600 text-purple-400"
              >
                Soon
              </Badge>
            </TabsTrigger>
            <TabsTrigger
              value="pvp"
              disabled
              className="relative text-slate-500"
            >
              PvP
              <Badge
                variant="outline"
                className="ml-2 text-xs py-0 h-5 border-purple-600 text-purple-400"
              >
                Soon
              </Badge>
            </TabsTrigger>
          </TabsList>

          <TabsContent value="mythic" className="mt-6">
            <div className="bg-slate-800/30 rounded-lg border border-slate-700 shadow-lg p-6 relative overflow-hidden">
              <div className="absolute top-0 right-0 w-32 h-32 bg-purple-500/10 rounded-full -mr-16 -mt-16"></div>
              <div className="absolute bottom-0 left-0 w-24 h-24 bg-indigo-500/10 rounded-full -ml-12 -mb-12"></div>

              <h3 className="text-xl font-bold mb-4 text-white">
                Top Performing Specs This Week
              </h3>

              <div className="flex items-center justify-center py-10">
                <Badge
                  variant="outline"
                  className="text-purple-400 border-purple-600 text-lg py-2 px-4"
                >
                  Coming Soon
                </Badge>
              </div>

              <div className="flex justify-center mt-6">
                <Link href="/performance-analysis">
                  <Button
                    variant="outline"
                    className="group border-purple-700 text-purple-300 hover:bg-purple-900/30 hover:text-purple-200"
                  >
                    Explore Full Performance Analysis
                    <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
                  </Button>
                </Link>
              </div>
            </div>
          </TabsContent>
        </Tabs>

        {/* Carousel for Spec Spotlight */}
        <div className="mt-12 max-w-4xl mx-auto">
          <div className="bg-slate-800/30 rounded-lg border border-slate-700 shadow-lg p-6 relative overflow-hidden">
            <div className="absolute top-0 right-0 w-32 h-32 bg-purple-500/10 rounded-full -mr-16 -mt-16"></div>
            <div className="absolute bottom-0 left-0 w-24 h-24 bg-indigo-500/10 rounded-full -ml-12 -mb-12"></div>

            <div className="flex justify-between items-center mb-6">
              <button
                onClick={prevSpotlight}
                className="bg-slate-700/50 hover:bg-purple-700/30 p-2 rounded-full text-white"
                aria-label="Previous spec"
              >
                <ArrowLeft className="h-5 w-5" />
              </button>

              <div className="flex flex-col md:flex-row gap-6 items-center flex-1 px-4">
                <div className="w-full md:w-1/3 flex justify-center">
                  <div className="relative w-48 h-48 rounded-full overflow-hidden border-4 border-purple-600 shadow-lg bg-purple-900/20 flex items-center justify-center">
                    <Badge className="text-purple-400 border-purple-600 text-lg py-2 px-4 bg-slate-900/70 z-10">
                      Coming Soon
                    </Badge>
                    <div className="absolute inset-0 opacity-50">
                      {/* Placeholder circle with purple glow */}
                      <div className="absolute inset-0 flex items-center justify-center">
                        <div className="w-full h-full flex items-center justify-center">
                          <div className="w-24 h-24 rounded-full bg-purple-500/20"></div>
                          <div className="absolute w-16 h-16 rounded-full bg-purple-500/40"></div>
                          <div className="absolute w-8 h-8 rounded-full bg-purple-500/60"></div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="w-full md:w-2/3 text-center md:text-left">
                  <div className="flex items-center mb-4 justify-center md:justify-start">
                    <h3 className="text-xl font-bold text-white">
                      {spotlightItems[activeSpotlight].name}
                    </h3>
                    <Badge className="ml-3 bg-purple-600">Spotlight</Badge>
                  </div>

                  <p className="text-slate-400 text-base mb-4">
                    Dominating high keys with exceptional burst AoE and strong
                    utility. The current meta favors this spec for its
                    consistent damage profile and valuable crowd control.
                  </p>

                  <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-6">
                    {["Coming", "Soon", "Stay", "Tuned"].map(
                      (talent, index) => (
                        <div
                          key={index}
                          className="bg-purple-900/30 border border-purple-800/50 rounded-md p-2 text-center text-sm font-medium text-purple-300"
                        >
                          {talent}
                        </div>
                      )
                    )}
                  </div>
                </div>
              </div>

              <button
                onClick={nextSpotlight}
                className="bg-slate-700/50 hover:bg-purple-700/30 p-2 rounded-full text-white"
                aria-label="Next spec"
              >
                <ArrowRight className="h-5 w-5" />
              </button>
            </div>

            <div className="flex justify-center mt-2">
              <Button className="w-full sm:w-auto bg-purple-600 hover:bg-purple-700 group">
                Explore All {spotlightItems[activeSpotlight].name} Builds
                <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
              </Button>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
