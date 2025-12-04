'use client'

import { use } from 'react'
import { useRouter } from 'next/navigation'
import { BookOpen, ChevronRight } from 'lucide-react'
import Link from 'next/link'
import CurriculumLayout from '@/components/curriculum/CurriculumLayout'
import WeekNavigator from '@/components/curriculum/WeekNavigator'
import { useSessionContext } from '@/contexts/sessionContext'
import { useThemes } from '@/hooks/useThemes'
import { GetGameContentsCategory } from '@/lib/api/theSpecialStandardAPI.schemas'

interface PageProps {
  params: Promise<{ id: string; level: string }>
}

const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
]

const CATEGORIES = {
  [GetGameContentsCategory.receptive_language]: { 
    label: 'Receptive Language', 
    icon: '👂', 
    colorClass: 'bg-blue',
  },
  [GetGameContentsCategory.expressive_language]: { 
    label: 'Expressive Language', 
    icon: '💬', 
    colorClass: 'bg-pink',
  },
  [GetGameContentsCategory.social_pragmatic_language]: { 
    label: 'Social Pragmatic Language', 
    icon: '🤝', 
    colorClass: 'bg-orange',
  },
  [GetGameContentsCategory.speech]: { 
    label: 'Speech', 
    icon: '🗣️', 
    colorClass: 'bg-blue',
  },
}

export default function LevelPage({ params }: PageProps) {
  const { id, level } = use(params)
  const levelNumber = Number.parseInt(level, 10)
  const router = useRouter()
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
      if (currentMonth === 11) {
        setCurrentMonth(0)
        setCurrentYear(currentYear + 1)
      } else {
        setCurrentMonth(currentMonth + 1)
      }
      setCurrentWeek(1)
    }
  }

  const handleCategoryClick = (category: GetGameContentsCategory) => {
    setCurrentLevel(levelNumber)
    router.push(`/games?sessionId=${id}&category=${category}`)
  }

  const sessionDate = new Date(session.start_datetime)
  const formattedDate = sessionDate.toLocaleDateString('en-US', {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  })

  return (
    <CurriculumLayout
      title={`Level ${levelNumber} Materials`}
      subtitle={formattedDate}
      backHref={`/sessions/${id}/curriculum`}
      backLabel="Back to Curriculum"
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
      <div className="max-w-6xl mx-auto px-8 py-12">
        <div className="grid grid-cols-1 lg:grid-cols-2 items-start">
          {/* Left: Book with Reading Button */}
          <div className="flex flex-col items-center gap-6">
            <div className="w-full max-w-xs">
              <div className="relative w-full" style={{ aspectRatio: '201/262' }}>
                <svg
                  width="201"
                  height="262"
                  viewBox="0 0 201 262"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                  className="absolute inset-0 w-full h-full"
                >
                  <path
                    d="M50.8421 0H170.528C179.177 0 187.851 1.46012 193.907 7.94362C199.872 14.3467 201 23.2147 201 31.4798V189.548C201 196.822 200.066 203.976 196.462 209.762C194.382 213.1 191.47 215.798 188.033 217.572V229.735C188.033 238.657 186.633 247.619 180.37 253.874C174.171 260.063 165.587 261.215 157.586 261.215H19.4619C14.3033 261.215 9.35603 259.098 5.70835 255.33C2.06067 251.562 0.0114299 246.451 0.0114299 241.122V50.4614C-0.0274712 43.7904 -0.0663707 35.4583 1.69714 27.863C3.86263 18.5932 8.85493 9.60468 19.4749 4.01869C23.8577 1.71464 28.474 0.803739 33.4015 0.401869C38.1344 -3.7427e-07 43.9047 0 50.8421 0ZM19.4619 241.122H157.586C164.679 241.122 166.443 239.836 166.845 239.434C167.169 239.112 168.582 237.358 168.582 229.735V221.028H38.9124C33.7538 221.028 28.8065 223.145 25.1589 226.913C21.5112 230.682 19.4619 235.793 19.4619 241.122Z"
                    fill="#f4b860"
                  />
                </svg>
                <div className="absolute inset-0 flex flex-col items-center justify-center p-8">
                  <Link
                    href={`/sessions/${id}/curriculum/reading`}
                    onClick={() => setCurrentLevel(levelNumber)}
                    className="w-full max-w-xs h-16 bg-pink hover:bg-pink-hover text-white gap-3 flex items-center justify-center rounded-full font-semibold transition-all hover:scale-105 text-lg"
                  >
                    <BookOpen className="w-5 h-5" />
                    Open Level {levelNumber} Reading
                  </Link>
                </div>
              </div>
            </div>
          </div>

          {/* Right: Exercise Categories */}
          <div className="flex flex-col gap-4">
            {Object.entries(CATEGORIES).map(([key, category]) => (
              <button
                key={key}
                onClick={() => handleCategoryClick(key as GetGameContentsCategory)}
                className="bg-pink cursor-pointer text-white rounded-xl shadow-md p-6 hover:scale-103 hover:shadow-lg transition-all duration-200 text-left group hover:bg-pink-hover border border-default hover:border-hover"
              >
                <div className="flex items-center gap-4">
                  <div className={`w-12 h-12 ${category.colorClass} rounded-full flex items-center justify-center text-xl text-white`}>
                    {category.icon}
                  </div>
                  <div className="flex-1">
                    <h3 className="font-semibold text-white">{category.label}</h3>
                  </div>
                  <ChevronRight className="w-5 h-5 text-muted group-hover:text-primary transition-colors" />
                </div>
              </button>
            ))}
          </div>
        </div>
      </div>
    </CurriculumLayout>
  )
}