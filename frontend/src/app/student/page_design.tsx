// DESIGN SYSTEM

'use client'

import { CirclePlus, PencilLine, Save, Trash2, X } from 'lucide-react'
import { useState } from 'react'
import StudentSchedule from '@/components/schedule/StudentSchedule'
import { Avatar } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'

// Mock RecentSession component for preview
function MockRecentSession() {
  const mockSessions = [
    {
      id: '1',
      notes: 'Articulation Practice',
      start_datetime: '2024-11-01T10:00:00Z',
      end_datetime: '2024-11-01T10:45:00Z',
      present: true,
    },
    {
      id: '2',
      notes: 'Vocabulary Building',
      start_datetime: '2024-10-29T14:00:00Z',
      end_datetime: '2024-10-29T14:45:00Z',
      present: true,
    },
    {
      id: '3',
      notes: 'Following Directions',
      start_datetime: '2024-10-25T10:00:00Z',
      end_datetime: '2024-10-25T10:45:00Z',
      present: false,
    },
  ]

  return (
    <div className="h-full overflow-y-auto space-y-2 text-background w-full">
      {mockSessions.map(session => (
        <div
          key={session.id}
          className={`p-4 border-b border-background/20 rounded-2xl flex flex-col justify-center 
        h-20 bg-background w-full text-primary last:border-b-0`}
        >
          <div className="w-full flex justify-between items-center">
            <div className="font-semibold text-base">
              {session.notes}
            </div>
            <div className="text-sm opacity-90">
              {new Date(session.start_datetime).toLocaleDateString()}
            </div>
          </div>
          <div className="text-sm mb-1">
            {session.present
              ? <span className="font-medium">Present ✓</span>
              : <span className="font-medium">Absent</span>}
          </div>
          <div className="text-sm opacity-75">
            {new Date(session.start_datetime).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' })}
            {' - '}
            {new Date(session.end_datetime).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' })}
          </div>
        </div>
      ))}
    </div>
  )
}

function StudentPage() {
  const [edit, setEdit] = useState(false)
  const [iepGoals, setIepGoals] = useState<string[]>([
    'Improve articulation of /r/ sound in conversation',
    'Increase expressive vocabulary by 20 words',
    'Follow 3-step directions with 80% accuracy',
  ])
  // const [additionalGoals] = useState<string[]>([])

  const CORNER_ROUND = 'rounded-4xl'
  const PADDING = 'p-5'

  // Mock sessions for Session Notes (will be replaced by real data)
  const mockSessionNotes = [
    {
      id: '1',
      title: 'Articulation Practice',
      notes: 'Worked on /r/ in 2-word phrases and during short conversation practice. This is a lot of text to test wrapping functionality in the session notes section. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.This is a lot of text to test wrapping functionality in the session notes section. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.',
      start_datetime: '2024-11-01T10:00:00Z',
      end_datetime: '2024-11-01T10:45:00Z',
    },
    {
      id: '2',
      title: 'Vocabulary Building',
      notes: 'Focused on thematic vocabulary (grocery items); introduced 8 new words.',
      start_datetime: '2024-10-29T14:00:00Z',
      end_datetime: '2024-10-29T14:45:00Z',
    },
    {
      id: '3',
      title: 'Following Directions',
      notes: 'Practiced 3-step directions with visual supports. 75% accuracy today.',
      start_datetime: '2024-10-25T10:00:00Z',
      end_datetime: '2024-10-25T10:45:00Z',
    },
  ]

  const handleSave = () => {
    // In real implementation, this would call updateStudent
    setEdit(false)
  }

  const handleCancel = () => {
    // Reset to mock data
    setIepGoals([
      'Improve articulation of /r/ sound in conversation',
      'Increase expressive vocabulary by 20 words',
      'Follow 3-step directions with 80% accuracy',
    ])
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

  return (
    <div className="min-h-screen h-screen flex items-center justify-center bg-background">
      <div className={`w-full h-full grid grid-rows-2 gap-8 ${PADDING} relative`}>
        {/* Edit toggle button */}
        <div className="absolute top-1/2 right-5 z-20 flex gap-2">
          {edit
            ? (
                <>
                  <Button
                    onClick={handleSave}
                    className="w-12 h-12 p-0 bg-green-600 hover:bg-green-700"
                    size="icon"
                  >
                    <Save size={20} />
                  </Button>
                  <Button
                    onClick={handleCancel}
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

        <div className="flex gap-8">
          {/* pfp and initials - placeholder */}
          <div className="flex flex-col items-center justify-between gap-2 w-1/5">
            <div className="w-full aspect-square border-2 border-accent rounded-full">
              <Avatar name="Student" variant="lorelei" className="w-full h-full" />
            </div>
            <div
              className={`w-full h-1/5 text-3xl font-bold flex items-center 
                justify-center bg-background border-border border-2 ${PADDING}
                 ${CORNER_ROUND}`}
            >
              S.S.
            </div>
          </div>
          {/* student schedule - no studentId, will show all sessions or empty */}
          <div className={`flex-[3] ${CORNER_ROUND} overflow-hidden bg-accent flex flex-col justify-between ${PADDING}`}>
            <StudentSchedule className="h-3/4" />
            <Button className="h-1/5 rounded-2xl text-lg font-bold " variant="secondary">
              View Student Schedule
            </Button>
          </div>
          <div className={`bg-accent flex-[2] flex flex-col items-center justify-between ${CORNER_ROUND} ${PADDING}`}>
            <div className="w-full h-3/4 text-3xl font-bold flex items-center rounded-2xl bg-primary">
              <MockRecentSession />
            </div>
            <Button className="w-full h-1/5 rounded-2xl text-lg font-bold " variant="secondary">
              View Student Attendance
            </Button>
          </div>
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
            <div className="w-full text-2xl text-primary flex items-baseline font-semibold">
              Session Notes
            </div>
            <div className="flex-1 overflow-y-auto">
              {mockSessionNotes.map(session => (
                <div
                  key={session.id}
                  className="w-full p-4 bg-background text-primary border-border border-b-2"
                >
                  <div className="w-full flex justify-between items-start">
                    <div className="font-semibold text-base">{session.title}</div>
                    <div className="text-sm opacity-90">
                      {new Date(session.start_datetime).toLocaleDateString()}
                    </div>
                  </div>
                  <div className="text-sm text-muted-foreground mt-2">{session.notes}</div>
                  <div className="text-xs opacity-75 mt-2">
                    {new Date(session.start_datetime).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' })}
                    {' - '}
                    {new Date(session.end_datetime).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' })}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default StudentPage
