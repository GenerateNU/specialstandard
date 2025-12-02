// File: hooks/useStudentGameResults.ts
import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getGameResult } from '@/lib/api/game-result'
import { getGameContent } from '@/lib/api/game-content'
import type { GameContent, GameResult } from '@/lib/api/theSpecialStandardAPI.schemas'

// Merged type combining GameResult and GameContent
export interface GameResultWithContent extends GameResult {
  category?: string
  exercise_type?: string
  difficulty_level?: number
  game_type?: string
}

export interface GamePerformanceStats {
  totalGames: number
  completedGames: number
  completionRate: number
  avgTimeSeconds: number
  totalErrors: number
  avgErrors: number
}

export interface PerformanceByType {
  type: string
  label: string
  total: number
  completed: number
  completionRate: number
  avgTimeSeconds: number
  avgErrors: number
}

export interface PerformanceTrend {
  date: string
  completed: number
  failed: number
  avgTimeSeconds: number
  totalErrors: number
}

export interface CategoryPerformance {
  category: string
  total: number
  completed: number
  completionRate: number
  avgErrors: number
  avgTimeSeconds: number
}

export function useStudentGameResults(studentId: string) {
  const gameResultApi = getGameResult()
  const gameContentApi = getGameContent()

  // Fetch game results
  const { data: gameResults = [], isLoading: resultsLoading } = useQuery({
    queryKey: ['game-results', studentId],
    queryFn: async () => {
      const response = await gameResultApi.getGameResults({
        student_id: studentId,
        limit: 1000,
      })
      return Array.isArray(response) ? response : []
    },
    enabled: !!studentId && studentId !== 'test-student',
  })

  // Fetch all game contents
  const { data: gameContents = [], isLoading: contentsLoading } = useQuery({
    queryKey: ['game-contents'],
    queryFn: async () => {
      const response = await gameContentApi.getGameContents({
        question_count: 1000,
        words_count: 2,
      })
      return Array.isArray(response) ? response : []
    },
  })

  // Merge game results with game contents
  const mergedGameResults = useMemo(() => {
    const contentMap = new Map<string, GameContent>()
    gameContents.forEach((content) => {
      contentMap.set(content.id, content)
    })

    return gameResults.map((result) => {
      const content = contentMap.get(result.content_id)
      return {
        ...result,
        category: content?.category,
        exercise_type: content?.exercise_type,
        difficulty_level: content?.difficulty_level,
        game_type: content?.applicable_game_types?.[0], // Use first applicable game type
      }
    })
  }, [gameResults, gameContents])

  // Overall stats
  const overallStats = useMemo(() => {
    if (mergedGameResults.length === 0) {
      return {
        totalGames: 0,
        completedGames: 0,
        completionRate: 0,
        avgTimeSeconds: 0,
        totalErrors: 0,
        avgErrors: 0,
      }
    }

    const completed = mergedGameResults.filter((r) => r.completed).length
    const totalTime = mergedGameResults.reduce((sum, r) => sum + r.time_taken_sec, 0)
    const totalErrors = mergedGameResults.reduce(
      (sum, r) => sum + r.count_of_incorrect_attempts,
      0
    )

    return {
      totalGames: mergedGameResults.length,
      completedGames: completed,
      completionRate: (completed / mergedGameResults.length) * 100,
      avgTimeSeconds: Math.round(totalTime / mergedGameResults.length),
      totalErrors,
      avgErrors: Math.round(totalErrors / mergedGameResults.length),
    }
  }, [mergedGameResults])

  // Performance by exercise type
  const performanceByExerciseType = useMemo(() => {
    const typeMap = new Map<string, GameResultWithContent[]>()

    mergedGameResults.forEach((result) => {
      const type = result.exercise_type || 'unknown'
      if (!typeMap.has(type)) {
        typeMap.set(type, [])
      }
      typeMap.get(type)!.push(result)
    })

    return Array.from(typeMap.entries()).map(([type, results]) => {
      const completed = results.filter((r) => r.completed).length
      const totalTime = results.reduce((sum, r) => sum + r.time_taken_sec, 0)
      const totalErrors = results.reduce(
        (sum, r) => sum + r.count_of_incorrect_attempts,
        0
      )

      return {
        type,
        label: type === 'game' ? 'Interactive Games' : 'PDFs',
        total: results.length,
        completed,
        completionRate: (completed / results.length) * 100,
        avgTimeSeconds: Math.round(totalTime / results.length),
        avgErrors: Math.round(totalErrors / results.length),
      }
    })
  }, [mergedGameResults])

  // Performance by game type
  const performanceByGameType = useMemo(() => {
    const typeMap = new Map<string, GameResultWithContent[]>()

    mergedGameResults.forEach((result) => {
      const type = result.game_type || 'unknown'
      if (!typeMap.has(type)) {
        typeMap.set(type, [])
      }
      typeMap.get(type)!.push(result)
    })

    return Array.from(typeMap.entries())
      .map(([type, results]) => {
        const completed = results.filter((r) => r.completed).length
        const totalTime = results.reduce((sum, r) => sum + r.time_taken_sec, 0)
        const totalErrors = results.reduce(
          (sum, r) => sum + r.count_of_incorrect_attempts,
          0
        )

        const labels: Record<string, string> = {
          drag_and_drop: 'Drag & Drop',
          'drag and drop': 'Drag & Drop',
          spinner: 'Spinner',
          'word/image_matching': 'Word/Image Match',
          'word/image matching': 'Word/Image Match',
          flashcards: 'Flashcards',
        }

        return {
          type,
          label: labels[type] || type,
          total: results.length,
          completed,
          completionRate: (completed / results.length) * 100,
          avgTimeSeconds: Math.round(totalTime / results.length),
          avgErrors: Math.round(totalErrors / results.length),
        }
      })
      .sort((a, b) => b.total - a.total)
  }, [mergedGameResults])

  // Performance trend over time
  const performanceTrend = useMemo(() => {
    const dateMap = new Map<string, GameResultWithContent[]>()

    mergedGameResults.forEach((result) => {
      const date = new Date(result.created_at).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
      })
      if (!dateMap.has(date)) {
        dateMap.set(date, [])
      }
      dateMap.get(date)!.push(result)
    })

    return Array.from(dateMap.entries())
      .sort(
        (a, b) =>
          new Date(a[0]).getTime() - new Date(b[0]).getTime()
      )
      .map(([date, results]) => {
        const completed = results.filter((r) => r.completed).length
        const failed = results.length - completed
        const totalTime = results.reduce((sum, r) => sum + r.time_taken_sec, 0)
        const totalErrors = results.reduce(
          (sum, r) => sum + r.count_of_incorrect_attempts,
          0
        )

        return {
          date,
          completed,
          failed,
          avgTimeSeconds: Math.round(totalTime / results.length),
          totalErrors,
        }
      })
  }, [mergedGameResults])

  // Category performance (language categories)
  const categoryPerformance = useMemo(() => {
    const categoryMap = new Map<string, GameResultWithContent[]>()

    mergedGameResults.forEach((result) => {
      const category = result.category || 'unknown'
      if (!categoryMap.has(category)) {
        categoryMap.set(category, [])
      }
      categoryMap.get(category)!.push(result)
    })

    const labels: Record<string, string> = {
      receptive_language: 'Receptive Language',
      expressive_language: 'Expressive Language',
      social_pragmatic_language: 'Social/Pragmatic',
      speech: 'Speech',
    }

    return Array.from(categoryMap.entries())
      .map(([category, results]) => {
        const completed = results.filter((r) => r.completed).length
        const totalTime = results.reduce((sum, r) => sum + r.time_taken_sec, 0)
        const totalErrors = results.reduce(
          (sum, r) => sum + r.count_of_incorrect_attempts,
          0
        )

        return {
          category,
          label: labels[category] || category,
          total: results.length,
          completed,
          completionRate: (completed / results.length) * 100,
          avgErrors: Math.round(totalErrors / results.length),
          avgTimeSeconds: Math.round(totalTime / results.length),
        }
      })
      .sort((a, b) => b.total - a.total)
  }, [mergedGameResults])

  // Difficulty level performance
  const difficultyPerformance = useMemo(() => {
    const diffMap = new Map<number, GameResultWithContent[]>()

    mergedGameResults.forEach((result) => {
      const diff = result.difficulty_level || 1
      if (!diffMap.has(diff)) {
        diffMap.set(diff, [])
      }
      diffMap.get(diff)!.push(result)
    })

    return Array.from(diffMap.entries())
      .sort((a, b) => a[0] - b[0])
      .map(([level, results]) => {
        const completed = results.filter((r) => r.completed).length
        const totalTime = results.reduce((sum, r) => sum + r.time_taken_sec, 0)
        const totalErrors = results.reduce(
          (sum, r) => sum + r.count_of_incorrect_attempts,
          0
        )

        return {
          level,
          total: results.length,
          completed,
          completionRate: (completed / results.length) * 100,
          avgTimeSeconds: Math.round(totalTime / results.length),
          avgErrors: Math.round(totalErrors / results.length),
        }
      })
  }, [mergedGameResults])

  return {
    gameResults: mergedGameResults,
    isLoading: resultsLoading || contentsLoading,
    error: null,
    overallStats,
    performanceByExerciseType,
    performanceByGameType,
    performanceTrend,
    categoryPerformance,
    difficultyPerformance,
  }
}