import React, { useState, useRef, useEffect } from "react";
import Image from "next/image";
import type { WowClass } from "@/types/warcraftlogs/dungeonRankings";
import type { ClassIconsMapping } from "@/utils/classandspecicons";

interface ClassSpecSelectorProps {
  selectedClass: WowClass | null;
  selectedSpec: string | null;
  onClassChange: (className: WowClass | null) => void;
  onSpecChange: (specName: string | null) => void;
  classMapping: ClassIconsMapping;
}

const ClassSpecSelector: React.FC<ClassSpecSelectorProps> = ({
  selectedClass,
  selectedSpec,
  onClassChange,
  onSpecChange,
  classMapping,
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [activeClass, setActiveClass] = useState<string | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const getClassColor = (className: string) => {
    const formattedClass = className
      .replace(/([A-Z])/g, "-$1")
      .toLowerCase()
      .replace(/^-/, "");
    return `class-color--${formattedClass}`;
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
        setActiveClass(null);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleClassHover = (className: string) => {
    setActiveClass(className);
  };

  const handleSpecClick = (className: WowClass, specName: string) => {
    onClassChange(className);
    onSpecChange(specName);
    setIsOpen(false);
    setActiveClass(null);
  };

  const getCurrentSelection = () => {
    if (!selectedClass) return "All Classes";
    if (!selectedSpec) return selectedClass;
    return `${selectedClass} - ${selectedSpec}`;
  };

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        className="w-[200px] h-10 px-3 bg-gradient-blue text-left text-white rounded-lg flex items-center justify-between"
        onClick={() => setIsOpen(!isOpen)}
      >
        <div className="flex items-center gap-2">
          {selectedClass && (
            <Image
              src={classMapping[selectedClass].classIcon}
              alt={selectedClass}
              width={24}
              height={24}
              className="rounded-sm"
              unoptimized
            />
          )}
          <span className={selectedClass ? getClassColor(selectedClass) : ""}>
            {getCurrentSelection()}
          </span>
        </div>
        <span>▶</span>
      </button>

      {isOpen && (
        <div className="absolute z-50 mt-1 w-[200px] bg-black border border-gray-700 rounded-lg shadow-lg">
          <div
            className="p-2 hover:bg-gradient-purple cursor-pointer flex items-center gap-2"
            onClick={() => {
              onClassChange(null);
              onSpecChange(null);
              setIsOpen(false);
            }}
          >
            <span className="w-6"></span>
            <span>All Classes</span>
            <span className="ml-auto">▶</span>
          </div>

          {Object.entries(classMapping).map(([className, data]) => (
            <div
              key={className}
              className="relative"
              onMouseEnter={() => handleClassHover(className)}
              onMouseLeave={() => !selectedClass && setActiveClass(null)}
            >
              <div className="p-2 hover:bg-gradient-purple cursor-pointer flex items-center gap-2">
                <Image
                  src={data.classIcon}
                  alt={className}
                  width={24}
                  height={24}
                  className="rounded-sm"
                  unoptimized
                />
                <span className={getClassColor(className)}>{className}</span>
                <span className="ml-auto">▶</span>
              </div>

              {activeClass === className && (
                <div className="absolute left-full top-0 ml-1 w-[200px] bg-black border border-gray-700 rounded-lg shadow-lg">
                  {Object.entries(data.spec).map(([specName, iconUrl]) => (
                    <div
                      key={specName}
                      className="p-2 hover:bg-gradient-purple cursor-pointer flex items-center gap-2"
                      onClick={() =>
                        handleSpecClick(className as WowClass, specName)
                      }
                    >
                      <Image
                        src={iconUrl}
                        alt={specName}
                        width={24}
                        height={24}
                        className="rounded-sm"
                        unoptimized
                      />
                      <span className={getClassColor(className)}>
                        {specName}
                      </span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ClassSpecSelector;
