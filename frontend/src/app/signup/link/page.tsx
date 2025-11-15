'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Dropdown } from '@/components/ui/dropdown'
import { MultiSelect } from '@/components/ui/multiselect'
import { ArrowLeft, Loader2 } from 'lucide-react'
import { useSchools } from '@/hooks/useSchools'
import { useTherapists } from '@/hooks/useTherapists'
import CustomAlert from '@/components/ui/CustomAlert'

export default function ProfilePage() {
  const router = useRouter()
  const { schools, districts, isLoading: isLoadingSchools } = useSchools()
  const { updateTherapist } = useTherapists()
  
  const [selectedDistrict, setSelectedDistrict] = useState<string>('')
  const [selectedSchools, setSelectedSchools] = useState<string[]>([])
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [showError, setShowError] = useState(false)
  const [therapistData, setTherapistData] = useState({
    firstName: '',
    lastName: '',
    email: '',
  })

  // Load saved data from localStorage (from signup step)
  useEffect(() => {
    const savedData = localStorage.getItem('onboardingData')
    const userId = localStorage.getItem('userId')
    
    if (savedData) {
      const parsed = JSON.parse(savedData)
      setTherapistData({
        firstName: parsed.firstName || '',
        lastName: parsed.lastName || '',
        email: parsed.email || '',
      })
    }
    
    if (!userId) {
      // If no userId, redirect back to signup
      router.push('/signup/welcome')
    }
  }, [router])

  // Filter schools based on selected district
  const filteredSchools = selectedDistrict
    ? schools.filter(school => school.district_id === Number(selectedDistrict))
    : schools

  // Convert data for dropdown components
  const districtOptions = districts.map(district => ({
    label: district.name ?? '',
    value: district.id!.toString(),
  }))

  const schoolOptions = filteredSchools.map(school => ({
    label: school.name ?? '',
    value: school.id!.toString(),
  }))

  const handleBack = () => {
    router.push('/signup/welcome')
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!selectedDistrict) {
      setError('Please select a district')
      setShowError(true)
      return
    }

    if (selectedSchools.length === 0) {
      setError('Please select at least one school')
      setShowError(true)
      return
    }

    setIsSubmitting(true)
    setError(null)

    try {
      const userId = localStorage.getItem('userId')
      
      if (!userId) {
        throw new Error('User ID not found. Please sign up again.')
      }

      // Create therapist profile
      await updateTherapist(
        userId,
        { first_name: therapistData.firstName,
        last_name: therapistData.lastName,
        email: therapistData.email,
        district_id: Number(selectedDistrict),
        schools: selectedSchools.map(id => Number(id)),
      })

      // Save to localStorage for next steps
      const profileData = {
        ...therapistData,
        districtId: selectedDistrict,
        schoolIds: selectedSchools,
      }
      localStorage.setItem('therapistProfile', JSON.stringify(profileData))
      
      // Move to next step
      router.push('/signup/students/add')
    } catch (err: any) {
      console.error('Profile setup error:', err)
      setError(err?.message || 'Failed to set up profile. Please try again.')
      setShowError(true)
    } finally {
      setIsSubmitting(false)
    }
  }

  // When district changes, clear selected schools if they're not in the new district
  const handleDistrictChange = (value: string) => {
    setSelectedDistrict(value)
    
    // Clear schools that don't belong to the new district
    const newDistrictSchoolIds = schools
      .filter(school => school.district_id === Number(value))
      .map(school => school.id!.toString())
    
    setSelectedSchools(prev => 
      prev.filter(schoolId => newDistrictSchoolIds.includes(schoolId))
    )
  }

  if (isLoadingSchools) {
    return (
      <div className="flex items-center justify-center min-h-screen p-8">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="flex items-center justify-center min-h-screen p-8">
      <div className="max-w-md w-full">
        <button
          onClick={handleBack}
          className="mb-6 flex items-center text-secondary hover:text-primary transition-colors"
        >
          <ArrowLeft className="w-4 h-4 mr-1" />
          Back
        </button>

        <h1 className="text-3xl font-bold text-primary mb-2">
          Almost There!
        </h1>
        
        <p className="text-secondary mb-8">
          Registered District
        </p>
        
        {showError && error && (
          <div className="mb-4">
            <CustomAlert
              variant="destructive"
              title="Setup Error"
              description={error}
              onClose={() => {
                setShowError(false)
                setError(null)
              }}
            />
          </div>
        )}
        
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-primary mb-2">
              Select District
            </label>
            <Dropdown
              items={districtOptions.map(opt => ({
                label: opt.label,
                value: opt.value,
                onClick: () => handleDistrictChange(opt.value),
              }))}
              value={selectedDistrict}
              placeholder="Select District"
              className="w-full min-w-30[rem]"
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-primary mb-2">
              Registered School
            </label>
            <MultiSelect
              options={schoolOptions}
              value={selectedSchools}
              onValueChange={setSelectedSchools}
              placeholder="Select School"
              showTags={false}
              showCount={false}
              className="w-full"
            />
            {selectedDistrict && selectedSchools.length > 0 && (
              <p className="text-xs text-secondary mt-2">
                {selectedSchools.length} school{selectedSchools.length !== 1 ? 's' : ''} selected
              </p>
            )}
          </div>
          
          {(selectedDistrict || selectedSchools.length > 0) && (
            <div className="mt-2">
              <button
                type="button"
                onClick={() => {
                  setSelectedSchools([])
                  setSelectedDistrict('')
                }}
                className="text-xs text-accent hover:text-accent-hover underline"
              >
                Reset Districts and Schools
              </button>
            </div>
          )}
          
          <div className="pt-4">
            <Button
              type="submit"
              size="long"
              className="w-full text-white"
              disabled={isSubmitting || !selectedDistrict || selectedSchools.length === 0}
            >
              {isSubmitting ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin mr-2" />
                  <span>Setting up profile...</span>
                </>
              ) : (
                <span>Continue</span>
              )}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}