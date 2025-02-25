// app/mythic-plus/spec-analysis/[slug]/page.tsx
"use client";

import React from "react";
import { useParams } from "next/navigation";
import SpecDetailView from "@/components/MythicPlus/PerformanceStatistics/SpecBestPlayerDetails";

export default function SpecDetailPage() {
  const params = useParams();
  const slug = params?.slug as string;

  return <SpecDetailView slug={slug} />;
}
