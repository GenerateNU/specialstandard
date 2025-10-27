import { useAuthContext } from '@/contexts/authContext'
import { getTherapists as getTherapistsApi } from '@/lib/api/therapists'
import type {
  CreateTherapistInput,
  Therapist,
  UpdateTherapistInput,
} from '@/lib/api/theSpecialStandardAPI.schemas'
import type { QueryObserverResult } from '@tanstack/react-query'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

interface UseTherapistsReturn {
  therapists: Therapist[]
  isLoading: boolean
  error: string | null
  refetch: () => Promise<QueryObserverResult<Therapist[], Error>>
  addTherapist: (therapist: CreateTherapistInput) => void
  updateTherapist: (id: string, updatedTherapist: UpdateTherapistInput) => void
  deleteTherapist: (id: string) => void
}

export function useTherapists(): UseTherapistsReturn {
  const queryClient = useQueryClient()
  const api = getTherapistsApi()
  const { userId: therapistId } = useAuthContext()

  const {
    data: therapistsResponse,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['therapists', therapistId],
    queryFn: () => api.getTherapists({ therapistId: therapistId! }),
    // we technically dont need this line but it is just defensive programming!!  
    enabled: !!therapistId, 
  })

  const therapists = therapistsResponse ?? []

  const addTherapistMutation = useMutation({
    mutationFn: (input: CreateTherapistInput) => api.postTherapists(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['therapists', therapistId] })
    },
  })

  const updateTherapistMutation = useMutation({
    mutationFn: ({ id, data }: { id: string, data: UpdateTherapistInput }) =>
      api.patchTherapistsId(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['therapists', therapistId] })
    },
  })

  const deleteTherapistMutation = useMutation({
    mutationFn: (id: string) => api.deleteTherapistsId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['therapists', therapistId] })
    },
  })

  return {
    therapists,
    isLoading,
    error: error?.message || null,
    refetch,
    addTherapist: (therapist: CreateTherapistInput) =>
      addTherapistMutation.mutate(therapist),
    updateTherapist: (id: string, data: UpdateTherapistInput) =>
      updateTherapistMutation.mutate({ id, data }),
    deleteTherapist: (id: string) => deleteTherapistMutation.mutate(id),
  }
}
