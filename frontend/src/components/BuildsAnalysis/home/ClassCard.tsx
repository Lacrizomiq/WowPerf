import Link from "next/link";
import SpecButton from "./SpecButton";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import {
  getClassIcon,
  formatDisplayClassName,
  classNameToPascalCase,
} from "@/utils/classandspecicons";
import Image from "next/image";

// Mapping of classes to available specs
const CLASS_SPECS: Record<WowClassParam, WowSpecParam[]> = {
  warrior: ["arms", "fury", "protection"],
  paladin: ["holy", "protection", "retribution"],
  hunter: ["beastmastery", "marksmanship", "survival"],
  rogue: ["assassination", "outlaw", "subtlety"],
  priest: ["discipline", "holy", "shadow"],
  deathknight: ["blood", "frost", "unholy"],
  shaman: ["elemental", "enhancement", "restoration"],
  mage: ["arcane", "fire", "frost"],
  warlock: ["affliction", "demonology", "destruction"],
  monk: ["brewmaster", "mistweaver", "windwalker"],
  druid: ["balance", "feral", "guardian", "restoration"],
  demonhunter: ["havoc", "vengeance"],
  evoker: ["devastation", "preservation", "augmentation"],
};

interface ClassCardProps {
  className: WowClassParam;
}

export default function ClassCard({ className }: ClassCardProps) {
  // Protection contre les valeurs undefined
  if (!className) {
    console.error("ClassCard received undefined className");
    return null;
  }

  // Get the specs for this class
  const specs = CLASS_SPECS[className] || [];

  // Get the formatted display name of the class (avec protection)
  const displayName =
    typeof className === "string"
      ? formatDisplayClassName(className)
      : className;

  // Style for the class color
  const classColorStyle = { color: `var(--color-${className})` };
  const classBorderStyle = { borderColor: `var(--color-${className})` };

  // Get the URL of the class icon (avec protection)
  const classIconName =
    typeof className === "string"
      ? classNameToPascalCase(className)
      : className;
  const classIconUrl = getClassIcon(classIconName);

  return (
    <div
      className={`bg-slate-800/50 rounded-lg border overflow-hidden`}
      style={classBorderStyle}
    >
      <div className="p-4 border-b border-slate-700">
        <h3
          className="text-xl font-bold flex items-center gap-3"
          style={classColorStyle}
        >
          <span className="w-8 h-8 flex-shrink-0 rounded-full overflow-hidden bg-slate-700">
            <Image
              src={classIconUrl}
              alt={displayName}
              width={32}
              height={32}
              className="w-full h-full object-cover"
            />
          </span>
          {displayName}
        </h3>
      </div>
      <div className="p-2 space-y-2">
        {specs.map((spec) => (
          <SpecButton key={spec} className={className} spec={spec} />
        ))}
      </div>
    </div>
  );
}
