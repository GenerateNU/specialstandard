import type { QueryObserverResult } from '@tanstack/react-query'
import type { GameContent } from '@/lib/api/theSpecialStandardAPI.schemas'
import { useQuery } from '@tanstack/react-query'
import { getGameContent } from '@/lib/api/game-content'

interface UseGameContentsReturn {
  gameContents: GameContent[]
  isLoading: boolean
  error: string | null
  refetch: () => Promise<QueryObserverResult<GameContent[], Error>>
}

export function useGameContents(): UseGameContentsReturn {
  const api = getGameContent()

  const {
    data: gameContentsResponse,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['game-contents'],
    queryFn: () => api.getGameContents(),
  })

  const gameContents = gameContentsResponse ?? []

  return {
    gameContents,
    isLoading,
    error: error?.message || null,
    refetch,
  }
}
