import { useAuthContext } from '@/contexts/authContext'

import type { QueryObserverResult } from '@tanstack/react-query'
import { useQuery } from '@tanstack/react-query'
import { getSchools as getSchoolsApi } from '@/lib/api/schools'
import { getDistricts as getDistrictsApi } from '@/lib/api/districts'
import type { District, School } from '@/lib/api/theSpecialStandardAPI.schemas'

interface UseSchoolsReturn {
  schools: School[]
  districts: District[]
  isLoading: boolean
  error: string | null
  refetch: () => Promise<QueryObserverResult<School[], Error>>
}

export function useSchools(): UseSchoolsReturn {
  const schoolsApi = getSchoolsApi()
  const districtsApi = getDistrictsApi()
  const { userId: therapistId } = useAuthContext()

  const {
    data: schoolsResponse,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['schools', therapistId],
    queryFn: () => schoolsApi.getSchools(),
  })

  const {
    data: districtsResponse,
  } = useQuery({
    queryKey: ['districts', therapistId],
    queryFn: () => districtsApi.getDistricts(),
  })

  const schools = schoolsResponse ?? []
  const districts = districtsResponse ?? []

  return {
    schools,
    districts,
    isLoading,
    error: error?.message || null,
    refetch,
  }
}

