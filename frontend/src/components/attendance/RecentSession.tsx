'use client'
import moment from 'moment'
import { useStudentSessions } from '@/hooks/useStudentSessions'

interface RecentSessionProps {
  studentId?: string
}

export default function RecentSession({ studentId }: RecentSessionProps) {
  const { sessions, isLoading, error } = studentId
    ? useStudentSessions(studentId)
    : { sessions: [], isLoading: false, error: null }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-sm text-background">Loading...</div>
      </div>
    )
  }

  if (error || sessions.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-2 text-background">
        <div className="text-sm">No recent sessions</div>
      </div>
    )
  }

  // Sort sessions by date (most recent first)
  const sortedSessions = [...sessions].sort((a, b) =>
    new Date(b.start_datetime).getTime() - new Date(a.start_datetime).getTime(),
  )

  return (
    <div className="h-full  overflow-y-auto space-y-2 text-background w-full">
      {sortedSessions.map((session) => {
        const startMoment = moment(session.start_datetime)
        const endMoment = moment(session.end_datetime)

        return (
          <div
            key={session.id}
            className="p-4 bg-background border-b border-background/20 rounded-2xl flex flex-col justify-center h-20 w-full text-primary last:border-b-0"
          >
            {/* Session title and date */}
            <div className="w-full flex justify-between items-center">
              <div className="font-semibold text-base">
                Session #
                {sortedSessions.length - sortedSessions.indexOf(session)}
              </div>
              <div className="text-sm opacity-90">
                {startMoment.format('MM/DD/YYYY')}
              </div>
            </div>
            {/* Present/Absent status */}
            <div className="text-sm mb-1">
              {session.present
                ? (
                    <span className="font-medium">Present ✓</span>
                  )
                : (
                    <span className="font-medium">Absent</span>
                  )}
            </div>
            {/* Time range */}
            <div className="text-sm opacity-75">
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
