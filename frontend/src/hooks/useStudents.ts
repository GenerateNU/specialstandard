import type { QueryObserverResult } from "@tanstack/react-query";
import type { Student } from "@/types/student";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createStudent,
  deleteStudent,
  getStudents as fetchStudents,
  updateStudent,
} from "@/lib/api/students";

interface UseStudentsReturn {
  students: Student[];
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<QueryObserverResult<Student[], Error>>;
  addStudent: (student: Omit<Student, "id">) => void;
  updateStudent: (id: string, updatedStudent: Partial<Student>) => void;
  deleteStudent: (id: string) => void;
}

export function useStudents(): UseStudentsReturn {
  const queryClient = useQueryClient();

  // Fetch students
  const {
    data: students = [],
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["students"],
    queryFn: fetchStudents,
  });

  // Create student
  const addStudentMutation = useMutation({
    mutationFn: createStudent,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["students"] });
    },
  });

  // Update student
  const updateStudentMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Student> }) =>
      updateStudent(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["students"] });
    },
  });

  // Delete student
  const deleteStudentMutation = useMutation({
    mutationFn: deleteStudent,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["students"] });
    },
  });

  return {
    students,
    isLoading,
    error: error instanceof Error ? error.message : null,
    refetch,
    addStudent: addStudentMutation.mutate,
    updateStudent: (id, updatedStudent) =>
      updateStudentMutation.mutate({ id, data: updatedStudent }),
    deleteStudent: deleteStudentMutation.mutate,
  };
}
