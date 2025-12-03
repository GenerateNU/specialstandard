import { useMutation } from '@tanstack/react-query'
import { getGameResult } from '@/lib/api/game-result'
import type { 
  PostGameResultInput 
} from '@/lib/api/theSpecialStandardAPI.schemas'

export function useManualGameResult() {
  const api = getGameResult();
  
  const mutation = useMutation({
    mutationFn: async (input: PostGameResultInput) => {
      return await api.postGameResults(input);
    },
  });

  return {
    submitResult: mutation.mutate,
    submitResultAsync: mutation.mutateAsync,
    isSubmitting: mutation.isPending,
    error: mutation.error,
  };
}
