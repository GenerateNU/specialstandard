// src/lib/api/auth.ts
import apiClient from './apiClient'

export interface LoginCredentials {
  email: string
  password: string
  remember_me?: boolean
}

export interface SignupCredentials {
  email: string
  password: string
  first_name: string
  last_name: string
}

export interface AuthResponse {
  access_token: string
  user: {
    id: string
    email?: string
  }
}

export async function login(
  credentials: LoginCredentials,
): Promise<AuthResponse> {
  const response = await apiClient.post<AuthResponse>(
    '/api/v1/auth/login',
    credentials,
  )
  return response.data
}

export async function signup(
  credentials: SignupCredentials,
): Promise<AuthResponse> {
  const response = await apiClient.post<AuthResponse>(
    '/api/v1/signup',
    credentials,
  )
  return response.data
}

export async function logout(): Promise<void> {
  // Since auth is cookie-based, we just need to clear cookies client-side
  // Backend doesn't have a logout endpoint - cookies expire naturally
  document.cookie.split(';').forEach((cookie) => {
    const eqPos = cookie.indexOf('=')
    const name = eqPos > -1 ? cookie.substring(0, eqPos).trim() : cookie.trim()
    document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`
  })
}
