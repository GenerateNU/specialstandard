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
  const [currentView, setCurrentView] = React.useState<View>(initialView)

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

  // Dynamic header text based on view
  const getHeaderText = () => {
    switch (currentView) {
      case 'day':
        return 'Today'
      case 'week':
        return 'This Week'
      case 'month':
        return 'This Month'
      default:
        return 'Schedule'
    }
  }

  return (
    <div className={`bg-card rounded-2xl overflow-hidden flex flex-col ${className}`}>
      {/* Custom header */}
      <div className="px-4 py-2 border-b border-border">
        <div className="font-semibold text-lg text-primary">{getHeaderText()}</div>
      </div>
      {/* Calendar */}
      <div className="flex-1 overflow-hidden">
        <style>
          {`
            .rbc-time-view .rbc-time-header {
              ${currentView === 'day' ? 'display: none;' : ''}
            }
            .rbc-today {
              background-color: rgba(59, 130, 246, 0.3) !important;
            }
            .rbc-current-time-indicator {
              background-color: #3b82f6 !important;
            }
            .rbc-day-bg.rbc-today {
              background-color: rgba(59, 130, 246, 0.05) !important;
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
          view={currentView}
          onView={setCurrentView}
          views={['day', 'week', 'month']}
          toolbar={false}
        />
      </div>
    </div>
  )
}
