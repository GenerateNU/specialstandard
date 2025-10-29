import type {
  CreateSessionStudentInput,
  DeleteSessionStudentsBody,
} from '@/lib/api/theSpecialStandardAPI.schemas'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { getSessionStudents as getSessionStudentsApi } from '@/lib/api/session-students'

export function useSessionStudents() {
  const queryClient = useQueryClient()
  const api = getSessionStudentsApi()

  const addStudentToSessionMutation = useMutation({
    mutationFn: (input: CreateSessionStudentInput) =>
      api.postSessionStudents(input),
    onSuccess: (_, variables) => {
      if (variables.session_ids) {
        variables.session_ids.forEach((id: string) => {
          queryClient.invalidQueries({
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

  return {
    addStudentToSession: (input: CreateSessionStudentInput) =>
      addStudentToSessionMutation.mutate(input),
    removeStudentFromSession: (input: DeleteSessionStudentsBody) =>
      removeStudentFromSessionMutation.mutate(input),
    isAdding: addStudentToSessionMutation.isPending,
    isRemoving: removeStudentFromSessionMutation.isPending,
    addError: addStudentToSessionMutation.error?.message || null,
    removeError: removeStudentFromSessionMutation.error?.message || null,
  }
}
