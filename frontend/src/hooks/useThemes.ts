import { useAuthContext } from '@/contexts/authContext'
import { getThemes as getThemesApi } from '@/lib/api/themes'
import type {
  CreateThemeInput,
  GetThemesParams,
  Theme,
} from '@/lib/api/theSpecialStandardAPI.schemas'
import type { QueryObserverResult } from '@tanstack/react-query'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

interface UseThemesOptions {
  month?: number // 1-12
  year?: number
}

interface UseThemesReturn {
  themes: Theme[]
  isLoading: boolean
  error: string | null
  refetch: () => Promise<QueryObserverResult<Theme[], Error>>
  addTheme: (theme: CreateThemeInput) => void
}

export function useThemes(options?: UseThemesOptions): UseThemesReturn {
  const queryClient = useQueryClient()
  const api = getThemesApi()
  const { userId: therapistId } = useAuthContext()

  const params: GetThemesParams | undefined = options?.month || options?.year 
    ? { month: options.month, year: options.year }
    : undefined

  const {
    data: themesResponse,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['themes', therapistId, options?.month, options?.year],
    queryFn: () => api.getThemes(params),
  })

  const themes = themesResponse ?? []

  const addThemeMutation = useMutation({
    mutationFn: (input: CreateThemeInput) => api.postThemes(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['themes', therapistId] })
    },
  })

  return {
    themes,
    isLoading,
    error: error?.message || null,
    refetch,
    addTheme: (theme: CreateThemeInput) => addThemeMutation.mutate(theme),
  }
}
