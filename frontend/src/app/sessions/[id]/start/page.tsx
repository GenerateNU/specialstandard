'use client'

import { ArrowLeft, ArrowRight, ChevronLeft, ChevronRight } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { use, useEffect, useRef, useState } from 'react'
import { useSessionContext } from '@/contexts/sessionContext'
import { useSession } from '@/hooks/useSessions'
import { useSessionStudentsForSession } from '@/hooks/useSessionStudents'
import { useThemes } from '@/hooks/useThemes'
import { getAvatarName, getAvatarVariant } from '@/lib/avatarUtils'
import { Avatar } from '@/components/ui/avatar'

interface PageProps {
  params: Promise<{ id: string }>
}

const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
]

export default function StartSessionPage({ params }: PageProps) {
  const { id } = use(params)
  const router = useRouter()
  const { session, isLoading: sessionLoading } = useSession(id)
  const { students: sessionStudents, isLoading: studentsLoading } = useSessionStudentsForSession(id)
  const { setSession, setStudents, setCurrentWeek, setCurrentMonth, setCurrentYear } = useSessionContext()
  const [selectedWeek, setSelectedWeek] = useState(1)
  const initializedRef = useRef(false)

  // Initialize month and year from current date
  const now = new Date()
  const [selectedMonth, setSelectedMonth] = useState(now.getMonth())
  const [selectedYear, setSelectedYear] = useState(now.getFullYear())

  // Fetch themes for the selected month and year
  // Note: useThemes expects month as 1-12, but selectedMonth is 0-11
  const { themes, isLoading: themesLoading, error: themesError } = useThemes({
    month: selectedMonth + 1,
    year: selectedYear,
  })

  // Available weeks - 4 weeks per month
  const availableWeeks = [1, 2, 3, 4]
  const studentTuples = sessionStudents?.map(s => ({
    studentId: s.id,
    sessionStudentId: s.session_student_id,
  })) || []

  useEffect(() => {
    if (session && studentTuples.length > 0 && !initializedRef.current) {
      initializedRef.current = true
      setSession(session)
      setStudents(studentTuples)
    }
  }, [session, studentTuples, setSession, setStudents])

  const handlePreviousMonth = () => {
    if (selectedMonth === 0) {
      setSelectedMonth(11)
      setSelectedYear(selectedYear - 1)
    } else {
      setSelectedMonth(selectedMonth - 1)
    }
  }

  const handleNextMonth = () => {
    if (selectedMonth === 11) {
      setSelectedMonth(0)
      setSelectedYear(selectedYear + 1)
    } else {
      setSelectedMonth(selectedMonth + 1)
    }
  }

  const handleStartCurriculum = () => {
    setCurrentWeek(selectedWeek)
    setCurrentMonth(selectedMonth)
    setCurrentYear(selectedYear)
    router.push(`/sessions/${id}/curriculum`)
  }

  if (sessionLoading || studentsLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div>Loading session...</div>
      </div>
    )
  }

  if (!session) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div>Session not found</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background p-8">
      {/* Back button */}
      <Link
        href={`/sessions/${id}`}
        className="inline-flex items-center gap-2 text-secondary hover:text-primary mb-6 transition-colors group"
      >
        <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
        <span className="text-sm font-medium">Back to Session</span>
      </Link>

      {/* Main content */}
      <div className="max-w-4xl mx-auto">
        {/* Week Selection */}
        <div className="flex flex-col gap-8">
          <h1 className=' w-full text-center '>Choose a week to begin!</h1>
          
          {/* Month/Year Selector */}
          <div className="flex flex-col items-center justify-center gap-6">
            <div className="flex items-center justify-center gap-4">
              <button
                onClick={handlePreviousMonth}
                aria-label="Previous month"
                title="Previous month"
                className="w-12 h-12 cursor-pointer rounded-full hover:bg-card-hover flex items-center justify-center transition-colors"
              >
                <ChevronLeft className="w-6 h-6" />
              </button>
              
              <h1 className="text-4xl font-bold text-center min-w-[300px]">
                {MONTHS[selectedMonth]} {selectedYear}
              </h1>
              
              <button
                onClick={handleNextMonth}
                aria-label="Next month"
                title="Next month"
                className="w-12 cursor-pointer h-12 rounded-full hover:bg-card-hover flex items-center justify-center transition-colors"
              >
                <ChevronRight className="w-6 h-6" />
              </button>
            </div>

            {/* Theme Display */}
            <div className="w-full text-center">
              {themesLoading ? (
                <p className="text-secondary text-sm">Loading theme...</p>
              ) : themesError ? (
                <p className="text-red-500 text-sm">Error loading theme</p>
              ) : themes.length > 0 ? (
                <div className="inline-block bg-card-hover rounded-xl p-4 px-6">
                  <p className="text-xs text-secondary uppercase tracking-wide mb-2">Theme for this month</p>
                  <p className="text-2xl font-semibold text-primary">{themes[0].name}</p>
                </div>
              ) : (
                <p className="text-secondary text-sm">No theme set for this month</p>
              )}
            </div>
          </div>

          <div className="grid grid-cols-2 gap-8">
            {availableWeeks.map(week => (
              <button
                key={week}
                onClick={() => setSelectedWeek(week)}
                className={`
                  p-6 rounded-2xl font-semibold text-lg transition-all cursor-pointer
                  hover:scale-102
                  ${selectedWeek === week
                    ? 'bg-pink text-white shadow-lg'
                    : 'bg-card-hover text-primary hover:bg-blue-light '
                  }
                `}
              >
                Week {week}
              </button>
            ))}
          </div>

           <div className="flex justify-center">
            <button
              type="button"
              onClick={handleStartCurriculum}
              className="bg-orange cursor-pointer text-white px-12 py-6 text-xl rounded-2xl hover:bg-orange-hover transition-all hover:scale-105 flex items-center gap-3 font-semibold"
            >
              Start Week {selectedWeek} – {MONTHS[selectedMonth]} {selectedYear}
              <ArrowRight className="w-6 h-6" />
            </button>
          </div>

          {/* Student List */}
          <div className="mb-8">
            <h3 className="text-xl font-semibold mb-4">Students in this session:</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {sessionStudents.map(student => (
                <div
                  key={student.id}
                  className="bg-card-hover rounded-xl p-4 flex items-center gap-3"
                >
                  <Avatar
                    name={getAvatarName(student.first_name, student.last_name, student.id)}
                    variant={getAvatarVariant(student.id)}
                    className="w-12 h-12"
                  />
                  <div>
                    <p className="font-medium">
                      {student.first_name} {student.last_name}
                    </p>
                    <p className="text-sm text-secondary">
                      {student.present ? 'Present' : 'Absent'}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}