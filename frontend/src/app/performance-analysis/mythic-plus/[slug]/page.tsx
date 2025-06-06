import SpecDetailView from "@/components/PerformanceStatistics/mythicplus/specdetails/SpecDetailView";

interface PageProps {
  params: {
    slug: string;
  };
}

export default function SpecAnalysisPage({ params }: PageProps) {
  return <SpecDetailView slug={params.slug} />;
}
