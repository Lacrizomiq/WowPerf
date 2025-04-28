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
  // Get the specs for this class
  const specs = CLASS_SPECS[className];

  // Get the formatted display name of the class
  const displayName = formatDisplayClassName(className);

  // Style for the class color
  const classColorStyle = { color: `var(--color-${className})` };
  const classBorderStyle = { borderColor: `var(--color-${className})` };

  // Get the URL of the class icon
  const classIconName = classNameToPascalCase(className);
  const classIconUrl = getClassIcon(classIconName);

  return (
    <div
      className="bg-[#112240] rounded-lg overflow-hidden shadow-lg border-l-4 hover:shadow-xl"
      style={classBorderStyle}
    >
      <div className="p-5">
        <h2
          className="text-xl font-bold mb-4 flex items-center gap-3"
          style={classColorStyle}
        >
          <span className="w-8 h-8 flex-shrink-0 rounded-full overflow-hidden">
            <Image
              src={classIconUrl}
              alt={displayName}
              width={32}
              height={32}
              className="w-full h-full object-cover"
            />
          </span>
          {displayName}
        </h2>

        <div className="space-y-2">
          {specs.map((spec) => (
            <SpecButton key={spec} className={className} spec={spec} />
          ))}
        </div>
      </div>
    </div>
  );
}
