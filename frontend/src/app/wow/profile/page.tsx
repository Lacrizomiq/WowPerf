"use client";

import { WoWProfile } from "@/components/UserProfile/WoWAccount/AccountProfile";
import { useRequireAuth } from "@/providers/AuthContext";

export default function WoWProfilePage() {
  const { isLoading } = useRequireAuth();

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="container mx-auto px-4 py-8">
      <WoWProfile />
    </div>
  );
}
