import { useState, useEffect } from "react";
import api from "@/libs/api";

// useCSRFToken is a hook that fetches the CSRF token from the backend
export const useCSRFToken = () => {
  const [csrfToken, setCSRFToken] = useState<string | null>(null);

  useEffect(() => {
    const fetchCSRFToken = async () => {
      try {
        const response = await api.get("/auth/csrf");
        setCSRFToken(response.data.csrf_token);
      } catch (error) {
        console.error("Failed to fetch CSRF token:", error);
      }
    };

    fetchCSRFToken();
  }, []);

  return csrfToken;
};
