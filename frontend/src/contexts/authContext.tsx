'use client'
import type {
  PostAuthLoginBody,
  PostAuthSignupBody,
} from '@/lib/api/theSpecialStandardAPI.schemas'
import { useRouter } from 'next/navigation'
import { createContext, useContext, useEffect, useState } from 'react'
import { useAuth } from '@/hooks/useAuth'

interface AuthContextType {
  userId: string | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (credentials: PostAuthLoginBody) => Promise<void>
  signup: (credentials: PostAuthSignupBody) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [userId, setUserId] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()
  const { userLogin, userLogout, userSignup } = useAuth()

  // Check if user is authenticated on mount (check localStorage instead of cookies)
  useEffect(() => {
    const checkAuth = () => {
      // Check localStorage for auth data
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

      // Store auth data in localStorage
      if (response.access_token) {
        localStorage.setItem('jwt', response.access_token)
      }
      if (response.user?.id) {
        localStorage.setItem('userId', response.user.id)
        setUserId(response.user.id)
      }

      router.push('/students')
    }
    catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  const signup = async (credentials: PostAuthSignupBody) => {
    try {
      const response = await userSignup(credentials)

      // Store auth data in localStorage
      if (response.access_token) {
        localStorage.setItem('jwt', response.access_token)
      }
      if (response.user?.id) {
        localStorage.setItem('userId', response.user.id)
        setUserId(response.user.id)
      }

      router.push('/students')
    }
    catch (error) {
      console.error('Signup failed:', error)
      throw error
    }
  }

  const logout = () => {
    // Clear localStorage
    localStorage.removeItem('jwt')
    localStorage.removeItem('userId')

    // Clear any remaining cookies (for backwards compatibility)
    document.cookie.split(';').forEach((cookie) => {
      const eqPos = cookie.indexOf('=')
      const name = eqPos > -1 ? cookie.substring(0, eqPos).trim() : cookie.trim()
      document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`
    })

    userLogout()
    setUserId(null)
    router.push('/login')
  }

  return (
    <AuthContext.Provider
      value={{
        userId,
        isAuthenticated: !!userId,
        isLoading,
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
