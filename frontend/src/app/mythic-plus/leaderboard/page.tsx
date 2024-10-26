// Dans ta page principale
import { RoleLeaderboards } from "@/components/MythicPlus/Leaderboard/RoleLeaderboard/RoleLeaderboards";

export default function HomePage() {
  return (
    <div className="p-6 bg-black w-full h-full">
      <h1 className="text-2xl font-bold text-center mb-6">
        Mythic+ Leaderboards
      </h1>
      <RoleLeaderboards />
    </div>
  );
}
