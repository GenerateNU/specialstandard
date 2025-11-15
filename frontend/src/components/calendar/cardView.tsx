import type { CalendarEvent } from '@/hooks/useCalendar'
import type { Session } from '@/lib/api/theSpecialStandardAPI.schemas'

import { Clock } from 'lucide-react'
import moment from 'moment'
import CardBookBg from './CardBookBg'
import CardViewSessionStudents from './cardViewSessionStudents'

interface CardViewProps {
  date: Date
  events: CalendarEvent[]
  onSelectSession: (session: Session, position: { x: number, y: number }) => void
}

// Hash function to consistently map session ID to a color
function hashIdToColor(id: string | number): 'blue' | 'yellow' | 'pink' {
  const colors: Array<'blue' | 'yellow' | 'pink'> = ['blue', 'yellow', 'pink']
  
  // Convert string to a numeric hash
  let hash = 0
  const idStr = String(id)
  for (let i = 0; i < idStr.length; i++) {
    hash = ((hash << 5) - hash) + idStr.charCodeAt(i)
    hash = hash & hash // Convert to 32-bit integer
  }
  
  return colors[Math.abs(hash) % colors.length]
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
                    // Map through sessions for the day
                    daySessions.map(event => (
                      <div key={event.id} className="flex flex-col gap-1 items-center">
                        <CardBookBg 
                          className="hover:scale-102 transition" 
                          size="md" 
                          color={hashIdToColor(event.id)}
                        >
                          <button
                            type="button"
                            onClick={() => {
                              onSelectSession(event.resource, {
                                x: window.innerWidth / 2,
                                y: window.innerHeight / 2,
                              })
                            }}
                            className="flex flex-col items-start gap-1 p-4 rounded-lg border-0 cursor-pointer transition-transform text-black transition-inherit focus:outline-none bg-transparent shadow-none w-full h-full min-h-[120px]"
                            style={{ background: 'none' }}
                          >
                            <div className="text-sm font-semibold">Session Name</div>
                            <div className="text-sm">School A</div>
                            <br />
                            <div className="text-sm">{moment(event.start).format('D MMM YYYY')}</div>
                            <div className="flex items-center gap-1.5 text-xs font-medium">
                              <Clock size={14} />
                              {moment(event.start).format('h:mm A')}
                            </div>
                          </button>
                        </CardBookBg>
                        <CardViewSessionStudents sessionId={event.id} />
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