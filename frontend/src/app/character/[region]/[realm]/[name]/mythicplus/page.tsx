import MythicDungeon from "@/components/MythicPlus/MythicDungeon";

const MythicPlusPage = (params: { params: { seasonSlug: string } }) => {
  const { seasonSlug } = params.params;
  return (
    <div className="p-6 bg-gradient-dark shadow-lg rounded-lg glow-effect m-12 max-w-6xl mx-auto">
      <MythicDungeon seasonSlug={seasonSlug} />
    </div>
  );
};

export default MythicPlusPage;
