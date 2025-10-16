'use client'

import type { View } from 'react-big-calendar'
import { ArrowLeft } from 'lucide-react'
import moment from 'moment'
import Link from 'next/link'
import { useState } from 'react'
import { Calendar, momentLocalizer } from 'react-big-calendar'
import { useSessions } from '@/hooks/useSessions'
import 'react-big-calendar/lib/css/react-big-calendar.css'
import './override-calendar.css'

const localizer = momentLocalizer(moment)

// Here, we are defining a calendar event type
interface CalendarEvent {
  id: string
  title: string
  start: Date
  end: Date
}

export default function MyCalendar() {
  const { sessions, isLoading, error } = useSessions()
  const [date, setDate] = useState(new Date())
  const [view, setView] = useState<View>('week')

  // Transform sessions into calendar events
  const events: CalendarEvent[] = sessions.map(session => ({
    id: session.id,
    title: 'session',
    start: new Date(session.start_datetime),
    end: new Date(session.end_datetime),
  }))

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
          views={['week', 'day', 'month']}
        />
      </div>
    </div>

  )
}
