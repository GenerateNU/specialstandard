import { useAuthContext } from '@/contexts/authContext'
import { getStudents } from '@/lib/api/students'
import type {
  CreateStudentInput,
  Student,
  UpdateStudentInput,
} from '@/lib/api/theSpecialStandardAPI.schemas'
import { gradeToDisplay } from '@/lib/gradeUtils'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

export type StudentBody = Omit<Student, 'grade'> & {
  grade: string
}

// interface UseStudentsReturn {
//   students: StudentBody[]
//   isLoading: boolean
//   error: string | null
//   refetch: () => Promise<QueryObserverResult<Student[], Error>>
//   addStudent: (student: Omit<Student, "id">) => void
//   updateStudent: (id: string, updatedStudent: Partial<Student>) => void
//   deleteStudent: (id: string) => void
// }

export function useStudents() {
  const queryClient = useQueryClient()
  const api = getStudents()
  const { userId: therapistId } = useAuthContext()

  const {
    data: studentsData = [],
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['students'],
    queryFn: () => api.getStudents({ limit: 100, therapist_id: therapistId! }), //TODO: add this, get rid of queryKey, and update get endpoints that dont have this, sessions, sessionstudents, student, session resources
    // we technically dont need this line but it is just defensive programming!!  
    enabled: !!therapistId,
  })

  // get students/id/sessions

  const students = studentsData.map(student => ({
    ...student,
    grade: gradeToDisplay(student.grade),
  }))

  const addStudentMutation = useMutation({
    mutationFn: (input: CreateStudentInput) => api.postStudents(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['students', therapistId] })
    },
  })

  const updateStudentMutation = useMutation({
    mutationFn: ({ id, data }: { id: string, data: UpdateStudentInput }) =>
      api.patchStudentsId(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['students', therapistId] })
    },
  })

  const deleteStudentMutation = useMutation({
    mutationFn: (id: string) => api.deleteStudentsId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['students', therapistId] })
    },
  })

  return {
    students,
    isLoading,
    error: error?.message || null,
    refetch,
    addStudent: (student: CreateStudentInput) =>
      addStudentMutation.mutate(student),
    updateStudent: (id: string, data: UpdateStudentInput) =>
      updateStudentMutation.mutate({ id, data }),
    deleteStudent: (id: string) => deleteStudentMutation.mutate(id),
  }
}
