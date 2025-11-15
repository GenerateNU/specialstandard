import { useQuery } from '@tanstack/react-query'
import { getResources } from '@/lib/api/resources'
import type { GetResourcesParams } from '@/lib/api/theSpecialStandardAPI.schemas'

export interface UseResourcesReturn {
  resources: any[]
  isLoading: boolean
  error: string | null
  refetch: () => void
}

export function useResources(params?: GetResourcesParams): UseResourcesReturn {
  const api = getResources()
  
  const {
    data: resources,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['resources', params],
    queryFn: async () => {
      const response = await api.getResources(params)
      return response
    },
    enabled: !!(params && Object.keys(params).length > 0),
  })
  
  return {
    resources: resources ?? [],
    isLoading,
    error: error?.message || null,
    refetch,
  }
}
