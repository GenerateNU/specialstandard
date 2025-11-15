'use client'

import { use, useState } from 'react'
import { BookOpen, Dumbbell } from 'lucide-react'
import Link from 'next/link'
import CurriculumLayout from '@/components/curriculum/CurriculumLayout'
import LevelButton from '@/components/curriculum/LevelButton'
import WeekNavigator from '@/components/curriculum/WeekNavigator'
import { useSessionContext } from '@/contexts/sessionContext'
import { Button } from '@/components/ui/button'

interface PageProps {
  params: Promise<{ id: string }>
}

const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
]

export default function CurriculumPage({ params }: PageProps) {
  const { id } = use(params)
  const { 
    session, 
    currentWeek, 
    currentMonth, 
    currentYear,
    setCurrentWeek,
    setCurrentMonth,
    setCurrentYear,
    setCurrentLevel,
  } = useSessionContext()
  const [selectedLevel, setSelectedLevel] = useState<number | null>(null)

  if (!session) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div>Session not found. Please start the session first.</div>
      </div>
    )
  }

  const handlePreviousWeek = () => {
    if (currentWeek > 1) {
      setCurrentWeek(currentWeek - 1)
      setSelectedLevel(null)
    } else {
      // Go to previous month, week 4
      if (currentMonth === 0) {
        setCurrentMonth(11)
        setCurrentYear(currentYear - 1)
      } else {
        setCurrentMonth(currentMonth - 1)
      }
      setCurrentWeek(4)
      setSelectedLevel(null)
    }
  }

  const handleNextWeek = () => {
    if (currentWeek < 4) {
      setCurrentWeek(currentWeek + 1)
      setSelectedLevel(null)
    } else {
      // Go to next month, week 1
      if (currentMonth === 11) {
        setCurrentMonth(0)
        setCurrentYear(currentYear + 1)
      } else {
        setCurrentMonth(currentMonth + 1)
      }
      setCurrentWeek(1)
      setSelectedLevel(null)
    }
  }

  const handleLevelClick = (level: number) => {
    setSelectedLevel(level)
  }

  // Format date from session
  const sessionDate = new Date(session.start_datetime)
  const formattedDate = sessionDate.toLocaleDateString('en-US', {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  })

  return (
    <CurriculumLayout
      title="Curriculum"
      subtitle={formattedDate}
      backHref={`/sessions/${id}`}
      backLabel="Back to Session"
      headerContent={(
        <div className="flex items-center gap-4">
          <span className="text-lg font-medium opacity-90">
            {MONTHS[currentMonth]} {currentYear}
          </span>
          <WeekNavigator
            currentWeek={currentWeek}
            onPreviousWeek={handlePreviousWeek}
            onNextWeek={handleNextWeek}
          />
        </div>
      )}
    >
      <div className="max-w-5xl mx-auto py-8">
        {/* Level Selection Grid */}
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6 mb-12">
          {[1, 2, 3, 4, 5].map(level => (
            <div key={level} className="flex justify-center">
              <LevelButton 
                level={level} 
                onClick={() => handleLevelClick(level)} 
                isSelected={selectedLevel === level} 
              />
            </div>
          ))}
        </div>

        {/* Book Component - Shows when a level is selected */}
        {selectedLevel && (
          <div className="flex justify-center mt-12">
            <div className="bg-card rounded-3xl shadow-2xl border border-default px-12 py-16 max-w-2xl w-full">
              <div className="text-center mb-2">
                <p className="text-lg font-medium text-secondary">Week {currentWeek}</p>
              </div>
              <div className="flex items-center justify-center gap-4 mb-8">
                <BookOpen className="w-12 h-12 text-pink" />
                <h2 className="text-3xl font-bold text-primary">
                  Level {selectedLevel} Materials
                </h2>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Link
                  href={`/sessions/${id}/curriculum/reading`}
                  onClick={() => setCurrentLevel(selectedLevel)}
                  className="h-24 text-xl bg-blue hover:bg-blue-hover text-white gap-3 flex items-center justify-center rounded-lg font-semibold transition-all hover:scale-105"
                >
                  <BookOpen className="w-6 h-6" />
                  Open Reading
                </Link>
                
                <Button
                  size="lg"
                  className="h-24 text-xl bg-pink hover:bg-pink-hover text-white gap-3"
                  onClick={() => {
                    // TODO: Navigate to exercises page
                  }}
                >
                  <Dumbbell className="w-6 h-6" />
                  Exercises
                </Button>
              </div>
            </div>
          </div>
        )}
      </div>
    </CurriculumLayout>
  )
}

