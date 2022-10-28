import { useState, useEffect, useRef } from "react";

import axios from "../lib/axios";

const useFetchSession = (apiUrl: string) => {
  const [csrfToken, setCsrfToken] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const didEffect = useRef(false);

  useEffect(() => {
    if (!didEffect.current) {
      didEffect.current = true;
      const getSession = async () => {
        try {
          const { headers } = await axios.get(`${apiUrl}/session`);
          if (headers["x-csrf-token"] !== undefined) {
            setCsrfToken(headers["x-csrf-token"]);
          }
          setIsLoading(false);
        } catch (error) {
          setError(error);
          setIsLoading(false);
        }
      };
      getSession();
    }
  }, [apiUrl]);
  return { csrfToken, isLoading, error };
};

export default useFetchSession;
