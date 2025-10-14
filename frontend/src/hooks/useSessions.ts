import type { QueryObserverResult } from "@tanstack/react-query";
import type {
  Session,
  PostSessionsBody,
  UpdateSessionInput,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getSessions as getSessionsApi } from "@/lib/api/sessions";

interface UseSessionsReturn {
  sessions: Session[];
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<QueryObserverResult<Session[], Error>>;
  addSession: (session: PostSessionsBody) => void;
  updateSession: (id: string, updatedSession: UpdateSessionInput) => void;
  deleteSession: (id: string) => void;
}

export function useSessions(): UseSessionsReturn {
  const queryClient = useQueryClient();
  const api = getSessionsApi();

  const {
    data: sessionsResponse,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["sessions"],
    queryFn: () => api.getSessions(),
  });

  const sessions = sessionsResponse ?? [];

  const addSessionMutation = useMutation({
    mutationFn: (input: PostSessionsBody) => api.postSessions(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
    },
  });

  const updateSessionMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateSessionInput }) =>
      api.patchSessionsId(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["ssions"] });
    },
  });

  const deleteSessionMutation = useMutation({
    mutationFn: (id: string) => api.deleteSessionsId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
    },
  });

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
  };
}
