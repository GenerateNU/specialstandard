import type { View } from 'react-big-calendar'
import { ArrowLeft, ArrowRight } from 'lucide-react'

interface CustomToolbarProps {
  label: string
  onNavigate: (action: 'PREV' | 'NEXT' | 'TODAY') => void
  onView: (view: View) => void
  view: View
  showViewSelector?: boolean // Add this
}

export default function CustomToolbar({
  label,
  onNavigate,
  onView,
  view,
  showViewSelector = true,
}: CustomToolbarProps) {
  return (
    <div className="rbc-toolbar">
      <span className="rbc-btn-group">
        <button type="button" onClick={() => onNavigate('TODAY')}>
          Today
        </button>
      </span>
      <span className="rbc-toolbar-label-container">
        <button type="button" className="rbc-toolbar-nav-btn" onClick={() => onNavigate('PREV')}>
          <ArrowLeft size={20} />
        </button>
        <span className="rbc-toolbar-label">{label}</span>
        <button type="button" className="rbc-toolbar-nav-btn " onClick={() => onNavigate('NEXT')}>
          <ArrowRight size={20} />
        </button>
      </span>
      {showViewSelector
        ? (
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
          )
        : (
            <span className="rbc-btn-group w-1/6" />
          )}
    </div>
  )
}
