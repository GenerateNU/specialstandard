'use client'
import moment from 'moment'
import { useStudentSessions } from '@/hooks/useStudentSessions'

interface UpcomingSessionProps {
  studentId?: string
}

export default function UpcomingSession({ studentId }: UpcomingSessionProps) {
  const { sessions, isLoading, error } = useStudentSessions(studentId || '')

  // Handle no studentId case
  if (!studentId) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-2">
        <div className="text-sm text-muted-foreground">No student selected</div>
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-sm text-muted-foreground">Loading...</div>
      </div>
    )
  }

  // Filter for upcoming sessions (future dates only)
  const now = moment()
  const upcomingSessionsWithStudentInfo = sessions
    .map(item => ({
      ...item,
      sessionData: (item as any).session
    }))
    .filter(item => moment(item.sessionData.start_datetime).isAfter(now))
    .sort((a, b) => 
      new Date(a.sessionData.start_datetime).getTime() - new Date(b.sessionData.start_datetime).getTime()
    )

  if (error || upcomingSessionsWithStudentInfo.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-2">
        <div className="text-sm text-muted-foreground">No upcoming sessions</div>
      </div>
    )
  }

  return (
    <div className="space-y-2 overflow-y-scroll h-[140px]">
      {upcomingSessionsWithStudentInfo.map((item, index) => {
        const { sessionData } = item
        const startMoment = moment(sessionData.start_datetime)
        const endMoment = moment(sessionData.end_datetime)

        return (
          <div
            key={sessionData.id}
            className="p-4 bg-card h-[90px] border-2 border-default rounded-xl flex flex-col gap-1"
          >
            <span className="font-semibold text-sm">
              {sessionData.session_name || `Session #${index + 1}`}
            </span>
            <span className="text-xs">
              {startMoment.format('dddd, MMMM D, YYYY')}
            </span>
            <span className="text-xs text-muted-foreground">
              {startMoment.format('h:mm A')} - {endMoment.format('h:mm A')}
            </span>
          </div>
        )
      })}
    </div>
  )
}
