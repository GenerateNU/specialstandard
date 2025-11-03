'use client'

import { ArrowLeft, Calendar, Clock, MapPin, Pencil, Plus, X } from 'lucide-react'
import Link from 'next/link'
import { use, useState } from 'react'
import { Avatar } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { useSession, useSessions } from '@/hooks/useSessions'
import {
  useSessionStudents,
  useSessionStudentsForSession,
} from '@/hooks/useSessionStudents'
import { useStudents } from '@/hooks/useStudents'
import { getAvatarName, getAvatarVariant, getStudentInitials } from '@/lib/avatarUtils'

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
  const { updateSession } = useSessions()
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

  const [mode, setMode] = useState<'view' | 'attendance' | 'editStudents'>('view')
  const [studentToRemove, setStudentToRemove] = useState<{ id: string, name: string } | null>(null)
  const [isEditingSession, setIsEditingSession] = useState(false)
  const [editedSession, setEditedSession] = useState({
    startDate: '',
    startTime: '',
    endDate: '',
    endTime: '',
    notes: '',
  })

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

  const handleRemoveStudent = (studentId: string, studentName: string) => {
    setStudentToRemove({ id: studentId, name: studentName })
  }

  const confirmRemoveStudent = () => {
    if (studentToRemove) {
      removeStudentFromSession({
        session_id: id,
        student_id: studentToRemove.id,
      })
    }
  }

  const handleAddStudent = (studentId: string) => {
    addStudentToSession({
      session_ids: [id],
      student_ids: [studentId],
      present: false,
    })
  }

  const handleToggleAttendance = (studentId: string, present: boolean) => {
    updateSessionStudent({
      session_id: id,
      student_id: studentId,
      present,
    })
  }

  const handleEditClick = () => {
    const start = new Date(session.start_datetime)
    const end = new Date(session.end_datetime)

    setEditedSession({
      startDate: start.toISOString().split('T')[0],
      startTime: start.toTimeString().slice(0, 5),
      endDate: end.toISOString().split('T')[0],
      endTime: end.toTimeString().slice(0, 5),
      notes: session.notes || '',
    })
    setIsEditingSession(true)
  }

  const handleSaveSession = () => {
    const startDatetime = new Date(`${editedSession.startDate}T${editedSession.startTime}`).toISOString()
    const endDatetime = new Date(`${editedSession.endDate}T${editedSession.endTime}`).toISOString()

    updateSession(id, {
      start_datetime: startDatetime,
      end_datetime: endDatetime,
      notes: editedSession.notes,
    })
    setIsEditingSession(false)
  }

  const handleCancelEdit = () => {
    setIsEditingSession(false)
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
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-4xl font-bold">Session Details</h1>
          {!isEditingSession
            ? (
                <Button
                  onClick={handleEditClick}
                  variant="outline"
                  className="flex items-center gap-2"
                >
                  <Pencil className="w-4 h-4" />
                  Edit
                </Button>
              )
            : (
                <div className="flex gap-2">
                  <Button
                    onClick={handleCancelEdit}
                    variant="outline"
                  >
                    Cancel
                  </Button>
                  <Button
                    onClick={handleSaveSession}
                    variant="default"
                  >
                    Save
                  </Button>
                </div>
              )}
        </div>

        {!isEditingSession
          ? (
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
            )
          : (
              <div className="space-y-4 mb-6">
                {/* Date & Time Inputs */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      <Calendar className="w-4 h-4 inline mr-2" />
                      Start Date & Time
                    </label>
                    <div className="flex gap-2">
                      <input
                        type="date"
                        value={editedSession.startDate}
                        onChange={e => setEditedSession({ ...editedSession, startDate: e.target.value })}
                        className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <input
                        type="time"
                        value={editedSession.startTime}
                        onChange={e => setEditedSession({ ...editedSession, startTime: e.target.value })}
                        className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      <Clock className="w-4 h-4 inline mr-2" />
                      End Date & Time
                    </label>
                    <div className="flex gap-2">
                      <input
                        type="date"
                        value={editedSession.endDate}
                        onChange={e => setEditedSession({ ...editedSession, endDate: e.target.value })}
                        className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <input
                        type="time"
                        value={editedSession.endTime}
                        onChange={e => setEditedSession({ ...editedSession, endTime: e.target.value })}
                        className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                    </div>
                  </div>
                </div>

                {/* Location Input */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    <MapPin className="w-4 h-4 inline mr-2" />
                    Location
                  </label>
                  <input
                    type="text"
                    value={editedSession.notes}
                    onChange={e => setEditedSession({ ...editedSession, notes: e.target.value })}
                    placeholder="e.g., Boston Latin Academy"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>
            )}
      </div>

      {/* Students section */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-3xl font-semibold">Students</h2>
          <div className="flex gap-3">
            {mode === 'view' && (
              <>
                <Button
                  onClick={() => setMode('attendance')}
                  variant="default"
                  size="lg"
                >
                  Attendance
                </Button>
                <Button
                  onClick={() => setMode('editStudents')}
                  variant="outline"
                  size="lg"
                >
                  Edit Students
                </Button>
              </>
            )}
            {mode !== 'view' && (
              <Button
                onClick={() => setMode('view')}
                variant="secondary"
                size="lg"
              >
                Done
              </Button>
            )}
          </div>
        </div>

        {/* Current students list */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
          {sessionStudents.map((student, index) => (
            <div
              key={student.id || `student-${index}`}
              className="bg-card rounded-2xl p-6 shadow-sm border border-default flex items-center justify-between"
            >
              <div className="flex items-center gap-4">
                <Avatar
                  name={getAvatarName(
                    student.first_name || 'Unknown',
                    student.last_name || 'Student',
                    student.id,
                  )}
                  variant={getAvatarVariant(student.id)}
                  className="w-16 h-16 ring-2 ring-accent-light"
                />
                <div>
                  <p className="text-lg font-medium">
                    {getStudentInitials(student.first_name, student.last_name)}
                  </p>
                  <p className="text-sm text-secondary">
                    {student.first_name || 'Unknown'}
                    {' '}
                    {student.last_name || 'Student'}
                  </p>
                </div>
              </div>

              {/* Attendance mode: Present/Absent buttons */}
              {mode === 'attendance' && (
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

              {/* Edit mode: Remove button */}
              {mode === 'editStudents' && (
                <Button
                  onClick={() => handleRemoveStudent(
                    student.id,
                    `${student.first_name || 'Unknown'} ${student.last_name || 'Student'}`,
                  )}
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

        {/* Add students section (only in edit mode) */}
        {mode === 'editStudents' && availableStudents.length > 0 && (
          <div>
            <h3 className="text-xl font-semibold mb-4">Add Students</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {availableStudents.map((student, index) => (
                <div
                  key={student.id || `available-student-${index}`}
                  className="bg-card rounded-2xl p-6 shadow-sm border border-default flex items-center justify-between"
                >
                  <div className="flex items-center gap-4">
                    <Avatar
                      name={getAvatarName(
                        student.first_name || 'Unknown',
                        student.last_name || 'Student',
                        student.id,
                      )}
                      variant={getAvatarVariant(student.id)}
                      className="w-16 h-16 ring-2 ring-accent-light"
                    />
                    <div>
                      <p className="text-lg font-medium">
                        {getStudentInitials(student.first_name, student.last_name)}
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
                    <Plus className="w-4 h-4 mr-1" />
                    Add
                  </Button>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Confirmation dialog for removing student */}
      <ConfirmDialog
        isOpen={!!studentToRemove}
        onClose={() => setStudentToRemove(null)}
        onConfirm={confirmRemoveStudent}
        title="Remove Student"
        description={`Are you sure you want to remove ${studentToRemove?.name} from this session?`}
        confirmText="Remove"
        cancelText="Cancel"
        variant="danger"
        isLoading={isRemoving}
      />
    </div>
  )
}
