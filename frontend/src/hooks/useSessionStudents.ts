import type {
  CreateSessionStudentInput,
  DeleteSessionStudentsBody,
  StudentWithSessionInfo,
  UpdateSessionStudentInput,
} from '@/lib/api/theSpecialStandardAPI.schemas'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { getSessionStudents as getSessionStudentsApi } from '@/lib/api/session-students'
import { getSessions } from '@/lib/api/sessions'
import { gradeToDisplay } from '@/lib/gradeUtils'

export type SessionStudentBody = Omit<StudentWithSessionInfo, 'grade'> & {
  grade: string
}

export function useSessionStudentsForSession(sessionId: string) {
  const sessionsApi = getSessions()

  const {
    data: studentsData,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['sessions', sessionId, 'students'],
    queryFn: () => sessionsApi.getSessionsSessionIdStudents(sessionId),
    enabled: !!sessionId,
  })

  // Flatten the nested student structure
  const students = (studentsData || []).map((item: any) => ({
    // Spread the nested student object first to get id, first_name, last_name, etc.
    ...(item.student || {}),
    // Then add session-specific fields
    session_id: item.session_id,
    present: item.present,
    notes: item.notes,
    // Override grade with display format
    grade: gradeToDisplay(item.student?.grade ?? item.grade),
  }))

  return {
    students,
    isLoading,
    error: error?.message || null,
    refetch,
  }
}

export function useSessionStudents() {
  const queryClient = useQueryClient()
  const api = getSessionStudentsApi()

  const addStudentToSessionMutation = useMutation({
    mutationFn: (input: CreateSessionStudentInput) =>
      api.postSessionStudents({
        // Backend expects arrays, wrap single values
        session_ids: [input.session_id],
        student_ids: [input.student_id],
        present: input.present,
        notes: input.notes,
      } as any),
    onSuccess: (_, variables) => {
      if (variables.session_ids) {
        variables.session_ids.forEach((id: string) => {
          queryClient.invalidateQueries({
            queryKey: ['sessions', id, 'students'],
          })
        })
      }

      queryClient.invalidateQueries({ queryKey: ['sessions'] })
    },
  })

  const removeStudentFromSessionMutation = useMutation({
    mutationFn: (input: DeleteSessionStudentsBody) =>
      api.deleteSessionStudents(input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['sessions', variables.session_id, 'students'],
      })
      queryClient.invalidateQueries({ queryKey: ['sessions'] })
    },
  })

  const updateSessionStudentMutation = useMutation({
    mutationFn: (input: UpdateSessionStudentInput) =>
      api.patchSessionStudents(input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['sessions', variables.session_id, 'students'],
      })
    },
  })

  return {
    addStudentToSession: (input: CreateSessionStudentInput) =>
      addStudentToSessionMutation.mutate(input),
    removeStudentFromSession: (input: DeleteSessionStudentsBody) =>
      removeStudentFromSessionMutation.mutate(input),
    updateSessionStudent: (input: UpdateSessionStudentInput) =>
      updateSessionStudentMutation.mutate(input),
    isAdding: addStudentToSessionMutation.isPending,
    isRemoving: removeStudentFromSessionMutation.isPending,
    isUpdating: updateSessionStudentMutation.isPending,
    addError: addStudentToSessionMutation.error?.message || null,
    removeError: removeStudentFromSessionMutation.error?.message || null,
    updateError: updateSessionStudentMutation.error?.message || null,
  }
}
