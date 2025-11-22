import { getSchools } from "@/lib/api/schools";
import { useQuery } from "@tanstack/react-query";

export interface SchoolOption {
  id: number;
  name: string;
}

/**
 * Hook to fetch all schools from the API
 */
export function useSchools() {
  const api = getSchools();

  const {
    data: schoolsData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["schools"],
    queryFn: () => api.getSchools(),
  });

  const schools: SchoolOption[] = (schoolsData || []).map(school => ({
    id: school.id!,
    name: school.name!,
  }));

  return {
    schools,
    isLoading,
    error: error?.message || null,
  };
}

