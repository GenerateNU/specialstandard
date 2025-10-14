import type { ModifySessionResource } from '@/lib/api/theSpecialStandardAPI.schemas'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { getSessionResource as getSessionResourceApi } from '@/lib/api/session-resource'

export function useSessionResource() {
  const queryClient = useQueryClient()
  const api = getSessionResourceApi()

  const addResourceToSessionMutation = useMutation({
    mutationFn: (input: ModifySessionResource) =>
      api.postSessionResource(input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['sessions', variables.session_id, 'resources'],
      })
      queryClient.invalidateQueries({ queryKey: ['sessions'] })
    },
  })

  const removeResourceFromSessionMutation = useMutation({
    mutationFn: (input: ModifySessionResource) =>
      api.deleteSessionResource(input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['sessions', variables.session_id, 'resources'],
      })
      queryClient.invalidateQueries({ queryKey: ['sessions'] })
    },
  })

  return {
    addResourceToSession: (input: ModifySessionResource) =>
      addResourceToSessionMutation.mutate(input),
    removeResourceFromSession: (input: ModifySessionResource) =>
      removeResourceFromSessionMutation.mutate(input),
    isAdding: addResourceToSessionMutation.isPending,
    isRemoving: removeResourceFromSessionMutation.isPending,
    addError: addResourceToSessionMutation.error?.message || null,
    removeError: removeResourceFromSessionMutation.error?.message || null,
  }
}
