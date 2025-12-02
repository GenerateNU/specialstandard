'use client'
import moment from 'moment'
import { useStudentSessions } from '@/hooks/useStudentSessions'
import {Badge} from "@/components/ui/badge";
import {getSchoolColor} from "@/lib/utils";
import {yellow} from "ansis";
import {useSessionStudents} from "@/hooks/useSessionStudents";
import {useSession} from "@/hooks/useSessions";
import {useStudents} from "@/hooks/useStudents";

interface UpcomingSessionProps {
  studentId?: string
  count?: number
  latest: boolean
  individuality?: boolean
}
//
// function determineIndividuality(sessionID: string) {
//   // const { sessions } = useStudents({ sessionID })
// }

export default function UpcomingSession({ studentId, count, latest, individuality = false }: UpcomingSessionProps) {
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
            className="p-4 bg-card border-2 border-default rounded-[32px] flex flex-col justify-center min-h-20 w-full shadow-md"
          >
            <div className="ml-2 my-4">
              <div className="flex flex-row">
                  <div>
                      {/* Session title and relative time */}
                      <div className="w-full flex justify-between items-center mb-1">
                        <div className="font-semibold text-base">
                          <h3>{sessionData.session_name || `Session #${index + 1}`}</h3>
                        </div>
                        {/* <div className="text-sm text-muted-foreground">
                          {daysUntil === 0 ? 'Today' : daysUntil === 1 ? 'Tomorrow' : `In ${daysUntil} days`}
                        </div> */}
                      </div>

                      {/* Date & Time Range */}
                      <div className="text-md text-muted-foreground">
                        {startMoment.format('dddd, MMMM D, YYYY')}
                        {' | '}
                        {startMoment.format('h:mm A')}
                        {' - '}
                        {endMoment.format('h:mm A')}
                      </div>
                  </div>
                  <div>
                      {/*<Badge className={}>*/}

                      {/*</Badge>*/}



                      <Badge className="">
                          {sessionData.id}
                      </Badge>
                  </div>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
