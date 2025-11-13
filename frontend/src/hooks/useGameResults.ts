// File: hooks/useGameResults.ts

import { useCallback, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { getGameResult } from '@/lib/api/game-result'
import type { 
  GameContent, 
  GetGameResultsParams,
  PostGameResultInput 
} from '@/lib/api/theSpecialStandardAPI.schemas'

interface GameResultTracker {
  contentId: string
  startTime: number
  endTime?: number
  completed: boolean
  incorrectAttempts: string[]
}

interface UseGameResultsProps {
  sessionStudentId: number
  sessionId?: string
  studentId?: string
}

export function useGameResults({ sessionStudentId, sessionId, studentId }: UseGameResultsProps) {
  const queryClient = useQueryClient()
  const api = getGameResult()
  const [results, setResults] = useState<Map<string, GameResultTracker>>(new Map())

  // Query to get existing game results
  const queryParams: GetGameResultsParams = {}
  if (sessionId) queryParams.session_id = sessionId
  if (studentId) queryParams.student_id = studentId

  const { 
    data: existingResults, 
    isLoading: isLoadingResults,
    refetch: refetchResults 
  } = useQuery({
    queryKey: ['game-results', queryParams],
    queryFn: async () => {
      const response = await api.getGameResults(queryParams)
      return response
    },
    enabled: !!sessionId || !!studentId,
  })

  // Start tracking a card
  const startCard = useCallback((content: GameContent) => {
    setResults(prev => {
      const newMap = new Map(prev)
      newMap.set(content.id, {
        contentId: content.id,
        startTime: Date.now(),
        completed: false,
        incorrectAttempts: []
      })
      return newMap
    })
  }, [])

  // Complete a card
  const completeCard = useCallback((contentId: string) => {
    setResults(prev => {
      const newMap = new Map(prev)
      const result = newMap.get(contentId)
      if (result) {
        newMap.set(contentId, {
          ...result,
          endTime: Date.now(),
          completed: true
        })
      }
      return newMap
    })
  }, [])

  // Track incorrect attempt (for quiz mode)
  const trackIncorrectAttempt = useCallback((contentId: string, attempt: string) => {
    setResults(prev => {
      const newMap = new Map(prev)
      const result = newMap.get(contentId)
      if (result) {
        newMap.set(contentId, {
          ...result,
          incorrectAttempts: [...result.incorrectAttempts, attempt]
        })
      }
      return newMap
    })
  }, [])

  // Mutation to save a single result
  const saveResultMutation = useMutation({
    mutationFn: async (input: PostGameResultInput) => {
      const response = await api.postGameResults(input)
      return response
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['game-results'] })
    }
  })

  // Save all results
  const saveAllResults = useCallback(async () => {
    const allResults = Array.from(results.values())
    const savePromises = allResults.map(result => {
      const timeTaken = result.endTime 
        ? Math.floor((result.endTime - result.startTime) / 1000)
        : 0
      
      const input: PostGameResultInput = {
        session_student_id: sessionStudentId,
        content_id: result.contentId,
        time_taken_sec: timeTaken,
        completed: result.completed,
        count_of_incorrect_attempts: result.incorrectAttempts.length,
        incorrect_attempts: result.incorrectAttempts
      }
      
      return saveResultMutation.mutateAsync(input)
    })
    
    try {
      await Promise.all(savePromises)
      setResults(new Map()) // Clear tracked results after successful save
    } catch (error) {
      console.error('Error saving game results:', error)
    }
  }, [results, sessionStudentId, saveResultMutation])

  // Get result for a specific content
  const getResultForContent = useCallback((contentId: string) => {
    return existingResults?.find(result => result.content_id === contentId)
  }, [existingResults])

  return {
    // Tracking functions
    startCard,
    completeCard,
    trackIncorrectAttempt,
    saveAllResults,
    
    // Query data
    existingResults: existingResults || [],
    isLoadingResults,
    refetchResults,
    getResultForContent,
    
    // Current session tracking
    currentResults: results,
    
    // Save status
    isSaving: saveResultMutation.isPending,
    saveError: saveResultMutation.error
  }
}
