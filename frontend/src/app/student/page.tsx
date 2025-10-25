'use client'

import RecentSession from '@/components/attendance/RecentSession'
import StudentSchedule from '@/components/schedule/StudentSchedule'
import { Avatar } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'

function StudentPage() {
  const CORNER_ROUND = 'rounded-4xl'
  const PADDING = 'p-5'

  return (
    <div className="min-h-screen h-screen flex items-center justify-center bg-background">
      <div className={`w-full h-full grid grid-rows-2 gap-8 ${PADDING}`}>
        <div className="flex gap-8">
          {/* pfp and initials - placeholder */}
          <div className="flex flex-col items-center  justify-between gap-2 w-1/5">
            <div className="w-full aspect-square border-2 border-accent rounded-full">
              <Avatar name="Student" variant="lorelei" className="w-full h-full" />
            </div>
            <div
              className={`w-full h-1/5 text-3xl font-bold flex items-center 
                justify-center bg-background border-border border-2 ${PADDING}
                overflow-y-auto ${CORNER_ROUND}`}
            >
              --
            </div>
          </div>
          {/* student schedule - no studentId, will show all sessions or empty */}
          <div className={`flex-[3] ${CORNER_ROUND} overflow-hidden bg-accent flex flex-col justify-between ${PADDING}`}>
            <StudentSchedule className="h-3/4" />
            <Button className="h-1/5 rounded-2xl text-lg font-bold " variant="secondary">
              View Student Schedule
            </Button>
          </div>
          <div className={`bg-accent flex-[2] flex flex-col items-center justify-between ${CORNER_ROUND} ${PADDING}`}>
            <div className={`w-full h-3/4 text-3xl font-bold flex items-center rounded-2xl bg-white ${PADDING}`}>
              <RecentSession />
            </div>
            <Button className="w-full h-1/5 rounded-2xl text-lg font-bold " variant="secondary">
              View Student Attendance
            </Button>
          </div>
        </div>
        {/* Student attendance */}
        <div className="grid grid-cols-2 gap-8 ">
          <div className={`bg-accent ${CORNER_ROUND}`}></div>
          <div className={`bg-accent ${CORNER_ROUND}`}></div>
        </div>
      </div>
    </div>
  )
}

export default StudentPage
