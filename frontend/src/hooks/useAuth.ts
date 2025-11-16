import type {
  PostAuthLoginBody,
  PostAuthSignupBody,
} from '@/lib/api/theSpecialStandardAPI.schemas'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { getAuth as getAuthApi } from '@/lib/api/auth'

export function useAuth() {
  const api = getAuthApi()

  const queryClient = useQueryClient()

  const loginMutation = useMutation({
    mutationFn: (credentials: PostAuthLoginBody) =>
      api.postAuthLogin(credentials),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user'] })
    },
  })

  const signupMutation = useMutation({
    mutationFn: (credentials: PostAuthSignupBody) =>
      api.postAuthSignup(credentials),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user'] })
    },
  })

  const logoutMutation = useMutation({
    mutationFn: async () => {
      // Clear localStorage (important for your auth setup!)
      localStorage.removeItem('jwt')
      localStorage.removeItem('userId')

      // This can be an API call if you have a logout endpoint
      // For cookie-based auth, clearing cookies client-side is common
      document.cookie.split(';').forEach((cookie) => {
        const eqPos = cookie.indexOf('=')
        const name
          = eqPos > -1 ? cookie.substring(0, eqPos).trim() : cookie.trim()
        document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`
      })
    },
    onSuccess: () => {
      queryClient.setQueryData(['user'], null)
      queryClient.invalidateQueries({ queryKey: ['user'] })
    },
  })

  return {
    userLogin: loginMutation.mutateAsync,
    userLogout: logoutMutation.mutateAsync,
    userSignup: signupMutation.mutateAsync,
    loginMutation,
    logoutMutation,
    signupMutation,
  }
}
