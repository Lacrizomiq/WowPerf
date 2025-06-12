// app/mythic-plus/builds/[class]/[spec]/enchants-gems/page.tsx
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import EnchantGemsContent from "@/components/BuildsAnalysis/mythicplus/enchants-gems/EnchantGemsContent";

export default async function BuildEnchantsGemsPage({
  params,
}: {
  params: Promise<{ class: string; spec: string }>;
}) {
  const resolvedParams = await params;
  const className = resolvedParams.class as WowClassParam;
  const spec = resolvedParams.spec as WowSpecParam;

  return <EnchantGemsContent className={className} spec={spec} />;
}
