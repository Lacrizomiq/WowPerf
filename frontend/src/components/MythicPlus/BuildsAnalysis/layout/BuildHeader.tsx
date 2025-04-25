import { useState } from "react";
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
import { getClassIcon, getSpecIcon } from "@/utils/classandspecicons";
import { useRouter } from "next/navigation";

interface BuildHeaderProps {
  className: WowClassParam;
  spec: WowSpecParam;
}

export default function BuildHeader({ className, spec }: BuildHeaderProps) {
  const router = useRouter();
  const [selectedClass, setSelectedClass] = useState<WowClassParam>(className);
  const [selectedSpec, setSelectedSpec] = useState<WowSpecParam>(spec);

  // Format class and spec names for display
  const displayClassName =
    className.charAt(0).toUpperCase() + className.slice(1);
  const displaySpecName = spec.charAt(0).toUpperCase() + spec.slice(1);

  // Get class and spec icons
  const classIconName =
    className === "deathknight"
      ? "DeathKnight"
      : className === "demonhunter"
      ? "DemonHunter"
      : displayClassName;
  const specIconName =
    spec === "beastmastery" ? "BeastMastery" : displaySpecName;

  const classIconUrl = getClassIcon(classIconName);
  const specIconUrl = getSpecIcon(classIconName, specIconName);

  // Handle class change
  const handleClassChange = (value: string) => {
    setSelectedClass(value as WowClassParam);
    router.push(`/builds/${value}/${selectedSpec}`);
  };

  // Handle spec change
  const handleSpecChange = (value: string) => {
    setSelectedSpec(value as WowSpecParam);
    router.push(`/builds/${selectedClass}/${value}`);
  };

  return (
    <div className="flex flex-col md:flex-row items-start md:items-center justify-between mb-6 gap-4">
      <div className="flex items-center gap-4">
        <div className="w-16 h-16 rounded-md bg-slate-800 overflow-hidden flex items-center justify-center">
          <Image
            src={specIconUrl}
            alt={`${displaySpecName} ${displayClassName}`}
            width={64}
            height={64}
            className="w-full h-full object-cover"
          />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-white">
            {displaySpecName} {displayClassName}
          </h1>
          <p className="text-slate-400">
            Talents, Hero Specs, and Gear for Mythic+
          </p>
        </div>
      </div>

      <div className="flex flex-wrap gap-2">
        <Select value={selectedClass} onValueChange={handleClassChange}>
          <SelectTrigger className="w-36 bg-slate-800 border-slate-700">
            <SelectValue placeholder="Class" />
          </SelectTrigger>
          <SelectContent className="bg-slate-800 border-slate-700 text-white">
            <SelectItem value="priest">Priest</SelectItem>
            <SelectItem value="mage">Mage</SelectItem>
            <SelectItem value="warrior">Warrior</SelectItem>
            <SelectItem value="paladin">Paladin</SelectItem>
            <SelectItem value="hunter">Hunter</SelectItem>
            <SelectItem value="rogue">Rogue</SelectItem>
            <SelectItem value="deathknight">Death Knight</SelectItem>
            <SelectItem value="shaman">Shaman</SelectItem>
            <SelectItem value="warlock">Warlock</SelectItem>
            <SelectItem value="monk">Monk</SelectItem>
            <SelectItem value="druid">Druid</SelectItem>
            <SelectItem value="demonhunter">Demon Hunter</SelectItem>
            <SelectItem value="evoker">Evoker</SelectItem>
          </SelectContent>
        </Select>

        <Select value={selectedSpec} onValueChange={handleSpecChange}>
          <SelectTrigger className="w-36 bg-slate-800 border-slate-700">
            <SelectValue placeholder="Spec" />
          </SelectTrigger>
          <SelectContent className="bg-slate-800 border-slate-700 text-white">
            {className === "priest" && (
              <>
                <SelectItem value="discipline">Discipline</SelectItem>
                <SelectItem value="holy">Holy</SelectItem>
                <SelectItem value="shadow">Shadow</SelectItem>
              </>
            )}
          </SelectContent>
        </Select>

        <Select defaultValue="mythicplus" disabled>
          <SelectTrigger className="w-36 bg-slate-800 border-slate-700">
            <SelectValue placeholder="Content" />
          </SelectTrigger>
          <SelectContent className="bg-slate-800 border-slate-700 text-white">
            <SelectItem value="mythicplus">Mythic+</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}
