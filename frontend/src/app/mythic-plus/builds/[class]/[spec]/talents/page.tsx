// # Homepage/talents of the spec
export default function BuildTalentsPage({
  params,
}: {
  params: { class: string; spec: string };
}) {
  return (
    <div>
      Hello ! This is the talentspage for {params.class}/{params.spec}
    </div>
  );
}
