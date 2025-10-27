import { useAuthContext } from '@/contexts/authContext'
import { getResources as getResourcesApi } from '@/lib/api/resources'
import type {
  CreateResourceBody,
  ResourceWithPresignURL,
  UpdateResourceBody,
} from '@/lib/api/theSpecialStandardAPI.schemas'
import type { QueryObserverResult } from '@tanstack/react-query'
import {
  useMutation,
  useQuery,
  useQueryClient,
} from '@tanstack/react-query'

interface UseResourcesReturn {
  resources: ResourceWithPresignURL[]
  isLoading: boolean
  error: string | null
  refetch: () => Promise<QueryObserverResult<ResourceWithPresignURL[], Error>>
  addResource: (resource: CreateResourceBody) => void
  updateResource: (id: string, updatedResource: UpdateResourceBody) => void
  deleteResource: (id: string) => void
}

export function useResources(): UseResourcesReturn {
  const queryClient = useQueryClient()
  const api = getResourcesApi()
  const { userId: therapistId } = useAuthContext()

  const {
    data: resources = [],
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['resources', therapistId],
    queryFn: () => api.getResources(),
  })

  const addResourceMutation = useMutation({
    mutationFn: (input: CreateResourceBody) => api.postResources(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['resources', therapistId] })
    },
  })

  const updateResourceMutation = useMutation({
    mutationFn: ({ id, data }: { id: string, data: UpdateResourceBody }) =>
      api.patchResourcesId(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['resources', therapistId] })
    },
  })

  const deleteResourceMutation = useMutation({
    mutationFn: (id: string) => api.deleteResourcesId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['resources', therapistId] })
    },
  })

  return {
    resources,
    isLoading,
    error: error?.message || null,
    refetch,
    addResource: (resource: CreateResourceBody) =>
      addResourceMutation.mutate(resource),
    updateResource: (id: string, data: UpdateResourceBody) =>
      updateResourceMutation.mutate({ id, data }),
    deleteResource: (id: string) => deleteResourceMutation.mutate(id),
  }
}
