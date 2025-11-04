'use client'

import type { Session } from '@/lib/api/theSpecialStandardAPI.schemas'
import { ChevronDown, Loader2 } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import AppLayout from '@/components/AppLayout'
import StudentSchedule from '@/components/schedule/StudentSchedule'
import { Button } from '@/components/ui/button'
import UpcomingSessionCard from '@/components/UpcomingSessionCard'
import { useAuthContext } from '@/contexts/authContext'
import { useSessions } from '@/hooks/useSessions'
import { useSessionStudentsForSession } from '@/hooks/useSessionStudents'
import { useTherapists } from '@/hooks/useTherapists'

export default function Home() {
  const { isAuthenticated, isLoading, userId } = useAuthContext()
  const router = useRouter()
  const CORNER_ROUND = 'rounded-2xl'

  // Fetch therapists data
  const { therapists } = useTherapists()
  const currentTherapist = therapists.find(t => t.id === userId)

  // Fetch all sessions for today (backend doesn't support therapist_id filtering yet)
  // TODO: Add therapist_id query param to backend API for better performance
  const { sessions: allSessions, isLoading: sessionsLoading } = useSessions({
    startdate: new Date(new Date().setHours(0, 0, 0, 0)).toISOString(), // Today at midnight
  })

  // Filter sessions by current therapist and limit to 5
  const sessions = currentTherapist
    ? allSessions.filter(s => s.therapist_id === currentTherapist.id).slice(0, 5)
    : allSessions.slice(0, 5)

  const [selectedSession, setSelectedSession] = useState<Session | null>(null)
  const [openSection, setOpenSection] = useState<'students' | 'curriculum' | 'calendar' | null>('students')

  // Fetch students for the selected session
  const { students: sessionStudents, isLoading: studentsLoading } = useSessionStudentsForSession(
    selectedSession?.id || '',
  )

  // Helper function to get therapist name
  const getTherapistName = (therapistId: string) => {
    const therapist = therapists.find(t => t.id === therapistId)
    return therapist ? `${therapist.first_name} ${therapist.last_name}` : 'Unknown Therapist'
  }

  // Helper function to format time
  const formatTime = (datetime: string) => {
    return new Date(datetime).toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true,
    })
  }

  // Helper function to format date
  const formatDate = (datetime: string) => {
    return new Date(datetime).toLocaleDateString('en-US', {
      month: '2-digit',
      day: '2-digit',
      year: 'numeric',
    })
  }

  // Set first session as selected by default
  useEffect(() => {
    if (sessions.length > 0 && !selectedSession) {
      setSelectedSession(sessions[0])
    }
  }, [sessions, selectedSession])

  // Redirect to login if not authenticated
  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login')
    }
  }, [isAuthenticated, isLoading, router])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!isAuthenticated) {
    return null
  }

  return (
    <AppLayout>
      <div className="grow bg-background flex flex-row h-screen">
        {/* Main Content */}
        <div className="w-full md:w-2/3 p-10 flex flex-col gap-10 overflow-y-scroll">
          <div className="flex flex-row justify-between items-end shrink-0">
            <div className="flex flex-col items-left justify-start">
              <div className="text-4xl font-serif font-bold">Your Educator Dashboard</div>
              <div className="text-2xl">
                {new Date().toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric' })}
              </div>

            </div>
            <div className="flex gap-2">
              <Button variant="outline">
                Download Newsletter
              </Button>
            </div>
          </div>
          <div className={`w-full bg-card p-6 gap-4 font-semibold text-xl ${CORNER_ROUND} flex flex-col`}>
            <div className="w-full flex items-center justify-between">
              <div>Upcoming Sessions</div>
              <Button
                size="sm"
                variant="default"
                onClick={() => router.push('/sessions')}
              >
                View All Sessions
              </Button>
            </div>
            <div className="grid grid-cols-[1fr_2fr] w-full gap-3 items-start">
              <div className="gap-2 flex flex-col">
                {sessionsLoading
                  ? (
                      <div className="flex items-center justify-center py-8">
                        <Loader2 className="w-6 h-6 animate-spin text-primary" />
                      </div>
                    )
                  : sessions.length > 0
                    ? (
                        sessions.map((session, index) => (
                          <div
                            key={session.id}
                            onClick={() => setSelectedSession(session)}
                            className="cursor-pointer animate-in fade-in slide-in-from-left-3 duration-300"
                            style={{ animationDelay: `${index * 75}ms` }}
                          >
                            <UpcomingSessionCard
                              className={`transition-all duration-200 ${selectedSession?.id === session.id ? 'ring-2 ring-offset-1 ring-blue-disabled scale-[1.02]' : 'hover:scale-[1.01]'}`}
                              sessionName="Session Name"
                              startTime={formatTime(session.start_datetime)}
                              endTime={formatTime(session.end_datetime)}
                              date={formatDate(session.start_datetime)}
                            />
                          </div>
                        ))
                      )
                    : (
                        <div className="flex items-center justify-center py-8 text-muted-foreground">
                          No upcoming sessions
                        </div>
                      )}
              </div>
              <div className={`p-4 text-sm font-normal ${CORNER_ROUND} justify-start flex flex-col gap-4 self-start transition-all duration-300`}>
                {selectedSession
                  ? (
                      <div className="flex flex-col gap-6 animate-in fade-in duration-300">
                        <div className="flex flex-row flex-1 ">
                          <div className="flex flex-col flex-1">
                            <strong>
                              Session Name
                            </strong>
                            <strong>
                              {getTherapistName(selectedSession.therapist_id)}
                            </strong>
                          </div>
                          <div className="flex flex-col flex-1">
                            <span>
                              {formatTime(selectedSession.start_datetime)}
                              {' '}
                              –
                              {formatTime(selectedSession.end_datetime)}
                            </span>
                            <span>
                              {formatDate(selectedSession.start_datetime)}
                            </span>
                          </div>
                        </div>
                        {selectedSession.notes && (
                          <div className="text-sm">
                            <strong>Notes:</strong>
                            {' '}
                            {selectedSession.notes}
                          </div>
                        )}
                      </div>
                    )
                  : (
                      <div className="flex items-center justify-center py-8">
                        <span>Select a session to view details</span>
                      </div>
                    )}
                {selectedSession && (
                  <div className="w-full gap-3 flex flex-col mt-2">
                    {/* Students Section */}
                    <div className="flex flex-col text-lg w-full border-b-2 border-border overflow-hidden">
                      <button
                        className="flex items-center justify-between px-2 py-2 cursor-pointer group hover:bg-card-hover transition-colors"
                        onClick={() => setOpenSection(openSection === 'students' ? null : 'students')}
                      >
                        <span className="font-semibold">Students</span>
                        <ChevronDown
                          size={18}
                          className={`transition-transform duration-300 ease-in-out ${openSection === 'students' ? 'rotate-180' : ''}`}
                        />
                      </button>
                      <div
                        className={`grid transition-all duration-300 ease-in-out ${
                          openSection === 'students' ? 'grid-rows-[1fr] opacity-100' : 'grid-rows-[0fr] opacity-0'
                        }`}
                      >
                        <div className="overflow-hidden">
                          <div className="mt-1 mb-3 text-sm font-normal space-y-1 px-2">
                            {studentsLoading
                              ? (
                                  <div className="flex items-center gap-2 text-muted">
                                    <Loader2 className="w-4 h-4 animate-spin" />
                                    Loading students...
                                  </div>
                                )
                              : sessionStudents.length > 0
                                ? (
                                    <div className="flex flex-wrap gap-x-2 gap-y-0.5">
                                      {sessionStudents.map((student, index) => (
                                        <span
                                          key={student.id}
                                          className="text-sm animate-in fade-in duration-300"
                                          style={{ animationDelay: `${index * 50}ms` }}
                                        >
                                          {student.first_name}
                                          {' '}
                                          {student.last_name}
                                          {index < sessionStudents.length - 1 && ','}
                                        </span>
                                      ))}
                                    </div>
                                  )
                                : (
                                    <div className="text-muted">No students assigned</div>
                                  )}
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Curriculum Section */}
                    <div className="flex flex-col text-lg w-full border-b-2 border-border overflow-hidden">
                      <button
                        className="flex items-center justify-between px-2 py-2 cursor-pointer group hover:bg-card-hover transition-colors"
                        onClick={() => setOpenSection(openSection === 'curriculum' ? null : 'curriculum')}
                      >
                        <span className="font-semibold">Curriculum</span>
                        <ChevronDown
                          size={18}
                          className={`transition-transform duration-300 ease-in-out ${openSection === 'curriculum' ? 'rotate-180' : ''}`}
                        />
                      </button>
                      <div
                        className={`grid transition-all duration-300 ease-in-out ${
                          openSection === 'curriculum' ? 'grid-rows-[1fr] opacity-100' : 'grid-rows-[0fr] opacity-0'
                        }`}
                      >
                        <div className="overflow-hidden">
                          <div className="mt-1 mb-3 text-sm font-normal px-2">
                            <div className="text-muted">Curriculum content coming soon...</div>
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Calendar Section */}
                    <div className="flex flex-col text-lg w-full border-b-2 border-border overflow-hidden">
                      <button
                        className="flex items-center justify-between px-2 py-2 cursor-pointer group hover:bg-card-hover transition-colors"
                        onClick={() => setOpenSection(openSection === 'calendar' ? null : 'calendar')}
                      >
                        <span className="font-semibold">Weekly Calendar</span>
                        <ChevronDown
                          size={18}
                          className={`transition-transform duration-300 ease-in-out ${openSection === 'calendar' ? 'rotate-180' : ''}`}
                        />
                      </button>
                      <div
                        className={`grid transition-all duration-300 ease-in-out ${
                          openSection === 'calendar' ? 'grid-rows-[1fr] opacity-100' : 'grid-rows-[0fr] opacity-0'
                        }`}
                      >
                        <div className="overflow-hidden">
                          <div className="mt-1 mb-3 text-sm font-normal px-2">
                            <div className="text-muted">Calendar view coming soon...</div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                )}

              </div>
            </div>
          </div>
          <div className={`w-full shrink-0 p-6 bg-card flex items-start flex-col gap-4 ${CORNER_ROUND}`}>
            <div className="w-full flex items-center justify-between font-bold text-xl">
              <div>
                Schedule
              </div>
              <Button
                size="sm"
                variant="default"
                onClick={() => router.push('/calendar')}
              >
                View Full Schedule
              </Button>
            </div>
            {/* Calendar */}
            <div className="w-full h-[500px]">
              <StudentSchedule
                initialView="day"
                className="h-full"
              />
            </div>
          </div>
        </div>
        {/* Sidebar */}
        <div className="hidden md:block w-1/3 h-screen bg-orange sticky top-0">

        </div>
      </div>
    </AppLayout>
  )
}
