import { useState, useEffect } from "react";
import { fetchData } from "../services/api";

export function useApi<T>(endpoint: string) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    async function loadData() {
      try {
        const result = await fetchData<T>(endpoint);
        setData(result);
      } catch (e) {
        setError(e as Error);
      } finally {
        setLoading(false);
      }
    }
    loadData();
  }, [endpoint]);

  return { data, loading, error };
}
