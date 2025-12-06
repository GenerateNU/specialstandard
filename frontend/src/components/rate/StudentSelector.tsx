// --- RateStudentSelector.tsx (Updated) ---

'use client'

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
} from '@/components/ui/select'
import { Avatar } from '@/components/ui/avatar'
import { useRouter } from 'next/navigation'
import { useStudents } from '@/hooks/useStudents'
import type { StudentTuple } from '@/contexts/sessionContext'
import { useSessionContext } from '@/contexts/sessionContext' // <-- Import context
import { getAvatarName, getAvatarVariant } from '@/lib/avatarUtils'

interface RateStudentSelectorProps {
  students: StudentTuple[]
  currentSessionStudentId: number
  sessionId: string
}

export default function RateStudentSelector({
  students,
  currentSessionStudentId,
  sessionId,
}: RateStudentSelectorProps) {
  const router = useRouter()
  const { attendance } = useSessionContext() 
  
  const presentStudents = students.filter(s => attendance[s.sessionStudentId] !== false) 
  
  // Fetch full student data for all *present* students in this session
  // Extract the isLoading state
  const { students: studentData, isLoading } = useStudents({ // <-- Get isLoading
    ids: presentStudents.map(s => s.studentId)
  })
  
  // 1. ADD LOADING CHECK HERE
  if (isLoading) {
    // Return a simple loading state or null while the data fetches
    return <div className="w-[280px] h-14 bg-gray-100 animate-pulse rounded-full mt-3.5" />;
  }

  // Map present student tuples to full student objects
  const studentMap = new Map(studentData?.map(s => [s.id, s]) ?? [])
  
  // ... (rest of enrichedStudents calculation remains the same)
  const enrichedStudents = presentStudents 
    .map(s => ({
      sessionStudentId: s.sessionStudentId,
      student: studentMap.get(s.studentId),
    }))
    .filter((s): s is { sessionStudentId: number; student: NonNullable<typeof s.student> } => 
      s.student !== undefined
    )

  const handleValueChange = (value: string) => {
    router.push(`/sessions/${sessionId}/rate/${value}`)
  }

  // Ensure the current student is still available (they should be present if they are being rated)
  const currentStudent = enrichedStudents.find(
    (s) => s.sessionStudentId === currentSessionStudentId
  )?.student

  if (!currentStudent) {
    // This case should ideally not happen if the URL is pointing to a present student
    // but we can fall back or return null if the current student is filtered out.
    return null
  }

  const avatarVariant = getAvatarVariant(currentStudent.id)

  return (
    <Select value={currentSessionStudentId.toString()} onValueChange={handleValueChange}>
      <SelectTrigger className="w-[280px] h-14 rounded-full bg-white border border-border mt-3.5 shadow-sm hover:bg-gray-50 transition-colors group">
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
      </SelectTrigger>
      <SelectContent className="bg-white border border-border">
        {/* Only present students are mapped into enrichedStudents */}
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