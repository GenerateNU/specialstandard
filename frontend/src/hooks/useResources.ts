import type {
  CreateResourceBody,
  Resource,
  UpdateResourceBody,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import {
  useMutation,
  useQuery,
  useQueryClient,
  type QueryObserverResult,
} from "@tanstack/react-query";
import { getResources as getResourcesApi } from "@/lib/api/resources";

interface UseResourcesReturn {
  resources: Resource[];
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<QueryObserverResult<Resource[], Error>>;
  addResource: (resource: CreateResourceBody) => void;
  updateResource: (id: string, updatedResource: UpdateResourceBody) => void;
  deleteResource: (id: string) => void;
}

export function useResources(): UseResourcesReturn {
  const queryClient = useQueryClient();
  const api = getResourcesApi();

  const {
    data: resources = [],
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["resources"],
    queryFn: () => api.getResources(),
  });

  const addResourceMutation = useMutation({
    mutationFn: (input: CreateResourceBody) => api.postResources(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["resources"] });
    },
  });

  const updateResourceMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateResourceBody }) =>
      api.patchResourcesId(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["resources"] });
    },
  });

  const deleteResourceMutation = useMutation({
    mutationFn: (id: string) => api.deleteResourcesId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["resources"] });
    },
  });

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
  };
}
