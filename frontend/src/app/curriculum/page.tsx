'use client'

import { AlertCircle, ArrowLeft, BookOpen, File, FileText, Gamepad2, Loader2, NotebookPen, RefreshCcw } from 'lucide-react'
import Link from 'next/link'
import AppLayout from '@/components/AppLayout'
import { useResources } from '@/hooks/useResources'
import { ResourceButton } from '@/components/curriculum/resourceButton'
import { Button } from '@/components/ui/button'

// Group resources by theme, then by week
function groupByThemeAndWeek(resources: any[]) {
  const themes = new Map()
  
  resources.forEach(resource => {
    // Filter out resources without theme or week
    if (!resource.theme || !resource.theme_id || resource.week === null || resource.week === undefined) return
    
    const themeId = resource.theme_id
    
    if (!themes.has(themeId)) {
      themes.set(themeId, {
        themeId,
        themeName: resource.theme.theme_name,
        themeMonth: resource.theme.theme_month,
        themeYear: resource.theme.theme_year,
        weeks: new Map()
      })
    }
    
    const themeData = themes.get(themeId)
    const weekNumber = resource.week
    
    if (!themeData.weeks.has(weekNumber)) {
      themeData.weeks.set(weekNumber, {
        weekNumber,
        resources: []
      })
    }
    
    themeData.weeks.get(weekNumber).resources.push(resource)
  })
  
  // Convert to array and sort themes by date (newest to oldest)
  const themesArray = Array.from(themes.values()).sort((a, b) => {
    const dateA = new Date(a.themeYear, a.themeMonth - 1)
    const dateB = new Date(b.themeYear, b.themeMonth - 1)
    return dateB.getTime() - dateA.getTime()
  })
  
  // Sort weeks within each theme
  themesArray.forEach(themeData => {
  themeData.weeks = Array.from(themeData.weeks.values()).sort((a: any, b: any) => a.weekNumber - b.weekNumber)
})
  
  return themesArray
}

function getMonthName(month: number): string {
  const months = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December'
  ]
  return months[month - 1] || ''
}

export default function Curriculum() {
  const { resources, isLoading, error, refetch } = useResources()

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

  const themeGroups = groupByThemeAndWeek(resources)

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
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center justify-between w-full">
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
            </div>
            <p className="text-secondary">
              View and access all available learning materials.
            </p>
          </header>

          <div className='w-full flex flex-col gap-10'>
            {themeGroups.length === 0 ? (
              <div className='bg-orange-disabled rounded-2xl p-6 text-center text-muted'>
                No resources scheduled yet
              </div>
            ) : (
              themeGroups.map((themeGroup, themeIndex) => (
                <div key={themeIndex} className='flex flex-col gap-2'>
                  {/* Theme Header */}
                  <h3 className='text-2xl font-semibold px-2'>
                    {getMonthName(themeGroup.themeMonth)} {themeGroup.themeYear} - {themeGroup.themeName}
                  </h3>

                  {/* Weeks within this theme */}
                  {themeGroup.weeks.map((week: any, weekIndex: number) => {
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

                    return (
                      <div key={weekIndex} className='bg-orange-blue rounded-2xl flex flex-col p-6'>
                        <h4 className='text-black'>Week {week.weekNumber}</h4>
                        <div className='grid grid-cols-2 w-full gap-6 mt-4'>
                          {/* Readings */}
                          <div className='bg-card h-full w-full flex flex-col rounded-2xl gap-3 p-6'>
                            <h4>Readings</h4>
                            {readings.length === 0 ? (
                              <div className='text-muted text-sm px-2'>No readings available</div>
                            ) : (
                              readings.map((reading: any) => (
                                <ResourceButton key={reading.id} resource={reading} icon={BookOpen} />
                              ))
                            )}
                          </div>

                          {/* Exercises and Games */}
                          <div className='bg-card h-full w-full flex flex-col rounded-2xl gap-3 p-6'>
                            <h4>Exercises and Games</h4>
                            {exercises.length === 0 && games.length === 0 && other.length === 0 ? (
                              <div className='text-muted text-sm px-2'>No exercises or games available</div>
                            ) : (
                              <>
                                {exercises.map((exercise: any) => (
                                  <ResourceButton key={exercise.id} resource={exercise} icon={NotebookPen} />
                                ))}
                                {games.map((game: any) => (
                                  <ResourceButton key={game.id} resource={game} icon={Gamepad2} />
                                ))}
                                {other.map((resource: any) => (
                                  <ResourceButton key={resource.id} resource={resource} icon={File} />
                                ))}
                              </>
                            )}
                          </div>
                        </div>
                      </div>
                    )
                  })}
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </AppLayout>
  )
}