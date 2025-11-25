'use client'

import type { SlotInfo } from 'react-big-calendar'
import { motion } from 'motion/react'
import { useSearchParams } from 'next/navigation'
import React, { Suspense, useEffect } from 'react'
import AppLayout from '@/components/AppLayout'
import CalendarHeader from '@/components/calendar/calendarHeader'
import CalendarView from '@/components/calendar/calendarView'
import CardView from '@/components/calendar/cardView'
import { CreateSessionDialog } from '@/components/calendar/NewSessionModal'
import SessionPreviewModal from '@/components/SessionPreviewModal'
import { useCalendarData, useCalendarState } from '@/hooks/useCalendar'
import { useAuthContext } from '@/contexts/authContext'

function CalendarPage() {
  const { userId, isLoading: authLoading } = useAuthContext()
  const searchParams = useSearchParams()
  
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

  // Set view mode from URL query parameter
  useEffect(() => {
    const viewParam = searchParams.get('view')
    if (viewParam === 'card') {
      setViewMode('card')
    }
  }, [searchParams, setViewMode])

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

  // Show nothing while auth is loading - let Suspense handle the loading state
  if (authLoading || !userId) {
    return null
  }

  return (
    <div className="flex justify-center">
      <div style={{ width: '90vw', maxHeight: 'calc(115vh - 100px)' }} className="pt-10 flex flex-col gap-6 overflow-hidden">
        <CreateSessionDialog
          open={newSessionOpen}
          therapistId={userId}
          students={students}
          setOpen={handleCloseModal}
          onSubmit={async data => {
            await addSession(data)
          }}
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
              <motion.div
                key="card-view"
                initial={{ opacity: 0, y: 60 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: 60 }}
                transition={{ type: 'spring', damping: 30 }}
              >
                <CardView
                  date={date}
                  events={events}
                  onSelectSession={(session, position) => {
                    setSelectedSession(session)
                    setModalPosition(position)
                  }}
                />
              </motion.div>
            )
          : (
              <motion.div
                key="calendar-view"
                initial={{ opacity: 0, y: 60 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: 60 }}
                transition={{ type: 'spring', damping: 30 }}
              >
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
              </motion.div>
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
  )
}

export default function MyCalendar() {
  const [mounted, setMounted] = React.useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted) {
    return (
      <AppLayout>
        <div className="flex justify-center">
          <div style={{ width: '90vw', maxHeight: 'calc(115vh - 100px)' }} className="pt-10 flex flex-col gap-6 overflow-hidden">
            <div className="flex justify-center items-center" style={{ minHeight: '400px' }}>
              <div>Loading...</div>
            </div>
          </div>
        </div>
      </AppLayout>
    )
  }

  return (
    <AppLayout>
      <Suspense fallback={
        <div className="flex justify-center">
          <div style={{ width: '90vw', maxHeight: 'calc(115vh - 100px)' }} className="pt-10 flex flex-col gap-6 overflow-hidden">
            <div className="flex justify-center items-center" style={{ minHeight: '400px' }}>
              <div>Loading...</div>
            </div>
          </div>
        </div>
      }>
        <CalendarPage />
      </Suspense>
    </AppLayout>
  )
}