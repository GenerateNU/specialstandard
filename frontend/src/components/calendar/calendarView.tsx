import type { SlotInfo, View } from 'react-big-calendar'
import type { CalendarEvent } from '@/hooks/useCalendar'
import moment from 'moment'
import { Calendar, momentLocalizer } from 'react-big-calendar'
import CustomToolbar from './customToolbar'

import 'react-big-calendar/lib/css/react-big-calendar.css'
import '@/app/calendar/override-calendar.css'

const localizer = momentLocalizer(moment)

interface CalendarViewProps {
  date: Date
  view: View
  events: CalendarEvent[]
  isLoading: boolean
  error: string | null
  onNavigate: (date: Date) => void
  onViewChange: (view: View) => void
  onSelectEvent: (event: CalendarEvent, e: React.SyntheticEvent) => void
  onSelectSlot: (slotInfo: SlotInfo) => void
}

export default function CalendarView({
  date,
  view,
  events,
  isLoading,
  error,
  onNavigate,
  onViewChange,
  onSelectEvent,
  onSelectSlot,
}: CalendarViewProps) {
  if (isLoading) {
    return (
      <div
        className="flex items-center justify-center"
        style={{ height: '70vh', backgroundColor: 'var(--card-bg)', borderRadius: '16px' }}
      >
        <div className="text-primary">Loading sessions...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div
        className="flex items-center justify-center"
        style={{ height: '70vh', backgroundColor: 'var(--card-bg)', borderRadius: '16px' }}
      >
        <div className="text-error">
          Error loading sessions:
          {error}
        </div>
      </div>
    )
  }

  return (
    <Calendar
      localizer={localizer}
      events={events}
      startAccessor="start"
      endAccessor="end"
      style={{ height: '80vh', width: '90vw' }}
      date={date}
      view={view}
      onNavigate={onNavigate}
      onView={onViewChange}
      onSelectEvent={onSelectEvent}
      views={['day', 'work_week', 'month']}
      selectable
      onSelectSlot={onSelectSlot}
      toolbar={false}
      eventPropGetter={() => ({
        style: {
          borderRadius: '8px',
        },
      })}
      formats={{
        eventTimeRangeFormat: () => '',
        timeGutterFormat: (date, culture, localizer) =>
          localizer?.format(date, 'h A', culture) || '',
      }}
      components={getCalendarComponents()}
    />
  )
}

export function getCalendarComponents() {
  return {
    toolbar: CustomToolbar,
    month: {
      header: ({ date, localizer }: any) => (
        <div
          style={{
            fontSize: '0.875rem',
            fontWeight: 600,
            textAlign: 'center',
            padding: '12px 8px',
          }}
        >
          {localizer.format(date, 'ddd', 'en').toUpperCase()}
        </div>
      ),
    },
    work_week: {
      header: ({ date, localizer }: any) => {
        const isToday = moment(date).isSame(new Date(), 'day')
        return (
          <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
            <div style={{ fontSize: '1.00rem', fontWeight: 600 }}>
              {localizer.format(date, 'ddd', 'en').toUpperCase()}
            </div>
            <div
              style={{
                fontSize: '1.00rem',
                fontWeight: 600,
                marginTop: '4px',
                backgroundColor: isToday ? 'var(--color-orange)' : 'transparent',
                borderRadius: '50%',
                width: '32px',
                height: '32px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: isToday ? 'var(--color-white)' : 'inherit',
              }}
            >
              {localizer.format(date, 'D', 'en')}
            </div>
          </div>
        )
      },
    },
    day: {
      header: ({ date, localizer }: any) => {
        const isToday = moment(date).isSame(new Date(), 'day')
        return (
          <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
            <div style={{ fontSize: '0.75rem', fontWeight: 600 }}>
              {localizer.format(date, 'ddd', 'en').toUpperCase()}
            </div>
            <div
              style={{
                fontSize: '1.25rem',
                fontWeight: 600,
                marginTop: '4px',
                backgroundColor: isToday ? 'var(--color-orange)' : 'transparent',
                borderRadius: '50%',
                width: '40px',
                height: '40px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: isToday ? 'var(--color-black)' : 'inherit',
              }}
            >
              {localizer.format(date, 'D', 'en')}
            </div>
          </div>
        )
      },
    },
  }
}
