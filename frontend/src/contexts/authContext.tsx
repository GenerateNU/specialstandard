'use client'
import { supabase } from '@/components/MFA/EnrollMFA'
import { useAuth } from '@/hooks/useAuth'
import type {
  PostAuthLoginBody,
  PostAuthSignupBody,
} from '@/lib/api/theSpecialStandardAPI.schemas'
import { useRouter } from 'next/navigation'
import { createContext, useContext, useEffect, useState } from 'react'

interface AuthContextType {
  userId: string | null
  isAuthenticated: boolean
  isLoading: boolean
  showMFAEnroll: boolean
  setShowMFAEnroll: (show: boolean) => void
  login: (credentials: PostAuthLoginBody) => Promise<void>
  signup: (credentials: PostAuthSignupBody) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [userId, setUserId] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [showMFAEnroll, setShowMFAEnroll] = useState(false)
  const router = useRouter()
  const { userLogin, userLogout, userSignup } = useAuth()

  // Check if user is authenticated on mount
  useEffect(() => {
    const checkAuth = () => {
      const storedUserId = localStorage.getItem('userId')
      const storedJwt = localStorage.getItem('jwt')

      if (storedJwt && storedUserId) {
        setUserId(storedUserId)
      }
      else {
        setUserId(null)
      }
      setIsLoading(false)
    }
    checkAuth()
  }, [])

  const login = async (credentials: PostAuthLoginBody) => {
    try {
      const response = await userLogin(credentials)
      console.warn('Login response:', response)

      // Store auth data in localStorage
      if (response.access_token) {
        localStorage.setItem('jwt', response.access_token)
        await supabase.auth.setSession({
          access_token: response.access_token,
          refresh_token: response.refresh_token || '',
        })
      }

      if (response.user?.id) {
        localStorage.setItem('userId', response.user.id)
        setUserId(response.user.id)
      }

      // Check if user has MFA enrolled
      const { data: factors } = await supabase.auth.mfa.listFactors()

      if (!factors || factors.totp.length === 0) {
        // No MFA enrolled, show enrollment
        setShowMFAEnroll(true)
      } else {
        // MFA already enrolled, proceed to students
        router.push('/students')
      }
    }
    catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  const signup = async (credentials: PostAuthSignupBody) => {
    try {
      const response = await userSignup(credentials)

      if (response.access_token) {
        localStorage.setItem('jwt', response.access_token)
        await supabase.auth.setSession({
          access_token: response.access_token,
          refresh_token: response.refresh_token || '',
        })
      }

      if (response.user?.id) {
        localStorage.setItem('userId', response.user.id)
        setUserId(response.user.id)
      }

      // Show MFA enrollment for new users
      setShowMFAEnroll(true)
    }
    catch (error) {
      console.error('Signup failed:', error)
      throw error
    }
  }

  const logout = () => {
    localStorage.removeItem('jwt')
    localStorage.removeItem('userId')

    document.cookie.split(';').forEach((cookie) => {
      const eqPos = cookie.indexOf('=')
      const name = eqPos > -1 ? cookie.substring(0, eqPos).trim() : cookie.trim()
      document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`
    })

    userLogout()
    setUserId(null)
    setShowMFAEnroll(false)
    router.push('/login')
  }

  return (
    <AuthContext.Provider
      value={{
        userId,
        isAuthenticated: !!userId,
        isLoading,
        showMFAEnroll,
        setShowMFAEnroll,
        login,
        signup,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuthContext() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
