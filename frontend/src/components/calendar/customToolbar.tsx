import type { View } from 'react-big-calendar'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '../ui/button'

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
    <div className="rbc-toolbar flex items-center justify-between w-full gap-4">
      <span className="rbc-btn-group">
        <button type="button" onClick={() => onNavigate('TODAY')}>
          Today
        </button>
      </span>
      <span className="rbc-toolbar-label-container flex items-center gap-2">
        <button type="button" className="rbc-toolbar-nav-btn" onClick={() => onNavigate('PREV')}>
          <ArrowLeft size={20} />
        </button>
        <span className="rbc-toolbar-label text-lg font-semibold">{label}</span>
        <button type="button" className="rbc-toolbar-nav-btn" onClick={() => onNavigate('NEXT')}>
          <ArrowRight size={20} />
        </button>
      </span>
      <span
        className="rbc-btn-group flex gap-2"
        style={showViewSelector ? {} : { opacity: 0, pointerEvents: 'none', userSelect: 'none', color: 'transparent', backgroundColor: 'transparent', borderColor: 'transparent', boxShadow: 'none' }}
      >
        <Button
          type="button"
          className={view === 'day' ? 'rbc-active' : ''}
          onClick={() => onView('day')}
        >
          Day
        </Button>
        <Button
          type="button"
          className={view === 'work_week' ? 'rbc-active' : ''}
          onClick={() => onView('work_week')}
        >
          Week
        </Button>
        <Button
          type="button"
          className={view === 'month' ? 'rbc-active' : ''}
          onClick={() => onView('month')}
        >
          Month
        </Button>
      </span>
    </div>
  )
}
