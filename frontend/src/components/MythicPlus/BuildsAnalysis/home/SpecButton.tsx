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
      href={`/mythic-plus/builds/${className}/${spec}`}
      className="flex items-center gap-3 py-2 px-3 rounded-md bg-slate-700/50 hover:bg-slate-600 transition-colors w-full group relative overflow-hidden"
    >
      {/* Effet de soulignement au hover */}
      <span className="absolute bottom-0 left-0 w-full h-0.5 bg-slate-500 transform scale-x-0 group-hover:scale-x-100 transition-transform origin-left"></span>

      <span className="w-6 h-6 flex-shrink-0 rounded-full overflow-hidden">
        <Image
          src={specIconUrl}
          alt={displayName}
          width={24}
          height={24}
          className="w-full h-full object-cover"
        />
      </span>
      <span>{displayName}</span>
    </Link>
  );
}
