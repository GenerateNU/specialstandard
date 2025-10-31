'use client'

import { ArrowLeft, Calendar, Clock, MapPin, Plus, X } from 'lucide-react'
import Link from 'next/link'
import { use, useState } from 'react'
import { Avatar } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { useSession } from '@/hooks/useSessions'
import {
  useSessionStudents,
  useSessionStudentsForSession,
} from '@/hooks/useSessionStudents'
import { useStudents } from '@/hooks/useStudents'

interface PageProps {
  params: Promise<
    {
      id: string
    }
  >
}

export default function SessionPage({ params }: PageProps) {
  const { id } = use(params)
  const { session, isLoading: sessionLoading } = useSession(id)
  const { students: sessionStudents, isLoading: studentsLoading }
    = useSessionStudentsForSession(id)
  const { students: allStudents } = useStudents()
  const {
    addStudentToSession,
    removeStudentFromSession,
    updateSessionStudent,
    isAdding,
    isRemoving,
  } = useSessionStudents()

  const [isAttendanceMode, setIsAttendanceMode] = useState(false)
  const [showAddStudents, setShowAddStudents] = useState(false)

  if (sessionLoading || studentsLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div>Loading session...</div>
      </div>
    )
  }

  if (!session) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div>Session not found</div>
      </div>
    )
  }

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric',
    })
  }

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true,
    })
  }

  const formatTimeRange = () => {
    return `${formatTime(session.start_datetime)} - ${formatTime(session.end_datetime)}`
  }

  const getInitials = (firstName?: string, lastName?: string) => {
    const first = firstName?.charAt(0) || '?'
    const last = lastName?.charAt(0) || '?'
    return `${first}${last}`.toUpperCase()
  }

  // Function to deterministically select avatar variant based on student ID
  function getAvatarVariant(id?: string): 'avataaars' | 'lorelei' | 'micah' | 'miniavs' | 'big-smile' | 'personas' {
    const variants = ['avataaars', 'lorelei', 'micah', 'miniavs', 'big-smile', 'personas'] as const

    // Default to first variant if no ID
    if (!id) {
      return variants[0]
    }

    // Simple hash function to get consistent index
    let hash = 0
    for (let i = 0; i < id.length; i++) {
      hash = ((hash << 5) - hash) + id.charCodeAt(i)
      hash = hash & hash // Convert to 32-bit integer
    }

    return variants[Math.abs(hash) % variants.length]
  }

  const handleRemoveStudent = (studentId: string) => {
    removeStudentFromSession({
      session_id: id,
      student_id: studentId,
    })
  }

  const handleAddStudent = (studentId: string) => {
    addStudentToSession({
      session_id: id,
      student_id: studentId,
      present: false,
    })
    setShowAddStudents(false)
  }

  const handleToggleAttendance = (studentId: string, present: boolean) => {
    updateSessionStudent({
      session_id: id,
      student_id: studentId,
      present,
    })
  }

  // Filter out students already in session
  const availableStudents = allStudents.filter(
    student => !sessionStudents.some(s => s.id === student.id),
  )

  return (
    <div className="min-h-screen bg-background p-8">
      {/* Back button */}
      <Link
        href="/calendar"
        className="inline-flex items-center gap-2 text-secondary hover:text-primary mb-6 transition-colors group"
      >
        <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
        <span className="text-sm font-medium">Back to Calendar</span>
      </Link>

      {/* Session header */}
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-6">Session Details</h1>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          {/* Date */}
          <div className="bg-card-hover rounded-2xl px-6 py-4 flex items-center gap-3">
            <Calendar className="w-5 h-5 text-accent" />
            <span className="text-lg">{formatDateTime(session.start_datetime)}</span>
          </div>

          {/* Time */}
          <div className="bg-card-hover rounded-2xl px-6 py-4 flex items-center gap-3">
            <Clock className="w-5 h-5 text-accent" />
            <span className="text-lg">{formatTimeRange()}</span>
          </div>

          {/* Location */}
          <div className="bg-card-hover rounded-2xl px-6 py-4 flex items-center gap-3">
            <MapPin className="w-5 h-5 text-accent" />
            <span className="text-lg">{session.notes || 'No location'}</span>
          </div>
        </div>
      </div>

      {/* Students section */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-3xl font-semibold">Students</h2>
          <div className="flex gap-3">
            {!isAttendanceMode && (
              <Button
                onClick={() => setIsAttendanceMode(true)}
                variant="default"
                size="lg"
              >
                Attendance
              </Button>
            )}
            {isAttendanceMode && (
              <>
                <Button
                  onClick={() => setShowAddStudents(!showAddStudents)}
                  variant="default"
                  size="lg"
                >
                  {showAddStudents ? 'Close' : 'Swap'}
                </Button>
                <Button
                  onClick={() => setIsAttendanceMode(false)}
                  variant="secondary"
                  size="lg"
                >
                  Done
                </Button>
              </>
            )}
          </div>
        </div>

        {/* Student list */}
        {!showAddStudents && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {sessionStudents.map((student, index) => (
              <div
                key={student.id || `student-${index}`}
                className="bg-card rounded-2xl p-6 shadow-sm border border-default flex items-center justify-between"
              >
                <div className="flex items-center gap-4">
                  <Avatar
                    name={`${student.first_name || 'Unknown'} ${student.last_name || 'Student'}${student.id}`}
                    variant={getAvatarVariant(student.id)}
                    className="w-16 h-16 ring-2 ring-accent-light"
                  />
                  <div>
                    <p className="text-lg font-medium">
                      {getInitials(student.first_name, student.last_name)}
                    </p>
                    <p className="text-sm text-secondary">
                      {student.first_name || 'Unknown'}
                      {' '}
                      {student.last_name || 'Student'}
                    </p>
                  </div>
                </div>

                {isAttendanceMode && (
                  <div className="flex gap-2">
                    <Button
                      onClick={() => handleToggleAttendance(student.id, true)}
                      variant={student.present ? 'default' : 'outline'}
                      size="sm"
                    >
                      Present
                    </Button>
                    <Button
                      onClick={() => handleToggleAttendance(student.id, false)}
                      variant={!student.present ? 'default' : 'outline'}
                      size="sm"
                    >
                      Absent
                    </Button>
                  </div>
                )}

                {!isAttendanceMode && (
                  <Button
                    onClick={() => handleRemoveStudent(student.id)}
                    disabled={isRemoving}
                    variant="ghost"
                    size="icon"
                    aria-label="Remove student"
                  >
                    <X className="w-5 h-5" />
                  </Button>
                )}
              </div>
            ))}
          </div>
        )}

        {/* Add students mode */}
        {showAddStudents && (
          <div>
            <h3 className="text-xl font-semibold mb-4">All Students</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {availableStudents.map((student, index) => (
                <div
                  key={student.id || `available-student-${index}`}
                  className="bg-card rounded-2xl p-6 shadow-sm border border-default flex items-center justify-between"
                >
                  <div className="flex items-center gap-4">
                    <Avatar
                      name={`${student.first_name || 'Unknown'} ${student.last_name || 'Student'}${student.id}`}
                      variant={getAvatarVariant(student.id)}
                      className="w-16 h-16 ring-2 ring-accent-light"
                    />
                    <div>
                      <p className="text-lg font-medium">
                        {getInitials(student.first_name, student.last_name)}
                      </p>
                      <p className="text-sm text-secondary">
                        {student.first_name || 'Unknown'}
                        {' '}
                        {student.last_name || 'Student'}
                      </p>
                    </div>
                  </div>

                  <Button
                    onClick={() => handleAddStudent(student.id)}
                    disabled={isAdding}
                    variant="default"
                    size="sm"
                  >
                    Add
                  </Button>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Add student button (non-attendance mode) */}
        {!isAttendanceMode && (
          <button
            onClick={() => {
              setIsAttendanceMode(true)
              setShowAddStudents(true)
            }}
            className="mt-6 flex items-center gap-2 text-accent hover:text-accent-dark transition-colors"
          >
            <div className="w-12 h-12 rounded-full bg-card-hover flex items-center justify-center">
              <Plus className="w-6 h-6" />
            </div>
            <span className="text-lg font-medium">Add student</span>
          </button>
        )}
      </div>
    </div>
  )
}
