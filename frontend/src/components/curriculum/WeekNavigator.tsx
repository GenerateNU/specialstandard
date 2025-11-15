'use client'

import { ChevronLeft, ChevronRight } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface WeekNavigatorProps {
  currentWeek: number
  onPreviousWeek: () => void
  onNextWeek: () => void
  className?: string
}

export default function WeekNavigator({
  currentWeek,
  onPreviousWeek,
  onNextWeek,
  className,
}: WeekNavigatorProps) {
  return (
    <div className={`flex items-center gap-4 ${className}`}>
      <Button
        onClick={onPreviousWeek}
        variant="ghost"
        size="icon"
        className="w-10 h-10 rounded-full hover:bg-blue-light"
      >
        <ChevronLeft className="w-6 h-6" />
      </Button>
      
      <div className="bg-card px-8 py-3 rounded-full shadow-md border border-default">
        <span className="text-lg font-semibold text-primary">
          Week
          {' '}
          {currentWeek}
        </span>
      </div>

      <Button
        onClick={onNextWeek}
        variant="ghost"
        size="icon"
        className="w-10 h-10 rounded-full hover:bg-blue-light"
      >
        <ChevronRight className="w-6 h-6" />
      </Button>
    </div>
  )
}

