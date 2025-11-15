'use client'
import moment from 'moment'
import { useStudentSessions } from '@/hooks/useStudentSessions'

interface SessionNotesProps {
  studentId?: string
}

export default function SessionNotes({ studentId }: SessionNotesProps) {
  const { sessions, isLoading, error } = studentId
    ? useStudentSessions(studentId)
    : { sessions: [], isLoading: false, error: null }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-sm text-muted-foreground">Loading session notes...</div>
      </div>
    )
  }

  if (error || sessions.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-2">
        <div className="text-sm text-muted-foreground italic">
          No session notes available
        </div>
      </div>
    )
  }

  // Sort sessions by date (most recent first)
  const sortedSessions = [...sessions]
    .map(item => ({
      ...item,
      sessionData: (item as any).session
    }))
    .sort((a, b) =>
      new Date(b.sessionData.start_datetime).getTime() - new Date(a.sessionData.start_datetime).getTime()
    )


  return (
    <div className="h-full overflow-y-auto">
      {sortedSessions.map((item) => {
        const { sessionData } = item
        const startMoment = moment(sessionData.start_datetime)
        const endMoment = moment(sessionData.end_datetime)

        return (
          <div
            key={sessionData.id}
            className="w-full p-4 bg-background text-primary border-border border-b-2 last:border-b-0"
          >
            {/* Session title and date */}
            <div className="w-full flex justify-between items-start">
              <div className="font-semibold text-base">
                Session #
                {sortedSessions.length - sortedSessions.indexOf(item)}
              </div>
              <div className="text-sm opacity-90">
                {startMoment.format('MM/DD/YYYY')}
              </div>
            </div>

            {item.notes && (
              <div className="text-sm text-muted-foreground mt-2">
                {item.notes}
              </div>
            )}

            {/* Time range */}
            <div className="text-xs opacity-75 mt-2">
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
