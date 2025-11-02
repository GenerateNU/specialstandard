'use client'

import { CirclePlus, PencilLine } from 'lucide-react'
import { useState } from 'react'
import RecentSession from '@/components/attendance/RecentSession'
import StudentSchedule from '@/components/schedule/StudentSchedule'
import { Avatar } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'

function StudentPage() {
  const CORNER_ROUND = 'rounded-4xl'
  const PADDING = 'p-5'
  const goals = ['Goal 1', 'Goal 2', 'Goal 3', 'Goal 4', 'Goal 5']
  const [edit, setEdit] = useState(false)

  return (
    <div className="min-h-screen h-screen flex items-center justify-center bg-background">
      <div className={`w-full h-full grid grid-rows-2 gap-8 ${PADDING} relative`}>
        {/* Edit toggle button */}
        <Button
          onClick={() => setEdit(!edit)}
          className={`absolute top-1/2 right-5 z-20 w-12 h-12 p-0 ${edit ? 'bg-accent' : 'bg-accent-hover'}`}
          variant="secondary"
          size="icon"
        >
          <PencilLine size={20} />
        </Button>

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
        <div className="grid grid-cols-2 gap-8 h-1/2">
          <div className="gap-2 flex flex-col relative h-3/4">
            <div className="w-full text-3xl text-primary flex items-baseline">Goals</div>
            <div className="w-full h-full overflow-y-scroll overflow-x-hidden flex flex-col gap-2">
              {goals.map((goal, index) => (
                <div
                  key={index}
                  className={`w-full h-1/4 text-xl flex items-center 
                rounded-2xl hover:scale-99 transition bg-accent ${PADDING}`}
                >
                  {goal}
                </div>
              ))}
            </div>
            {edit && (
              <div className="flex items-center justify-center absolute -bottom-4
              -right-4 w-12 h-12 hover:scale-105 transition cursor-pointer z-10"
              >
                <CirclePlus size={36} strokeWidth={1.5} fill="white" className="text-primary" />
              </div>
            )}
          </div>
          <div className="gap-2 flex flex-col relative h-3/4">
            <div className="w-full text-3xl text-primary flex items-baseline">Session Notes</div>
            <div className="w-full flex-1 relative">
              <div className={`w-full h-full bg-accent rounded-2xl ${PADDING} overflow-y-auto`}>
                {/* Session notes content - will display notes from past sessions */}
              </div>
              {edit && (
                <div className="flex items-center justify-center absolute -bottom-4
                -right-4 w-12 h-12 hover:scale-105 transition cursor-pointer z-10"
                >
                  <CirclePlus size={36} strokeWidth={1.5} fill="white" className="text-primary" />
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default StudentPage
