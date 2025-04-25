import Image from "next/image";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  WowClassParam,
  WowSpecParam,
} from "@/types/warcraftlogs/builds/classSpec";
import { getSpecIcon } from "@/utils/classandspecicons";

// Mapping of classes to their available specs
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

// Format the display name of a spec
function formatSpecName(spec: WowSpecParam): string {
  // Special cases
  if (spec === "beastmastery") return "Beast Mastery";

  // For other specs, capitalize the first letter
  return spec.charAt(0).toUpperCase() + spec.slice(1);
}

// Convert class and spec names to PascalCase for icons
function toPascalCaseForIcon(name: string): string {
  if (name === "deathknight") return "DeathKnight";
  if (name === "demonhunter") return "DemonHunter";
  return name.charAt(0).toUpperCase() + name.slice(1);
}

// Convert spec name to PascalCase for icons
function specToPascalCaseForIcon(spec: string): string {
  if (spec === "beastmastery") return "BeastMastery";
  return spec.charAt(0).toUpperCase() + spec.slice(1);
}

interface SpecSelectorProps {
  selectedClass: WowClassParam;
  selectedSpec: WowSpecParam;
  onSpecChange: (spec: WowSpecParam) => void;
}

export default function SpecSelector({
  selectedClass,
  selectedSpec,
  onSpecChange,
}: SpecSelectorProps) {
  const specs = CLASS_SPECS[selectedClass] || [];

  const classIconName = toPascalCaseForIcon(selectedClass);

  return (
    <Select
      value={selectedSpec}
      onValueChange={(value) => onSpecChange(value as WowSpecParam)}
    >
      <SelectTrigger className="w-[180px] bg-slate-800 text-white border-slate-700">
        {selectedSpec && (
          <div className="flex items-center gap-2">
            <div className="w-5 h-5 rounded-full overflow-hidden">
              <Image
                src={getSpecIcon(
                  classIconName,
                  specToPascalCaseForIcon(selectedSpec)
                )}
                alt={formatSpecName(selectedSpec)}
                width={20}
                height={20}
                className="object-cover"
              />
            </div>
            <span>{formatSpecName(selectedSpec)}</span>
          </div>
        )}
        {!selectedSpec && <SelectValue placeholder="Select Spec" />}
      </SelectTrigger>

      <SelectContent className="bg-slate-900 border-slate-700 text-white">
        {specs.map((spec) => (
          <SelectItem key={spec} value={spec} className="hover:bg-slate-800">
            <div className="flex items-center gap-2">
              <div className="w-5 h-5 rounded-full overflow-hidden">
                <Image
                  src={getSpecIcon(
                    classIconName,
                    specToPascalCaseForIcon(spec)
                  )}
                  alt={formatSpecName(spec)}
                  width={20}
                  height={20}
                  className="object-cover"
                />
              </div>
              <span>{formatSpecName(spec)}</span>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
