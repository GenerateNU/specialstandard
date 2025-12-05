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
  content_id: string
  start_time: number
  end_time?: number
  time_taken_sec?: number
  completed: boolean
  incorrect_attempts: string[] 
  count_of_incorrect_attempts?: number
}

interface UseGameResultsProps {
  session_student_id: number
  session_id?: string
  student_id?: string
}

export function useGameResults({ session_student_id, session_id, student_id }: UseGameResultsProps) {
  const queryClient = useQueryClient()
  const api = getGameResult()
  const [results, setResults] = useState<Map<string, GameResultTracker>>(new Map())

  // Query existing results
  const queryParams: GetGameResultsParams = {}
  if (session_id) queryParams.session_id = session_id
  if (student_id) queryParams.student_id = student_id

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
    enabled: !!session_id && session_id !== 'test-session' && !!student_id && student_id !== 'test-student'
  })

  // Start tracking a game
  const startCard = useCallback((content: GameContent) => {
    setResults(prev => {
      const newMap = new Map(prev)
      newMap.set(content.id, {
        content_id: content.id,
        start_time: Date.now(),
        completed: false,
        incorrect_attempts: []
      })
      return newMap
    })
  }, [])

  // Mark game complete with optional time override
  const completeCard = useCallback((content_id: string, time_taken_sec?: number, count_of_incorrect_attempts?: number) => {
    setResults(prev => {
      const newMap = new Map(prev)
      const result = newMap.get(content_id)
      if (result) {
        const calculatedTime = time_taken_sec ?? Math.max(0, Math.floor((Date.now() - result.start_time) / 1000))
        newMap.set(content_id, {
          ...result,
          end_time: Date.now(),
          time_taken_sec: calculatedTime,
          completed: true,
          count_of_incorrect_attempts
        })
      }
      return newMap
    })
  }, [])

  // Track incorrect attempts
  const trackIncorrectAttempt = useCallback((content_id: string, attempt: string) => {
    setResults(prev => {
      const newMap = new Map(prev)
      const result = newMap.get(content_id)
      if (result) {
        newMap.set(content_id, {
          ...result,
          incorrect_attempts: [...result.incorrect_attempts, attempt]
        })
      }
      return newMap
    })
  }, [])

  // Mutation for saving single result
  const saveResultMutation = useMutation({
    mutationFn: async (input: PostGameResultInput) => {
      const response = await api.postGameResults(input)
      return response
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['game-results'] })
    }
  })

  // Save all tracked results
  const saveAllResults = useCallback(async () => {
    const allResults = Array.from(results.values())
    console.warn("🧠 Saving all results:", allResults)

    const savePromises = allResults.map(result => {
      const time_taken_sec = result.time_taken_sec ?? (result.end_time
        ? Math.max(0, Math.floor((result.end_time - result.start_time) / 1000))
        : 0)

      const incorrectAttempts = result.incorrect_attempts && result.incorrect_attempts.length > 0 
        ? result.incorrect_attempts 
        : undefined

      const input: PostGameResultInput = {
        session_student_id,
        content_id: result.content_id,
        time_taken_sec,
        completed: true,
        count_of_incorrect_attempts: result.count_of_incorrect_attempts ?? result.incorrect_attempts.length,
        incorrect_attempts: incorrectAttempts
      }

      console.warn("📤 Sending result:", input)
      console.warn("📤 JSON would be:", JSON.stringify(input))
      return saveResultMutation.mutateAsync(input)
    })

    try {
      await Promise.all(savePromises)
      setResults(new Map()) // Clear after success
    } catch (error) {
      console.error('Error saving game results:', error)
    }
  }, [results, session_student_id, saveResultMutation])

  // Get existing result for specific content
  const getResultForContent = useCallback((content_id: string) => {
    return existingResults?.find(result => result.content_id === content_id)
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
    
    // Current in-session results
    currentResults: results,
    
    // Save status
    isSaving: saveResultMutation.isPending,
    saveError: saveResultMutation.error
  }
}