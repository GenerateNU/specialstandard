'use client'

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
} from '@/components/ui/select'
import { Avatar } from '@/components/ui/avatar'
import { ChevronDown } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useStudents } from '@/hooks/useStudents'
import type { StudentTuple } from '@/contexts/sessionContext'
import { getAvatarName, getAvatarVariant } from '@/lib/avatarUtils'

interface StudentSelectorProps {
  students: StudentTuple[]
  currentSessionStudentId: number
  sessionId: string
}

export default function StudentSelector({
  students,
  currentSessionStudentId,
  sessionId,
}: StudentSelectorProps) {
  const router = useRouter()
  
  // Fetch full student data for all students in this session
  const { students: studentData } = useStudents({ 
    ids: students.map(s => s.studentId) 
  })

  // Map student tuples to full student objects
  const studentMap = new Map(studentData?.map(s => [s.id, s]) ?? [])
  
  const enrichedStudents = students
    .map(s => ({
      sessionStudentId: s.sessionStudentId,
      student: studentMap.get(s.studentId),
    }))
    .filter((s): s is { sessionStudentId: number; student: NonNullable<typeof s.student> } => 
      s.student !== undefined
    )

  const handleValueChange = (value: string) => {
    router.push(`/sessions/${sessionId}/rate?id=${sessionId}&sessionStudentId=${value}`)
  }

  const currentStudent = enrichedStudents.find(
    (s) => s.sessionStudentId === currentSessionStudentId
  )?.student

  if (!currentStudent) {
    return null
  }

  const avatarVariant = getAvatarVariant(currentStudent.id)

  return (
    <Select value={currentSessionStudentId.toString()} onValueChange={handleValueChange}>
      <SelectTrigger className="w-[280px] h-14 rounded-full bg-white border border-border shadow-sm hover:bg-gray-50 transition-colors group">
        <div className="flex items-center gap-3 flex-1">
          <Avatar
            name={getAvatarName(currentStudent.first_name, currentStudent.last_name, currentStudent.id)}
            variant={avatarVariant}
            className="w-8 h-8"
          />
          <span className="font-semibold text-sm text-foreground">
            {currentStudent.first_name}
          </span>
        </div>
        <ChevronDown className="h-4 w-4 text-muted-foreground transition-transform duration-200 group-data-[state=open]:rotate-180 shrink-0 ml-2" />
      </SelectTrigger>
      <SelectContent className="bg-white border border-border">
        {enrichedStudents.map((s) => (
          <SelectItem 
            key={s.sessionStudentId} 
            value={s.sessionStudentId.toString()}
            className="hover:bg-gray-50"
          >
            <div className="flex items-center gap-2">
              <Avatar
                name={getAvatarName(s.student.first_name, s.student.last_name, s.student.id)}
                variant={getAvatarVariant(s.student.id)}
                className="h-6 w-6"
              />
              <span>{s.student.first_name}</span>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}