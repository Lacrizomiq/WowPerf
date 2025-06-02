import Link from "next/link";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import {
  getSpecIcon,
  formatDisplaySpecName,
  classNameToPascalCase,
  specNameToPascalCase,
} from "@/utils/classandspecicons";
import Image from "next/image";

interface SpecButtonProps {
  className: WowClassParam;
  spec: WowSpecParam;
}

export default function SpecButton({ className, spec }: SpecButtonProps) {
  // Format the name of the spec for display
  const displayName = formatDisplaySpecName(spec);

  // Get the URL of the spec icon
  const classIconName = classNameToPascalCase(className);
  const specIconName = specNameToPascalCase(spec);
  const specIconUrl = getSpecIcon(classIconName, specIconName);

  return (
    <Link
      href={`/builds/mythic-plus/${className}/${spec}`}
      className="flex items-center gap-3 p-3 rounded-md bg-slate-800/70 hover:bg-slate-700 transition-colors mb-2"
    >
      <div className="relative w-6 h-6 rounded-full overflow-hidden bg-slate-700">
        <Image
          src={specIconUrl}
          alt={displayName}
          width={24}
          height={24}
          className="w-full h-full object-cover"
        />
      </div>
      <span className="text-slate-200">{displayName}</span>
    </Link>
  );
}
