'use client'

import type { SlotInfo, View } from 'react-big-calendar'
import type { Session } from '@/lib/api/theSpecialStandardAPI.schemas'
import { ArrowLeft } from 'lucide-react'
import moment from 'moment'
import Link from 'next/link'
import { useState } from 'react'
import { Calendar, momentLocalizer } from 'react-big-calendar'
import { CreateSessionDialog } from '@/components/calendar/NewSessionModal'
import SessionPreviewModal from '@/components/SessionPreviewModal'
import { Button } from '@/components/ui/button'
import { useSessions } from '@/hooks/useSessions'
import { useStudents } from '@/hooks/useStudents'
import 'react-big-calendar/lib/css/react-big-calendar.css'
import './override-calendar.css'

const localizer = momentLocalizer(moment)

interface CalendarEvent {
  id: string
  title: string
  start: Date
  end: Date
  resource: Session
}

// Hard coded therapist ID as requested (using existing therapist from database)
const HARDCODED_THERAPIST_ID = '9dad94d8-6534-4510-90d7-e4e97c175a65' // John Doe

export default function MyCalendar() {
  const { students } = useStudents()
  const [date, setDate] = useState(new Date())
  const [view, setView] = useState<View>('week')
  const [newSessionOpen, setNewSessionOpen] = useState(false)
  const [selectedSlot, setSelectedSlot] = useState<{ start: Date, end: Date } | null>(null)
  const [selectedSession, setSelectedSession] = useState<Session | null>(null)
  const [modalPosition, setModalPosition] = useState<{ x: number, y: number } | null>(null)

  const getViewRange = () => {
    const startOfView = moment(date).startOf(view === 'day' ? 'day' : view === 'week' ? 'week' : 'month').toDate()
    const endOfView = moment(date).endOf(view === 'day' ? 'day' : view === 'week' ? 'week' : 'month').toDate()
    return { startdate: startOfView, enddate: endOfView }
  }

  const viewRange = getViewRange()
  const { sessions, isLoading, error, addSession } = useSessions({
    startdate: viewRange.startdate.toISOString(),
    enddate: viewRange.enddate.toISOString(),
  })

  // Transform sessions into calendar events
  const events: CalendarEvent[] = sessions.map(session => ({
    id: session.id,
    title: 'session',
    start: new Date(session.start_datetime),
    end: new Date(session.end_datetime),
    resource: session,
  }))

  const handleSelectSlot = (slotInfo: SlotInfo) => {
    setSelectedSlot({
      start: slotInfo.start as Date,
      end: slotInfo.end as Date,
    })
    setNewSessionOpen(true)
  }

  const handleCloseModal = () => {
    setNewSessionOpen(false)
    setSelectedSlot(null)
  }

  const handleSelectEvent = (event: CalendarEvent, e: React.SyntheticEvent) => {
    const target = e.target as HTMLElement
    const rect = target.getBoundingClientRect()
    
    // Position modal to the right of the clicked event
    setModalPosition({
      x: rect.right + 10,
      y: rect.top,
    })
    setSelectedSession(event.resource)
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div>Loading sessions...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div>
          Error loading sessions:
          {error}
        </div>
      </div>
    )
  }

  return (
    <div>
      {/* Back button */}
      <Link
        href="/"
        className="inline-flex items-center gap-2 text-secondary hover:text-primary mb-4 transition-colors group"
      >
        <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
        <span className="text-sm font-medium">Back to Home</span>
      </Link>

      <Button variant="default" onClick={() => setNewSessionOpen(true)}>
        + New Session
      </Button>

      <CreateSessionDialog
        open={newSessionOpen}
        therapistId={HARDCODED_THERAPIST_ID}
        students={students}
        setOpen={handleCloseModal}
        onSubmit={async data => addSession(data)}
        initialDateTime={selectedSlot ?? undefined}
      />

      <div className="flex items-center justify-center h-screen">
        <Calendar
          localizer={localizer}
          events={events}
          startAccessor="start"
          endAccessor="end"
          style={{ height: '80vh', width: '90vw' }}
          date={date}
          view={view}
          onNavigate={setDate}
          onView={setView}
          onSelectEvent={handleSelectEvent}
          views={['week', 'day', 'month']}
          selectable
          onSelectSlot={handleSelectSlot}
        />
      </div>

      {/* Session preview modal */}
      {selectedSession && modalPosition && (
        <SessionPreviewModal
          session={selectedSession}
          position={modalPosition}
          onClose={() => {
            setSelectedSession(null)
            setModalPosition(null)
          }}
        />
      )}
    </div>
  )
}
