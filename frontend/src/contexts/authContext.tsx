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
  // Check if user is authenticated on mount (by checking cookies)
  useEffect(() => {
    const checkAuth = () => {
      // Check if jwt cookie exists
      const cookies = document.cookie.split(';')
      const jwtCookie = cookies.find(c => c.trim().startsWith('jwt='))
      const userIdCookie = cookies.find(c => c.trim().startsWith('userID='))

      if (jwtCookie && userIdCookie) {
        const userId = userIdCookie.split('=')[1]
        setUserId(userId)
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
      setUserId(response.user.id ?? null)
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
      setUserId(response.user.id ?? null)
      router.push('/students')
    }
    catch (error) {
      console.error('Signup failed:', error)
      throw error
    }
  }

  const logout = () => {
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
