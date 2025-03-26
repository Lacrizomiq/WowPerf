// components/Footer/Footer.tsx
import Image from "next/image";
import { Github, Twitter } from "lucide-react";
import Link from "next/link";

export default function Footer() {
  return (
    <footer className="relative">
      {/* Background Image */}
      <div className="absolute inset-0 z-0">
        <Image
          src="/tww.png"
          alt="Background"
          fill
          className="object-cover object-top brightness-[0.35]"
          quality={100}
        />
      </div>

      {/* Content */}
      <div className="relative z-10 bg-gradient-to-b from-transparent to-slate-900/90">
        <div className="container mx-auto px-4 pt-20 pb-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-12">
            {/* Features Section */}
            <div>
              <h3 className="text-blue-400 font-semibold mb-4">Features</h3>
              <ul className="space-y-2">
                <li>
                  <Link
                    href="/"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors"
                  >
                    Character Search
                  </Link>
                </li>
                <li>
                  <Link
                    href="/mythic-plus/leaderboard"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors"
                  >
                    Mythic+ Leaderboard
                  </Link>
                </li>
                <li>
                  <Link
                    href="/mythic-plus/statistics"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors"
                  >
                    Mythic+ Statistics
                  </Link>
                </li>
                <li>
                  <Link
                    href="/raids"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors"
                  >
                    Raid Progress
                  </Link>
                </li>
              </ul>
            </div>

            {/* Resources Section */}
            <div>
              <h3 className="text-blue-400 font-semibold mb-4">Resources</h3>
              <ul className="space-y-2">
                <li>
                  <Link
                    href="#"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors"
                  >
                    Documentation
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors"
                  >
                    FAQ
                  </Link>
                </li>
              </ul>
            </div>

            {/* Community Section */}
            <div>
              <h3 className="text-blue-400 font-semibold mb-4">Community</h3>
              <ul className="space-y-2">
                <li>
                  <Link
                    href="#"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors inline-flex items-center gap-2"
                  >
                    Report a bug
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors inline-flex items-center gap-2"
                  >
                    Suggest a feature
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors inline-flex items-center gap-2"
                  >
                    <Twitter size={16} />
                    Twitter
                  </Link>
                </li>
              </ul>
            </div>

            {/* Legal Section */}
            <div>
              <h3 className="text-blue-400 font-semibold mb-4">Legal</h3>
              <ul className="space-y-2">
                <li>
                  <Link
                    href="#"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors"
                  >
                    Terms of Service
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors"
                  >
                    Privacy Policy
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-slate-300 hover:text-blue-400 text-sm transition-colors"
                  >
                    Cookie Policy
                  </Link>
                </li>
              </ul>
            </div>
          </div>

          {/* Bottom Section */}
          <div className="border-t border-slate-700/50 pt-8">
            <div className="flex flex-col md:flex-row items-center justify-between gap-4">
              <div className="text-slate-400 text-sm w-1/2">
                <p className="mb-2">© 2024 WoW Perf. All rights reserved.</p>
                <p>
                  World of Warcraft, Warcraft and Blizzard Entertainment are
                  trademarks or registered trademarks of Blizzard Entertainment,
                  Inc. in the U.S. and/or other countries.
                </p>
              </div>

              <div className="flex items-center gap-4">
                <span className="text-slate-400 text-sm">Powered by</span>
                <div className="flex items-center gap-2">
                  <Link
                    href="https://develop.battle.net/"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="px-3 py-1 bg-blue-900/30 rounded-full text-blue-400 text-xs hover:bg-blue-900/50 transition-colors"
                  >
                    Blizzard API
                  </Link>
                  <Link
                    href="https://www.warcraftlogs.com/"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="px-3 py-1 bg-blue-900/30 rounded-full text-blue-400 text-xs hover:bg-blue-900/50 transition-colors"
                  >
                    Warcraft Logs
                  </Link>
                  <Link
                    href="https://raider.io/"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="px-3 py-1 bg-blue-900/30 rounded-full text-blue-400 text-xs hover:bg-blue-900/50 transition-colors"
                  >
                    Raider.io
                  </Link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
