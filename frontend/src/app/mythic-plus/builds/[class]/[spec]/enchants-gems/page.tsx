// # Homepage/enchants-gems of the spec
export default function BuildEnchantsGemsPage({
  params,
}: {
  params: { class: string; spec: string };
}) {
  return (
    <div>
      Hello ! This is the enchants & gems page for {params.class}/{params.spec}
    </div>
  );
}
