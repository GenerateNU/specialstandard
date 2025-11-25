import { getSchools } from "@/lib/api/schools";
import { getDistricts } from "@/lib/api/districts";
import { useQuery } from "@tanstack/react-query";

export interface SchoolOption {
  id: number;
  name: string;
  district_id?: number;
}

export interface DistrictOption {
  id: number;
  name: string;
}

/**
 * Hook to fetch all schools and districts from the API
 */
export function useSchools() {
  const schoolsApi = getSchools();
  const districtsApi = getDistricts();

  const {
    data: schoolsData,
    isLoading: isLoadingSchools,
    error: schoolsError,
  } = useQuery({
    queryKey: ["schools"],
    queryFn: () => schoolsApi.getSchools(),
  });

  const {
    data: districtsData,
    isLoading: isLoadingDistricts,
    error: districtsError,
  } = useQuery({
    queryKey: ["districts"],
    queryFn: () => districtsApi.getDistricts(),
  });

  const schools: SchoolOption[] = (schoolsData || []).map(school => ({
    id: school.id!,
    name: school.name!,
    district_id: school.district_id,
  }));

  const districts: DistrictOption[] = (districtsData || []).map(district => ({
    id: district.id!,
    name: district.name!,
  }));

  return {
    schools,
    districts,
    isLoading: isLoadingSchools || isLoadingDistricts,
    error: schoolsError?.message || districtsError?.message || null,
  };
}

