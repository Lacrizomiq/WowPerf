// # Homepage/gear of the spec
export default function BuildGearPage({
  params,
}: {
  params: { class: string; spec: string };
}) {
  return (
    <div>
      Hello ! This is the gea rpage for {params.class}/{params.spec}
    </div>
  );
}
