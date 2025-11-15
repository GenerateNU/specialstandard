'use client'

import { ArrowLeft, BookOpen } from 'lucide-react'
import Link from 'next/link'
import { use, useState } from 'react'
import { Button } from '@/components/ui/button'
import { useSessionContext } from '@/contexts/sessionContext'
import { useResources } from '@/hooks/useResources'

interface PageProps {
  params: Promise<{ id: string }>
}

const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
]

export default function ReadingPage({ params }: PageProps) {
  const { id } = use(params)
  const { session, currentWeek, currentMonth, currentYear, currentLevel } = useSessionContext()
  const [viewMode, setViewMode] = useState<'reading' | 'images'>('reading')

  if (!session) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div>Session not found. Please start the session first.</div>
      </div>
    )
  }

  if (!currentLevel) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div>Please select a level from the curriculum page first.</div>
      </div>
    )
  }

  // Format session date
  const sessionDate = new Date(session.start_datetime)
  const formattedDate = sessionDate.toLocaleDateString('en-US', {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  })

  // Fetch reading PDF (type: "Passage")
  const { resources: readingResources, isLoading: readingLoading } = useResources({
    grade_level: currentLevel,
    week: currentWeek,
    theme_month: currentMonth + 1, // API expects 1-12, context has 0-11
    theme_year: currentYear,
    type: 'Passage',
  })

  // Fetch images PDF (type: "Visual")
  const { resources: imagesResources, isLoading: imagesLoading } = useResources({
    grade_level: currentLevel,
    week: currentWeek,
    theme_month: currentMonth + 1, // API expects 1-12, context has 0-11
    theme_year: currentYear,
    type: 'Visual',
  })

  const readingPdfUrl = readingResources[0]?.presigned_url || null
  const imagesPdfUrl = imagesResources[0]?.presigned_url || null
  const isLoadingPdfs = readingLoading || imagesLoading

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="bg-blue text-white px-8 py-6">
        <div>
          <h1 className="text-5xl font-bold mb-2">Curriculum</h1>
          <p className="text-xl opacity-90">
            {formattedDate} | Week {currentWeek}
          </p>
        </div>
      </div>

      {/* Content Area */}
      <div className="px-8 py-6">
        {/* Navigation Bar */}
        <div className="flex items-center justify-between mb-6">
          <Link
            href={`/sessions/${id}/curriculum`}
            className="inline-flex items-center gap-2 text-primary hover:text-secondary transition-colors"
          >
            <ArrowLeft className="w-5 h-5" />
            <span className="text-lg font-medium">Back to Map</span>
          </Link>

          <div className="flex items-center gap-4">
            <BookOpen className="w-6 h-6 text-primary" />
            <h2 className="text-2xl font-bold text-primary">
              Week {currentWeek} {viewMode === 'reading' ? 'Reading' : 'Reading Images'}
            </h2>
          </div>

          <Button
            onClick={() => setViewMode(viewMode === 'reading' ? 'images' : 'reading')}
            className="bg-pink hover:bg-pink-hover text-white px-6 py-3 text-base"
          >
            {viewMode === 'reading' ? 'View Pictures' : 'View Reading'}
          </Button>
        </div>

        {/* PDF Viewer Container */}
        <div className="bg-card rounded-3xl shadow-lg border border-default p-6 max-w-7xl mx-auto">
          {isLoadingPdfs ? (
            <div className="flex items-center justify-center h-[800px] bg-card-hover rounded-xl">
              <div className="text-center">
                <BookOpen className="w-16 h-16 text-secondary mx-auto mb-4 animate-pulse" />
                <p className="text-xl text-primary font-medium mb-2">
                  Loading {viewMode === 'reading' ? 'Reading Material' : 'Reading Images'}...
                </p>
                <p className="text-secondary">
                  Level {currentLevel} • Week {currentWeek} • {MONTHS[currentMonth]} {currentYear}
                </p>
              </div>
            </div>
          ) : viewMode === 'reading' ? (
            readingPdfUrl ? (
              <iframe
                src={readingPdfUrl}
                className="w-full h-[800px] rounded-xl"
                title={`Week ${currentWeek} Reading - Level ${currentLevel}`}
              />
            ) : (
              <div className="flex items-center justify-center h-[800px] bg-card-hover rounded-xl">
                <div className="text-center">
                  <BookOpen className="w-16 h-16 text-secondary mx-auto mb-4" />
                  <p className="text-xl text-primary font-medium mb-2">
                    No Reading Material Available
                  </p>
                  <p className="text-secondary">
                    Level {currentLevel} • Week {currentWeek} • {MONTHS[currentMonth]} {currentYear}
                  </p>
                  <p className="text-sm text-muted mt-4">
                    Please add PDF to database with type: "Passage"
                  </p>
                </div>
              </div>
            )
          ) : (
            imagesPdfUrl ? (
              <iframe
                src={imagesPdfUrl}
                className="w-full h-[800px] rounded-xl"
                title={`Week ${currentWeek} Reading Images - Level ${currentLevel}`}
              />
            ) : (
              <div className="flex items-center justify-center h-[800px] bg-card-hover rounded-xl">
                <div className="text-center">
                  <BookOpen className="w-16 h-16 text-secondary mx-auto mb-4" />
                  <p className="text-xl text-primary font-medium mb-2">
                    No Reading Images Available
                  </p>
                  <p className="text-secondary">
                    Level {currentLevel} • Week {currentWeek} • {MONTHS[currentMonth]} {currentYear}
                  </p>
                  <p className="text-sm text-muted mt-4">
                    Please add PDF to database with type: "Visual"
                  </p>
                </div>
              </div>
            )
          )}
        </div>
      </div>
    </div>
  )
}

