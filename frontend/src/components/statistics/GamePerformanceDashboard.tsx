// File: components/statistics/GamePerformanceDashboard.tsx
import { useMemo, useState } from 'react'
import { Button } from '@/components/ui/button'
import type {
  CategoryPerformance,
  GameResultWithContent,
  PerformanceByType,
  PerformanceTrend,
  } from '@/hooks/useStudentGameResults'
import {
  CategoryPerformanceChart,
  CompletionRateChart,
  DifficultyPerformanceChart,
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
  const [gameTypeFilter, setGameTypeFilter] = useState<string | null>(null)

  // Filter results based on active filters
  const filteredResults = useMemo(() => {
    return gameResults.filter((result) => {
      if (exerciseTypeFilter && result.exercise_type !== exerciseTypeFilter) return false
      if (categoryFilter && result.category !== categoryFilter) return false
      if (gameTypeFilter && result.game_type !== gameTypeFilter) return false
      return true
    })
  }, [gameResults, exerciseTypeFilter, categoryFilter, gameTypeFilter])

  // Calculate filtered stats
  const filteredStats = useMemo(() => {
    const gameResults = filteredResults.filter((r) => r.exercise_type === 'game')
    const pdfResults = filteredResults.filter((r) => r.exercise_type === 'pdf')

    const calculateStats = (results: typeof filteredResults) => {
      if (results.length === 0)
        return {
          total: 0,
          completed: 0,
          completionRate: 0,
          avgTime: 0,
          totalErrors: 0,
          avgErrors: 0,
        }

      const completed = results.filter((r) => r.completed).length
      const totalTime = results.reduce((sum, r) => sum + r.time_taken_sec, 0)
      const totalErrors = results.reduce(
        (sum, r) => sum + r.count_of_incorrect_attempts,
        0
      )

      return {
        total: results.length,
        completed,
        completionRate: Math.round((completed / results.length) * 100),
        avgTime: Math.round(totalTime / results.length),
        totalErrors,
        avgErrors: Math.round(totalErrors / results.length),
      }
    }

    return {
      games: calculateStats(gameResults),
      pdfs: calculateStats(pdfResults),
      overall: calculateStats(filteredResults),
    }
  }, [filteredResults])

  const filteredPerformanceByExerciseType: PerformanceByType[] = useMemo(() => {
    const byType: { [key: string]: { total: number; completed: number; timeSum: number; errorSum: number } } = {}
    for (const result of filteredResults) {
      const type = result.exercise_type ?? 'unknown'
      if (!byType[type]) {
        byType[type] = { total: 0, completed: 0, timeSum: 0, errorSum: 0 }
      }
      byType[type].total++
      if (result.completed) {
        byType[type].completed++
      }
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

  const filteredPerformanceByGameType: PerformanceByType[] = useMemo(() => {
    const byType: { [key: string]: { total: number; completed: number; timeSum: number; errorSum: number } } = {}
    const gameResultsOnly = filteredResults.filter(
      (r) => r.exercise_type === 'game' && r.game_type
    )
    for (const result of gameResultsOnly) {
      const type = result.game_type!
      if (!byType[type]) {
        byType[type] = { total: 0, completed: 0, timeSum: 0, errorSum: 0 }
      }
      byType[type].total++
      if (result.completed) {
        byType[type].completed++
      }
      byType[type].timeSum += result.time_taken_sec
      byType[type].errorSum += result.count_of_incorrect_attempts
    }
    const labels: Record<string, string> = {
      drag_and_drop: 'Drag & Drop',
      'drag and drop': 'Drag & Drop',
      spinner: 'Spinner',
      'word/image_matching': 'Word/Image Match',
      'word/image matching': 'Word/Image Match',
      flashcards: 'Flashcards',
    }
    return Object.entries(byType).map(([type, data]) => ({
      type,
      label: labels[type] || type,
      total: data.total,
      completed: data.completed,
      completionRate: data.total > 0 ? (data.completed / data.total) * 100 : 0,
      avgTimeSeconds: data.total > 0 ? Math.round(data.timeSum / data.total) : 0,
      avgErrors: data.total > 0 ? Math.round(data.errorSum / data.total) : 0,
    }))
  }, [filteredResults])

  const filteredPerformanceTrend: PerformanceTrend[] = useMemo(() => {
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

  const filteredCategoryPerformance: CategoryPerformance[] = useMemo(() => {
    const byCategory: { [key: string]: { total: number; completed: number; timeSum: number; errorSum: number } } = {}
    for (const result of filteredResults) {
      const category = result.category ?? 'uncategorized'
      if (!byCategory[category]) {
        byCategory[category] = { total: 0, completed: 0, timeSum: 0, errorSum: 0 }
      }
      byCategory[category].total++
      if (result.completed) {
        byCategory[category].completed++
      }
      byCategory[category].timeSum += result.time_taken_sec
      byCategory[category].errorSum += result.count_of_incorrect_attempts
    }
    const labels: Record<string, string> = {
      receptive_language: 'Receptive Language',
      expressive_language: 'Expressive Language',
      social_pragmatic_language: 'Social/Pragmatic',
      speech: 'Speech',
    }
    return Object.entries(byCategory).map(([category, data]) => ({
      category,
      label: labels[category] || category,
      total: data.total,
      completed: data.completed,
      completionRate: data.total > 0 ? (data.completed / data.total) * 100 : 0,
      avgErrors: data.total > 0 ? Math.round(data.errorSum / data.total) : 0,
      avgTimeSeconds: data.total > 0 ? Math.round(data.timeSum / data.total) : 0,
    }))
  }, [filteredResults])

  const filteredDifficultyPerformance = useMemo(() => {
    const byLevel: {
      [key: number]: { total: number; completed: number; timeSum: number; errorSum: number }
    } = {}
    for (const result of filteredResults) {
      const level = result.difficulty_level
      if (level === null || level === undefined) continue

      if (!byLevel[level]) {
        byLevel[level] = { total: 0, completed: 0, timeSum: 0, errorSum: 0 }
      }
      byLevel[level].total++
      if (result.completed) {
        byLevel[level].completed++
      }
      byLevel[level].timeSum += result.time_taken_sec
      byLevel[level].errorSum += result.count_of_incorrect_attempts
    }
    return Object.entries(byLevel)
      .map(([levelStr, data]) => ({
        level: Number.parseInt(levelStr, 10),
        total: data.total,
        completed: data.completed,
        completionRate: data.total > 0 ? (data.completed / data.total) * 100 : 0,
        avgTimeSeconds: data.total > 0 ? Math.round(data.timeSum / data.total) : 0,
        avgErrors: data.total > 0 ? Math.round(data.errorSum / data.total) : 0,
      }))
      .sort((a, b) => a.level - b.level)
  }, [filteredResults])

  const StatCard = ({
    label,
    value,
    unit = '',
    variant = 'default',
  }: {
    label: string
    value: number | string
    unit?: string
    variant?: 'default' | 'games' | 'pdfs'
  }) => {
    const bgColors = {
      default: 'bg-card',
      games: 'bg-blue/5',
      pdfs: 'bg-orange/5',
    }

    return (
      <div
        className={`flex flex-col gap-2 p-4 rounded-3xl ${bgColors[variant]} border-2 border-default hover:border-primary/40 transition-all`}
      >
        <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          {label}
        </span>
        <div className="flex items-baseline gap-2">
          <span className="text-2xl font-bold text-primary">{value}</span>
          {unit && <span className="text-sm text-muted-foreground">{unit}</span>}
        </div>
      </div>
    )
  }

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
    new Set(gameResults.filter((r) => r.exercise_type).map((r) => r.exercise_type))
  ).sort()

  return (
    <div className="w-full flex flex-col gap-8">
      {/* Title and Filter Section */}
      <div className="flex flex-col gap-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-primary">Game Performance Analytics</h1>
          <span className="text-sm text-muted-foreground">
            {filteredResults.length} game sessions
          </span>
        </div>

        {/* Filter Buttons */}
        <div className="flex flex-col gap-4">
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
                  {type === 'game' ? '🎮 Games' : '📄 PDFs'}
                </Button>
              ))}
            </div>
          </div>

          {/* Category Filter */}
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
                    {type}
                  </Button>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Overall Stats */}
      <div className="flex flex-col gap-3">
        <h2 className="text-lg font-semibold text-primary">Overall Performance</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3">
          <StatCard label="Total" value={filteredStats.overall.total} />
          <StatCard label="Completed" value={filteredStats.overall.completed} />
          <StatCard label="Rate" value={filteredStats.overall.completionRate} unit="%" />
          <StatCard label="Avg Time" value={filteredStats.overall.avgTime} unit="s" />
          <StatCard label="Total Errors" value={filteredStats.overall.totalErrors} />
          <StatCard label="Avg Errors" value={filteredStats.overall.avgErrors} />
        </div>
      </div>

      {/* Separated Stats by Exercise Type */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Games Stats */}
        <div className="flex flex-col gap-3">
          <h2 className="text-lg font-semibold text-blue">🎮 Interactive Games</h2>
          <div className="grid grid-cols-3 gap-2">
            <StatCard
              label="Total"
              value={filteredStats.games.total}
              variant="games"
            />
            <StatCard
              label="Completed"
              value={filteredStats.games.completed}
              variant="games"
            />
            <StatCard
              label="Rate"
              value={filteredStats.games.completionRate}
              unit="%"
              variant="games"
            />
            <StatCard
              label="Avg Time"
              value={filteredStats.games.avgTime}
              unit="s"
              variant="games"
            />
            <StatCard
              label="Total Errors"
              value={filteredStats.games.totalErrors}
              variant="games"
            />
            <StatCard
              label="Avg Errors"
              value={filteredStats.games.avgErrors}
              variant="games"
            />
          </div>
        </div>

        {/* PDFs Stats */}
        <div className="flex flex-col gap-3">
          <h2 className="text-lg font-semibold text-orange">📄 PDF Content</h2>
          <div className="grid grid-cols-3 gap-2">
            <StatCard
              label="Total"
              value={filteredStats.pdfs.total}
              variant="pdfs"
            />
            <StatCard
              label="Completed"
              value={filteredStats.pdfs.completed}
              variant="pdfs"
            />
            <StatCard
              label="Rate"
              value={filteredStats.pdfs.completionRate}
              unit="%"
              variant="pdfs"
            />
            <StatCard
              label="Avg Time"
              value={filteredStats.pdfs.avgTime}
              unit="s"
              variant="pdfs"
            />
            <StatCard
              label="Total Errors"
              value={filteredStats.pdfs.totalErrors}
              variant="pdfs"
            />
            <StatCard
              label="Avg Errors"
              value={filteredStats.pdfs.avgErrors}
              variant="pdfs"
            />
          </div>
        </div>
      </div>

      {/* Charts Grid */}
      {filteredPerformanceByExerciseType.length > 0 && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <PerformanceByTypeChart
            data={filteredPerformanceByExerciseType}
            title="Performance by Exercise Type"
          />
          <CompletionRateChart data={filteredPerformanceByExerciseType} title="Completion Rate" />
        </div>
      )}

      {filteredPerformanceTrend.length > 0 && (
        <PerformanceTrendChart data={filteredPerformanceTrend} />
      )}

      {filteredPerformanceByGameType.length > 0 && (
        <PerformanceByTypeChart data={filteredPerformanceByGameType} title="Performance by Game Type" />
      )}

      {filteredCategoryPerformance.length > 0 && (
        <CategoryPerformanceChart data={filteredCategoryPerformance} />
      )}

      {filteredDifficultyPerformance.length > 0 && (
        <DifficultyPerformanceChart data={filteredDifficultyPerformance} />
      )}
    </div>
  )
}