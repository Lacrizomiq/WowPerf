import Image from "next/image";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { WowClassParam } from "@/types/warcraftlogs/builds/classSpec";
import { getClassIcon } from "@/utils/classandspecicons";

// Mapping of classes to their available specs
const CLASS_SPECS: Record<WowClassParam, string[]> = {
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

// Format the display name of the class
function formatDisplayName(className: WowClassParam): string {
  if (className === "deathknight") return "Death Knight";
  if (className === "demonhunter") return "Demon Hunter";
  return className.charAt(0).toUpperCase() + className.slice(1);
}

// Convert the class name to PascalCase for icons
function toPascalCaseForIcon(name: string): string {
  if (name === "deathknight") return "DeathKnight";
  if (name === "demonhunter") return "DemonHunter";
  return name.charAt(0).toUpperCase() + name.slice(1);
}

interface ClassSelectorProps {
  selectedClass: WowClassParam;
  onClassChange: (className: WowClassParam) => void;
}

export default function ClassSelector({
  selectedClass,
  onClassChange,
}: ClassSelectorProps) {
  const classes = Object.keys(CLASS_SPECS) as WowClassParam[];

  return (
    <Select
      value={selectedClass}
      onValueChange={(value) => onClassChange(value as WowClassParam)}
    >
      <SelectTrigger className="w-[180px] bg-slate-800 text-white border-slate-700">
        {selectedClass && (
          <div className="flex items-center gap-2">
            <div className="w-5 h-5 rounded-full overflow-hidden">
              <Image
                src={getClassIcon(toPascalCaseForIcon(selectedClass))}
                alt={formatDisplayName(selectedClass)}
                width={20}
                height={20}
                className="object-cover"
              />
            </div>
            <span>{formatDisplayName(selectedClass)}</span>
          </div>
        )}
        {!selectedClass && <SelectValue placeholder="Select Class" />}
      </SelectTrigger>

      <SelectContent className="bg-slate-900 border-slate-700 text-white max-h-80">
        {classes.map((className) => (
          <SelectItem
            key={className}
            value={className}
            className="hover:bg-slate-800"
          >
            <div className="flex items-center gap-2">
              <div className="w-5 h-5 rounded-full overflow-hidden">
                <Image
                  src={getClassIcon(toPascalCaseForIcon(className))}
                  alt={formatDisplayName(className)}
                  width={20}
                  height={20}
                  className="object-cover"
                />
              </div>
              <span>{formatDisplayName(className)}</span>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
