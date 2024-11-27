// components/Home/Statistics.tsx
import React from "react";

export default function Statistics() {
  return (
    <div className="bg-slate-800/50 py-20">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold text-white mb-4">
            Live Statistics
          </h2>
          <p className="text-slate-300">
            Real-time insights from the top players and teams
          </p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="bg-slate-900/50 rounded-xl p-6 text-center">
            <div className="text-4xl font-bold text-blue-400 mb-2">17-19+</div>
            <div className="text-slate-300">Average M+ Key Level</div>
          </div>
          <div className="bg-slate-900/50 rounded-xl p-6 text-center">
            <div className="text-4xl font-bold text-blue-400 mb-2">3500+</div>
            <div className="text-slate-300">Top Rating Score</div>
          </div>
          <div className="bg-slate-900/50 rounded-xl p-6 text-center">
            <div className="text-4xl font-bold text-blue-400 mb-2">24hr</div>
            <div className="text-slate-300">Data Update Frequency</div>
          </div>
        </div>
      </div>
    </div>
  );
}
