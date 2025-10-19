import React from 'react'
import { cn } from '@/lib/utils'
import { Input } from './input'

interface TimeRangeProps {
  value?: { startTime: string, endTime: string }
  onChange: (value: { startTime: string, endTime: string }) => void
  className?: string
}

const TimeRange = React.forwardRef<HTMLDivElement, TimeRangeProps>(
  ({ value, onChange, className }, ref) => {
    const handleStartTimeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const newStartTime = e.target.value
      const currentEndTime = value?.endTime || '10:00'

      // If new start time is after end time, adjust end time to be 1 hour after start
      if (newStartTime >= currentEndTime) {
        const [hours, minutes] = newStartTime.split(':').map(Number)
        const adjustedEndTime = `${String(hours + 1).padStart(2, '0')}:${String(minutes).padStart(2, '0')}`
        onChange({
          startTime: newStartTime,
          endTime: adjustedEndTime,
        })
      }
      else {
        onChange({
          startTime: newStartTime,
          endTime: currentEndTime,
        })
      }
    }

    const handleEndTimeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const newEndTime = e.target.value
      const currentStartTime = value?.startTime || '09:00'

      // Don't allow end time to be before or equal to start time
      if (newEndTime <= currentStartTime) {
        return
      }

      onChange({
        startTime: currentStartTime,
        endTime: newEndTime,
      })
    }

    return (
      <div ref={ref} className={cn('space-y-2', className)}>
        <div className="rounded-lg p-3 flex items-center gap-3 border border-gray-200">
          <div className="relative flex-1 w-min">
            <Input
              type="time"
              value={value?.startTime || '09:00'}
              onChange={handleStartTimeChange}
              className="text-white"
            />
          </div>

          <span className="text-gray-400">–</span>

          <div className="relative flex-1">
            <Input
              type="time"
              value={value?.endTime || '10:00'}
              onChange={handleEndTimeChange}
              // className="[&::-webkit-calendar-picker-indicator]:hidden"
              min={value?.startTime || '09:00'}
            />
          </div>
        </div>
      </div>
    )
  },
)

TimeRange.displayName = 'TimeRange'
export { TimeRange }
