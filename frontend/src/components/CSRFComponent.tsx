import { useEffect } from "react";
import { useCSRFToken } from "@/hooks/useCSRFToken";

export default function CSRFComponent() {
  const { fetchCSRFToken } = useCSRFToken();

  useEffect(() => {
    fetchCSRFToken();
  }, [fetchCSRFToken]);

  return null;
}
