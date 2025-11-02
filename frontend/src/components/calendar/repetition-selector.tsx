import { Calendar, Repeat } from 'lucide-react'
import React from 'react'
import {
  Card,
  CardContent,
  CardHeader,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { cn } from '@/lib/utils'
import { FormLabel } from '../ui/form'

export interface RepetitionConfig {
  enabled: boolean
  recur_start: string
  recur_end: string
  every_n_weeks: number
}

interface RepetitionFieldProps {
  value?: Partial<RepetitionConfig>
  onChange: (value: RepetitionConfig | undefined) => void
  sessionDate: string
  sessionTime?: string
}

export function RepetitionSelector({
  value,
  onChange,
  sessionDate,
  sessionTime = '09:00',
}: RepetitionFieldProps) {
  const [isEnabled, setIsEnabled] = React.useState(value?.enabled || false)
  const [recurStart, setRecurStart] = React.useState(value?.recur_start || sessionDate)
  const [recurEnd, setRecurEnd] = React.useState(value?.recur_end || '')
  const [everyNWeeks, setEveryNWeeks] = React.useState(value?.every_n_weeks || 1)

  React.useEffect(() => {
    // Update recur_start when sessionDate changes
    if (!isEnabled) {
      setRecurStart(sessionDate)
    }
  }, [sessionDate, isEnabled])

  React.useEffect(() => {
    if (isEnabled) {
      // Combine date and time for datetime values
      const startDateTime = new Date(`${recurStart}T${sessionTime}:00`)
      const endDateTime = recurEnd ? new Date(`${recurEnd}T${sessionTime}:00`) : null
      if (endDateTime && startDateTime <= endDateTime) {
        onChange({
          enabled: true,
          recur_start: startDateTime.toISOString(),
          recur_end: endDateTime.toISOString(),
          every_n_weeks: everyNWeeks,
        })
      }
    }
    else {
      onChange(undefined)
    }
  }, [isEnabled, recurStart, recurEnd, everyNWeeks, sessionTime, onChange])

  const handleToggle = (checked: boolean) => {
    setIsEnabled(checked)
    if (!checked) {
      // Reset values when disabled
      setRecurStart(sessionDate)
      setRecurEnd('')
      setEveryNWeeks(1)
    }
  }

  // Calculate default end date (3 months from start)
  const getDefaultEndDate = () => {
    const date = new Date(recurStart || sessionDate)
    date.setMonth(date.getMonth() + 3)
    return date.toISOString().split('T')[0]
  }

  // Calculate number of sessions
  const calculateSessionCount = () => {
    if (!isEnabled || !recurEnd || !recurStart) {
      return 0
    }
    const start = new Date(recurStart)
    const end = new Date(recurEnd)
    const diffTime = Math.abs(end.getTime() - start.getTime())
    const diffWeeks = Math.ceil(diffTime / (1000 * 60 * 60 * 24 * 7))
    return Math.floor(diffWeeks / everyNWeeks) + 1
  }

  return (
    <Card className="border-none bg-transparent shadow-none">
      <CardHeader className="pb-4 -mt-3 px-0">
        <div className="flex items-center justify-between">
          <FormLabel className="flex items-center gap-3 ">
            <>
              <Repeat className="w-4 h-4" />
              Repetition
            </>
            <Switch
              checked={isEnabled}
              onCheckedChange={handleToggle}
              aria-label="Enable recurring sessions"
              className={cn(isEnabled ? 'bg-accent' : 'bg-border', 'transition-colors')}
            />
          </FormLabel>
        </div>
      </CardHeader>
      {isEnabled && (
        <CardContent className="space-y-4 pt-0 border border-border rounded-md">
          {/* Frequency Selector */}
          <div className="space-y-2 rounded-md pt-3">
            <Label htmlFor="frequency" className="text-sm font-medium">
              Frequency
            </Label>
            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">Repeat every</span>
              <Input
                type="number"
                min="1"
                max="52"
                value={everyNWeeks}
                onChange={(e) => {
                  const val = Number.parseInt(e.target.value)
                  if (!Number.isNaN(val) && val > 0) {
                    setEveryNWeeks(val)
                  }
                }}
                className="w-min max-w-[4rem]"
              />
              <span className="text-sm text-muted-foreground">
                {everyNWeeks === 1 ? 'week' : 'weeks'}
              </span>
            </div>
          </div>

          {/* Date Range */}
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-2">
              <Label htmlFor="recur-start" className="text-sm font-medium flex items-center gap-1">
                <Calendar className="w-3 h-3" />
                Starts on
              </Label>
              <Input
                id="recur-start"
                type="date"
                value={recurStart}
                onChange={e => setRecurStart(e.target.value)}
                min={sessionDate}
                className="text-sm"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="recur-end" className="text-sm font-medium flex items-center gap-1">
                <Calendar className="w-3 h-3" />
                Ends on
              </Label>
              <Input
                id="recur-end"
                type="date"
                value={recurEnd}
                onChange={e => setRecurEnd(e.target.value)}
                min={recurStart || sessionDate}
                placeholder={getDefaultEndDate()}
                className="text-sm"
              />
            </div>
          </div>

          {/* Session Count Preview */}
          {recurEnd && (
            <div className="bg-muted/50 rounded-md p-3">
              <p className="text-sm text-muted-foreground">
                This will create approximately
                {' '}
                <span className="font-semibold text-foreground">
                  {calculateSessionCount()}
                  {' '}
                  sessions
                </span>
                {everyNWeeks === 1
                  ? ' (weekly)'
                  : ` (every ${everyNWeeks} weeks)`}
              </p>
            </div>
          )}
        </CardContent>
      )}
    </Card>
  )
}
