// components/Home/Statistics.tsx
import {
  ArrowRight,
  BarChart3,
  PieChart,
  Trophy,
  Settings,
  TrendingUp,
  LayoutDashboard,
} from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

export default function Statistics() {
  return (
    <section className="py-2 pb-8 bg-[#1A1D21]">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold mb-4 text-white">
            Powerful Tools at Your Fingertips
          </h2>
          <p className="text-slate-400 text-base max-w-2xl mx-auto">
            Dive into our comprehensive suite of analytics and optimization
            tools designed for serious players.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 auto-rows-fr">
          <FeatureCard
            icon={<BarChart3 className="h-10 w-10 text-purple-500" />}
            title="Deep Dive into Performance Rankings"
            description="Analyze class and spec performance across all contents with detailed breakdowns and filters."
            link="/performance-analysis"
            linkText="Performance Analysis"
          />

          <FeatureCard
            icon={<PieChart className="h-10 w-10 text-purple-500" />}
            title="Uncover Trends & Statistics"
            description="Stay ahead of the meta with weekly trends and statistics."
            link="/statistics"
            linkText="Statistics"
          />

          <FeatureCard
            icon={<Trophy className="h-10 w-10 text-purple-500" />}
            title="Study the Anatomy of Top Runs"
            description="Learn from the best by examining the gear, talents, and composition of top-ranking players."
            link="/best-runs"
            linkText="Best Runs"
          />

          <FeatureCard
            icon={<Settings className="h-10 w-10 text-purple-500" />}
            title="Find Your Optimal Character Build"
            description="Discover the most effective talent builds, gear setups, and stat priorities for your class and spec."
            link="/builds"
            linkText="Builds"
          />

          <FeatureCard
            icon={<TrendingUp className="h-10 w-10 text-purple-500" />}
            title="Monitor Your Character's Ascent"
            description="Track your progress over time with detailed performance metrics and improvement suggestions."
            link="/character-progress"
            linkText=""
            comingSoon
          />

          <FeatureCard
            icon={<LayoutDashboard className="h-10 w-10 text-purple-500" />}
            title="Your Personal Command Center"
            description="A customizable dashboard with all your important metrics and tools in one place."
            link="/dashboard"
            linkText=""
            comingSoon
          />
        </div>
      </div>
    </section>
  );
}
function FeatureCard({
  icon,
  title,
  description,
  link,
  linkText,
  comingSoon = false,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
  link: string;
  linkText: string;
  comingSoon?: boolean;
}) {
  return (
    <div className="flex flex-col h-full overflow-hidden border border-slate-700 bg-slate-800/30 rounded-lg transition-all hover:shadow-lg">
      {/* Contenu de la carte - utilise flex-grow-1 pour prendre l'espace disponible */}
      <div className="p-6 flex-grow">
        <div className="mb-4">{icon}</div>
        <h3 className="text-xl font-bold text-white mb-3">{title}</h3>
        <p className="text-slate-400 text-base">{description}</p>
      </div>

      {/* Conteneur du bouton - toujours en bas de la carte */}
      <div className="p-6 pt-0 mt-auto">
        <Button
          variant={comingSoon ? "outline" : "default"}
          className={`w-full justify-between group ${
            comingSoon
              ? "border-purple-700/50 text-purple-300 hover:bg-purple-900/30 hover:text-purple-200"
              : "bg-purple-600 hover:bg-purple-700"
          }`}
          asChild
        >
          <Link href={link}>
            {comingSoon && (
              <Badge
                variant="outline"
                className="mr-2 border-purple-600 text-purple-400"
              >
                Coming Soon
              </Badge>
            )}
            {linkText}
            <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
          </Link>
        </Button>
      </div>
    </div>
  );
}
