import SpecDetailView from "@/components/PerformanceStatistics/mythicplus/specdetails/SpecDetailView";

interface PageProps {
  params: Promise<{
    slug: string;
  }>;
}

export default async function SpecAnalysisPage({ params }: PageProps) {
  const resolvedParams = await params;
  return <SpecDetailView slug={resolvedParams.slug} />;
}
