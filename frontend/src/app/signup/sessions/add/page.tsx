'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { ArrowLeft } from 'lucide-react'
import { CreateSessionDialog } from '@/components/calendar/NewSessionModal'
import { useStudents } from '@/hooks/useStudents'
import { useSessions } from '@/hooks/useSessions'
import type { PostSessionsBody } from '@/lib/api/theSpecialStandardAPI.schemas'

export default function AddSessionsPage() {
  const router = useRouter()
  const { students, isLoading: loadingStudents } = useStudents()
  const { addSession } = useSessions()
  
  const [isModalOpen, setIsModalOpen] = useState(true)
  const [therapistId, setTherapistId] = useState<string>('')
  
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
  
  
  const handleSessionSubmit = async (data: PostSessionsBody) => {
    try {
      await addSession(data)
      
      // Store session info in localStorage for summary
      const sessionsData = localStorage.getItem('onboardingSessions')
      const existingSessions = sessionsData ? JSON.parse(sessionsData) : []
      existingSessions.push({
        name: data.session_name,
        date: data.start_datetime,
        students: data.student_ids,
      })
      localStorage.setItem('onboardingSessions', JSON.stringify(existingSessions))
      
      setIsModalOpen(false)
      router.push('/signup/sessions')
    } catch (error) {
      console.error('Failed to create session:', error)
      // The dialog component should handle its own error display
    }
  }
  
  const handleModalClose = (open: boolean) => {
    setIsModalOpen(open)
    if (!open) {
      router.push('/signup/sessions')
    }
  }

  if (loadingStudents || !therapistId) {
    return (
      <div className="flex items-center justify-center min-h-screen p-8">
        <div className="animate-pulse">Loading...</div>
      </div>
    )
  }

  return (
    <div className="flex items-center justify-center min-h-screen p-8">
      <div className="max-w-2xl w-full">
        <button
          onClick={handleBack}
          className="mb-6 flex items-center text-secondary hover:text-primary transition-colors"
        >
          <ArrowLeft className="w-4 h-4 mr-1" />
          Back
        </button>

        <h1 className="text-3xl font-bold text-primary mb-2">
          Add Sessions
        </h1>
        <p className="text-secondary mb-8">
          Schedule your first therapy session
        </p>

        {therapistStudents.length === 0 ? (
          <div className="bg-card rounded-lg border border-default p-8 text-center">
            <p className="text-secondary mb-4">
              You need to add students before scheduling sessions.
            </p>
            <Button
              onClick={() => router.push('/signup/students-add')}
              className="hover:bg-accent-hover text-white"
            >
              Add Students First
            </Button>
          </div>
        ) : (
          <>
            <div className="bg-accent-light rounded-lg p-4 mb-6">
              <p className="text-sm text-primary">
                💡 <strong>Tip:</strong> You can schedule recurring sessions and add more sessions later from your dashboard.
              </p>
            </div>

            <CreateSessionDialog
              open={isModalOpen}
              setOpen={handleModalClose}
              therapistId={therapistId}
              students={therapistStudents}
              onSubmit={handleSessionSubmit}
            />

            {!isModalOpen && (
              <div className="space-y-4">
                <Button
                  onClick={() => setIsModalOpen(true)}
                  className="w-full bg-accent hover:bg-accent-hover text-white"
                >
                  Add Another Session
                </Button>
                
                <Button
                  onClick={() => router.push('/signup/complete')}
                  variant="outline"
                  className="w-full"
                >
                  Continue to Finish
                </Button>
                
                <div className="text-center">
                  <button
                    onClick={() => router.push('/signup/complete')}
                    className="text-sm text-secondary hover:text-primary underline"
                  >
                    Skip this step
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}