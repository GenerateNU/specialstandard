import type { CalendarEvent } from '@/hooks/useCalendar'
import type { Session } from '@/lib/api/theSpecialStandardAPI.schemas'
import { Clock } from 'lucide-react'
import moment from 'moment'
import SessionStudents from './sessionStudents'

interface CardViewProps {
  date: Date
  events: CalendarEvent[]
  onSelectSession: (session: Session, position: { x: number, y: number }) => void
}

export default function CardView({ date, events, onSelectSession }: CardViewProps) {
  return (
    <div
      className="grid grid-cols-5 p-6 rounded-2xl overflow-y-auto"
      style={{
        backgroundColor: '',
        width: '90vw',
        minHeight: '600px',
        maxHeight: '80vh',
      }}
    >
      {Array.from({ length: 5 }).map((_, dayIndex) => {
        const currentDay = moment(date).startOf('isoWeek').add(dayIndex, 'days')
        const isToday = currentDay.isSame(new Date(), 'day')
        const daySessions = events.filter(event => moment(event.start).isSame(currentDay, 'day'))

        return (
          <div
            key={dayIndex}
            className={`flex flex-col min-h-[400px] px-6 ${dayIndex < 4 ? 'border-r border-black' : ''}`}
          >
            <div className="flex flex-col items-center mb-3 pt-3">
              <div className="text-base font-semibold" style={{ color: 'var(--text-primary)' }}>
                {currentDay.format('ddd').toUpperCase()}
              </div>
              <div
                className="text-base font-semibold mt-1 rounded-full w-8 h-8 flex items-center justify-center"
                style={{
                  backgroundColor: isToday ? 'var(--color-orange)' : 'transparent',
                  color: isToday ? 'var(--color-white)' : 'var(--text-primary)',
                }}
              >
                {currentDay.format('D')}
              </div>
            </div>

            <div className="flex flex-col gap-2">
              {daySessions.length === 0
                ? (
                    <div className="text-sm text-center mt-5" style={{ color: 'var(--text-muted)' }}>
                      No sessions
                    </div>
                  )
                : (
                    daySessions.map(event => (
                      <div key={event.id} className="flex flex-col gap-1">
                        <button
                          type="button"
                          onClick={() => {
                            onSelectSession(event.resource, {
                              x: window.innerWidth / 2,
                              y: window.innerHeight / 2,
                            })
                          }}
                          className="flex flex-col items-start gap-1 p-3 rounded-lg border-0 cursor-pointer transition-all hover:scale-103"
                          style={{
                            backgroundColor: 'var(--color-blue)',
                            color: 'var(--color-white)',
                          }}
                        >
                          <div className="text-sm font-semibold">Session</div>
                          <div className="flex items-center gap-1.5 text-xs font-medium">
                            <Clock size={14} />
                            {moment(event.start).format('h:mm A')}
                          </div>
                        </button>
                        <SessionStudents sessionId={event.id} />
                      </div>
                    ))
                  )}
            </div>
          </div>
        )
      })}
    </div>
  )
}
