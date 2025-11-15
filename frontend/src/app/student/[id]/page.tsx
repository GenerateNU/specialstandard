'use client'

import { ChevronLeft, CirclePlus, PencilLine, Save, Trash2, X } from 'lucide-react'
import { useParams } from 'next/navigation'
import { useEffect, useState } from 'react'
import AppLayout from '@/components/AppLayout'
import UpcomingSession from '@/components/sessions/UpcomingSession'

import SessionNotes from '@/components/sessions/sessionNotes'
import { Avatar } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { useRecentlyViewedStudents } from '@/hooks/useRecentlyViewedStudents'
import { useStudents } from '@/hooks/useStudents'
import { getAvatarVariant } from '@/lib/utils'
import SchoolTag from '@/components/school/schoolTag'

import { useStudentAttendance } from '@/hooks/useStudentAttendance'
import CustomPieChart from '@/components/statistics/PieChart'

function StudentPage() {
  const params = useParams()
  const studentId = params.id as string

  const { students, isLoading, updateStudent } = useStudents()
  const student = students.find(s => s.id === studentId)
  const { addRecentStudent } = useRecentlyViewedStudents()

  const { attendance, isLoading: attendanceLoading } = useStudentAttendance({
  studentId: studentId || '',
})

  // Track this student as recently viewed - only when studentId changes
  useEffect(() => {
    if (student && studentId) {
      addRecentStudent(student)
    }
  }, [studentId]) // Only track when navigating to a new student

  const [edit, setEdit] = useState(false)
  const [iepGoals, setIepGoals] = useState<string[]>([])
  const [isSaving, setIsSaving] = useState(false)

  // Initialize IEP goals from student data
  useEffect(() => {
    if (student?.iep && Array.isArray(student.iep)) {
      setIepGoals(student.iep)
    }
    else {
      setIepGoals([])
    }
  }, [student?.iep])

  const CORNER_ROUND = 'rounded-4xl'
  const PADDING = 'p-5'

  const fullName = student ? `${student.first_name} ${student.last_name}` : ''
  const initials = student ? `${student.first_name[0]}.${student.last_name[0]}.` : ''
  const avatarVariant = student ? getAvatarVariant(student.id) : 'lorelei'

  const handleSave = async () => {
    if (!student)
      return

    // Filter out empty goals before saving
    const filteredGoals = iepGoals.filter(goal => goal.trim() !== '')

    setIsSaving(true)
    try {
      // Save IEP goals as array
      await updateStudent(student.id, { iep: filteredGoals })
      setIepGoals(filteredGoals)
      setEdit(false)
    }
    catch (error) {
      console.error('Failed to save IEP goals:', error)
    }
    finally {
      setIsSaving(false)
    }
  }

  const handleCancel = () => {
    // Reset to original data
    if (student?.iep && Array.isArray(student.iep)) {
      setIepGoals(student.iep)
    }
    else {
      setIepGoals([])
    }
    setEdit(false)
  }

  const addGoal = () => {
    setIepGoals([...iepGoals, ''])
  }

  const updateGoal = (index: number, value: string) => {
    const newGoals = [...iepGoals]
    newGoals[index] = value
    setIepGoals(newGoals)
  }

  const deleteGoal = (index: number) => {
    setIepGoals(iepGoals.filter((_, i) => i !== index))
  }

  if (isLoading) {
    return (
      <div className="min-h-screen h-screen flex items-center justify-center bg-background">
        <div className="text-primary">Loading student...</div>
      </div>
    )
  }

  if (!student) {
    return (
      <div className="min-h-screen h-screen flex items-center justify-center bg-background">
        <div className="text-error">Student not found</div>
      </div>
    )
  }

  return (
    <AppLayout>
      <div className="w-full h-screen bg-background">
        <div className={`w-full h-full flex flex-col gap-8 ${PADDING} relative overflow-y-auto`}>
          {/* Edit toggle button */}
          <div className="absolute top-1/2 right-5 z-20 flex gap-2">
            {edit
              ? (
                  <>
                    <Button
                      onClick={handleSave}
                      disabled={isSaving}
                      className="w-12 h-12 p-0 bg-green-600 hover:bg-green-700"
                      size="icon"
                    >
                      <Save size={20} />
                    </Button>
                    <Button
                      onClick={handleCancel}
                      disabled={isSaving}
                      className="w-12 h-12 p-0 bg-red-600 hover:bg-red-700"
                      size="icon"
                    >
                      <X size={20} />
                    </Button>
                  </>
                )
              : (
                  <Button
                    onClick={() => setEdit(!edit)}
                    className="w-12 h-12 p-0"
                    variant="secondary"
                    size="icon"
                  >
                    <PencilLine size={20} />
                  </Button>
                )}
          </div>

          <div className="flex flex-col gap-4 flex-shrink-0">
            {/* Back button */}
            <Button
              variant="outline"
              className={`w-fit p-4 flex flex-row items-center gap-2 ${CORNER_ROUND} flex-shrink-0`}
              onClick={() => window.history.back()}
            >
              <ChevronLeft />
              Back
            </Button>

            {/* Profile and Upcoming Sessions row */}
            <div className="flex gap-8 flex-1 min-h-0">
              {/* Student Profile */}
              <div className={`flex-1 bg-card border-2 border-default ${CORNER_ROUND} overflow-hidden flex flex-col relative`}>
                {/* Edit Profile Button - Separate Section */}
                <div className="flex justify-end p-3 flex-shrink-0 relative z-10">
                  <Button
                    onClick={() => {/* Navigate to edit page */}}
                    variant="secondary"
                    className="flex items-center gap-2 hover:bg-accent h-8"
                    size="sm"
                  >
                    <span className="text-base font-medium">Edit Profile</span>
                    <PencilLine size={18} />
                  </Button>
                </div>
                
                {/* Content - positioned absolutely to center in entire card */}
                <div className="absolute inset-0 flex items-center justify-center gap-8 p-5 pointer-events-none">
                  <div className="max-h-full aspect-square border-2 border-default rounded-full flex-shrink-0 pointer-events-auto">
                    <Avatar
                      name={fullName + student.id}
                      variant={avatarVariant}
                      className="w-full h-full"
                    />
                  </div>
                  
                  <div className="flex flex-col gap-3 flex-1 pointer-events-auto">
                    <div className="text-4xl font-bold text-primary">{initials}</div>
                    
                    <div className="flex items-center gap-3 flex-wrap">
                      <span className="text-xl font-medium text-primary">
                        Grade {student.grade}
                      </span>
                      {student.school_name && <SchoolTag schoolName={student.school_name} />}
                    </div>
                  </div>
                </div>
              </div>

              {/* Upcoming Sessions */}
              <div className={`flex-1 bg-card border-2 border-default ${CORNER_ROUND} ${PADDING} flex flex-col gap-4`}>
                <div className="text-2xl font-semibold text-primary flex-shrink-0">Upcoming Sessions</div>
                <div className="flex-1 min-h-0">
                  {/* Upcoming sessions content will go here */}
                  <UpcomingSession studentId={studentId}/>
                </div>
              </div>
            </div>
          </div>
          {/* Attendance */}
          <div className={`bg-card border-2 border-default ${CORNER_ROUND} ${PADDING} flex-shrink-0 ${attendance?.total_count === 0 ? 'opacity-50 pointer-events-none' : ''}`}>
            {attendanceLoading ? (
              <div className="flex items-center justify-center h-full">
                <div className="text-sm text-muted-foreground">Loading attendance...</div>
              </div>
            ) : attendance && attendance.total_count > 0 ? (
              <CustomPieChart
                percentage={Math.round((attendance.present_count / attendance.total_count) * 100)}
                title="Attendance"
              />
            ) : (
              <div className="flex items-center justify-center h-full">
                <div className="text-sm text-muted-foreground">No attendance data</div>
              </div>
            )}
          </div>
          {/* Goals and Session Notes */}
          <div className="grid grid-cols-2 gap-8 overflow-hidden">
            <div className="gap-2 flex flex-col overflow-hidden">
              <div className="w-full text-2xl text-primary flex items-baseline font-semibold">
                IEP Goals
              </div>
              <div className="flex-1 overflow-y-auto flex flex-col gap-2">
                {iepGoals.length === 0 && !edit
                  ? (
                      <div className="text-muted-foreground italic">No IEP goals set</div>
                    )
                  : (
                      iepGoals.map((goal, index) => (
                        <div
                          key={index}
                          className={`w-full text-lg flex items-center gap-2
                    rounded-2xl transition bg-background select-none border-2 border-border ${PADDING} ${!edit && 'hover:scale-99'}`}
                        >
                          {edit
                            ? (
                                <>
                                  <input
                                    value={goal}
                                    onChange={e => updateGoal(index, e.target.value)}
                                    onBlur={() => goal.trim() === '' && deleteGoal(index)}
                                    className="flex-1 bg-transparent outline-none py-1 leading-normal"
                                    placeholder="Enter IEP goal..."
                                  />
                                  <Button
                                    onClick={() => deleteGoal(index)}
                                    variant="ghost"
                                    size="icon"
                                    className="text-red-600 hover:text-red-700 hover:bg-red-100 flex-shrink-0"
                                  >
                                    <Trash2 size={18} />
                                  </Button>
                                </>
                              )
                            : (
                                <span className="py-1 leading-normal">{goal}</span>
                              )}
                        </div>
                      ))
                    )}
              </div>
              {edit && (
                <Button
                  onClick={addGoal}
                  variant="outline"
                  className="w-full border-2 rounded-full gap-2"
                >
                  <CirclePlus size={20} />
                  Add IEP Goal
                </Button>
              )}
            </div>

            <div className="gap-2 flex flex-col overflow-hidden">
              <div className="gap-2 flex flex-col overflow-hidden">
                <div className="w-full text-2xl text-primary flex items-baseline font-semibold">
                  Session Notes
                </div>
                <div className="flex-1 overflow-y-auto">
                  <SessionNotes studentId={studentId} />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AppLayout>
  )
}

export default StudentPage
