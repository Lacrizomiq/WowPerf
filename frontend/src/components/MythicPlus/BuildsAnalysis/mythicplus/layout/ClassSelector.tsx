// ClassSelector.tsx - Version complète harmonisée
import Image from "next/image";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  formatDisplayClassName,
  classNameToPascalCase,
  getClassIcon,
} from "@/utils/classandspecicons";
import { WowClassParam } from "@/types/warcraftlogs/builds/classSpec";
import { useRouter } from "next/navigation";

// Mapping of classes to their available specs
const CLASS_SPECS: Record<WowClassParam, string[]> = {
  deathknight: ["blood", "frost", "unholy"],
  demonhunter: ["havoc", "vengeance"],
  druid: ["balance", "feral", "guardian", "restoration"],
  evoker: ["devastation", "preservation", "augmentation"],
  hunter: ["beastmastery", "marksmanship", "survival"],
  mage: ["arcane", "fire", "frost"],
  monk: ["brewmaster", "mistweaver", "windwalker"],
  paladin: ["holy", "protection", "retribution"],
  priest: ["discipline", "holy", "shadow"],
  rogue: ["assassination", "outlaw", "subtlety"],
  shaman: ["elemental", "enhancement", "restoration"],
  warlock: ["affliction", "demonology", "destruction"],
  warrior: ["arms", "fury", "protection"],
};

interface ClassSelectorProps {
  selectedClass: WowClassParam;
  onClassChange: (className: WowClassParam) => void;
}

export default function ClassSelector({
  selectedClass,
  onClassChange,
}: ClassSelectorProps) {
  const router = useRouter();
  const classes = Object.keys(CLASS_SPECS) as WowClassParam[];

  const handleClassChange = (value: string) => {
    const newClass = value as WowClassParam;

    // if onClassChange is provided, call it
    if (onClassChange) {
      onClassChange(newClass);
    }

    // Get the first spec of the selected class
    const firstSpec = CLASS_SPECS[newClass][0];

    router.push(`/builds/mythic-plus/${newClass}/${firstSpec}`);
  };

  return (
    <Select value={selectedClass} onValueChange={handleClassChange}>
      <SelectTrigger className="w-[180px] bg-slate-800/50 text-white border-slate-700 focus:ring-purple-600">
        {selectedClass && (
          <div className="flex items-center gap-2">
            <div className="w-5 h-5 rounded-full overflow-hidden">
              <Image
                src={getClassIcon(classNameToPascalCase(selectedClass))}
                alt={formatDisplayClassName(selectedClass)}
                width={20}
                height={20}
                className="object-cover"
              />
            </div>
            <span>{formatDisplayClassName(selectedClass)}</span>
          </div>
        )}
        {!selectedClass && <SelectValue placeholder="Select Class" />}
      </SelectTrigger>

      <SelectContent className="bg-slate-900 border-slate-700 text-white max-h-80">
        {classes.map((className) => (
          <SelectItem
            key={className}
            value={className}
            className="hover:bg-slate-800 focus:bg-purple-600 focus:text-white"
          >
            <div className="flex items-center gap-2">
              <div className="w-5 h-5 rounded-full overflow-hidden">
                <Image
                  src={getClassIcon(classNameToPascalCase(className))}
                  alt={formatDisplayClassName(className)}
                  width={20}
                  height={20}
                  className="object-cover"
                />
              </div>
              <span>{formatDisplayClassName(className)}</span>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
