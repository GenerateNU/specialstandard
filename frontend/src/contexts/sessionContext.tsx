'use client'

import type { Session } from '@/lib/api/theSpecialStandardAPI.schemas'
import { createContext, useCallback, useContext, useState } from 'react'

export interface StudentTuple {
  studentId: string
  sessionStudentId: number
}

interface SessionContextType {
  session: Session | null
  students: StudentTuple[]
  currentWeek: number
  setSession: (session: Session) => void
  setStudents: (students: StudentTuple[]) => void
  setCurrentWeek: (week: number) => void
  clearSession: () => void
}

const SessionContext = createContext<SessionContextType | undefined>(undefined)

export function SessionProvider({ children }: { children: React.ReactNode }) {
  const [session, setSessionState] = useState<Session | null>(null)
  const [students, setStudentsState] = useState<StudentTuple[]>([])
  const [currentWeek, setCurrentWeek] = useState<number>(1)

  const setSession = useCallback((newSession: Session) => {
    setSessionState(newSession)
  }, [])

  const setStudents = useCallback((newStudents: StudentTuple[]) => {
    setStudentsState(newStudents)
  }, [])

  const clearSession = useCallback(() => {
    setSessionState(null)
    setStudentsState([])
    setCurrentWeek(1)
  }, [])

  return (
    <SessionContext.Provider
      value={{
        session,
        students,
        currentWeek,
        setSession,
        setStudents,
        setCurrentWeek,
        clearSession,
      }}
    >
      {children}
    </SessionContext.Provider>
  )
}

export function useSessionContext() {
  const context = useContext(SessionContext)
  if (context === undefined) {
    throw new Error('useSessionContext must be used within a SessionProvider')
  }
  return context
}

