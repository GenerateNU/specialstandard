import type { PostGameResultInput } from '@/lib/api/theSpecialStandardAPI.schemas'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { getGameResult } from '@/lib/api/game-result'

export function useGameResults() {
  const queryClient = useQueryClient()
  const api = getGameResult()

  const {
    data: gameResultsData = [],
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['game-results'],
    queryFn: () => api.getGameResults(),
  })

  const addGameResultMutation = useMutation({
    mutationFn: (input: PostGameResultInput) => api.postGameResults(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['game-results'] })
    },
  })

  return {
    gameResultsData,
    isLoading,
    error: error?.message || null,
    refetch,
    addGameResult: (gameResult: PostGameResultInput) =>
      addGameResultMutation.mutate(gameResult),
  }
}
