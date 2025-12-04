'use client'
import type { Session } from '@/lib/api/theSpecialStandardAPI.schemas'
import { createContext, useCallback, useContext, useEffect, useState } from 'react'

export interface StudentTuple {
  studentId: string
  sessionStudentId: number
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
  const [currentMonth, setCurrentMonth] = useState<number>(now.getMonth())
  const [currentYear, setCurrentYear] = useState<number>(now.getFullYear())
  const [currentLevel, setCurrentLevel] = useState<number | null>(null)

  // Load from localStorage on mount
  useEffect(() => {
    try {
      const savedSession = localStorage.getItem('session')
      const savedStudents = localStorage.getItem('students')
      const savedCurrentWeek = localStorage.getItem('currentWeek')
      const savedCurrentMonth = localStorage.getItem('currentMonth')
      const savedCurrentYear = localStorage.getItem('currentYear')
      const savedCurrentLevel = localStorage.getItem('currentLevel')

      if (savedSession) setSessionState(JSON.parse(savedSession))
      if (savedStudents) setStudentsState(JSON.parse(savedStudents))
      if (savedCurrentWeek) setCurrentWeek(Number(savedCurrentWeek))
      if (savedCurrentMonth) setCurrentMonth(Number(savedCurrentMonth))
      if (savedCurrentYear) setCurrentYear(Number(savedCurrentYear))
      if (savedCurrentLevel) setCurrentLevel(JSON.parse(savedCurrentLevel))
    } catch (error) {
      console.error('Failed to load session from localStorage:', error)
    }
  }, [])

  // Sync state to localStorage whenever it changes
  useEffect(() => {
    localStorage.setItem('session', JSON.stringify(session))
  }, [session])

  useEffect(() => {
    localStorage.setItem('students', JSON.stringify(students))
  }, [students])

  useEffect(() => {
    localStorage.setItem('currentWeek', String(currentWeek))
  }, [currentWeek])

  useEffect(() => {
    localStorage.setItem('currentMonth', String(currentMonth))
  }, [currentMonth])

  useEffect(() => {
    localStorage.setItem('currentYear', String(currentYear))
  }, [currentYear])

  useEffect(() => {
    localStorage.setItem('currentLevel', JSON.stringify(currentLevel))
  }, [currentLevel])

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

    // Clear localStorage
    localStorage.removeItem('session')
    localStorage.removeItem('students')
    localStorage.removeItem('currentWeek')
    localStorage.removeItem('currentMonth')
    localStorage.removeItem('currentYear')
    localStorage.removeItem('currentLevel')
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