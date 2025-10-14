'use client'

import type { LoginCredentials, SignupCredentials } from '@/lib/api/auth'
import { useRouter } from 'next/navigation'
import { createContext, useContext, useEffect, useState } from 'react'
import {
  login as apiLogin,
  logout as apiLogout,
  signup as apiSignup,
} from '@/lib/api/auth'

interface User {
  id: string
  email?: string
}

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (credentials: LoginCredentials) => Promise<void>
  signup: (credentials: SignupCredentials) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()

  // Check if user is authenticated on mount (by checking cookies)
  useEffect(() => {
    const checkAuth = () => {
      // Check if jwt cookie exists
      const cookies = document.cookie.split(';')
      const jwtCookie = cookies.find(c => c.trim().startsWith('jwt='))
      const userIdCookie = cookies.find(c => c.trim().startsWith('userID='))

      if (jwtCookie && userIdCookie) {
        const userId = userIdCookie.split('=')[1]
        setUser({ id: userId })
      }
      else {
        setUser(null)
      }

      setIsLoading(false)
    }

    checkAuth()
  }, [])

  const login = async (credentials: LoginCredentials) => {
    try {
      const response = await apiLogin(credentials)
      setUser({ id: response.user.id, email: response.user.email })
      router.push('/students')
    }
    catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  const signup = async (credentials: SignupCredentials) => {
    try {
      const response = await apiSignup(credentials)
      setUser({ id: response.user.id, email: response.user.email })
      router.push('/students')
    }
    catch (error) {
      console.error('Signup failed:', error)
      throw error
    }
  }

  const logout = () => {
    apiLogout()
    setUser(null)
    router.push('/login')
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
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

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
