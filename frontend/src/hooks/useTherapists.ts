import { getTherapists as getTherapistsApi } from "@/lib/api/therapists";
import type {
  CreateTherapistInput,
  Therapist,
  UpdateTherapistInput,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import type { QueryObserverResult } from "@tanstack/react-query";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

interface UseTherapistsReturn {
  therapists: Therapist[];
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<QueryObserverResult<Therapist[], Error>>;
  addTherapist: (therapist: CreateTherapistInput) => Promise<Therapist>;
  updateTherapist: (id: string, updatedTherapist: UpdateTherapistInput) => Promise<Therapist>;
  deleteTherapist: (id: string) => Promise<void>;
  isAddingTherapist: boolean;
  isUpdatingTherapist: boolean;
  isDeletingTherapist: boolean;
  addError: string | null;
  updateError: string | null;
  deleteError: string | null;
}

interface UseTherapistsOptions {
  fetchOnMount: boolean;
}

export function useTherapists(
  _options: UseTherapistsOptions = { fetchOnMount: true }
): UseTherapistsReturn {
  const queryClient = useQueryClient();
  const api = getTherapistsApi();

  const {
    data: therapistsResponse,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["therapists"],
    queryFn: () => api.getTherapists({}),
    enabled: _options.fetchOnMount,
  });

  const therapists = therapistsResponse ?? [];

  const addTherapistMutation = useMutation({
    mutationFn: (input: CreateTherapistInput) => api.postTherapists(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["therapists"] });
      queryClient.invalidateQueries({ queryKey: ["therapist"] });
    },
  });

  const updateTherapistMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTherapistInput }) =>
      api.patchTherapistsId(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["therapists"] });
      queryClient.invalidateQueries({ queryKey: ["therapist"] });
    },
  });

  const deleteTherapistMutation = useMutation({
    mutationFn: (id: string) => api.deleteTherapistsId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["therapists"] });
      queryClient.invalidateQueries({ queryKey: ["therapist"] });
    },
  });

  return {
    therapists,
    isLoading,
    error: error?.message || null,
    refetch,
    addTherapist: async (therapist: CreateTherapistInput) => {
      return new Promise((resolve, reject) => {
        addTherapistMutation.mutate(therapist, {
          onSuccess: (data) => resolve(data),
          onError: (err) => reject(err),
        });
      });
    },
    updateTherapist: async (id: string, data: UpdateTherapistInput) => {
      return new Promise((resolve, reject) => {
        updateTherapistMutation.mutate({ id, data }, {
          onSuccess: (data) => resolve(data),
          onError: (err) => reject(err),
        });
      });
    },
    deleteTherapist: async (id: string) => {
      return new Promise((resolve, reject) => {
        deleteTherapistMutation.mutate(id, {
          onSuccess: () => resolve(),
          onError: (err) => reject(err),
        });
      });
    },
    isAddingTherapist: addTherapistMutation.isPending,
    isUpdatingTherapist: updateTherapistMutation.isPending,
    isDeletingTherapist: deleteTherapistMutation.isPending,
    addError: addTherapistMutation.error?.message || null,
    updateError: updateTherapistMutation.error?.message || null,
    deleteError: deleteTherapistMutation.error?.message || null,
  };
}

// Hook for fetching a single therapist
interface UseTherapistReturn {
  therapist: Therapist | null;
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<QueryObserverResult<Therapist, Error>>;
}

export function useTherapist(therapistId: string | null): UseTherapistReturn {
  const api = getTherapistsApi();
  const {
    data: therapist,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["therapist", therapistId],
    queryFn: () => api.getTherapistsId(therapistId!),
    enabled: !!therapistId,
  });

  return {
    therapist: therapist ?? null,
    isLoading,
    error: error?.message || null,
    refetch,
  };
}