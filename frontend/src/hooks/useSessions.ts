// hooks/useSessions.ts

import { useAuthContext } from "@/contexts/authContext";
import { getSessions as getSessionsApi } from "@/lib/api/sessions";
import type {
  PostSessionsBody,
  Session,
  UpdateSessionInput,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import type { QueryObserverResult } from "@tanstack/react-query";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";

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
  const { userId: contextUserId } = useAuthContext();
  
  // Support onboarding flow by checking localStorage
  const [therapistId, setTherapistId] = useState<string | undefined>(contextUserId ?? undefined);

  useEffect(() => {
    // If we don't have a therapist ID from context, check localStorage (onboarding)
    if (!contextUserId && typeof window !== 'undefined') {
      const tempUserId = localStorage.getItem("temp_userId") || localStorage.getItem("userId");
      if (tempUserId) {
        setTherapistId(tempUserId);
      }
    } else if (contextUserId) {
      setTherapistId(contextUserId);
    }
  }, [contextUserId]);

  const {
    data: sessionsResponse,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["sessions", params, therapistId],
    queryFn: async () => {
      // Fetch all pages
      const allSessions: Session[] = [];
      let page = 1;
      const limit = params?.limit ?? 100;
      let hasMore = true;

      while (hasMore) {
        const pageData = await api.getSessions({
          limit,
          page,
          startdate: params?.startdate,
          enddate: params?.enddate,
          therapist_id: therapistId!,
        });

        // Check if pageData is valid before spreading
        if (!pageData || !Array.isArray(pageData)) {
          hasMore = false;
          break;
        }

        allSessions.push(...pageData);

        // If we got fewer results than the limit, we've reached the last page
        if (pageData.length < limit) {
          hasMore = false;
        } else {
          page++;
        }
      }

      return allSessions;
    },
    enabled: !!therapistId, // Only run when we have a valid therapist ID
  });

  const sessions = sessionsResponse ?? [];

  // Add Session Mutation
  const addSessionMutation = useMutation({
    mutationFn: (input: PostSessionsBody) => api.postSessions(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
      queryClient.invalidateQueries({
        queryKey: ["studentSessions"],
        exact: false,
      });
    },
  });

  // Update Session Mutation
  const updateSessionMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateSessionInput }) =>
      api.patchSessionsId(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
      queryClient.invalidateQueries({ queryKey: ["session"] });
      queryClient.invalidateQueries({
        queryKey: ["studentSessions"],
        exact: false,
      });
    },
  });

  // Delete Single Session Mutation
  const deleteSessionMutation = useMutation({
    mutationFn: (id: string) => api.deleteSessionsId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
      queryClient.invalidateQueries({ queryKey: ["session"] });
      queryClient.invalidateQueries({
        queryKey: ["studentSessions"],
        exact: false,
      });
    },
  });

  // Delete Recurring Sessions Group Mutation
  const deleteRecurringMutation = useMutation({
    mutationFn: (id: string) => api.deleteSessionsIdRecurring(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
      queryClient.invalidateQueries({ queryKey: ["session"] });
      queryClient.invalidateQueries({
        queryKey: ["studentSessions"],
        exact: false,
      });
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