'use client'

import { AlertCircle, Loader2, Plus, Users } from 'lucide-react'
import AddStudentModal from '@/components/students/AddStudentModal'
import StudentCard from '@/components/students/studentCard'
import { useStudents } from '@/hooks/useStudents'

export default function StudentsPage() {
  const { students, isLoading, error, refetch } = useStudents()

  // useEffect(() => {
  //   console.log("Students:", students);
  //   console.log("Resources:", resources);
  //   console.log("Sessions:", sessions);
  //   console.log("Themes:", themes);
  // }, [students, resources, sessions, themes]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <Loader2 className="w-8 h-8 animate-spin text-accent mx-auto mb-4" />
          <p className="text-secondary">Loading students...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center max-w-md">
          <AlertCircle className="w-12 h-12 text-error mx-auto mb-4" />
          <h2 className="text-xl font-semibold text-primary mb-2">
            Error Loading Students
          </h2>
          <p className="text-secondary mb-4">{error}</p>
          <button
            onClick={() => refetch()}
            className="px-4 py-2 bg-accent text-white rounded-lg hover:bg-accent-hover transition-colors"
          >
            Try Again
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background py-8">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <header className="mb-8">
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center space-x-3">
              <Users className="w-8 h-8 text-accent" />
              <h1 className="text-3xl font-bold text-primary">Students</h1>
            </div>
            <AddStudentModal />
          </div>
          <p className="text-secondary">
            View and manage all student information
          </p>
        </header>

        {students.length === 0
          ? (
              <div className="text-center py-12">
                <Users className="w-16 h-16 text-muted mx-auto mb-4 opacity-30" />
                <h2 className="text-xl font-semibold text-primary mb-2">
                  No Students Found
                </h2>
                <p className="text-secondary mb-6">
                  There are no students in the system yet. Get started by adding
                  your first student.
                </p>
                <AddStudentModal
                  trigger={(
                    <button className="px-6 py-3 bg-accent text-white rounded-lg hover:bg-accent-hover transition-colors flex items-center gap-2 mx-auto">
                      <Plus className="w-4 h-4" />
                      Add Your First Student
                    </button>
                  )}
                />
              </div>
            )
          : (
              <div className="space-y-4">
                <div className="flex justify-between items-center mb-4">
                  <p className="text-sm text-secondary">
                    Showing
                    {' '}
                    {students.length}
                    {' '}
                    student
                    {students.length !== 1 ? 's' : ''}
                  </p>
                </div>
                <div className="grid gap-4">
                  {students.map(student => (
                    <StudentCard key={student.id} student={student} />
                  ))}
                </div>
              </div>
            )}
      </div>
    </div>
  )
}
