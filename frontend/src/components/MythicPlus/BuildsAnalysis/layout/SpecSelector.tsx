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
import {
  formatDisplaySpecName,
  classNameToPascalCase,
  specNameToPascalCase,
  getSpecIcon,
} from "@/utils/classandspecicons";

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

  const classIconName = classNameToPascalCase(selectedClass);

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
                  specNameToPascalCase(selectedSpec)
                )}
                alt={formatDisplaySpecName(selectedSpec)}
                width={20}
                height={20}
                className="object-cover"
              />
            </div>
            <span>{formatDisplaySpecName(selectedSpec)}</span>
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
                  src={getSpecIcon(classIconName, specNameToPascalCase(spec))}
                  alt={formatDisplaySpecName(spec)}
                  width={20}
                  height={20}
                  className="object-cover"
                />
              </div>
              <span>{formatDisplaySpecName(spec)}</span>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
