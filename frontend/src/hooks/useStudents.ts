import { useAuthContext } from "@/contexts/authContext";
import { getStudents } from "@/lib/api/students";
import type {
  CreateStudentInput,
  Student,
  UpdateStudentInput,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import { gradeToDisplay } from "@/lib/gradeUtils";
import { useRecentlyViewedStudents } from "@/hooks/useRecentlyViewedStudents";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export type StudentBody = Omit<Student, "grade"> & {
  grade: string;
};

interface UseStudentsOptions {
  ids?: string[];
}

export function useStudents(options?: UseStudentsOptions) {
  const queryClient = useQueryClient();
  const api = getStudents();
  const { userId: therapistId } = useAuthContext();
  const { removeRecentStudent } = useRecentlyViewedStudents();
  const ids = options?.ids;
  
  console.warn("🔍 useStudents - therapistId:", therapistId, "ids:", ids);
  
  const {
    data: studentsData = [],
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["students", therapistId, ids],
    queryFn: () => {
      console.warn("🚀 Fetching students for therapist:", therapistId);
      return api.getStudents({ limit: 100, therapist_id: therapistId! })
        .then((data) => {
          console.warn("✅ Students fetched successfully:", data);
          // Filter by ids if provided
          if (ids && ids.length > 0) {
            return data.filter(student => ids.includes(student.id));
          }
          return data;
        })
        .catch((err) => {
          console.warn("❌ Error fetching students:", err);
          throw err;
        });
    },
    // we technically dont need this line but it is just defensive programming!!
    enabled: !!therapistId,
  });
  
  console.warn("📊 Query state - isLoading:", isLoading, "error:", error);
  console.warn("📋 Raw students data:", studentsData);
  
  // get students/id/sessions
  const students = studentsData.map((student) => {
    const transformed = {
      ...student,
      grade: gradeToDisplay(student.grade),
    };
    console.warn(`🔄 Transformed student ${student.id}: grade ${student.grade} → ${transformed.grade}`);
    return transformed;
  });
  
  console.warn("📋 Transformed students:", students);
  
  const addStudentMutation = useMutation({
    mutationFn: (input: CreateStudentInput) => {
      console.warn("➕ Adding student:", input);
      return api.postStudents(input);
    },
    onSuccess: (data) => {
      console.warn("✅ Student added successfully:", data);
      queryClient.invalidateQueries({ queryKey: ["students"] });
    },
    onError: (error) => {
      console.warn("❌ Error adding student:", error);
    },
  });
  
  const updateStudentMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateStudentInput }) => {
      console.warn("✏️ Updating student:", id, "with data:", data);
      return api.patchStudentsId(id, data);
    },
    onSuccess: (data) => {
      console.warn("✅ Student updated successfully:", data);
      queryClient.invalidateQueries({ queryKey: ["students"] });
    },
    onError: (error) => {
      console.warn("❌ Error updating student:", error);
    },
  });
  
  const deleteStudentMutation = useMutation({
    mutationFn: (id: string) => {
      console.warn("🗑️ Deleting student:", id);
      return api.deleteStudentsId(id);
    },
    onSuccess: (data, deletedStudentId) => {
      console.warn("✅ Student deleted successfully:", data);
      // Remove from recently viewed list
      removeRecentStudent(deletedStudentId);
      queryClient.invalidateQueries({ queryKey: ["students"] });
    },
    onError: (error) => {
      console.warn("❌ Error deleting student:", error);
    },
  });
  
  // Log mutation states
  console.warn("🔧 Mutation states:", {
    addStudent: {
      isLoading: addStudentMutation.isPending,
      isError: addStudentMutation.isError,
      error: addStudentMutation.error,
    },
    updateStudent: {
      isLoading: updateStudentMutation.isPending,
      isError: updateStudentMutation.isError,
      error: updateStudentMutation.error,
    },
    deleteStudent: {
      isLoading: deleteStudentMutation.isPending,
      isError: deleteStudentMutation.isError,
      error: deleteStudentMutation.error,
    },
  });
  
  return {
    students,
    isLoading,
    error: error?.message || null,
    refetch,
    addStudent: (student: CreateStudentInput) => {
      console.warn("🎯 addStudent called with:", student);
      return addStudentMutation.mutate(student);
    },
    updateStudent: (id: string, data: UpdateStudentInput) => {
      console.warn("🎯 updateStudent called with id:", id, "data:", data);
      return updateStudentMutation.mutate({ id, data });
    },
    deleteStudent: (id: string) => {
      console.warn("🎯 deleteStudent called with id:", id);
      return deleteStudentMutation.mutate(id);
    },
  };
}