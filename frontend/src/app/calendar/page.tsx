'use client'

import type { SlotInfo } from 'react-big-calendar'
import AppLayout from '@/components/AppLayout'
import CalendarHeader from '@/components/calendar/calendarHeader'
import CalendarView from '@/components/calendar/calendarView'
import CardView from '@/components/calendar/cardView'
import { CreateSessionDialog } from '@/components/calendar/NewSessionModal'
import SessionPreviewModal from '@/components/SessionPreviewModal'
import { useCalendarData, useCalendarState } from '@/hooks/useCalendar'

const HARDCODED_THERAPIST_ID = '9dad94d8-6534-4510-90d7-e4e97c175a65'

export default function MyCalendar() {
  const {
    date,
    setDate,
    view,
    setView,
    viewMode,
    setViewMode,
    selectedSlot,
    setSelectedSlot,
    selectedSession,
    setSelectedSession,
    modalPosition,
    setModalPosition,
    newSessionOpen,
    setNewSessionOpen,
  } = useCalendarState()

  const { students, events, isLoading, error, addSession } = useCalendarData(
    date,
    view,
  )

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

  const handleSelectEvent = (event: any, e: React.SyntheticEvent) => {
    const target = e.target as HTMLElement
    const eventElement = target.closest('.rbc-event') as HTMLElement
    const rect = eventElement?.getBoundingClientRect() || target.getBoundingClientRect()

    setModalPosition({
      x: rect.right + 10,
      y: rect.top,
    })
    setSelectedSession(event.resource)
  }

  return (
    <AppLayout>
      <div className="flex justify-center">
        <div style={{ width: '90vw', maxHeight: 'calc(100vh - 100px)' }} className="pt-10">
          <CreateSessionDialog
            open={newSessionOpen}
            therapistId={HARDCODED_THERAPIST_ID}
            students={students}
            setOpen={handleCloseModal}
            onSubmit={async data => addSession(data)}
            initialDateTime={selectedSlot ?? undefined}
          />

          <CalendarHeader
            viewMode={viewMode}
            onViewModeChange={setViewMode}
            onAddSession={() => setNewSessionOpen(true)}
            date={date}
            view={view}
            onNavigate={setDate}
            onViewChange={setView}
          />

          {viewMode === 'card'
            ? (
                <CardView
                  date={date}
                  events={events}
                  onSelectSession={(session, position) => {
                    setSelectedSession(session)
                    setModalPosition(position)
                  }}
                />
              )
            : (
                <CalendarView
                  date={date}
                  view={view}
                  events={events}
                  isLoading={isLoading}
                  error={error}
                  onNavigate={setDate}
                  onViewChange={setView}
                  onSelectEvent={handleSelectEvent}
                  onSelectSlot={handleSelectSlot}
                />
              )}

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
      </div>
    </AppLayout>
  )
}
