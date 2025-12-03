'use client'

import AppLayout from "@/components/AppLayout";
import {ArrowLeft} from "lucide-react";
import {useParams, useRouter} from "next/navigation";
import {useStudents} from "@/hooks/useStudents";
import UpcomingSession from "@/components/sessions/UpcomingSession";

export default function StudentSessionHistory() {
  const params = useParams()
  const studentID = params.studentID as string
  const { students } = useStudents()
  const student = students.find(s => s.id === studentID)
  const initials = student ? `${student.first_name[0]}${student.last_name[0]}` : ''
  const router = useRouter()

  return(
    <AppLayout>
      <div className="grow bg-background flex flex-row h-screen">
        <div className="w-[85%] p-10 flex flex-col overflow-y-scroll">
          <div className="flex flex-row">
            <ArrowLeft onClick={() => router.push(`/student/${studentID}`)}
                       className="mt-1 mr-1 w-8 h-8 cursor-pointer"/>
            <div className="flex items-center justify-between mb-8">
              <h1 className="text-3xl font-bold text-primary">{initials}'{initials.endsWith('S') ? '' : 's'} Sessions</h1>
            </div>
          </div>

          <UpcomingSession studentId={studentID} latest={false} individuality={true} />
        </div>
      </div>
    </AppLayout>
  )
}