import { useAuthContext } from '@/contexts/authContext'
import { getSessions as getSessionsApi } from '@/lib/api/sessions'
import type {
  PostSessionsBody,
  Session,
  UpdateSessionInput,
} from '@/lib/api/theSpecialStandardAPI.schemas'
import type { QueryObserverResult } from '@tanstack/react-query'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

interface UseSessionsReturn {
  sessions: Session[]
  isLoading: boolean
  error: string | null
  refetch: () => Promise<QueryObserverResult<Session[], Error>>
  addSession: (session: PostSessionsBody) => void
  updateSession: (id: string, updatedSession: UpdateSessionInput) => void
  deleteSession: (id: string) => void
}

interface UseSessionsParams {
  startdate?: string
  enddate?: string
  limit?: number
}

export function useSessions(params?: UseSessionsParams): UseSessionsReturn {
  const queryClient = useQueryClient()
  const api = getSessionsApi()
  const { userId: therapistId } = useAuthContext()

  const {
    data: sessionsResponse,
    isLoading,
    error,
    refetch,
  } = useQuery({
<<<<<<< HEAD
    queryKey: ['sessions', therapistId],
    queryFn: () => api.getSessions(),
    // we technically dont need this line but it is just defensive programming!!  
    enabled: !!therapistId, 
=======
    queryKey: ['sessions', params],
    queryFn: () => api.getSessions({
      limit: params?.limit ?? 100,
      startdate: params?.startdate,
      enddate: params?.enddate,
    }),
>>>>>>> main
  })

  const sessions = sessionsResponse ?? []

  const addSessionMutation = useMutation({
    mutationFn: (input: PostSessionsBody) => api.postSessions(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sessions', therapistId] })
    },
  })

  const updateSessionMutation = useMutation({
    mutationFn: ({ id, data }: { id: string, data: UpdateSessionInput }) =>
      api.patchSessionsId(id, data),
    onSuccess: () => {
<<<<<<< HEAD
      queryClient.invalidateQueries({ queryKey: ['sessions', therapistId] })
=======
      queryClient.invalidateQueries({ queryKey: ['sessions'] })
>>>>>>> main
    },
  })

  const deleteSessionMutation = useMutation({
    mutationFn: (id: string) => api.deleteSessionsId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sessions', therapistId] })
    },
  })

  return {
    sessions,
    isLoading,
    error: error?.message || null,
    refetch,
    addSession: (session: PostSessionsBody) =>
      addSessionMutation.mutate(session),
    updateSession: (id: string, data: UpdateSessionInput) =>
      updateSessionMutation.mutate({ id, data }),
    deleteSession: (id: string) => deleteSessionMutation.mutate(id),
  }
}
