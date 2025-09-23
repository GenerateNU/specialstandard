// src/hooks/useSessions.ts
import type { Session } from '@/types/session'
import { useCallback, useEffect, useState } from 'react'
import { fetchSessions } from '@/lib/api/sessions'

interface UseSessionsReturn {
  sessions: Session[]
  isLoading: boolean
  error: string | null
  refetch: () => Promise<void>
  addSession: (session: Session) => void
  updateSession: (id: string, updatedSession: Partial<Session>) => void
  deleteSession: (id: string) => void
}

export function useSessions(): UseSessionsReturn {
  const [sessions, setSessions] = useState<Session[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const loadSessions = useCallback(async () => {
    try {
      setIsLoading(true)
      setError(null)
      const data = await fetchSessions()
      setSessions(data)
    }
    catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load sessions'
      setError(errorMessage)
      console.error('Error loading sessions:', err)
    }
    finally {
      setIsLoading(false)
    }
  }, [])

  // Initial fetch
  useEffect(() => {
    loadSessions()
  }, [loadSessions])

  // Optimistic update functions
  const addSession = useCallback((session: Session) => {
    setSessions(prev => [...prev, session])
  }, [])

  const updateSession = useCallback((id: string, updatedSession: Partial<Session>) => {
    setSessions(prev =>
      prev.map(session =>
        session.id === id ? { ...session, ...updatedSession } : session,
      ),
    )
  }, [])

  const deleteSession = useCallback((id: string) => {
    setSessions(prev => prev.filter(session => session.id !== id))
  }, [])

  return {
    sessions,
    isLoading,
    error,
    refetch: loadSessions,
    addSession,
    updateSession,
    deleteSession,
  }
}
