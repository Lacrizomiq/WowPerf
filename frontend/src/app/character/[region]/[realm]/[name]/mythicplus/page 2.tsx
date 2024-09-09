import MythicDungeonOverview from "@/components/MythicPlus/MythicOverview";

const MythicPlusPage = ({
  params,
}: {
  params: { region: string; realm: string; name: string; seasonSlug: string };
}) => {
  const { region, realm, name, seasonSlug } = params;
  return (
    <div className="p-6 bg-gradient-dark shadow-lg rounded-lg glow-effect m-12 max-w-6xl mx-auto">
      <MythicDungeonOverview
        characterName={name}
        realmSlug={realm}
        region={region}
        namespace={`profile-${region}`}
        locale="en_GB"
        seasonSlug={seasonSlug}
      />
    </div>
  );
};

export default MythicPlusPage;
