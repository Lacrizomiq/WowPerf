export default function CharacterLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <div className="relative z-20">{children}</div>;
}
