import type { Student } from '@/types/student'
// src/hooks/useStudents.ts
import { useCallback, useEffect, useState } from 'react'
import { fetchStudents } from '@/lib/api/students'

interface UseStudentsReturn {
  students: Student[]
  isLoading: boolean
  error: string | null
  refetch: () => Promise<void>
  addStudent: (student: Student) => void
  updateStudent: (id: string, updatedStudent: Partial<Student>) => void
  deleteStudent: (id: string) => void
}

export function useStudents(): UseStudentsReturn {
  const [students, setStudents] = useState<Student[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const loadStudents = useCallback(async () => {
    try {
      setIsLoading(true)
      setError(null)
      const data = await fetchStudents()
      setStudents(data)
    }
    catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load students'
      setError(errorMessage)
      console.error('Error loading students:', err)
    }
    finally {
      setIsLoading(false)
    }
  }, [])

  // Initial fetch
  useEffect(() => {
    loadStudents()
  }, [loadStudents])

  // Optimistic update functions
  const addStudent = useCallback((student: Student) => {
    setStudents(prev => [...prev, student])
  }, [])

  const updateStudent = useCallback((id: string, updatedStudent: Partial<Student>) => {
    setStudents(prev =>
      prev.map(student =>
        student.id === id ? { ...student, ...updatedStudent } : student,
      ),
    )
  }, [])

  const deleteStudent = useCallback((id: string) => {
    setStudents(prev => prev.filter(student => student.id !== id))
  }, [])

  return {
    students,
    isLoading,
    error,
    refetch: loadStudents,
    addStudent,
    updateStudent,
    deleteStudent,
  }
}
