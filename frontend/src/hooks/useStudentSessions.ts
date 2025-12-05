import type { QueryObserverResult } from '@tanstack/react-query'
import type { SessionWithStudentInfo } from '@/lib/api/theSpecialStandardAPI.schemas'
import { useQuery } from '@tanstack/react-query'
import { getStudents } from '@/lib/api/students'

interface UseStudentSessionsReturn {
  sessions: SessionWithStudentInfo[]
  isLoading: boolean
  error: string | null
  refetch: () => Promise<QueryObserverResult<SessionWithStudentInfo[], Error>>
}

export function useStudentSessions(studentId: string, params?: Record<string, unknown>): UseStudentSessionsReturn {
  const api = getStudents()
  
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['studentSessions', studentId, ...(params ? Object.values(params) : [])],
    queryFn: () => api.getStudentsStudentIdSessions(studentId, params),
    enabled: Boolean(studentId),
  })

  return {
    sessions: data ?? [],
    isLoading,
    error: error?.message || null,
    refetch,
  }
}