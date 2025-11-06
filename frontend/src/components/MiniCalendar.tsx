import { cn } from '@/lib/utils'

interface MiniCalendarProps {
  className?: string
}

export default function MiniCalendar({ className }: MiniCalendarProps) {
  const today = new Date()
  const currentMonth = today.getMonth()
  const currentYear = today.getFullYear()
  const currentDate = today.getDate()

  // Get first day of month (0 = Sunday, 6 = Saturday)
  const firstDayOfMonth = new Date(currentYear, currentMonth, 1).getDay()

  // Get number of days in current month
  const daysInMonth = new Date(currentYear, currentMonth + 1, 0).getDate()

  // Month names
  const monthNames = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December',
  ]

  // Day abbreviations
  const dayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']

  // Create array of day numbers
  const days = []
  // Add empty slots for days before month starts
  for (let i = 0; i < firstDayOfMonth; i++) {
    days.push(null)
  }
  // Add actual days
  for (let i = 1; i <= daysInMonth; i++) {
    days.push(i)
  }

  return (
    <div className={cn('p-4 select-none', className)}>
      {/* Month and Year Header */}
      <div className="text-center mb-2">
        <h3 className="text-xl font-semibold text-foreground">
          {monthNames[currentMonth]}
          {' '}
          {currentYear}
        </h3>
      </div>

      {/* Day Names */}
      <div className="grid grid-cols-7 gap-2 mb-2">
        {dayNames.map((day, index) => (
          <div
            key={`${day}-${index}`}
            className="text-center text-sm font-medium text-muted-foreground"
          >
            {day.charAt(0)}
          </div>
        ))}
      </div>

      {/* Calendar Grid */}
      <div className="grid grid-cols-7 gap-1">
        {days.map((day, index) => (
          <div
            key={index}
            className={cn(
              'h-5 flex items-center justify-center text-sm rounded-full transition-colors select-none',
              day === null && 'invisible',
              day === currentDate && 'bg-pink text-background font-bold shadow-sm',
              day !== null && day !== currentDate && 'text-foreground hover:bg-pink-disabled',
            )}
          >
            {day}
          </div>
        ))}
      </div>
    </div>
  )
}
