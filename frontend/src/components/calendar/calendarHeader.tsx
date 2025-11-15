import type { View } from 'react-big-calendar'
import { Book, Calendar, Plus } from 'lucide-react'
import moment from 'moment'
import CustomToolbar from './customToolbar'

interface CalendarHeaderProps {
  viewMode: 'calendar' | 'card'
  onViewModeChange: (mode: 'calendar' | 'card') => void
  onAddSession: () => void
  date: Date
  view: View
  onNavigate: (date: Date) => void
  onViewChange: (view: View) => void
}

export default function CalendarHeader({
  viewMode,
  onViewModeChange,
  onAddSession,
  date,
  view,
  onNavigate,
  onViewChange,
}: CalendarHeaderProps) {
  const handleNavigate = (action: 'PREV' | 'NEXT' | 'TODAY') => {
    if (action === 'TODAY') {
      onNavigate(new Date())
    }
    else if (action === 'NEXT') {
      const unit = view === 'month' ? 'months' : view === 'work_week' ? 'weeks' : 'days'
      onNavigate(moment(date).add(1, unit).toDate())
    }
    else if (action === 'PREV') {
      const unit = view === 'month' ? 'months' : view === 'work_week' ? 'weeks' : 'days'
      onNavigate(moment(date).subtract(1, unit).toDate())
    }
  }

  const label
    = view === 'month'
      ? moment(date).format('MMMM YYYY')
      : view === 'work_week'
        ? `${moment(date).startOf('isoWeek').format('MMMM D')} - ${moment(date).endOf('isoWeek').format('D')}`
        : moment(date).format('dddd, MMMM D')

  return (
    <>
      <div
        className="flex justify-between items-end"
        style={{ paddingLeft: '24px', paddingRight: '24px' }}
      >
        <div className="flex flex-col items-left justify-start">
          <h1>Your Schedule</h1>
        </div>
        <div className="flex gap-4">
          <button
            type="button"
            onClick={() => onViewModeChange(viewMode === 'calendar' ? 'card' : 'calendar')}
            className="inline-flex items-center gap-2 text-pink hover:text-primary-hover cursor-pointer transition-colors group"
          >
            <span className="font-bold text-pink hover:inherit">
              {viewMode === 'calendar' ? 'Card View' : 'Calendar View'}
            </span>
            <span className="group-hover:scale-110 transition will-change-transform">
              {viewMode === 'card' ? <Calendar /> : <Book />}
            </span>
          </button>
          <button
            type="button"
            onClick={onAddSession}
            className="inline-flex items-center gap-2 text-pink hover:text-primary-hover cursor-pointer transition-colors group"
          >
            <span className="font-bold text-pink hover:inherit">Add Session</span>
            <span className="flex items-center justify-center w-6 h-6 bg-pink text-white rounded text-sm font-bold transition-transform cursor-pointer group-hover:scale-110 will-change-transform">
              <Plus strokeWidth={3} size={16} />
            </span>
          </button>
        </div>
      </div>

      <div className="" style={{ paddingLeft: '24px', paddingRight: '24px' }}>
        <CustomToolbar
          label={label}
          onNavigate={handleNavigate}
          onView={onViewChange}
          view={view}
          showViewSelector={viewMode === 'calendar'} // Add this line
        />
      </div>
    </>
  )
}
