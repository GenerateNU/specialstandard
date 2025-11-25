'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dropdown } from '@/components/ui/dropdown'
import { ArrowLeft, Loader2, X } from 'lucide-react'
import { useStudents } from '@/hooks/useStudents'
import { useSessions } from '@/hooks/useSessions'
import type { PostSessionsBody } from '@/lib/api/theSpecialStandardAPI.schemas'
import { RepetitionSelector } from '@/components/calendar/repetition-selector'
import type { RepetitionConfig } from '@/components/calendar/repetition-selector'
import { Avatar } from '@/components/ui/avatar'
import { getAvatarName, getAvatarVariant } from '@/lib/avatarUtils'
import CustomAlert from '@/components/ui/CustomAlert'

export default function AddSessionsPage() {
  const router = useRouter()
  const { students, isLoading: loadingStudents } = useStudents()
  const { addSession } = useSessions()
  
  const [therapistId, setTherapistId] = useState<string>('')
  const [repetitionConfig, setRepetitionConfig] = useState<RepetitionConfig | undefined>()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [showError, setShowError] = useState(false)
  const [selectedStudents, setSelectedStudents] = useState<string[]>([])

  const [formData, setFormData] = useState({
    session_name: '',
    sessionDate: new Date().toISOString().split('T')[0],
    startTime: '14:30',
    endTime: '15:30',
    location: '',
    notes: '',
  })

  useEffect(() => {
    const userId = localStorage.getItem('userId')
    if (userId) {
      setTherapistId(userId)
    } else {
      router.push('/signup/welcome')
    }
  }, [router])
  
  // Filter students for this therapist
  const therapistStudents = students.filter(
    student => student.therapist_id === therapistId
  )
  
  const handleBack = () => {
    router.push('/signup/students')
  }

  const handleSkip = () => {
    router.push('/signup/complete')
  }

  const addStudent = (studentId: string) => {
    if (!selectedStudents.includes(studentId)) {
      setSelectedStudents([...selectedStudents, studentId])
    }
  }

  const removeStudent = (studentId: string) => {
    setSelectedStudents(selectedStudents.filter(id => id !== studentId))
  }
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.session_name.trim()) {
      setError('Session name is required')
      setShowError(true)
      return
    }

    if (selectedStudents.length === 0) {
      setError('Please select at least one student')
      setShowError(true)
      return
    }

    setIsSubmitting(true)
    setError(null)

    try {
      const [year, month, day] = formData.sessionDate.split('-').map(Number)
      const sessionDate = new Date(year, month - 1, day)

      const [startHour, startMin] = formData.startTime.split(':').map(Number)
      const [endHour, endMin] = formData.endTime.split(':').map(Number)

      const startDateTime = new Date(sessionDate)
      startDateTime.setHours(startHour, startMin, 0, 0)

      const endDateTime = new Date(sessionDate)
      endDateTime.setHours(endHour, endMin, 0, 0)

      if (endDateTime <= startDateTime) {
        setError('End time must be after start time')
        setShowError(true)
        setIsSubmitting(false)
        return
      }

      const postBody: PostSessionsBody = {
        session_name: formData.session_name,
        start_datetime: startDateTime.toISOString(),
        end_datetime: endDateTime.toISOString(),
        therapist_id: therapistId!,
        notes: formData.notes || undefined,
        location: formData.location || undefined,
        student_ids: selectedStudents,
        repetition: repetitionConfig,
      }

      await addSession(postBody)
      
      // Store session info in localStorage for summary
      const sessionsData = localStorage.getItem('onboardingSessions')
      const existingSessions = sessionsData ? JSON.parse(sessionsData) : []
      existingSessions.push({
        name: formData.session_name,
        date: postBody.start_datetime,
        students: selectedStudents,
      })
      localStorage.setItem('onboardingSessions', JSON.stringify(existingSessions))
      
      router.push('/signup/sessions')
    } catch (error) {
      console.error('Failed to create session:', error)
      setError(error instanceof Error ? error.message : 'Failed to create session')
      setShowError(true)
    } finally {
      setIsSubmitting(false)
    }
  }

  const availableStudents = therapistStudents.filter(s => !selectedStudents.includes(s.id))

  if (loadingStudents || !therapistId) {
    return (
      <div className="flex items-center justify-center min-h-screen p-8">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="flex items-center justify-center min-h-screen p-8">
      <div className="max-w-md w-full">
        <button
          onClick={handleBack}
          className="mb-6 flex items-center text-secondary hover:text-primary cursor-pointer transition-colors"
        >
          <ArrowLeft className="w-4 h-4 mr-1" />
          Back
        </button>

        <h1 className="text-3xl font-bold text-primary mb-8">
          Add Sessions
        </h1>

        {showError && error && (
          <div className="mb-4">
            <CustomAlert
              variant="destructive"
              title="Error"
              description={error}
              onClose={() => {
                setShowError(false)
                setError(null)
              }}
            />
          </div>
        )}
        
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Session Name */}
          <div>
            <p className="text-xs text-secondary mb-1">Session Name</p>
            <Input
              value={formData.session_name}
              onChange={(e) => setFormData({ ...formData, session_name: e.target.value })}
              placeholder="e.g. Language Group"
              className="border border-gray-300"
              required
              disabled={isSubmitting}
            />
          </div>

          {/* Date and Time */}
          <div>
            <p className="text-xs text-secondary mb-1">Date & Time</p>
            <div className="flex gap-2 items-center">
              <Input
                type="date"
                value={formData.sessionDate}
                onChange={(e) => setFormData({ ...formData, sessionDate: e.target.value })}
                className="border border-gray-300 flex-1"
                disabled={isSubmitting}
              />
              <Input
                type="time"
                value={formData.startTime}
                onChange={(e) => setFormData({ ...formData, startTime: e.target.value })}
                className="border border-gray-300 w-24  [&::-webkit-calendar-picker-indicator]:hidden"
                disabled={isSubmitting}
              />
              <span className="text-secondary">–</span>
              <Input
                type="time"
                value={formData.endTime}
                onChange={(e) => setFormData({ ...formData, endTime: e.target.value })}
                className="border border-gray-300 w-24 [&::-webkit-calendar-picker-indicator]:hidden"
                disabled={isSubmitting}
                
              />
            </div>
          </div>

          {/* Repetition */}
          <div>
            <p className="text-xs text-secondary mb-1">Repeat Session (optional)</p>
            <RepetitionSelector
              value={repetitionConfig}
              onChange={setRepetitionConfig}
              sessionDate={formData.sessionDate}
              sessionTime={formData.startTime}
            />
          </div>

          {/* Location */}
          <div>
            <p className="text-xs text-secondary mb-1">Location (optional)</p>
            <Input
              value={formData.location}
              onChange={(e) => setFormData({ ...formData, location: e.target.value })}
              placeholder="e.g. Room 204"
              className="border border-gray-300"
              disabled={isSubmitting}
            />
          </div>

          {/* Students */}
          <div>
            <label className="block text-sm font-medium text-primary mb-2">
              Students
            </label>
            
            {/* Selected Students */}
            <div className="space-y-2 mb-3">
              {selectedStudents.map(id => {
                const student = therapistStudents.find(s => s.id === id)
                if (!student) return null
                return (
                  <div key={id} className="bg-gray-50 rounded p-3 flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <Avatar
                        name={getAvatarName(student.first_name, student.last_name, student.id)}
                        variant={getAvatarVariant(student.id)}
                        className="w-10 h-10"
                      />
                      <div>
                        <div className="font-medium">{student.first_name} {student.last_name}</div>
                        <div className="text-xs text-secondary">
                          Grade {student.grade}
                          {student.school_name && ` • ${student.school_name}`}
                        </div>
                      </div>
                    </div>
                    <button
                      type="button"
                      onClick={() => removeStudent(id)}
                      className="text-secondary hover:text-red-500 transition-colors"
                      disabled={isSubmitting}
                    >
                      <X className="w-4 h-4" />
                    </button>
                  </div>
                )
              })}
            </div>

            {/* Add Student Dropdown - Always visible if there are available students */}
            {availableStudents.length > 0 ? (
              <Dropdown
                items={availableStudents.map(student => ({
                  label: `${student.first_name} ${student.last_name}`,
                  value: student.id,
                  onClick: () => addStudent(student.id),
                }))}
                placeholder="+ Add Student"
                className="w-full border border-gray-300 border-dashed hover:border-gray-400 transition-colors"
              />
            ) : therapistStudents.length === 0 ? (
              <p className="text-xs text-secondary text-center py-3 border border-gray-300 border-dashed rounded">
                No students available. Add students first.
              </p>
            ) : (
              <p className="text-xs text-secondary text-center py-3 border border-gray-300 border-dashed rounded">
                All students have been added to this session
              </p>
            )}
            
            {/* Count indicator */}
            {selectedStudents.length > 0 && (
              <p className="text-xs text-secondary mt-2">
                {selectedStudents.length} student{selectedStudents.length !== 1 ? 's' : ''} selected
              </p>
            )}
          </div>

          {/* Notes */}
          <div>
            <p className="text-xs text-secondary mb-1">Notes (optional)</p>
            <Input
              value={formData.notes}
              onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
              placeholder="Add session notes"
              className="border border-gray-300"
              disabled={isSubmitting}
            />
          </div>

          <div className="pt-4">
            <Button
              type="submit"
              size="long"
              className="w-full text-white"
              disabled={isSubmitting}
            >
              {isSubmitting ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin mr-2" />
                  <span>Creating session...</span>
                </>
              ) : (
                <span>Save</span>
              )}
            </Button>
          </div>
          
          <div className="text-center">
            <button
              type="button"
              onClick={handleSkip}
              className="text-sm cursor-pointer text-secondary hover:text-primary underline"
            >
              Skip this step
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}