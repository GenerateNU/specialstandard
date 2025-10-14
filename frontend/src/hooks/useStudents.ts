import type { QueryObserverResult } from "@tanstack/react-query";
import type {
  CreateStudentInput,
  Student,
  UpdateStudentInput,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getStudents } from "@/lib/api/students";
import { gradeToDisplay, gradeToStorage } from "@/lib/gradeUtils";

export type StudentBody = Omit<Student, "grade"> & {
  grade: string;
};

interface UseStudentsReturn {
  students: StudentBody[];
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<QueryObserverResult<Student[], Error>>;
  addStudent: (student: Omit<Student, "id">) => void;
  updateStudent: (id: string, updatedStudent: Partial<Student>) => void;
  deleteStudent: (id: string) => void;
}

export function useStudents() {
  const queryClient = useQueryClient();
  const api = getStudents();

  const {
    data: studentsData = [],
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["students"],
    queryFn: () => api.getStudents({ limit: 100 }),
  });

  const students = studentsData.map((student) => ({
    ...student,
    grade: gradeToDisplay(student.grade),
  }));

  const addStudentMutation = useMutation({
    mutationFn: (input: CreateStudentInput) => api.postStudents(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["students"] });
    },
  });

  const updateStudentMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateStudentInput }) =>
      api.patchStudentsId(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["students"] });
    },
  });

  const deleteStudentMutation = useMutation({
    mutationFn: (id: string) => api.deleteStudentsId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["students"] });
    },
  });

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
  };
}
