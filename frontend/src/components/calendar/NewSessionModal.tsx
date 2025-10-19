import type { StudentBody } from '@/hooks/useStudents'
import type { PostSessionsBody } from '@/lib/api/theSpecialStandardAPI.schemas'
import { Calendar, Clock, FileText, User } from 'lucide-react'
import React, { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { TimeRange } from '@/components/ui/time-range'
import { MultiSelect } from '../ui/multiselect'

// Main Component
interface CreateSessionDialogProps {
  open: boolean
  setOpen: (open: boolean) => void
  therapistId: string
  students?: Array<StudentBody>
  onSubmit?: (data: PostSessionsBody) => Promise<void>
  initialDateTime?: { start: Date, end: Date }
}

export function CreateSessionDialog({
  open,
  setOpen,
  therapistId,
  students = [],
  onSubmit,
  initialDateTime,
}: CreateSessionDialogProps) {
  const form = useForm<{
    student_ids: string[]
    sessionDate: string
    timeRange: [string, string]
    duration: number
    notes?: string
  }>({
    defaultValues: {
      student_ids: [],
      sessionDate: initialDateTime?.start.toISOString().split('T')[0] || new Date().toISOString().split('T')[0], // Today's date
      timeRange: [initialDateTime?.start.toTimeString().split(' ')[0].slice(0, 5) || '09:00', initialDateTime?.end.toTimeString().split(' ')[0].slice(0, 5) || '10:00'],
      duration: 60,
      notes: '',
    },
  })

  const watchTimeRange = form.watch('timeRange')

  useEffect(() => {
    if (watchTimeRange) {
      const [start, end] = watchTimeRange
      const [startHour, startMin] = start.split(':').map(Number)
      const [endHour, endMin] = end.split(':').map(Number)
      const duration = (endHour * 60 + endMin) - (startHour * 60 + startMin)
      if (duration > 0) {
        form.setValue('duration', duration)
      }
    }
  }, [watchTimeRange, form])

  useEffect(() => {
    if (open && initialDateTime) {
      const startTime = initialDateTime.start.toTimeString().split(' ')[0].slice(0, 5)
      const endTime = initialDateTime.end.toTimeString().split(' ')[0].slice(0, 5)
      form.reset({
        student_ids: [],
        sessionDate: initialDateTime.start.toISOString().split('T')[0],
        timeRange: [startTime, endTime],
        duration: 60,
        notes: '',
      })
    }
    else if (open && !initialDateTime) {
      form.reset({
        student_ids: [],
        sessionDate: new Date().toISOString().split('T')[0],
        timeRange: ['09:00', '10:00'],
        duration: 60,
        notes: '',
      })
    }
  }, [open, initialDateTime, form])

  const handleSubmit = async (data: any) => {
    try {
      const sessionDate = new Date(data.sessionDate)

      const [startHour, startMin] = data.startTime.split(':')
      const startDateTime = new Date(sessionDate)
      startDateTime.setHours(Number.parseInt(startHour), Number.parseInt(startMin), 0, 0)

      const [endHour, endMin] = data.endTime.split(':')
      const endDateTime = new Date(sessionDate)
      endDateTime.setHours(Number.parseInt(endHour), Number.parseInt(endMin), 0, 0)

      const postBody: PostSessionsBody = {
        start_datetime: startDateTime.toISOString(),
        end_datetime: endDateTime.toISOString(),
        therapist_id: therapistId,
        notes: data.notes || undefined,
      }

      if (onSubmit) {
        await onSubmit(postBody)
      }

      // Reset form and close dialog on success
      form.reset()
      setOpen(false)
    }
    catch (error) {
      console.error('Error creating session:', error)
      // Handle error - you might want to show a toast notification here
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Create new session</DialogTitle>
          <DialogDescription>
            Schedule a new therapy session by filling out the details below.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            {/* Date Field */}
            <FormField
              control={form.control}
              name="sessionDate"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="flex items-center gap-1">
                    <Calendar className="w-4 h-4" />
                    Date
                  </FormLabel>
                  <FormControl>
                    <Input
                      type="date"
                      {...field}
                      value={field.value || ''}
                      className="[&::-webkit-calendar-picker-indicator]:hidden"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Time Range */}
            <FormField
              control={form.control}
              name="timeRange"
              render={() => (
                <FormItem>
                  <FormLabel className="flex items-center gap-1">
                    <Clock className="w-4 h-4" />
                    Time
                  </FormLabel>
                  <FormControl>
                    <TimeRange
                      value={{ startTime: form.watch('timeRange')[0], endTime: form.watch('timeRange')[1] }}
                      onChange={(value) => {
                        form.setValue('timeRange', [value.startTime, value.endTime])
                      }}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            {/* Session Notes */}
            <FormField
              control={form.control}
              name="notes"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="flex items-center gap-1">
                    <FileText className="w-4 h-4" />
                    Session Notes
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Goals, activities planned, or any special considerations for this session..."
                      rows={3}
                      {...field}
                      value={field.value ?? ''}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="student_ids"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="flex items-center gap-1">
                    <User className="w-4 h-4" />
                    Students
                  </FormLabel>
                  <FormControl>
                    <MultiSelect
                      options={students.map(student => ({
                        label: `${student.first_name} ${student.last_name}`,
                        value: student.id,
                      }))}
                      value={field.value}
                      onValueChange={field.onChange}
                      placeholder="Select students"
                      showTags={true}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter className="gap-2 sm:gap-0">
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  form.reset()
                  setOpen(false)
                }}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={form.formState.isSubmitting}>
                {form.formState.isSubmitting ? 'Creating...' : 'Create Session'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
