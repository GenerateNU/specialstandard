'use client'
import { use } from 'react'
import { useRouter } from 'next/navigation'
import CurriculumLayout from '@/components/curriculum/CurriculumLayout'
import LevelButton from '@/components/curriculum/LevelButton'
import WeekNavigator from '@/components/curriculum/WeekNavigator'
import CurriculumRoad from '@/components/curriculum/CurriculumRoad'
import { useSessionContext } from '@/contexts/sessionContext'
import { useThemes } from '@/hooks/useThemes'

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

  // Fetch theme for current month and year
  const { themes, isLoading: themesLoading } = useThemes({
    month: currentMonth + 1,
    year: currentYear,
  })

  const router = useRouter()

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
    } else {
      // Go to previous month, week 4
      if (currentMonth === 0) {
        setCurrentMonth(11)
        setCurrentYear(currentYear - 1)
      } else {
        setCurrentMonth(currentMonth - 1)
      }
      setCurrentWeek(4)
    }
  }

  const handleNextWeek = () => {
    if (currentWeek < 4) {
      setCurrentWeek(currentWeek + 1)
    } else {
      // Go to next month, week 1
      if (currentMonth === 11) {
        setCurrentMonth(0)
        setCurrentYear(currentYear + 1)
      } else {
        setCurrentMonth(currentMonth + 1)
      }
      setCurrentWeek(1)
    }
  }

  const handleLevelClick = (level: number) => {
    setCurrentLevel(level)
    router.push(`/sessions/${id}/curriculum/level/${level}`)
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
          <div className="flex flex-col gap-1">
            <span className="text-lg font-medium opacity-90">
              {MONTHS[currentMonth]} {currentYear}
            </span>
            {!themesLoading && themes.length > 0 && (
              <span className="text-sm font-semibold text-pink">
                {themes[0].name}
              </span>
            )}
          </div>
          <WeekNavigator
            currentWeek={currentWeek}
            onPreviousWeek={handlePreviousWeek}
            onNextWeek={handleNextWeek}
          />
        </div>
      )}
    >
      {/* Level Selection with Road - Full Width */}
      <div className="relative w-full h-[500px]">
        {/* SVG Road Background */}
        <CurriculumRoad />
        {/* Level Buttons Positioned on Road */}
        {[
          { level: 1, left: '30%', top: '94%' },
          { level: 2, left: '40%', top: '73%' },
          { level: 3, left: '50%', top: '51%' },
          { level: 4, left: '60%', top: '28%' },
          { level: 5, left: '70%', top: '6%' },
        ].map(({ level, left, top }) => (
          <div
            key={level}
            className="absolute transform -translate-x-1/2 -translate-y-1/2"
            style={{ left, top }}
          >
            <LevelButton
              level={level}
              onClick={() => handleLevelClick(level)}
            />
          </div>
        ))}
      </div>
    </CurriculumLayout>
  )
}