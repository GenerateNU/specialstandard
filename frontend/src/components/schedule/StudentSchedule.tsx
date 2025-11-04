'use client'

import type { View } from 'react-big-calendar'
import moment from 'moment'
import React from 'react'
import { Calendar, momentLocalizer } from 'react-big-calendar'
import { useSessions } from '@/hooks/useSessions'

import { useStudentSessions } from '@/hooks/useStudentSessions'
import 'react-big-calendar/lib/css/react-big-calendar.css'

const localizer = momentLocalizer(moment)

interface StudentScheduleProps {
  studentId?: string
  initialView?: View
  height?: string | number
  className?: string
}

export default function StudentSchedule({ studentId, initialView = 'day', className }: StudentScheduleProps) {
  // If a studentId is provided, call the student-specific endpoint, otherwise fall back to all sessions
  const studentHook = studentId ? useStudentSessions(studentId) : null
  const allHook = useSessions()

  const error = studentHook ? studentHook.error : allHook.error

  const sessions = studentHook ? studentHook.sessions : allHook.sessions

  // Always build events array (empty if error or no data)
  const events = error
    ? []
    : sessions.map(s => ({
        id: s.id,
        title: s.notes ? s.notes : 'Session',
        start: new Date(s.start_datetime),
        end: new Date(s.end_datetime),
        allDay: false,
      }))

  return (
    <div className={`bg-card rounded-2xl overflow-hidden flex flex-col ${className}`}>
      {/* Custom header */}
      <div className="px-4 py-2 border-b border-border">
        <div className="font-semibold text-lg text-primary">Today</div>
      </div>
      {/* Calendar - hide its built-in header */}
      <div className="flex-1 overflow-hidden">
        <style>
          {`
            .rbc-time-view .rbc-time-header {
              display: none;
            }
          `}
        </style>
        <Calendar
          localizer={localizer}
          events={events}
          startAccessor="start"
          endAccessor="end"
          style={{ height: '100%' }}
          defaultView={initialView}
          views={['day']}
          toolbar={false}
        />
      </div>
    </div>
  )
}
