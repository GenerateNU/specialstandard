import { useAuthContext } from "@/contexts/authContext";
import { getSessionStudents as getSessionStudentsApi } from "@/lib/api/session-students";
import { getSessions } from "@/lib/api/sessions";
import type {
  CreateSessionStudentInput,
  DeleteSessionStudentsBody,
  StudentWithSessionInfo,
  UpdateSessionStudentInput,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import { gradeToDisplay } from "@/lib/gradeUtils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export type SessionStudentBody = Omit<StudentWithSessionInfo, "grade"> & {
  grade: string;
};

export function useSessionStudentsForSession(sessionId: string) {
  const sessionsApi = getSessions();
  const { userId: therapistId } = useAuthContext();

  const {
    data: studentsData,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["sessions", sessionId, "students", therapistId],
    queryFn: () =>
      sessionsApi.getSessionsSessionIdStudents(sessionId, {
        therapist_id: therapistId!,
      }),
    enabled: !!sessionId && !!therapistId,
  });

  // Flatten the nested student structure
  // Note: API returns nested structure but TypeScript type doesn't reflect this
  const students = (studentsData || []).map((item: any) => ({
    // Spread the nested student object first to get id, first_name, last_name, etc.
    ...(item.student || {}),
    // Then add session-specific fields
    session_id: item.session_id,
    present: item.present,
    notes: item.notes,
    created_at: item.created_at,
    updated_at: item.updated_at,
    session_student_id: item.session_student_id,
    // Override grade with display format
    grade: gradeToDisplay(item.student?.grade ?? item.grade),
  }));

  return {
    students,
    isLoading,
    error: error?.message || null,
    refetch,
  };
}

export function useSessionStudent(sessionStudentId: number, sessionId?: string) {
  const { students, isLoading, error } = useSessionStudentsForSession(sessionId || '');
  
  const sessionStudent = students.find(
    (s: any) => s.session_student_id === sessionStudentId
  );

  return {
    sessionStudent,
    isLoading: isLoading && !sessionStudent,
    error,
  };
}

export function useSessionStudents() {
  const queryClient = useQueryClient();
  const api = getSessionStudentsApi();

  const addStudentToSessionMutation = useMutation({
    mutationFn: (input: CreateSessionStudentInput) =>
      api.postSessionStudents(input),
    onSuccess: (_, variables) => {
      if (variables.session_ids) {
        variables.session_ids.forEach((id: string) => {
          queryClient.invalidateQueries({
            queryKey: ["sessions", id, "students"],
          });
        });
      }

      queryClient.invalidateQueries({ queryKey: ["sessions"] });
    },
  });

  const removeStudentFromSessionMutation = useMutation({
    mutationFn: (input: DeleteSessionStudentsBody) =>
      api.deleteSessionStudents(input),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["sessions", variables.session_id, "students"],
      });
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
    },
  });

  const updateSessionStudentMutation = useMutation({
    mutationFn: (input: UpdateSessionStudentInput) =>
      api.patchSessionStudents(input),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["sessions", variables.session_id, "students"],
      });
    },
  });

  return {
    addStudentToSession: (input: CreateSessionStudentInput) =>
      addStudentToSessionMutation.mutate(input),
    removeStudentFromSession: (input: DeleteSessionStudentsBody) =>
      removeStudentFromSessionMutation.mutate(input),
    updateSessionStudent: (input: UpdateSessionStudentInput) =>
      updateSessionStudentMutation.mutateAsync(input),
    isAdding: addStudentToSessionMutation.isPending,
    isRemoving: removeStudentFromSessionMutation.isPending,
    isUpdating: updateSessionStudentMutation.isPending,
    addError: addStudentToSessionMutation.error?.message || null,
    removeError: removeStudentFromSessionMutation.error?.message || null,
    updateError: updateSessionStudentMutation.error?.message || null,
  };
}