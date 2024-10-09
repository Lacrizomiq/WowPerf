import { useState, useCallback, useEffect } from "react";
import { authService } from "@/libs/authService";

export function useCSRFToken() {
  const [csrfToken, setCSRFToken] = useState<string | null>(null);

  const fetchCSRFToken = useCallback(async () => {
    try {
      const token = await authService.getCSRFToken();
      setCSRFToken(token);
      localStorage.setItem("csrfToken", token);
    } catch (error) {
      console.error("Failed to fetch CSRF token:", error);
    }
  }, []);

  useEffect(() => {
    const storedToken = localStorage.getItem("csrfToken");
    if (storedToken) {
      setCSRFToken(storedToken);
    } else {
      fetchCSRFToken();
    }
  }, [fetchCSRFToken]);

  return { csrfToken, fetchCSRFToken };
}
