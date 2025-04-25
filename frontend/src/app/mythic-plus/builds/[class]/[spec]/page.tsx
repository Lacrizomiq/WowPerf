// # Page d'accueil/overview de la spec
export default function BuildPage({
  params,
}: {
  params: { class: string; spec: string };
}) {
  return (
    <div>
      Salut ! Ceci est la page pour {params.class}/{params.spec}
    </div>
  );
}
