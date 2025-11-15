import { useCallback, useEffect, useState } from 'react'

const STORAGE_KEY = 'recentlyViewedStudents'
const MAX_RECENT_STUDENTS = 3

interface RecentStudent {
  id: string
  first_name: string
  last_name: string
  grade?: number | null
  viewedAt: number
}

export function useRecentlyViewedStudents() {
  const [recentStudents, setRecentStudents] = useState<RecentStudent[]>([])

  // Load from localStorage on mount
  useEffect(() => {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      try {
        const parsed = JSON.parse(stored) as RecentStudent[]
        setRecentStudents(parsed)
      }
      catch (error) {
        console.error('Failed to parse recently viewed students:', error)
      }
    }
  }, [])

  // Add a student to recently viewed - memoized to prevent infinite loops
  const addRecentStudent = useCallback((student: { id: string, first_name: string, last_name: string, grade?: number | string | null }) => {
    const recentStudent: RecentStudent = {
      id: student.id,
      first_name: student.first_name,
      last_name: student.last_name,
      grade: typeof student.grade === 'number' ? student.grade : (student.grade ? Number.parseInt(student.grade as string, 10) : null),
      viewedAt: Date.now(),
    }

    setRecentStudents((prev) => {
      // Remove if already exists
      const filtered = prev.filter(s => s.id !== student.id)
      // Add to front
      const updated = [recentStudent, ...filtered].slice(0, MAX_RECENT_STUDENTS)
      // Save to localStorage
      localStorage.setItem(STORAGE_KEY, JSON.stringify(updated))
      return updated
    })
  }, [])

  // Clear all recent students
  const clearRecentStudents = useCallback(() => {
    setRecentStudents([])
    localStorage.removeItem(STORAGE_KEY)
  }, [])

  return {
    recentStudents,
    addRecentStudent,
    clearRecentStudents,
  }
}
