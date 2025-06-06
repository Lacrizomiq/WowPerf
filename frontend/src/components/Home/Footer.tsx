// components/Home/Footer.tsx
import Link from "next/link";
import { Twitter } from "lucide-react";

export default function Footer() {
  return (
    <footer className="bg-[#1A1D21] border-t border-slate-800/60 pt-16 pb-8">
      <div className="container mx-auto px-4">
        {/* Main Footer Links */}
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-8 mb-12">
          {/* Features Column */}
          <div>
            <h3 className="text-lg font-semibold text-purple-400 mb-4">
              Features
            </h3>
            <ul className="space-y-3">
              <li>
                <Link
                  href="/performance-analysis"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  Performance Analysis
                </Link>
              </li>
              <li>
                <Link
                  href="/builds"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  Builds
                </Link>
              </li>
              <li>
                <Link
                  href="/statistics"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  Statistics
                </Link>
              </li>
              <li>
                <Link
                  href="/leaderboards"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  Leaderboards
                </Link>
              </li>
            </ul>
          </div>

          {/* Resources Column */}
          <div>
            <h3 className="text-lg font-semibold text-purple-400 mb-4">
              Resources
            </h3>
            <ul className="space-y-3">
              <li>
                <Link
                  href="/documentation"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  Documentation
                </Link>
              </li>
              <li>
                <Link
                  href="/faq"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  FAQ
                </Link>
              </li>
            </ul>
          </div>

          {/* Community Column */}
          <div>
            <h3 className="text-lg font-semibold text-purple-400 mb-4">
              Community
            </h3>
            <ul className="space-y-3">
              <li>
                <Link
                  href="/report-bug"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  Report a bug
                </Link>
              </li>
              <li>
                <Link
                  href="/suggest-feature"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  Suggest a feature
                </Link>
              </li>
              <li>
                <Link
                  href="https://twitter.com/wowperf"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm inline-flex items-center"
                >
                  <Twitter className="mr-2 h-4 w-4" />
                  Twitter
                </Link>
              </li>
            </ul>
          </div>

          {/* Legal Column */}
          <div>
            <h3 className="text-lg font-semibold text-purple-400 mb-4">
              Legal
            </h3>
            <ul className="space-y-3">
              <li>
                <Link
                  href="/about"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  About
                </Link>
              </li>
              <li>
                <Link
                  href="/contact"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  Contact
                </Link>
              </li>
              <li>
                <Link
                  href="/privacy-policy"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  Privacy Policy
                </Link>
              </li>
              <li>
                <Link
                  href="/terms-of-service"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  Terms of Service
                </Link>
              </li>
              <li>
                <Link
                  href="/cookie-policy"
                  className="text-slate-300 hover:text-purple-400 transition-colors text-sm"
                >
                  Cookie Policy
                </Link>
              </li>
            </ul>
          </div>
        </div>

        {/* Divider */}
        <div className="border-t border-slate-800 mb-8"></div>

        {/* Bottom Footer with Copyright and Powered By */}
        <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
          {/* Copyright Section */}
          <div className="text-slate-400 text-sm max-w-2xl">
            <p className="mb-2">© 2025 WoW Perf. All rights reserved.</p>
            <p>
              World of Warcraft, Warcraft and Blizzard Entertainment are
              trademarks or registered trademarks of Blizzard Entertainment,
              Inc. in the U.S. and/or other countries.
            </p>
          </div>

          {/* Powered By Section */}
          <div className="flex items-center flex-wrap gap-3">
            <span className="text-slate-400 text-sm mr-2">Powered by</span>
            <div className="flex flex-wrap items-center gap-2">
              <Link
                href="https://develop.battle.net/"
                target="_blank"
                rel="noopener noreferrer"
                className="px-3 py-1 bg-slate-800 hover:bg-slate-700 rounded-full text-purple-400 text-xs transition-colors"
              >
                Blizzard API
              </Link>
              <Link
                href="https://www.warcraftlogs.com/"
                target="_blank"
                rel="noopener noreferrer"
                className="px-3 py-1 bg-slate-800 hover:bg-slate-700 rounded-full text-purple-400 text-xs transition-colors"
              >
                Warcraft Logs
              </Link>
              <Link
                href="https://raider.io/"
                target="_blank"
                rel="noopener noreferrer"
                className="px-3 py-1 bg-slate-800 hover:bg-slate-700 rounded-full text-purple-400 text-xs transition-colors"
              >
                Raider.io
              </Link>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
