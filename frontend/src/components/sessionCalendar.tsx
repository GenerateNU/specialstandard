// src/components/SessionCalendar.tsx
'use client'

import type { Session } from '@/types/session'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { useState } from 'react'

interface SessionCalendarProps {
  sessions: Session[]
  onSessionClick: (session: Session) => void
  onDateClick: () => void
}

export default function SessionCalendar({
  sessions,
  onSessionClick,
  onDateClick,
}: SessionCalendarProps) {
  const [currentDate, setCurrentDate] = useState(new Date())

  const monthNames = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December']
  const daysOfWeek = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']

  const getDaysInMonth = (date: Date) => {
    const year = date.getFullYear()
    const month = date.getMonth()
    const firstDay = new Date(year, month, 1)
    const lastDay = new Date(year, month + 1, 0)
    const daysInMonth = lastDay.getDate()
    const startingDayOfWeek = firstDay.getDay()

    const days: (Date | null)[] = []
    for (let i = 0; i < startingDayOfWeek; i++) {
      days.push(null)
    }
    for (let i = 1; i <= daysInMonth; i++) {
      days.push(new Date(year, month, i))
    }
    return days
  }

  const getSessionsForDate = (date: Date | null) => {
    if (!date)
      return []
    return sessions.filter((session) => {
      const sessionDate = new Date(session.start_datetime)
      return sessionDate.getDate() === date.getDate()
        && sessionDate.getMonth() === date.getMonth()
        && sessionDate.getFullYear() === date.getFullYear()
    })
  }

  const formatTime = (datetime: string) => {
    return new Date(datetime).toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true,
    })
  }

  const navigateMonth = (direction: number) => {
    const newDate = new Date(currentDate)
    newDate.setMonth(currentDate.getMonth() + direction)
    setCurrentDate(newDate)
  }

  return (
    <div>
      {/* Calendar Navigation */}
      <div className="flex justify-between items-center mb-4">
        <button
          onClick={() => navigateMonth(-1)}
          className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          title="Previous Month"
        >
          <ChevronLeft className="w-5 h-5" />
        </button>
        <h2 className="text-xl font-semibold">
          {monthNames[currentDate.getMonth()]}
          {' '}
          {currentDate.getFullYear()}
        </h2>
        <button
          onClick={() => navigateMonth(1)}
          className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          title="Next Month"
        >
          <ChevronRight className="w-5 h-5" />
        </button>
      </div>

      {/* Calendar Grid */}
      <div className="grid grid-cols-7 gap-1">
        {/* Day headers */}
        {daysOfWeek.map(day => (
          <div key={day} className="p-2 text-center text-sm font-semibold text-secondary">
            {day}
          </div>
        ))}

        {/* Calendar days */}
        {getDaysInMonth(currentDate).map((date, index) => {
          const daysSessions = getSessionsForDate(date)
          const isToday = date
            && date.toDateString() === new Date().toDateString()

          return (
            <div
              key={index}
              className={`min-h-24 p-2 border rounded-lg cursor-pointer ${
                !date
                  ? 'bg-gray-50'
                  : isToday
                    ? 'bg-accent-light border-accent'
                    : 'bg-card hover:bg-card-hover'
              } transition-colors`}
              onClick={() => date && onDateClick()}
            >
              {date && (
                <>
                  <div className="font-semibold text-sm mb-1">
                    {date.getDate()}
                  </div>
                  <div className="space-y-1">
                    {daysSessions.slice(0, 3).map(session => (
                      <div
                        key={session.id}
                        onClick={(e) => {
                          e.stopPropagation()
                          onSessionClick(session)
                        }}
                        className="text-xs p-1 bg-accent text-white rounded cursor-pointer hover:bg-accent-dark transition-colors"
                      >
                        {formatTime(session.start_datetime)}
                      </div>
                    ))}
                    {daysSessions.length > 3 && (
                      <div className="text-xs text-secondary">
                        +
                        {daysSessions.length - 3}
                        {' '}
                        more
                      </div>
                    )}
                  </div>
                </>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
