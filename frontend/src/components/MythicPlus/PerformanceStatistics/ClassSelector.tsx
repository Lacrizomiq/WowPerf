import React, { useState, useRef, useEffect } from "react";
import Image from "next/image";
import { ChevronsUpDown } from "lucide-react";
import { getClassIcon } from "@/utils/classandspecicons";

// Define the props interface
interface ClassSelectorProps {
  selectedClass: string | null;
  onClassChange: (className: string | null) => void;
  availableClasses: string[];
}

const ClassSelector: React.FC<ClassSelectorProps> = ({
  selectedClass,
  onClassChange,
  availableClasses,
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Helper to get class color CSS class
  const getClassColor = (className: string) => {
    const formattedClass = className
      .replace(/([A-Z])/g, "-$1")
      .toLowerCase()
      .replace(/^-/, "");
    return `class-color--${formattedClass} font-bold`;
  };

  // Handle clicks outside the dropdown to close it
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  // Handle selecting a class
  const handleClassClick = (className: string | null) => {
    onClassChange(className);
    setIsOpen(false);
  };

  // Get the current selection text
  const getCurrentSelection = () => {
    if (!selectedClass) return "All Classes";
    return selectedClass;
  };

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        className="w-[200px] h-9 px-3 bg-gradient-blue text-left text-white rounded-lg flex items-center justify-between"
        onClick={() => setIsOpen(!isOpen)}
        aria-haspopup="true"
        aria-expanded={isOpen}
      >
        <div className="flex items-center gap-2">
          {selectedClass && (
            <Image
              src={getClassIcon(selectedClass)}
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
        <span>
          <ChevronsUpDown className="w-3 h-3 text-gray-300" />
        </span>
      </button>

      {isOpen && (
        <div className="absolute z-50 mt-1 w-[200px] bg-black border border-gray-700 rounded-lg shadow-lg">
          <div
            className="p-2 hover:bg-gradient-purple cursor-pointer flex items-center gap-2"
            onClick={() => handleClassClick(null)}
          >
            <span className="w-6"></span>
            <span>All Classes</span>
          </div>

          {availableClasses.map((className) => (
            <div
              key={className}
              className="p-2 hover:bg-gray-800 cursor-pointer flex items-center gap-2"
              onClick={() => handleClassClick(className)}
            >
              <Image
                src={getClassIcon(className)}
                alt={className}
                width={24}
                height={24}
                className="rounded-sm"
                unoptimized
              />
              <span className={getClassColor(className)}>{className}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ClassSelector;
