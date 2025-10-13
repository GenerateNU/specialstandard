import type { QueryObserverResult } from "@tanstack/react-query";
import type {
  Theme,
  CreateThemeInput,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getThemes as getThemesApi } from "@/lib/api/themes";

interface UseThemesReturn {
  themes: Theme[];
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<QueryObserverResult<Theme[], Error>>;
  addTheme: (theme: CreateThemeInput) => void;
}

export function useThemes(): UseThemesReturn {
  const queryClient = useQueryClient();
  const api = getThemesApi();

  const {
    data: themesResponse,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["themes"],
    queryFn: () => api.getThemes(),
  });
  const themes = themesResponse? || [];

  const addThemeMutation = useMutation({
    mutationFn: (input: CreateThemeInput) => api.postThemes(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["themes"] });
    },
  });

  return {
    themes,
    isLoading,
    error: error?.message || null,
    refetch,
    addTheme: (theme: CreateThemeInput) => addThemeMutation.mutate(theme),
  };
}
