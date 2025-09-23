// src/lib/api/sessions.ts
import type {
  CreateSessionInput,
  Session,
  UpdateSessionInput,
} from '@/types/session'
import apiClient from './apiClient'

export async function fetchSessions(): Promise<Session[]> {
  try {
    const response = await apiClient.get<Session[]>('/api/v1/sessions')
    return response.data
  }
  catch (error) {
    console.error('Error fetching sessions:', error)
    throw error
  }
}

export async function fetchSession(sessionId: string): Promise<Session | null> {
  try {
    const response = await apiClient.get<Session>(`/api/v1/sessions/${sessionId}`)
    return response.data
  }
  catch (error) {
    console.error('Error fetching session:', error)
    return null
  }
}

export async function createSession(data: CreateSessionInput): Promise<Session | null> {
  try {
    const response = await apiClient.post<Session>('/api/v1/sessions', data)
    return response.data
  }
  catch (error) {
    console.error('Error creating session:', error)
    return null
  }
}

export async function updateSession(
  sessionId: string,
  data: UpdateSessionInput,
): Promise<Session | null> {
  try {
    const response = await apiClient.patch<Session>(`/api/v1/sessions/${sessionId}`, data)
    return response.data
  }
  catch (error) {
    console.error('Error updating session:', error)
    return null
  }
}

export async function deleteSession(sessionId: string): Promise<boolean> {
  try {
    await apiClient.delete(`/api/v1/sessions/${sessionId}`)
    return true
  }
  catch (error) {
    console.error('Error deleting session:', error)
    return false
  }
}
