'use client'

import { useParams } from 'next/navigation'
import RecentSession from '@/components/attendance/RecentSession'
import StudentSchedule from '@/components/schedule/StudentSchedule'
import { Avatar } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { useStudents } from '@/hooks/useStudents'
import { getAvatarVariant } from '@/lib/utils'

function StudentPage() {
  const params = useParams()
  const studentId = params.id as string

  const { students, isLoading } = useStudents()
  const student = students.find(s => s.id === studentId)

  const CORNER_ROUND = 'rounded-4xl'
  const PADDING = 'p-5'

  if (isLoading) {
    return (
      <div className="min-h-screen h-screen flex items-center justify-center bg-background">
        <div className="text-primary">Loading student...</div>
      </div>
    )
  }

  if (!student) {
    return (
      <div className="min-h-screen h-screen flex items-center justify-center bg-background">
        <div className="text-error">Student not found</div>
      </div>
    )
  }

  const fullName = `${student.first_name} ${student.last_name}`
  const initials = `${student.first_name[0]}.${student.last_name[0]}.`
  const avatarVariant = getAvatarVariant(student.id)

  return (
    <div className="min-h-screen h-screen flex items-center justify-center bg-background">
      <div className={`w-full h-full grid grid-rows-2 gap-8 ${PADDING}`}>
        <div className="flex gap-8">
          {/* pfp and initials */}
          <div className="flex flex-col items-center  justify-between gap-2 w-1/5">
            <div className="w-full aspect-square border-2 border-accent rounded-full">
              <Avatar
                name={fullName + student.id}
                variant={avatarVariant}
                className="w-full h-full"
              />
            </div>
            <div
              className={`w-full h-1/5 text-3xl font-bold flex items-center 
                justify-center bg-background border-border border-2 ${PADDING}
                overflow-y-auto ${CORNER_ROUND}`}
            >
              {initials}
            </div>
          </div>
          {/* student schedule */}
          <div className={`flex-[3] ${CORNER_ROUND} overflow-hidden bg-accent flex flex-col justify-between ${PADDING}`}>
            <StudentSchedule studentId={studentId} className="h-3/4" />
            <Button className="h-1/5 rounded-2xl text-lg font-bold " variant="secondary">
              View Student Schedule
            </Button>
          </div>
          <div className={`bg-accent flex-[2] flex flex-col items-center justify-between ${CORNER_ROUND} ${PADDING}`}>
            <div className={`w-full h-3/4 text-3xl font-bold flex items-center rounded-2xl bg-white ${PADDING}`}>
              <RecentSession studentId={studentId} />
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
