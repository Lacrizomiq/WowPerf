import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: ["class"], // Permet le mode sombre basé sur une classe (généralement sur <html>)
  content: [
    // Chemins où Tailwind cherchera les classes à utiliser
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      // Vous pouvez étendre le thème ici avec des valeurs personnalisées
      backgroundImage: {
        "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
        "gradient-conic":
          "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
      },
      borderRadius: {
        // Rayons de bordure basés sur la variable CSS --radius
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      colors: {
        // Vos couleurs personnalisées existantes.
        // Assurez-vous que les variables CSS correspondantes (--color-deep-blue, etc.)
        // sont définies dans votre globals.css si vous les utilisez.
        "deep-blue": "var(--color-deep-blue)",
        "midnight-blue": "var(--color-midnight-blue)",
        "royal-blue": "var(--color-royal-blue)",
        "electric-blue": "var(--color-electric-blue)",

        // Couleurs de base du thème (shadcn/ui et wow-perf-landing)
        // Ces couleurs sont définies par des variables CSS dans votre globals.css
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        border: "hsl(var(--border))",
        input: "hsl(var(--input))", // Utilisé pour les champs de saisie
        ring: "hsl(var(--ring))", // Utilisé pour l'anneau de focus

        primary: {
          DEFAULT: "hsl(var(--primary))", // Couleur principale (ex: boutons importants)
          foreground: "hsl(var(--primary-foreground))", // Texte sur la couleur primaire
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))", // Pour les actions destructrices (ex: suppression)
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))", // Pour les textes/éléments moins importants
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))", // Souvent utilisé pour les survols ou accents interactifs
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))", // Fond des popovers
          foreground: "hsl(var(--popover-foreground))", // Texte dans les popovers
        },
        card: {
          DEFAULT: "hsl(var(--card))", // Fond des cartes
          foreground: "hsl(var(--card-foreground))", // Texte dans les cartes
        },

        // Couleurs spécifiques à la Sidebar
        // Ces noms correspondent à ceux que nous avons utilisés dans AppSidebar.tsx
        // et pointent vers les variables CSS --sidebar-* de votre globals.css
        sidebar: {
          DEFAULT: "hsl(var(--sidebar-background))", // ex: bg-sidebar
          foreground: "hsl(var(--sidebar-foreground))", // ex: text-sidebar-foreground
          primary: "hsl(var(--sidebar-primary))", // ex: bg-sidebar-primary (pour le logo)
          "primary-foreground": "hsl(var(--sidebar-primary-foreground))", // ex: text-sidebar-primary-foreground
          accent: "hsl(var(--sidebar-accent))", // ex: hover:bg-sidebar-accent (pour les items de menu)
          "accent-foreground": "hsl(var(--sidebar-accent-foreground))", // ex: hover:text-sidebar-accent-foreground
          border: "hsl(var(--sidebar-border))", // ex: border-sidebar-border (pour les séparateurs)
          ring: "hsl(var(--sidebar-ring))", // Anneau de focus dans la sidebar
        },

        // Couleurs pour les graphiques (si utilisées)
        chart: {
          "1": "hsl(var(--chart-1))",
          "2": "hsl(var(--chart-2))",
          "3": "hsl(var(--chart-3))",
          "4": "hsl(var(--chart-4))",
          "5": "hsl(var(--chart-5))",
        },
      },
      keyframes: {
        // Animations pour les composants shadcn/ui (ex: Accordion)
        "accordion-down": {
          from: { height: "0" },
          to: { height: "var(--radix-accordion-content-height)" },
        },
        "accordion-up": {
          from: { height: "var(--radix-accordion-content-height)" },
          to: { height: "0" },
        },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
      },
    },
  },
  plugins: [
    require("tailwindcss-animate"), // Plugin pour les animations (utilisé par shadcn/ui)
  ],
};

export default config;
