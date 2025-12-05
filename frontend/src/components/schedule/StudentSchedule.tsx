'use client'
import type { View } from 'react-big-calendar'
import moment from 'moment'
import React from 'react'
import { Calendar, momentLocalizer } from 'react-big-calendar'
import { useSessions } from '@/hooks/useSessions'
import { useStudentSessions } from '@/hooks/useStudentSessions'
import 'react-big-calendar/lib/css/react-big-calendar.css'
import '@/app/calendar/override-calendar.css'

const localizer = momentLocalizer(moment)

// Color palette for events - using disabled variants from global styling
const eventColors = ['#f9aeda', '#bac0ff', '#ffdfb1', '#bac0ff', '#ffdfb1']
const getEventColor = (index: number) => eventColors[index % eventColors.length]

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
    : sessions.map((s, idx) => {
        // Check if this is from studentHook (nested) or allHook (flat)
        const sessionData = studentHook ? (s as any).session : s
        return {
          id: sessionData.id,
          title: sessionData.session_name ? sessionData.session_name : 'Session',
          start: new Date(sessionData.start_datetime),
          end: new Date(sessionData.end_datetime),
          allDay: false,
          resource: {
            color: getEventColor(idx),
            school: (sessionData as any).school || 'School',
          },
        }
      })

  const minTime = new Date()
  minTime.setHours(7, 0, 0)
  const maxTime = new Date()
  maxTime.setHours(19, 0, 0)

  const eventStyleGetter = (event: any) => {
    return {
      style: {
        backgroundColor: event.resource?.color || '#bac0ff',
        borderRadius: '0px',
        opacity: 0.95,
        color: '#000',
        border: 'none',
        display: 'block',
        fontWeight: 500,
        padding: '6px 8px',
        fontSize: '0.875rem',
        boxShadow: 'none',
      }
    }
  }

  return (
    <div className={`bg-card rounded-2xl overflow-hidden flex flex-col ${className}`}>

      {/* Calendar */}
      <div className="flex-1 overflow-hidden">
        <Calendar
          localizer={localizer}
          events={events}
          startAccessor="start"
          endAccessor="end"
          style={{ height: '100%' }}
          defaultView={initialView}
          view={currentView}
          onView={setCurrentView}
          views={['work_week',]}
          toolbar={false}
          min={minTime}
          max={maxTime}
          selectable={false}
          eventPropGetter={eventStyleGetter}
          step={30}
          timeslots={2}
          components={{
            header: ({ date }) => {
              const m = moment(date)
              return (
                <div className="rbc-custom-header" style={{ 
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  gap: '4px'
                }}>
                  <span style={{ 
                    fontSize: '0.75rem',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em'
                  }}>
                    {m.format("ddd")}
                  </span>

                  <div className="rbc-date-number" style={{
                    width: '38px',
                    height: '38px',
                    borderRadius: '50%',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontWeight: 700,
                    fontSize: '0.875rem',
                  }}>
                    {m.format("D")}
                  </div>
                </div>
              )
            }
          }}

        />
      </div>
    </div>
  )
}