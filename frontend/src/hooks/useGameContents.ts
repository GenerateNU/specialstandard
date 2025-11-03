import type { GetGameContentsCategory } from '@/lib/api/theSpecialStandardAPI.schemas'
import { useQuery } from '@tanstack/react-query'
import { getGameContent } from '@/lib/api/game-content'

export function useGameContents(category: GetGameContentsCategory, level: number, count: number) {
  const api = getGameContent()

  const {
    data: gameContentsData = [],
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['game-contents', category, level, count],
    queryFn: () => api.getGameContents({ category, level, count }),
    staleTime: 0,
    refetchOnMount: 'always',
    refetchOnWindowFocus: false,
  })

  return {
    gameContentsData,
    isLoading,
    error: error?.message || null,
    refetch,
  }
}
