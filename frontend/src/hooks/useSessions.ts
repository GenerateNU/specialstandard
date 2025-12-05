// hooks/useSessions.ts - FULLY UPDATED WITH REPETITION SUPPORT & AUTO-REFETCH

import { useAuthContext } from "@/contexts/authContext";
import { getSessions as getSessionsApi } from "@/lib/api/sessions";
import type {
  PostSessionsBody,
  Session,
  UpdateSessionInput,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import type { QueryObserverResult } from "@tanstack/react-query";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

interface UseSessionsReturn {
  sessions: Session[];
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<QueryObserverResult<Session[], Error>>;
  addSession: (session: PostSessionsBody) => Promise<Session[]>;
  updateSession: (id: string, updatedSession: UpdateSessionInput) => Promise<Session>;
  deleteSession: (id: string) => Promise<void>;
  deleteRecurringSessions: (id: string) => Promise<void>;
  isAddingSession: boolean;
  isDeletingSession: boolean;
  isDeletingRecurring: boolean;
  addSessionError: string | null;
  deleteError: string | null;
}

interface UseSessionsParams {
  startdate?: string;
  enddate?: string;
  limit?: number;
}

interface UseSessionReturn {
  session: Session | null;
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<QueryObserverResult<Session, Error>>;
  isRecurring: boolean;
}

export function useSessions(params?: UseSessionsParams): UseSessionsReturn {
  const queryClient = useQueryClient();
  const api = getSessionsApi();
  const { userId: therapistId } = useAuthContext();

  const {
    data: sessionsResponse,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["sessions", params],
    queryFn: () =>
      api.getSessions({
        limit: params?.limit ?? 100,
        startdate: params?.startdate,
        enddate: params?.enddate,
        therapist_id: therapistId!,
      }),
    enabled: !!therapistId,
  });

  const sessions = sessionsResponse ?? [];

  // Add Session Mutation
  const addSessionMutation = useMutation({
    mutationFn: (input: PostSessionsBody) => api.postSessions(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
    },
  });

  // Update Session Mutation
  const updateSessionMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateSessionInput }) =>
      api.patchSessionsId(id, data),
    onSuccess: () => {
      // Invalidate both query keys so both useSessions() and useSession(id) refetch
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
      queryClient.invalidateQueries({ queryKey: ["session"] });
    },
  });

  // Delete Single Session Mutation
  const deleteSessionMutation = useMutation({
    mutationFn: (id: string) => api.deleteSessionsId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
      queryClient.invalidateQueries({ queryKey: ["session"] });
    },
  });

  // Delete Recurring Sessions Group Mutation
  const deleteRecurringMutation = useMutation({
    mutationFn: (id: string) => api.deleteSessionsIdRecurring(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
      queryClient.invalidateQueries({ queryKey: ["session"] });
    },
  });

  return {
    sessions,
    isLoading,
    error: error?.message || null,
    refetch,
    addSession: async (session: PostSessionsBody) => {
      return new Promise((resolve, reject) => {
        addSessionMutation.mutate(session, {
          onSuccess: (data) => resolve(data),
          onError: (err) => reject(err),
        });
      });
    },
    updateSession: async (id: string, data: UpdateSessionInput) => {
      return new Promise((resolve, reject) => {
        updateSessionMutation.mutate({ 
          id, 
          data: { ...data, therapist_id: therapistId ?? undefined } 
        }, {
          onSuccess: (data) => resolve(data),
          onError: (err) => reject(err),
        });
      });
    },
    deleteSession: async (id: string) => {
      return new Promise((resolve, reject) => {
        deleteSessionMutation.mutate(id, {
          onSuccess: () => resolve(),
          onError: (err) => reject(err),
        });
      });
    },
    deleteRecurringSessions: async (id: string) => {
      return new Promise((resolve, reject) => {
        deleteRecurringMutation.mutate(id, {
          onSuccess: () => resolve(),
          onError: (err) => reject(err),
        });
      });
    },
    isAddingSession: addSessionMutation.isPending,
    isDeletingSession: deleteSessionMutation.isPending,
    isDeletingRecurring: deleteRecurringMutation.isPending,
    addSessionError: addSessionMutation.error?.message || null,
    deleteError:
      deleteSessionMutation.error?.message ||
      deleteRecurringMutation.error?.message ||
      null,
  };
}

export function useSession(id: string): UseSessionReturn {
  const api = getSessionsApi();
  const {
    data: session,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["session", id],
    queryFn: () => api.getSessionsId(id),
    enabled: !!id,
  });

  const isRecurring = !!session?.repetition;

  return {
    session: session || null,
    isLoading,
    error: error?.message || null,
    refetch,
    isRecurring,
  };
}

// ============================================
// UTILITY FUNCTION - Format recurrence text
// ============================================
export function formatRecurrence(repetition: any): string {
  if (!repetition) return "Single session";

  const DAYS_OF_WEEK = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
  const dayNames = repetition.days
    .map((d: number) => DAYS_OF_WEEK[d])
    .join(", ");

  const endDate = new Date(repetition.recur_end).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });

  return `Every ${repetition.every_n_weeks} week${repetition.every_n_weeks !== 1 ? "s" : ""} on ${dayNames} until ${endDate}`;
}