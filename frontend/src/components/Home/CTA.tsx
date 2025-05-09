// components/Home/CTA.tsx
"use client";

import { Button } from "@/components/ui/button";
import Image from "next/image";
import { useEffect, useState } from "react";
import Link from "next/link";

export default function CTA() {
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  return (
    <section className="relative py-20 overflow-hidden">
      {/* Gradient overlay */}
      <div className="absolute inset-0  z-10"></div>

      {/* Background image */}
      {isMounted && (
        <Image
          src="/xal.jpg"
          alt="World of Warcraft Background"
          layout="fill"
          objectFit="cover"
          quality={100}
          priority
          className="opacity-50 z-0"
        />
      )}

      <div className="container mx-auto px-4 text-center relative z-20">
        <h2 className="text-3xl font-bold mb-4 text-white">
          Ready to Elevate Your Game?
        </h2>
        <p className="max-w-2xl mx-auto mb-8 text-purple-100">
          Join thousands of players who are already using WoW Perf to optimize
          their gameplay and climb the ranks.
        </p>
        <Link href="/signup">
          <Button
            size="lg"
            className="bg-white text-purple-900 hover:bg-purple-100"
          >
            Get Started Now
          </Button>
        </Link>
      </div>
    </section>
  );
}
