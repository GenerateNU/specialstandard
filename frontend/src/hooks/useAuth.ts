import type {
    PostAuthForgotPasswordBody,
    PostAuthLoginBody,
    PostAuthSignupBody,
    PutAuthUpdatePasswordBody,
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
            // Clear localStorage
            localStorage.removeItem('jwt')
            localStorage.removeItem('userId')
            localStorage.removeItem('recentlyViewedStudents')     
            // Clear cookies
            document.cookie.split(';').forEach((cookie) => {
                const eqPos = cookie.indexOf('=')
                const name =
                    eqPos > -1 ? cookie.substring(0, eqPos).trim() : cookie.trim()
                document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`
            })
        },
        onSuccess: () => {
            queryClient.setQueryData(['user'], null)
            queryClient.invalidateQueries({ queryKey: ['user'] })
        },
    })

    const forgotPasswordMutation = useMutation({
        mutationFn: (body: PostAuthForgotPasswordBody) =>
            api.postAuthForgotPassword(body),
    })

    const updatePasswordMutation = useMutation({
        mutationFn: (data: { password: string; token?: string }) => {
            // Extract token and password
            const { token, password } = data
            const body: PutAuthUpdatePasswordBody = { password }
            const params = { token }

            // Call API with both body and params
            return api.putAuthUpdatePassword(body, params)
        },
        onSuccess: () => {
            // Clear auth data after successful password reset
            localStorage.removeItem('jwt')
            localStorage.removeItem('userId')
            localStorage.removeItem('recentlyViewedStudents')     
            queryClient.setQueryData(['user'], null)
            queryClient.invalidateQueries({ queryKey: ['user'] })
        },
    })

    const deleteAccountMutation = useMutation({
        mutationFn: (id: string) => api.deleteAuthDeleteAccountId(id),
        onSuccess: () => {
            localStorage.removeItem('jwt')
            localStorage.removeItem('userId')
            localStorage.removeItem('recentlyViewedStudents')   
            localStorage.removeItem("signup_requires_mfa");  

            queryClient.setQueryData(['user'], null)
            queryClient.invalidateQueries({ queryKey: ['user'] })
        },
    })

    return {
        userLogin: loginMutation.mutateAsync,
        userLogout: logoutMutation.mutateAsync,
        userSignup: signupMutation.mutateAsync,
        updatePassword: updatePasswordMutation.mutateAsync,
        forgotPassword: forgotPasswordMutation.mutateAsync,
        deleteAccount: deleteAccountMutation.mutateAsync,
        // Return mutation objects for isPending, isError, error states
        loginMutation,
        logoutMutation,
        signupMutation,
        updatePasswordMutation,
        forgotPasswordMutation,
        deleteAccountMutation,
    }
}
