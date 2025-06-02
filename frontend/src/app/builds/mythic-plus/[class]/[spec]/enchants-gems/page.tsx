// app/mythic-plus/builds/[class]/[spec]/enchants-gems/page.tsx
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import EnchantGemsContent from "@/components/BuildsAnalysis/mythicplus/enchants-gems/EnchantGemsContent";

export default function BuildEnchantsGemsPage({
  params,
}: {
  params: { class: string; spec: string };
}) {
  const className = params.class as WowClassParam;
  const spec = params.spec as WowSpecParam;

  return <EnchantGemsContent className={className} spec={spec} />;
}
