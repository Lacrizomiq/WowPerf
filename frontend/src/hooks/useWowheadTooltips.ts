import { useEffect } from "react";

export function useWowheadTooltips() {
  useEffect(() => {
    const loadWowheadTooltips = () => {
      const script = document.createElement("script");
      script.src = "https://wow.zamimg.com/widgets/power.js";
      script.async = true;
      document.body.appendChild(script);

      script.onload = () => {
        if (window.$WowheadPower) {
          window.$WowheadPower.refreshLinks();
        } else {
          console.error("Wowhead Power object not found");
        }
      };
    };

    if (
      !document.querySelector(
        'script[src="https://wow.zamimg.com/widgets/power.js"]'
      )
    ) {
      loadWowheadTooltips();
    } else if (window.$WowheadPower) {
      window.$WowheadPower.refreshLinks();
    }
  }, []);
}
