'use client'

import { AlertCircle, ArrowLeft, BookOpen, Download, File, FileText, Gamepad2, Loader2, NotebookPen, RefreshCcw } from 'lucide-react'
import Link from 'next/link'
import { useEffect, useMemo, useState } from 'react'
import AppLayout from '@/components/AppLayout'
import { useResources } from '@/hooks/useResources'
import { useGameContents } from '@/hooks/useGameContents'
import { ResourceButton } from '@/components/curriculum/resourceButton'
import { Button } from '@/components/ui/button'
import { Dropdown } from '@/components/ui/dropdown'
import type { GameContent } from '@/lib/api/theSpecialStandardAPI.schemas'

const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December'
]

function getAvailableYears(resources: any[]): number[] {
  const years = new Set<number>()
  resources.forEach(resource => {
    if (resource.theme?.theme_year) {
      years.add(resource.theme.theme_year)
    }
  })
  return Array.from(years).sort((a, b) => b - a)
}

function groupResourcesByMonthAndWeek(resources: any[], selectedYear: number) {
  const monthsMap = new Map()
  
  MONTHS.forEach((monthName, index) => {
    const month = index + 1
    const monthKey = `${month}`
    
    monthsMap.set(monthKey, {
      month,
      year: selectedYear,
      monthName,
      weeks: new Map()
    })
  })
  
  resources.forEach(resource => {
    if (!resource.theme?.theme_month || resource.week === null || resource.week === undefined) return
    if (resource.theme.theme_year !== selectedYear) return
    
    const month = resource.theme.theme_month
    const monthKey = `${month}`
    
    if (monthsMap.has(monthKey)) {
      const monthData = monthsMap.get(monthKey)
      const weekNumber = resource.week
      
      if (!monthData.weeks.has(weekNumber)) {
        monthData.weeks.set(weekNumber, {
          weekNumber,
          resources: [],
          themeId: resource.theme?.id
        })
      }
      
      monthData.weeks.get(weekNumber).resources.push(resource)
    }
  })
  
  const monthsArray = Array.from(monthsMap.values()).map(monthData => ({
    ...monthData,
    weeks: Array.from(monthData.weeks.values()).sort((a: any, b: any) => a.weekNumber - b.weekNumber)
  }))
  
  return monthsArray
}

// Convert GameContent to a format compatible with ResourceButton
function gameContentToResource(gameContent: GameContent) {
  // Format category and question_type for display
  const formatLabel = (str: string) => {
    return str
      .split('_')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')
  }
  
  return {
    id: gameContent.id,
    title: gameContent.question, // This will be the filename like 'apple.jpg'
    type: 'pdf',
    url: null,
    resource_name: `${formatLabel(gameContent.category)} - ${formatLabel(gameContent.question_type)}`,
    metadata: {
      category: gameContent.category,
      difficulty: gameContent.difficulty_level,
      question_type: gameContent.question_type
    }
  }
}

export default function Curriculum() {
  const { resources, isLoading, error, refetch } = useResources()
  const [selectedMonth, setSelectedMonth] = useState(1)
  const [weekGameContents, setWeekGameContents] = useState<Map<string, GameContent[]>>(new Map())
  const [isLoadingGameContents, setIsLoadingGameContents] = useState(false)
  
  const availableYears = useMemo(() => {
    const years = getAvailableYears(resources)
    return years.length > 0 ? years : [new Date().getFullYear()]
  }, [resources])
  
  const [selectedYear, setSelectedYear] = useState(availableYears[0])
  
  const monthGroups = useMemo(() => groupResourcesByMonthAndWeek(resources, selectedYear), [resources, selectedYear])
  const activeMonthData = monthGroups.find(m => m.month === selectedMonth) || monthGroups[0]

  // Fetch PDF exercises for all weeks in the active month
  useEffect(() => {
    const fetchPdfExercises = async () => {
      if (!activeMonthData?.weeks || activeMonthData.weeks.length === 0) {
        setWeekGameContents(new Map())
        return
      }

      setIsLoadingGameContents(true)
      const contentMap = new Map<string, GameContent[]>()
      
      try {
        const categories = ['receptive_language', 'expressive_language', 'social_pragmatic_language', 'speech']
        const questionTypes = ['sequencing', 'following_directions', 'wh_questions', 'true_false', 'concepts_sorting', 'fill_in_the_blank', 'categorical_language', 'emotions', 'teamwork_talk', 'express_excitement_interest', 'fluency', 'articulation_s', 'articulation_l']
        
        for (const week of activeMonthData.weeks) {
          if (!week.themeId) continue
          
          const weekKey = `${week.weekNumber}`
          const weekContents: GameContent[] = []
          
          // Fetch for each category/question_type combination
          for (const category of categories) {
            for (const questionType of questionTypes) {
              try {
                const api = await import('@/lib/api/game-content').then(m => m.getGameContent())
                const response = await api.getGameContents({
                  theme_id: week.themeId,
                  category: category as any,
                  question_type: questionType as any,
                  exercise_type: 'pdf',
                  difficulty_level: 1
                })
                
                if (response && Array.isArray(response)) {
                  weekContents.push(...response.filter((item: GameContent) => item.exercise_type === 'pdf'))
                }
              } catch {
                // Silently continue if a specific combination doesn't exist
                continue
              }
            }
          }
          
          // Remove duplicates by id
          const uniqueContents = Array.from(
            new Map(weekContents.map(item => [item.id, item])).values()
          )
          contentMap.set(weekKey, uniqueContents)
        }
        
        setWeekGameContents(contentMap)
      } catch (err) {
        console.error('Error fetching PDF exercises:', err)
      } finally {
        setIsLoadingGameContents(false)
      }
    }

    fetchPdfExercises()
  }, [activeMonthData])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <Loader2 className="w-8 h-8 animate-spin text-accent mx-auto mb-4" />
          <p className="text-secondary">Loading resources...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <AppLayout>
        <div className="min-h-screen flex items-center justify-center bg-background">
          <div className="text-center max-w-md">
            <AlertCircle className="w-12 h-12 text-error mx-auto mb-4" />
            <h2 className="text-xl font-semibold text-primary mb-2">Error Loading Resources</h2>
            <p className="text-secondary mb-4">{error}</p>
            <button
              onClick={() => refetch()}
              className="px-4 py-2 bg-accent text-white rounded-lg hover:bg-accent-hover transition-colors flex items-center gap-2 mx-auto"
            >
              <RefreshCcw className="w-4 h-4" />
              Try Again
            </button>
          </div>
        </div>
      </AppLayout>
    )
  }

  return (
    <AppLayout>
      <div className="grow bg-background flex flex-row h-screen">
        <div className="w-full p-10 flex flex-col overflow-y-scroll">
          {/* Header */}
          <header className="mb-8">
            <Link
              href="/"
              className="inline-flex items-center gap-2 text-secondary hover:text-primary mb-4 transition-colors group"
            >
              <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
              <span className="text-sm font-medium">Back to Home</span>
            </Link>
            <div className="flex items-center justify-between mb-8">
              <div className='flex flex-row items-center gap-2'>
                <FileText className="w-8 h-8 text-accent" />
                <h1 className="text-3xl font-bold text-primary">Curriculum Calendar</h1>
              </div>
              <Link href="/games">
                <Button variant={'outline'} className='px-10 py-5 items-center text-xl font-serif font-bold'>
                  <Gamepad2 size={36} className='!h-6 !w-6 text-xl' />
                  Games
                </Button>
              </Link>
            </div>

            {/* Year Selector */}
            <div className="flex items-center gap-2 mb-4">
              <label className="text-sm font-medium text-secondary">Year:</label>
              <Dropdown
                value={String(selectedYear)}
                onValueChange={(value) => setSelectedYear(Number.parseInt(value))}
                items={availableYears.map((year) => ({
                  label: String(year),
                  value: String(year)
                }))}
              />
            </div>

            {/* Month Selector + Download Button */}
            <div className="flex flex-col gap-3">

              {/* Month Selector (full width, scrollable if needed) */}
              <div className="flex gap-6 overflow-x-auto pb-2 w-full">
                {MONTHS.map((monthName, index) => {
                  const monthNum = index + 1
                  const isActive = monthNum === selectedMonth

                  return (
                    <button
                      key={monthNum}
                      onClick={() => setSelectedMonth(monthNum)}
                      className={`px-1 py-2 text-sm font-medium whitespace-nowrap transition-colors ${
                        isActive
                          ? 'text-primary border-b-2 border-accent'
                          : 'text-secondary hover:text-primary'
                      }`}
                    >
                      {monthName}
                    </button>
                  )
                })}
              </div>

              {/* Right-Aligned Download Button */}
              {activeMonthData && (
                <div className="flex justify-end w-full">
                  <button className="px-4 py-2 border border-accent text-accent rounded-full hover:bg-accent hover:text-white transition-colors flex items-center gap-2 text-sm font-medium whitespace-nowrap">
                    <Download className="w-4 h-4" />
                    Download {activeMonthData.monthName} Newsletter
                  </button>
                </div>
              )}
            </div>

            {/* Year Description */}
            <p className="text-secondary text-sm mt-4">
              View and access all available learning materials for {selectedYear}.
            </p>

          </header>

          {/* Content */}
          {activeMonthData && activeMonthData.weeks.length > 0 ? (
            <div className='w-full flex flex-col gap-6'>
              {activeMonthData.weeks.map((week: any, weekIndex: number) => {
                const readings = week.resources.filter((r: any) => 
                  r.type === 'reading' || r.type === 'Passage' || r.type === 'Video'
                )
                const exercises = week.resources.filter((r: any) => 
                  r.type === 'exercise' || r.type === 'Worksheet'
                )
                const games = week.resources.filter((r: any) => 
                  r.type === 'game' || r.type === 'Game'
                )
                const other = week.resources.filter((r: any) => 
                  !r.type || (!['reading', 'Passage', 'Video', 'exercise', 'Worksheet', 'game', 'Game'].includes(r.type))
                )
                
                // Get PDF exercises from game contents
                const pdfExercises = weekGameContents.get(String(week.weekNumber)) || []

                return (
                  <div key={weekIndex} className='bg-orange-blue rounded-2xl flex flex-col p-6'>
                    <div className="flex items-center gap-3 mb-4">
                      <h4 className='text-lg font-semibold text-black'>Week {week.weekNumber}</h4>
                      {week.resources[0]?.theme && (
                        <span className='bg-orange-300 text-black text-xs font-medium px-3 py-1 rounded-full'>
                          {week.resources[0].theme.theme_name}
                        </span>
                      )}
                    </div>
                    
                    <div className='grid grid-cols-3 w-full gap-6'>
                      {/* Readings */}
                      <div className='bg-card h-full w-full flex flex-col rounded-2xl gap-3 p-6'>
                        <h5 className='font-semibold text-black'>Readings</h5>
                        {readings.length === 0 ? (
                          <div className='text-muted text-sm'>No readings available</div>
                        ) : (
                          readings.map((reading: any) => (
                            <ResourceButton key={reading.id} resource={reading} icon={BookOpen} />
                          ))
                        )}
                      </div>

                      {/* Exercises */}
                      <div className='bg-card h-full w-full flex flex-col rounded-2xl gap-3 p-6'>
                        <h5 className='font-semibold text-black'>Exercises {isLoadingGameContents && <span className='text-xs text-secondary'>(loading...)</span>}</h5>
                        {exercises.length === 0 && other.length === 0 && pdfExercises.length === 0 ? (
                          <div className='text-muted text-sm'>No exercises available</div>
                        ) : (
                          <>
                            {exercises.map((exercise: any) => (
                              <ResourceButton key={exercise.id} resource={exercise} icon={NotebookPen} />
                            ))}
                            {other.map((resource: any) => (
                              <ResourceButton key={resource.id} resource={resource} icon={File} />
                            ))}
                            {pdfExercises.map((pdfExercise: GameContent) => {
                              const formatLabel = (str: string) => {
                                return str
                                  .split('_')
                                  .map(word => word.charAt(0).toUpperCase() + word.slice(1))
                                  .join(' ')
                              }
                              
                              return (
                                <div key={pdfExercise.id} className='flex flex-col gap-1'>
                                  <ResourceButton 
                                    resource={gameContentToResource(pdfExercise)} 
                                    icon={File} 
                                  />
                                  <div className='flex gap-2 ml-1'>
                                    <span className='text-xs bg-blue-100 text-blue-700 px-2 py-1 rounded'>
                                      {formatLabel(pdfExercise.category)}
                                    </span>
                                    <span className='text-xs bg-purple-100 text-purple-700 px-2 py-1 rounded'>
                                      {formatLabel(pdfExercise.question_type)}
                                    </span>
                                  </div>
                                </div>
                              )
                            })}
                          </>
                        )}
                      </div>

                      {/* Games */}
                      <div className='bg-card h-full w-full flex flex-col rounded-2xl gap-3 p-6'>
                        <h5 className='font-semibold text-black'>Games</h5>
                        {games.length === 0 ? (
                          <div className='text-muted text-sm'>No games available</div>
                        ) : (
                          games.map((game: any) => (
                            <ResourceButton key={game.id} resource={game} icon={Gamepad2} />
                          ))
                        )}
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          ) : (
            <div className='bg-orange-disabled rounded-2xl p-6 text-center text-muted'>
              No resources scheduled for {activeMonthData?.monthName || 'this month'}
            </div>
          )}
        </div>
      </div>
    </AppLayout>
  )
}