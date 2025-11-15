'use client'

import { ArrowLeft, ArrowRight } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { use, useEffect, useRef, useState } from 'react'
import { useSessionContext } from '@/contexts/sessionContext'
import { useSession } from '@/hooks/useSessions'
import { useSessionStudentsForSession } from '@/hooks/useSessionStudents'

interface PageProps {
  params: Promise<{ id: string }>
}

export default function StartSessionPage({ params }: PageProps) {
  const { id } = use(params)
  const router = useRouter()
  const { session, isLoading: sessionLoading } = useSession(id)
  const { students: sessionStudents, isLoading: studentsLoading } = useSessionStudentsForSession(id)
  const { setSession, setStudents, setCurrentWeek } = useSessionContext()
  const [selectedWeek, setSelectedWeek] = useState(1)
  const initializedRef = useRef(false)

  // Available weeks - 4 weeks per month
  const availableWeeks = [1, 2, 3, 4]
  
  // Get month name from session date
  const monthName = session ? new Date(session.start_datetime).toLocaleDateString('en-US', { month: 'long' }) : ''

  useEffect(() => {
    if (session && sessionStudents && !initializedRef.current) {
      initializedRef.current = true
      setSession(session)
      const studentTuples = sessionStudents.map(s => ({
        studentId: s.id,
        sessionStudentId: s.session_student_id,
      }))
      setStudents(studentTuples)
    }
  }, [session, sessionStudents, setSession, setStudents])

  const handleStartCurriculum = () => {
    setCurrentWeek(selectedWeek)
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
        <h1 className="text-4xl font-bold mb-8 text-center">{monthName}</h1>

        {/* Week Selection */}
        <div className="bg-card rounded-3xl p-8 shadow-lg border border-default">
          <h2 className="text-2xl font-semibold mb-6 text-center">Select Week</h2>
          
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
            {availableWeeks.map(week => (
              <button
                key={week}
                onClick={() => setSelectedWeek(week)}
                className={`
                  p-6 rounded-2xl font-semibold text-lg transition-all
                  ${selectedWeek === week
                    ? 'bg-blue text-white scale-105 shadow-lg'
                    : 'bg-card-hover text-primary hover:bg-blue-light hover:scale-102'
                  }
                `}
              >
                Week {week}
              </button>
            ))}
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
                  <div className="w-10 h-10 rounded-full bg-blue-light flex items-center justify-center font-semibold text-blue">
                    {student.first_name?.[0]}{student.last_name?.[0]}
                  </div>
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

          {/* Start Button */}
          <div className="flex justify-center">
            <button
              type="button"
              onClick={handleStartCurriculum}
              className="bg-blue text-white px-12 py-6 text-xl rounded-2xl hover:bg-blue-hover transition-all hover:scale-105 flex items-center gap-3 font-semibold"
            >
              Start Week {selectedWeek}
              <ArrowRight className="w-6 h-6" />
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

