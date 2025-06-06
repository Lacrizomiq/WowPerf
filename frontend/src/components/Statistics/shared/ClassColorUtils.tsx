// components/Statistics/shared/ClassColorUtils.tsx

/**
 * Utilitaires pour appliquer les couleurs des classes WoW
 * Gère automatiquement la conversion des noms d'API vers les classes CSS
 */

/**
 * Convertit un nom de classe API vers la classe CSS correspondante
 *
 * @param className - Nom de la classe depuis l'API (ex: "Death Knight", "Demon Hunter", "Warrior")
 * @returns Classe CSS correspondante (ex: "class-color--death-knight")
 *
 * @example
 * ```typescript
 * getClassColor("Death Knight") // → "class-color--death-knight"
 * getClassColor("Warrior") // → "class-color--warrior"
 * getClassColor("Demon Hunter") // → "class-color--demon-hunter"
 * ```
 */
export const getClassColor = (className: string): string => {
  if (!className) return "";

  const normalizedClass = className
    .toLowerCase()
    .replace(/\s+/g, "-") // Remplace les espaces par des tirets
    .trim();

  return `class-color--${normalizedClass}`;
};

/**
 * Applique directement la couleur de classe à un élément React
 *
 * @param className - Nom de la classe depuis l'API
 * @param additionalClasses - Classes CSS additionnelles à appliquer
 * @returns Chaîne de classes CSS complète
 *
 * @example
 * ```tsx
 * <span className={applyClassColor("Warrior", "font-bold text-lg")}>
 *   Protection Warrior
 * </span>
 * ```
 */
export const applyClassColor = (
  className: string,
  additionalClasses: string = ""
): string => {
  const colorClass = getClassColor(className);
  return `${colorClass} ${additionalClasses}`.trim();
};

/**
 * Composant React pour afficher du texte avec la couleur de classe
 *
 * @param className - Nom de la classe WoW
 * @param children - Contenu à afficher
 * @param additionalClasses - Classes CSS supplémentaires
 *
 * @example
 * ```tsx
 * <ClassColoredText className="Death Knight" additionalClasses="font-bold">
 *   Blood Death Knight
 * </ClassColoredText>
 * ```
 */
interface ClassColoredTextProps {
  className: string;
  children: React.ReactNode;
  additionalClasses?: string;
  as?: keyof JSX.IntrinsicElements;
}

export const ClassColoredText: React.FC<ClassColoredTextProps> = ({
  className,
  children,
  additionalClasses = "",
  as: Component = "span",
}) => {
  const colorClass = applyClassColor(className, additionalClasses);

  return <Component className={colorClass}>{children}</Component>;
};

/**
 * Hook personnalisé pour les couleurs de classe
 * Utile quand vous avez besoin de la classe CSS dans la logique du composant
 *
 * @param className - Nom de la classe WoW
 * @returns Objet contenant les utilitaires de couleur
 *
 * @example
 * ```tsx
 * const MyComponent = ({ spec }) => {
 *   const { colorClass, applyColor } = useClassColor(spec.class);
 *
 *   return (
 *     <div className={applyColor("p-4 rounded")}>
 *       <span className={colorClass}>{spec.display}</span>
 *     </div>
 *   );
 * };
 * ```
 */
export const useClassColor = (className: string) => {
  const colorClass = getClassColor(className);

  const applyColor = (additionalClasses: string = "") =>
    `${colorClass} ${additionalClasses}`.trim();

  return {
    colorClass,
    applyColor,
  };
};

/**
 * Mapping des classes WoW pour validation et debug
 * Peut être utilisé pour des listes déroulantes ou de la validation
 */
export const WOW_CLASSES = [
  "Death Knight",
  "Demon Hunter",
  "Druid",
  "Evoker",
  "Hunter",
  "Mage",
  "Monk",
  "Paladin",
  "Priest",
  "Rogue",
  "Shaman",
  "Warlock",
  "Warrior",
] as const;

export type WoWClass = (typeof WOW_CLASSES)[number];

/**
 * Vérifie si une chaîne est une classe WoW valide
 *
 * @param className - Nom à vérifier
 * @returns true si c'est une classe WoW valide
 */
export const isValidWoWClass = (className: string): className is WoWClass => {
  return WOW_CLASSES.includes(className as WoWClass);
};

/**
 * Mapping des rôles vers leurs couleurs/styles
 */
export const ROLE_STYLES = {
  tank: {
    color: "text-blue-400",
    bgColor: "bg-blue-500/20",
    borderColor: "border-blue-500",
  },
  healer: {
    color: "text-green-400",
    bgColor: "bg-green-500/20",
    borderColor: "border-green-500",
  },
  dps: {
    color: "text-red-400",
    bgColor: "bg-red-500/20",
    borderColor: "border-red-500",
  },
} as const;

/**
 * Applique les styles de rôle
 *
 * @param role - Rôle (tank, healer, dps)
 * @param styleType - Type de style à appliquer
 * @returns Classe CSS correspondante
 */
export const getRoleStyle = (
  role: keyof typeof ROLE_STYLES,
  styleType: keyof typeof ROLE_STYLES.tank
): string => {
  return ROLE_STYLES[role]?.[styleType] || "";
};
