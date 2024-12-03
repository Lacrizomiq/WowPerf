import React from "react";
import { Search, Trophy, LineChart, Users, Star, Layout } from "lucide-react";
import Link from "next/link";

export default function Feature() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 to-blue-900">
      {/* Features Grid */}
      <div className="bg-slate-900 py-20">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {/* Leaderboards */}
            <div className="bg-slate-800/50 rounded-xl p-6 backdrop-blur-sm hover:bg-slate-800/70 transition-all duration-300">
              <Link href="/mythic-plus/leaderboards">
                <div className="rounded-full bg-blue-600/10 w-12 h-12 flex items-center justify-center mb-4">
                  <Trophy className="text-blue-400" size={24} />
                </div>
                <h3 className="text-xl font-bold text-white mb-2">
                  Mythic+ Global Leaderboards
                </h3>
                <p className="text-slate-300">
                  Track the world&apos;s top players by role, spec, and dungeon.
                  Updated daily with the latest rankings.
                </p>
              </Link>
            </div>

            {/* Statistics */}
            <div className="bg-slate-800/50 rounded-xl p-6 backdrop-blur-sm">
              <div className="rounded-full bg-blue-600/10 w-12 h-12 flex items-center justify-center mb-4">
                <LineChart className="text-blue-400" size={24} />
              </div>
              <h3 className="text-xl font-bold text-white mb-2">
                Mythic+ Statistics
              </h3>
              <p className="text-slate-300">
                Detailed statistics on class distribution, team compositions,
                and success rates in high-level content.
              </p>
            </div>

            {/* Best Runs */}
            <div className="bg-slate-800/50 rounded-xl p-6 backdrop-blur-sm">
              <div className="rounded-full bg-blue-600/10 w-12 h-12 flex items-center justify-center mb-4">
                <Star className="text-blue-400" size={24} />
              </div>
              <h3 className="text-xl font-bold text-white mb-2">
                Best Mythic+ Runs
              </h3>
              <p className="text-slate-300">
                Explore the highest scoring Mythic+ runs with detailed team
                compositions and builds.
              </p>
            </div>

            {/* Character Progress */}
            <div className="bg-slate-800/50 rounded-xl p-6 backdrop-blur-sm">
              <div className="rounded-full bg-blue-600/10 w-12 h-12 flex items-center justify-center mb-4">
                <Users className="text-blue-400" size={24} />
              </div>
              <h3 className="text-xl font-bold text-white mb-2">
                Character Progress
              </h3>
              <p className="text-slate-300">
                Track your characters&apos; progression across raids and
                Mythic+.
              </p>
            </div>

            {/* Coming Soon: Dashboard */}
            <div className="bg-slate-800/50 rounded-xl p-6 backdrop-blur-sm">
              <div className="rounded-full bg-blue-600/10 w-12 h-12 flex items-center justify-center mb-4">
                <Layout className="text-blue-400" size={24} />
              </div>
              <h3 className="text-xl font-bold text-white mb-2">
                Personal Dashboard
                <span className="ml-2 text-xs bg-blue-600 px-2 py-1 rounded-full">
                  Coming Soon
                </span>
              </h3>
              <p className="text-slate-300">
                Get personalized recommendations and track multiple characters
                from one central dashboard.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
