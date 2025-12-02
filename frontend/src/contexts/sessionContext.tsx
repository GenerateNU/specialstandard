'use client'

import type { Session } from '@/lib/api/theSpecialStandardAPI.schemas'
import { createContext, useCallback, useContext, useState } from 'react'

export interface StudentTuple {
  studentId: string
  sessionStudentId: number
  firstName: string
  lastName: string
}

interface SessionContextType {
  session: Session | null
  students: StudentTuple[]
  currentWeek: number
  currentMonth: number
  currentYear: number
  currentLevel: number | null
  setSession: (session: Session) => void
  setStudents: (students: StudentTuple[]) => void
  setCurrentWeek: (week: number) => void
  setCurrentMonth: (month: number) => void
  setCurrentYear: (year: number) => void
  setCurrentLevel: (level: number | null) => void
  clearSession: () => void
}

const SessionContext = createContext<SessionContextType | undefined>(undefined)

export function SessionProvider({ children }: { children: React.ReactNode }) {
  const [session, setSessionState] = useState<Session | null>(null)
  const [students, setStudentsState] = useState<StudentTuple[]>([])
  const [currentWeek, setCurrentWeek] = useState<number>(1)
  const now = new Date()
  const [currentMonth, setCurrentMonth] = useState<number>(now.getMonth()) // 0-11
  const [currentYear, setCurrentYear] = useState<number>(now.getFullYear())
  const [currentLevel, setCurrentLevel] = useState<number | null>(null)

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
    const now = new Date()
    setCurrentMonth(now.getMonth())
    setCurrentYear(now.getFullYear())
    setCurrentLevel(null)
  }, [])

  return (
    <SessionContext.Provider
      value={{
        session,
        students,
        currentWeek,
        currentMonth,
        currentYear,
        currentLevel,
        setSession,
        setStudents,
        setCurrentWeek,
        setCurrentMonth,
        setCurrentYear,
        setCurrentLevel,
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

