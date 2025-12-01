import { useQuery } from '@tanstack/react-query';
import { getStudents } from '@/lib/api/students';
import type { GetStudentsStudentIdRatingsParams } from '@/lib/api/theSpecialStandardAPI.schemas';

interface UseStudentRatingsProps {
	studentId: string;
	params?: GetStudentsStudentIdRatingsParams;
}

export function useStudentRatings({ studentId, params }: UseStudentRatingsProps) {
	const api = getStudents();

	const {
		data: ratings,
		isLoading,
		error,
		refetch,
	} = useQuery({
		queryKey: ['student-ratings', studentId, params ?? {}],
		queryFn: () => api.getStudentsStudentIdRatings(studentId, params),
		enabled: Boolean(studentId),
	});

	return {
		ratings: ratings ?? [],
		isLoading,
		error,
		refetch,
	};
}

