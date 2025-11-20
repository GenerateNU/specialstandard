import type { QueryObserverResult } from '@tanstack/react-query'
import { useQuery } from '@tanstack/react-query'
import { getStudents } from '@/lib/api/students'
import type { AttendanceRecord } from '@/lib/api/theSpecialStandardAPI.schemas'

interface UseStudentAttendanceReturn {
  attendance: AttendanceRecord | null
  isLoading: boolean
  error: string | null
  refetch: () => Promise<QueryObserverResult<AttendanceRecord, Error>>
}
interface UseStudentAttendanceParams {
  studentId: string
  dateFrom?: string // YYYY-MM-DD format
  dateTo?: string   // YYYY-MM-DD format
}


export function useStudentAttendance({ 
  studentId, 
  dateFrom, 
  dateTo 
}: UseStudentAttendanceParams): UseStudentAttendanceReturn {
  const api = getStudents()

  const params: Record<string, unknown> = {}
  if (dateFrom) params.date_from = dateFrom
  if (dateTo) params.date_to = dateTo

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['studentAttendance', studentId, params],
    queryFn: () => api.getStudentsStudentIdAttendance(studentId, params),
    enabled: Boolean(studentId),
  })

  return {
    attendance: data ?? null,
    isLoading,
    error: error?.message || null,
    refetch,
  }
}
