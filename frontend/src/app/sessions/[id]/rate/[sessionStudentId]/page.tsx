'use client'

import CurriculumLayout from '@/components/curriculum/CurriculumLayout'
import RateStudent from '@/components/rate/RateStudent'
import StudentSelector from '@/components/rate/StudentSelector'
import { useSessionContext } from '@/contexts/sessionContext'
import { useStudents } from '@/hooks/useStudents'
import { useParams } from 'next/navigation'

export default function RateStudentPage() {
  const params = useParams()
  const id = params.id as string
  const sessionStudentId = params.sessionStudentId as string
  const { students, session } = useSessionContext()
  
  const { students: studentData } = useStudents({
    ids: students.map(s => s.studentId)
  })

  // map studentId to student details
  const studentMap = new Map(studentData?.map(s => [s.id, s]) ?? [])

  if (!id || !sessionStudentId) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">Missing parameters</h1>
          <p className="text-muted-foreground">
            Session id or student id is missing from the URL.
          </p>
        </div>
      </div>
    )
  }

  const currentSessionStudent = students.find(
    (s) => s.sessionStudentId === Number(sessionStudentId)
  )

  if (students.length === 0) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">Loading session data...</div>
      </div>
    )
  }


  if (!currentSessionStudent) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">Student Not Found</h1>
          <p className="text-muted-foreground">
            The requested student could not be found in this session.
          </p>
        </div>
      </div>
    )
  }

  const currentStudent = studentMap.get(currentSessionStudent.studentId)
  
  if (!currentStudent) {
    return <div>Loading student details...</div>
  }

  const sessionDate = session ? new Date(session.start_datetime) : new Date()
  const formattedDate = sessionDate.toLocaleDateString('en-US', {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  })

  return (
    <CurriculumLayout
      title="Session Rating"
      subtitle={formattedDate}
      backHref={`/`}
      backLabel="Back to Home"
      headerContent={
        <div className="flex items-center gap-4">
        <StudentSelector
          students={students}
          currentSessionStudentId={Number(sessionStudentId)}
          sessionId={id}
        />
        </div>
      }
    >
      <div className="flex flex-col items-center justify-center py-16">
        <RateStudent
          sessionId={id}
          studentId={currentStudent.id}
          sessionStudentId={currentSessionStudent.sessionStudentId}
          firstName={currentStudent.first_name}
          lastName={currentStudent.last_name}
        />
      </div>
    </CurriculumLayout>
  )
}