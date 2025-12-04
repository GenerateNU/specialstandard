// File: components/statistics/GamePerformanceDashboard.tsx
import { useMemo, useState } from 'react'
import { Button } from '@/components/ui/button'
import type { GameResultWithContent } from '@/hooks/useStudentGameResults'
import {
  PerformanceByTypeChart,
  PerformanceTrendChart,
} from './GamePerformanceCharts'

interface GamePerformanceDashboardProps {
  gameResults: GameResultWithContent[]
}

export function GamePerformanceDashboard({
  gameResults,
}: GamePerformanceDashboardProps) {
  const [exerciseTypeFilter, setExerciseTypeFilter] = useState<string | null>(null)
  const [categoryFilter, setCategoryFilter] = useState<string | null>(null)
  const [difficultyFilter, setDifficultyFilter] = useState<number | null>(null)
  const [gameTypeFilter, setGameTypeFilter] = useState<string | null>(null)

  // Filter results based on active filters
  const filteredResults = useMemo(() => {
    return gameResults.filter((result) => {
      if (exerciseTypeFilter && result.exercise_type !== exerciseTypeFilter) return false
      if (categoryFilter && result.category !== categoryFilter) return false
      if (gameTypeFilter && result.game_type !== gameTypeFilter) return false
      if (difficultyFilter !== null && result.difficulty_level !== difficultyFilter) return false
      return true
    })
  }, [gameResults, exerciseTypeFilter, categoryFilter, gameTypeFilter, difficultyFilter])

  // Calculate filtered stats
  const filteredStats = useMemo(() => {
    if (filteredResults.length === 0)
      return {
        total: 0,
        completed: 0,
        completionRate: 0,
        avgTime: 0,
        totalErrors: 0,
        avgErrors: 0,
      }

    const completed = filteredResults.filter((r) => r.completed).length
    const totalTime = filteredResults.reduce((sum, r) => sum + r.time_taken_sec, 0)
    const totalErrors = filteredResults.reduce(
      (sum, r) => sum + r.count_of_incorrect_attempts,
      0
    )

    return {
      total: filteredResults.length,
      completed,
      completionRate: Math.round((completed / filteredResults.length) * 100),
      avgTime: Math.round(totalTime / filteredResults.length),
      totalErrors,
      avgErrors: Math.round(totalErrors / filteredResults.length),
    }
  }, [filteredResults])

  // Performance by exercise type
  const performanceByExerciseType = useMemo(() => {
    const byType: { [key: string]: { total: number; completed: number; timeSum: number; errorSum: number } } = {}
    for (const result of filteredResults) {
      const type = result.exercise_type ?? 'unknown'
      if (!byType[type]) {
        byType[type] = { total: 0, completed: 0, timeSum: 0, errorSum: 0 }
      }
      byType[type].total++
      if (result.completed) byType[type].completed++
      byType[type].timeSum += result.time_taken_sec
      byType[type].errorSum += result.count_of_incorrect_attempts
    }
    return Object.entries(byType).map(([type, data]) => ({
      type,
      label: type === 'game' ? 'Interactive Games' : 'PDFs',
      total: data.total,
      completed: data.completed,
      completionRate: data.total > 0 ? (data.completed / data.total) * 100 : 0,
      avgTimeSeconds: data.total > 0 ? Math.round(data.timeSum / data.total) : 0,
      avgErrors: data.total > 0 ? Math.round(data.errorSum / data.total) : 0,
    }))
  }, [filteredResults])

  // Performance trend
  const performanceTrend = useMemo(() => {
    const byDate: { [key: string]: { completed: number; failed: number; timeSum: number; errorSum: number; count: number } } = {}
    for (const result of filteredResults) {
      const date = new Date(result.created_at).toISOString().split('T')[0]
      if (!byDate[date]) {
        byDate[date] = { completed: 0, failed: 0, timeSum: 0, errorSum: 0, count: 0 }
      }
      byDate[date].count++
      if (result.completed) {
        byDate[date].completed++
      } else {
        byDate[date].failed++
      }
      byDate[date].timeSum += result.time_taken_sec
      byDate[date].errorSum += result.count_of_incorrect_attempts
    }
    return Object.entries(byDate)
      .map(([date, data]) => ({
        date,
        completed: data.completed,
        failed: data.failed,
        avgTimeSeconds: data.count > 0 ? Math.round(data.timeSum / data.count) : 0,
        totalErrors: data.errorSum,
      }))
      .sort((a, b) => a.date.localeCompare(b.date))
  }, [filteredResults])

  const StatCard = ({
    label,
    value,
    unit = '',
  }: {
    label: string
    value: number | string
    unit?: string
  }) => (
    <div className="flex flex-col gap-2 p-4 rounded-3xl bg-card border-2 border-default hover:border-primary/40 transition-all">
      <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
        {label}
      </span>
      <div className="flex items-baseline gap-2">
        <span className="text-2xl font-bold text-primary">{value}</span>
        {unit && <span className="text-sm text-muted-foreground">{unit}</span>}
      </div>
    </div>
  )

  const uniqueCategories = Array.from(
    new Set(gameResults.filter((r) => r.category).map((r) => r.category))
  ).sort()
  const uniqueGameTypes: string[] = Array.from(
    new Set(
      gameResults
        .flatMap((r) => (r.game_type ? [r.game_type as string] : []))
    )
  ).sort()
  const uniqueExerciseTypes = Array.from(
    new Set(
      gameResults
        .map((r) => r.exercise_type ?? 'unknown') // Map null/undefined to 'unknown'
        .filter(Boolean) // Remove any falsy values
    )
  ).sort()

  const uniqueDifficulties: number[] = Array.from(
    new Set(
      gameResults
        .map((r) => r.difficulty_level)
        .filter((level): level is number => typeof level === 'number')
    )
  ).sort((a, b) => a - b)

  return (
    <div className="w-full flex flex-col gap-8">
      {/* Title and Filter Section */}
      <div className="flex flex-col gap-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-primary">Game Performance Analytics</h1>
          <span className="text-sm text-muted-foreground">
            {filteredResults.length} questions
          </span>
        </div>

        {/* Filters */}
        <div className="flex flex-col gap-4 p-4 bg-card border border-default rounded-lg">
          {/* Exercise Type Filter */}
          <div className="flex flex-col gap-2">
            <span className="text-sm font-semibold text-muted-foreground">Exercise Type:</span>
            <div className="flex gap-2 flex-wrap">
              <Button
                onClick={() => setExerciseTypeFilter(null)}
                variant={exerciseTypeFilter === null ? 'default' : 'outline'}
                size="sm"
                className="rounded-full"
              >
                All
              </Button>
              {uniqueExerciseTypes.map((type) => (
                <Button
                  key={type}
                  onClick={() => setExerciseTypeFilter(type ?? null)}
                  variant={exerciseTypeFilter === type ? 'default' : 'outline'}
                  size="sm"
                  className="rounded-full"
                >
                  {type === 'game' ? 'Interactive Games' : 'PDFs'}
                </Button>
              ))}
            </div>
          </div>

          {/* Category Filter */}
          {uniqueCategories.length > 0 && (
            <div className="flex flex-col gap-2">
              <span className="text-sm font-semibold text-muted-foreground">Category:</span>
              <div className="flex gap-2 flex-wrap">
                <Button
                  onClick={() => setCategoryFilter(null)}
                  variant={categoryFilter === null ? 'default' : 'outline'}
                  size="sm"
                  className="rounded-full"
                >
                  All
                </Button>
                {uniqueCategories.map((cat) => (
                  <Button
                    key={cat}
                    onClick={() => setCategoryFilter(cat ?? null)}
                    variant={categoryFilter === cat ? 'default' : 'outline'}
                    size="sm"
                    className="rounded-full"
                  >
                    {cat
                      ?.replace(/_/g, ' ')
                      .split(' ')
                      .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
                      .join(' ')}
                  </Button>
                ))}
              </div>
            </div>
          )}

          {/* Difficulty Filter */}
          {uniqueDifficulties.length > 0 && (
            <div className="flex flex-col gap-2">
              <span className="text-sm font-semibold text-muted-foreground">Difficulty Level:</span>
              <div className="flex gap-2 flex-wrap">
                <Button
                  onClick={() => setDifficultyFilter(null)}
                  variant={difficultyFilter === null ? 'default' : 'outline'}
                  size="sm"
                  className="rounded-full"
                >
                  All
                </Button>
                {uniqueDifficulties.map((level) => (
                  <Button
                    key={level}
                    onClick={() => setDifficultyFilter(level)}
                    variant={difficultyFilter === level ? 'default' : 'outline'}
                    size="sm"
                    className="rounded-full"
                  >
                    Level {level}
                  </Button>
                ))}
              </div>
            </div>
          )}

          {/* Game Type Filter */}
          {uniqueGameTypes.length > 0 && (
            <div className="flex flex-col gap-2">
              <span className="text-sm font-semibold text-muted-foreground">Game Type:</span>
              <div className="flex gap-2 flex-wrap">
                <Button
                  onClick={() => setGameTypeFilter(null)}
                  variant={gameTypeFilter === null ? 'default' : 'outline'}
                  size="sm"
                  className="rounded-full"
                >
                  All
                </Button>
                {uniqueGameTypes.map((type) => (
                  <Button
                    key={type}
                    onClick={() => setGameTypeFilter(type)}
                    variant={gameTypeFilter === type ? 'default' : 'outline'}
                    size="sm"
                    className="rounded-full"
                  >
                    {type
                      .replace(/_/g, ' ')
                      .split(' ')
                      .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
                      .join(' ')}
                  </Button>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Overall Stats */}
      {filteredResults.length > 0 && (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3">
          <StatCard label="Total" value={filteredStats.total} />
          <StatCard label="Completed" value={filteredStats.completed} />
          <StatCard label="Rate" value={filteredStats.completionRate} unit="%" />
          <StatCard label="Avg Time" value={filteredStats.avgTime} unit="s" />
          <StatCard label="Total Errors" value={filteredStats.totalErrors} />
          <StatCard label="Avg Errors" value={filteredStats.avgErrors} />
        </div>
      )}

      {/* Charts */}
      <div className="flex flex-col gap-6">
        {performanceByExerciseType.length > 0 && (
          <PerformanceByTypeChart
            data={performanceByExerciseType}
            title="Performance by Exercise Type"
          />
        )}

        {performanceTrend.length > 0 && (
          <PerformanceTrendChart data={performanceTrend} />
        )}
      </div>

      {filteredResults.length === 0 && (
        <div className="flex items-center justify-center h-64 text-muted-foreground">
          No data matches the selected filters
        </div>
      )}
    </div>
  )
}