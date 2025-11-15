// File: hooks/useGameContents.ts

import type { QueryObserverResult } from '@tanstack/react-query'
import type { GameContent, GetGameContentsParams } from '@/lib/api/theSpecialStandardAPI.schemas'
import { useQuery } from '@tanstack/react-query'
import { getGameContent } from '@/lib/api/game-content'

interface UseGameContentsReturn {
  gameContents: GameContent[]
  isLoading: boolean
  error: string | null
  refetch: () => Promise<QueryObserverResult<GameContent[], Error>>
}

export function useGameContents(params?: GetGameContentsParams): UseGameContentsReturn {
  const api = getGameContent()
  
  const {
    data: gameContentsResponse,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['game-contents', params],
    queryFn: async () => {
      const response = await api.getGameContents(params)
      return response
    },
    enabled: !!params?.theme_id && !!params?.category && !!params?.question_type,
  })
  
  const gameContents = gameContentsResponse ?? []
  
  return {
    gameContents,
    isLoading,
    error: error?.message || null,
    refetch,
  }
}