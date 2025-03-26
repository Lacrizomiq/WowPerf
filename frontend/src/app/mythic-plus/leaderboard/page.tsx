import LeaderboardTabs from "@/components/MythicPlus/Leaderboard/LeaderboardTabs";

export default function HomePage() {
  return (
    <div className=" w-full h-full bg-black py-6">
      <h1 className="text-2xl font-bold text-left px-8">
        Mythic+ Leaderboards
      </h1>
      <div className="mx-auto px-8 py-6 bg-black">
        <p className="text-left text-xl font-bold">
          Data are updated every 24 hours.
        </p>
        <p className="text-left text-sm">
          Due to technical limitations, the data is not updated in real-time and
          only includes the very best players of each role
        </p>
        <p className="text-left text-sm">
          The data is provided by{" "}
          <a
            href="https://www.warcraftlogs.com/about"
            target="_blank"
            className="text-blue-400 hover:text-blue-300 transition-colors"
          >
            Warcraft Logs
          </a>{" "}
          , check them out for more detailed data.
        </p>
      </div>
      <LeaderboardTabs />
    </div>
  );
}
