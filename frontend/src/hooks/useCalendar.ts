// hooks/useCalendar.ts - UPDATED WITH PARENT ID

import type { View } from 'react-big-calendar'
import type { Session } from '@/lib/api/theSpecialStandardAPI.schemas'
import moment from 'moment'
import { useMemo, useState } from 'react'
import { useSessions } from '@/hooks/useSessions'
import { useStudents } from '@/hooks/useStudents'

export function useCalendarState() {
  const [date, setDate] = useState(new Date())
  const [view, setView] = useState<View>('work_week')
  const [viewMode, setViewMode] = useState<'calendar' | 'card'>('calendar')
  const [newSessionOpen, setNewSessionOpen] = useState(false)
  const [selectedSlot, setSelectedSlot] = useState<{ start: Date, end: Date } | null>(null)
  const [selectedSession, setSelectedSession] = useState<Session | null>(null)
  const [modalPosition, setModalPosition] = useState<{ x: number, y: number } | null>(null)

  return {
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
  }
}

export interface CalendarEvent {
  id: string
  parentId: string // Add parent ID for color hashing
  title: string
  start: Date
  end: Date
  resource: Session
}

export function useCalendarData(date: Date, view: View) {
  const { students } = useStudents()

  const viewRange = useMemo(() => {
    const startOfView = moment(date)
      .startOf(view === 'day' ? 'day' : view === 'work_week' ? 'isoWeek' : 'month')
      .toDate()
    const endOfView = moment(date)
      .endOf(view === 'day' ? 'day' : view === 'work_week' ? 'isoWeek' : 'month')
      .toDate()
    return { startdate: startOfView, enddate: endOfView }
  }, [date, view])

  const { sessions, isLoading, error, addSession } = useSessions({
    startdate: viewRange.startdate.toISOString(),
    enddate: viewRange.enddate.toISOString(),
  })

  const events: CalendarEvent[] = useMemo(
    () =>
      sessions.map(session => ({
        id: session.id,
        parentId: session.session_parent_id || session.id, // Use parent ID if available, else fall back to session ID
        title: session.session_name,
        location: session.location,
        start: new Date(session.start_datetime),
        end: new Date(session.end_datetime),
        resource: session,
      })),
    [sessions],
  )

  return {
    sessions,
    students,
    events,
    isLoading,
    error,
    addSession,
    viewRange,
  }
}