import { useQuery } from "@tanstack/react-query";
import { getHealth as getHealthApi } from "@/lib/api/health";
import type { QueryObserverResult } from "@tanstack/react-query";

interface UseHealthReturn {
  isHealthy: boolean;
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<QueryObserverResult<boolean, Error>>;
}

export function useHealth(): UseHealthReturn {
  const api = getHealthApi();

  const {
    data: isHealthy = false,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["health"],
    queryFn: async () => {
      try {
        await api.getHealth();
        return true;
      } catch {
        return false;
      }
    },
  });

  return {
    isHealthy,
    isLoading,
    error: error?.message || null,
    refetch,
  };
}
