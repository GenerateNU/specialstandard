'use client'
import moment from 'moment'
import { useStudentSessions } from '@/hooks/useStudentSessions'

interface UpcomingSessionProps {
  studentId?: string
  count?: number
  latest: boolean
}

export default function UpcomingSession({ studentId, count, latest }: UpcomingSessionProps) {
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
  let upcomingSessionsWithStudentInfo
  if (latest) {
    upcomingSessionsWithStudentInfo = sessions
      .map(item => ({ ...item, sessionData: (item as any).session }))
      .filter(item => moment(item.sessionData.start_datetime).isAfter(now))
      .sort((a, b) =>
        new Date(a.sessionData.start_datetime).getTime() - new Date(b.sessionData.start_datetime).getTime()
      )
  } else {
    upcomingSessionsWithStudentInfo = sessions
      .map(item => ({ ...item, sessionData: (item as any).session }))
      .filter(item => moment(item.sessionData.start_datetime).isBefore(now))
      .sort((a, b) =>
        new Date(b.sessionData.start_datetime).getTime() - new Date(a.sessionData.start_datetime).getTime()
      )
  }
  if (count && upcomingSessionsWithStudentInfo.length > count) {
    upcomingSessionsWithStudentInfo = upcomingSessionsWithStudentInfo.slice(0, count)
  }

  if (error || upcomingSessionsWithStudentInfo.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-2">
        <div className="text-sm text-muted-foreground">No upcoming sessions</div>
      </div>
    )
  }

  return (
    <div className="h-full overflow-y-auto space-y-2 w-full p-2">
      {upcomingSessionsWithStudentInfo.map((item, index) => {
        const { sessionData } = item
        const startMoment = moment(sessionData.start_datetime)
        const endMoment = moment(sessionData.end_datetime)
        // const daysUntil = startMoment.diff(now, 'days')

        return (
          <div
            key={sessionData.id}
            className="p-4 bg-card border-2 border-default rounded-[32px] flex flex-col justify-center min-h-20 w-full shadow-card"
          >
            {/* Session title and relative time */}
            <div className="w-full flex justify-between items-center mb-2">
              <div className="font-semibold text-base">
                {sessionData.session_name || `Session #${index + 1}`}
              </div>
              {/* <div className="text-sm text-muted-foreground">
                {daysUntil === 0 ? 'Today' : daysUntil === 1 ? 'Tomorrow' : `In ${daysUntil} days`}
              </div> */}
            </div>
            
            {/* Date */}
            <div className="text-sm mb-1">
              {startMoment.format('dddd, MMMM D, YYYY')}
            </div>
            
            {/* Time range */}
            <div className="text-sm text-muted-foreground">
              {startMoment.format('h:mm A')}
              {' - '}
              {endMoment.format('h:mm A')}
            </div>
          </div>
        )
      })}
    </div>
  )
}
