'use client'

import { AlertCircle, ArrowLeft, BookOpen, FileText, Gamepad2, Loader2, NotebookPen, RefreshCcw } from 'lucide-react'
import Link from 'next/link'
import AppLayout from '@/components/AppLayout'
import { useResources } from '@/hooks/useResources'
import { ResourceButton } from '@/components/curriculum/resourceButton'

// Group resources by week number
function groupByWeek(resources: any[]) {
  const weeks = new Map()
  
  resources.forEach(resource => {
    if (resource.week === null || resource.week === undefined) return
    
    const weekNumber = resource.week
    
    if (!weeks.has(weekNumber)) {
      weeks.set(weekNumber, {
        weekNumber,
        resources: []
      })
    }
    
    weeks.get(weekNumber).resources.push(resource)
  })
  
  // Sort weeks by week number
  return Array.from(weeks.values()).sort((a, b) => a.weekNumber - b.weekNumber)
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

  const weeklyResources = groupByWeek(resources)

  return (
    <AppLayout>
      <div className="grow bg-background flex flex-row h-screen">
        <div className="w-full p-10 flex flex-col gap-10 overflow-y-scroll">
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
              <div className="flex items-center space-x-3">
                <FileText className="w-8 h-8 text-accent" />
                <h1 className="text-3xl font-bold text-primary">Curriculum Calendar</h1>
              </div>
            </div>
            <p className="text-secondary">
              View and access all available learning materials.
            </p>
          </header>

          <div className='w-full flex flex-col items-left gap-6'>

            {weeklyResources.length === 0 ? (
              <div className='bg-orange-disabled rounded-2xl p-6 text-center text-muted'>
                No resources scheduled yet
              </div>
            ) : (
              weeklyResources.map((week, index) => {
                // Map actual database types to display categories
                const readings = week.resources.filter((r: any) => 
                  r.type === 'reading' || r.type === 'Passage' || r.type === 'Video'
                )
                const exercises = week.resources.filter((r: any) => 
                  r.type === 'exercise' || r.type === 'Worksheet'
                )
                const games = week.resources.filter((r: any) => 
                  r.type === 'game' || r.type === 'Game'
                )
                
                // Get theme from first resource (assuming all in same week share theme)
                const theme = week.resources[0]?.category || week.resources[0]?.theme?.name || 'General'

                return (
                  <div key={index} className='bg-orange-disabled rounded-2xl flex flex-col p-6'>
                    <h4>Week {week.weekNumber}</h4>
                    <span>{theme}</span>
                    <div className='grid grid-cols-2 w-full gap-6 mt-4'>
                      {/* Readings */}
                      <div className='bg-white h-full w-full flex flex-col rounded-2xl gap-3 p-6'>
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
                      <div className='bg-white h-full w-full flex flex-col rounded-2xl gap-3 p-6'>
                        <h4>Exercises and Games</h4>
                        {exercises.length === 0 && games.length === 0 ? (
                          <div className='text-muted text-sm px-2'>No exercises or games available</div>
                        ) : (
                          <>
                            {exercises.map((exercise: any) => (
                              <ResourceButton key={exercise.id} resource={exercise} icon={NotebookPen} />
                            ))}
                            {games.map((game: any) => (
                              <ResourceButton key={game.id} resource={game} icon={Gamepad2} />
                            ))}
                          </>
                        )}
                      </div>
                    </div>
                  </div>
                )
              })
            )}
          </div>
        </div>
      </div>
    </AppLayout>
  )
}