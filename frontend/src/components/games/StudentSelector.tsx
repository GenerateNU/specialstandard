'use client'

import { useSessionContext } from '@/contexts/sessionContext'
import { useStudents } from '@/hooks/useStudents'
import { Users } from 'lucide-react'
import { useState } from 'react'
import { getAvatarName, getAvatarVariant } from '@/lib/avatarUtils'
import { Avatar } from '../ui/avatar'

interface StudentSelectorProps {
  onBack: () => void
  onStudentsSelected: (studentIds: string[]) => void
  gameTitle: string
}

export function StudentSelector({ 
  onBack, 
  onStudentsSelected,
  gameTitle 
}: StudentSelectorProps) {
  const { students: sessionStudents, attendance } = useSessionContext() // <-- Get attendance
  const [selectedStudentIds, setSelectedStudentIds] = useState<string[]>([])
  const { students: allStudents, isLoading } = useStudents()

  // 1. FILTER: Only keep students who are marked as present
  const presentSessionStudents = sessionStudents.filter(s => attendance[s.sessionStudentId] !== false)

  // 2. Get full student details for *present* session students
  const studentsInSession = presentSessionStudents // <-- Use filtered list
    .map(({ studentId, sessionStudentId }) => {
      const student = allStudents?.find(s => s.id === studentId)
      return student ? {
        ...student,
        sessionStudentId
      } : null
    })
    .filter((s): s is NonNullable<typeof s> => s !== null)

  const toggleStudent = (sessionStudentId: number) => {
    const id = sessionStudentId.toString()
    setSelectedStudentIds(prev => 
      prev.includes(id) 
        ? prev.filter(sid => sid !== id)
        : [...prev, id]
    )
  }

  const handleStartGame = () => {
    if (selectedStudentIds.length > 0) {
      // NOTE: This component passes back the sessionStudentId (as a string), 
      // which is correctly the unique identifier needed for tracking within the session.
      onStudentsSelected(selectedStudentIds)
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue mx-auto mb-4"></div>
          <p className="text-secondary">Loading students...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-2xl mx-auto">
        <button
          onClick={onBack}
          className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
        >
          ← Back
        </button>

        <div className="bg-card rounded-lg shadow-lg p-8 border border-default">
          <div className="flex items-center gap-3 mb-6">
            <Users className="w-8 h-8 text-blue" />
            <div>
              <h1 className="text-2xl font-bold">{gameTitle}</h1>
              <p className="text-secondary text-sm">Select students who will play</p>
            </div>
          </div>

          {studentsInSession.length === 0 ? (
            <div className="text-center py-8">
              <p className="text-secondary mb-4">No students are marked as Present for this session.</p>
              <button
                onClick={onBack}
                className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover transition-colors"
              >
                Go Back
              </button>
            </div>
          ) : (
            <>
              <div className="space-y-3 mb-6">
                {/* This map automatically only includes PRESENT students */}
                {studentsInSession.map((student) => ( 
                  <label
                    key={student.sessionStudentId}
                    className={`flex items-center gap-3 p-4 rounded-lg border-2 cursor-pointer transition-all ${
                      selectedStudentIds.includes(student.sessionStudentId.toString())
                        ? 'border-blue bg-blue/5'
                        : 'border-default hover:border-hover hover:bg-card-hover'
                    }`}
                  >
                    <input
                      type="checkbox"
                      checked={selectedStudentIds.includes(student.sessionStudentId.toString())}
                      onChange={() => toggleStudent(student.sessionStudentId)}
                      className="w-5 h-5 text-blue rounded focus:ring-blue focus:ring-2"
                    />
                    <Avatar
                      name={getAvatarName(student.first_name, student.last_name, student.id)}
                      variant={getAvatarVariant(student.id)}
                      className="w-12 h-12"
                    />
                    <div className="flex-1">
                      <p className="font-semibold text-primary">
                        {student.first_name} {student.last_name}
                      </p>
                      {student.grade && (
                        <p className="text-sm text-secondary">Grade {student.grade}</p>
                      )}
                    </div>
                  </label>
                ))}
              </div>

              <div className="border-t border-default pt-4">
                <div className="flex items-center justify-between mb-4">
                  <p className="text-sm text-secondary">
                    {selectedStudentIds.length} student{selectedStudentIds.length !== 1 ? 's' : ''} selected
                  </p>
                </div>
                
                <button
                  onClick={handleStartGame}
                  disabled={selectedStudentIds.length === 0}
                  className={`w-full py-3 rounded-lg font-semibold transition-colors ${
                    selectedStudentIds.length === 0
                      ? 'bg-blue-disabled text-white cursor-not-allowed'
                      : 'bg-blue text-white hover:bg-blue-hover'
                  }`}
                >
                  Start Game
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  )
}