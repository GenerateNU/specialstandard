'use client'

import AppLayout from "@/components/AppLayout";
import Link from "next/link";
import {ArrowLeft} from "lucide-react";
import {useParams} from "next/navigation";
import {useStudents} from "@/hooks/useStudents";
import UpcomingSession from "@/components/sessions/UpcomingSession";

export default function StudentSessionHistory() {
  const params = useParams()
  const studentID = params.id as string
  const { students } = useStudents()
  const student = students.find(s => s.id === studentID)

  const initials = student ? `${student.first_name[0]}${student.last_name[0]}` : ''

  return(
    <AppLayout>
      <div className="grow bg-background flex flex-row h-screen">
        <div className="w-[85%] p-10 flex flex-col overflow-y-scroll">
          <Link href={`/student/${studentID}`}
                className="inline-flex items-center gap-2 text-secondary hover:text-primary mb-4
                           transition-colors group">
            <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
            <span className="text-sm font-medium">Back to Student Profile</span>
          </Link>
          <div className="flex items-center justify-between mb-8">
            <h1 className="text-3xl font-bold text-primary">{initials}'{initials.endsWith('S') ? '' : 's'} Sessions</h1>
          </div>

          <UpcomingSession studentId={studentID} latest={false} individuality={true} />
        </div>
      </div>
    </AppLayout>
  )
}