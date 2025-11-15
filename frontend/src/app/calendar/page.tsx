'use client'

import type { SlotInfo, View } from 'react-big-calendar'
import type { Session } from '@/lib/api/theSpecialStandardAPI.schemas'
import { ArrowLeft, ArrowRight, Book, Plus } from 'lucide-react'
import moment from 'moment'
import Link from 'next/link'
import { useState } from 'react'
import { Calendar, momentLocalizer } from 'react-big-calendar'
import AppLayout from '@/components/AppLayout'
import { CreateSessionDialog } from '@/components/calendar/NewSessionModal'
import SessionPreviewModal from '@/components/SessionPreviewModal'

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

// Custom toolbar component with arrows on sides
function CustomToolbar({ label, onNavigate, onView, view }: {
  label: string
  onNavigate: (action: 'PREV' | 'NEXT' | 'TODAY') => void
  onView: (view: View) => void
  view: View
}) {
  return (
    <div className="rbc-toolbar">
      <span className="rbc-btn-group">
        <button type="button" onClick={() => onNavigate('TODAY')}>
          Today
        </button>
      </span>
      <span className="rbc-toolbar-label-container">
        <button
          type="button"
          className="rbc-toolbar-nav-btn"
          onClick={() => onNavigate('PREV')}
        >
          <ArrowLeft size={20} />
        </button>
        <span className="rbc-toolbar-label">{label}</span>
        <button
          type="button"
          className="rbc-toolbar-nav-btn "
          onClick={() => onNavigate('NEXT')}
        >
          <ArrowRight size={20} />
        </button>
      </span>
      <span className="rbc-btn-group">
        <button
          type="button"
          className={view === 'day' ? 'rbc-active' : ''}
          onClick={() => onView('day')}
        >
          Day
        </button>
        <button
          type="button"
          className={view === 'work_week' ? 'rbc-active' : ''}
          onClick={() => onView('work_week')}
        >
          Week
        </button>
        <button
          type="button"
          className={view === 'month' ? 'rbc-active' : ''}
          onClick={() => onView('month')}
        >
          Month
        </button>

      </span>
    </div>
  )
}

export default function MyCalendar() {
  const { students } = useStudents()
  const [date, setDate] = useState(new Date())
  const [view, setView] = useState<View>('work_week')
  const [newSessionOpen, setNewSessionOpen] = useState(false)
  const [selectedSlot, setSelectedSlot] = useState<{ start: Date, end: Date } | null>(null)
  const [selectedSession, setSelectedSession] = useState<Session | null>(null)
  const [modalPosition, setModalPosition] = useState<{ x: number, y: number } | null>(null)

  const getViewRange = () => {
    const startOfView = moment(date).startOf(view === 'day' ? 'day' : view === 'work_week' ? 'isoWeek' : 'month').toDate()
    const endOfView = moment(date).endOf(view === 'day' ? 'day' : view === 'work_week' ? 'isoWeek' : 'month').toDate()
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
    title: session.session_name,
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
    // Find the actual event container for consistent positioning
    const eventElement = target.closest('.rbc-event') as HTMLElement
    const rect = eventElement?.getBoundingClientRect() || target.getBoundingClientRect()

    // Position modal to the right of the event block with 10px gap
    setModalPosition({
      x: rect.right + 10,
      y: rect.top,
    })
    setSelectedSession(event.resource)
  }

  return (
    <AppLayout>
      <div className="flex justify-center">
        <div style={{ width: '90vw', maxHeight: 'calc(100vh - 100px)' }} className="pt-4">
          <CreateSessionDialog
            open={newSessionOpen}
            therapistId={HARDCODED_THERAPIST_ID}
            students={students}
            setOpen={handleCloseModal}
            onSubmit={async data => addSession(data)}
            initialDateTime={selectedSlot ?? undefined}
          />

          {/* Header with back button and action buttons */}
          <div className="mb-3 flex justify-between items-end" style={{ paddingLeft: '24px', paddingRight: '24px' }}>
            <div className="flex flex-col items-left justify-start">
              <Link
                href="/"
                className="inline-flex items-center gap-2 text-secondary hover:text-primary mb-2 transition-colors group self-start"
              >
                <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
                <span className="text-sm font-medium">Back</span>
              </Link>
              <h1>Your Schedule</h1>
            </div>
            <div className="flex gap-4">
              <button
                type="button"
                onClick={() => { /* TODO: Implement card view */ }}
                className="inline-flex items-center gap-2 px-4 py-2 text-secondary hover:text-primary border border-secondary/20 hover:border-primary/30 rounded-lg cursor-pointer transition-all"
              >
                <span className="font-medium">Card View</span>
                <Book size={18} />
              </button>
              <button
                type="button"
                onClick={() => setNewSessionOpen(true)}
                className="inline-flex items-center gap-2 text-pink hover:text-primary-hover cursor-pointer transition-colors group"
              >
                <span className="font-bold text-pink hover:inherit">Add Session</span>
                <span className="flex items-center justify-center w-6 h-6 bg-pink text-white rounded text-sm font-bold transition-transform cursor-pointer group-hover:scale-110 will-change-transform">
                  <Plus strokeWidth={3} size={16} />
                </span>
              </button>
            </div>
          </div>

          {/* Custom toolbar outside calendar */}
          <div className="mb-4" style={{ paddingLeft: '24px', paddingRight: '24px' }}>
            <CustomToolbar
              label={
                view === 'month'
                  ? moment(date).format('MMMM YYYY')
                  : view === 'work_week'
                    ? `${moment(date).startOf('isoWeek').format('MMMM D')} - ${moment(date).endOf('isoWeek').format('D')}`
                    : moment(date).format('dddd, MMMM D')
              }
              onNavigate={(action) => {
                if (action === 'TODAY') {
                  setDate(new Date())
                }
                else if (action === 'NEXT') {
                  const unit = view === 'month' ? 'months' : view === 'work_week' ? 'weeks' : 'days'
                  setDate(moment(date).add(1, unit).toDate())
                }
                else if (action === 'PREV') {
                  const unit = view === 'month' ? 'months' : view === 'work_week' ? 'weeks' : 'days'
                  setDate(moment(date).subtract(1, unit).toDate())
                }
              }}
              onView={setView}
              view={view}
            />
          </div>

          {isLoading
            ? (
                <div className="flex items-center justify-center" style={{ height: '70vh', backgroundColor: 'var(--card-bg)', borderRadius: '16px' }}>
                  <div className="text-primary">Loading sessions...</div>
                </div>
              )
            : error
              ? (
                  <div className="flex items-center justify-center" style={{ height: '70vh', backgroundColor: 'var(--card-bg)', borderRadius: '16px' }}>
                    <div className="text-error">
                      Error loading sessions:
                      {' '}
                      {error}
                    </div>
                  </div>
                )
              : (
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
                    views={['day', 'work_week', 'month']}
                    selectable
                    onSelectSlot={handleSelectSlot}
                    toolbar={false}
                    // Better event display
                    eventPropGetter={() => ({
                      style: {
                        borderRadius: '8px',
                      },
                    })}
                    // Show time range in day/week view
                    formats={{
                      eventTimeRangeFormat: () => '',
                      timeGutterFormat: (date, culture, localizer) =>
                        localizer?.format(date, 'h A', culture) || '',
                    }}
                    // Custom toolbar and header components
                    components={{
                      toolbar: CustomToolbar,
                      month: {
                        header: ({ date, localizer }) => (
                          <div style={{ fontSize: '0.875rem', fontWeight: 600, textAlign: 'center', padding: '12px 8px' }}>
                            {localizer.format(date, 'ddd', 'en').toUpperCase()}
                          </div>
                        ),
                      },
                      work_week: {
                        header: ({ date, localizer }) => {
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
                        header: ({ date, localizer }) => {
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
                    }}
                  />
                )}

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
      </div>
    </AppLayout>
  )
}
