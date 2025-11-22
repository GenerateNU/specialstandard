'use client'

import type { StudentBody } from '@/hooks/useStudents'
import {
  GraduationCap,
} from 'lucide-react'

import { Avatar } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { getAvatarName, getAvatarVariant } from '@/lib/avatarUtils'
import { getSchoolColor } from '@/lib/utils'

interface StudentCardProps {
  student: StudentBody
}
// getAvatarVariant is provided by the shared utils module
export default function StudentCard({ student }: StudentCardProps) {
  const getFullName = () => `${student.first_name} ${student.last_name}`

  const getAge = () => {
    if (!student.dob)
      return null
    const birthDate = new Date(student.dob)
    const today = new Date()
    let age = today.getFullYear() - birthDate.getFullYear()
    const monthDiff = today.getMonth() - birthDate.getMonth()

    if (
      monthDiff < 0
      || (monthDiff === 0 && today.getDate() < birthDate.getDate())
    ) {
      age--
    }

    return age
  }

  // Get avatar variant based on student ID for variety
  const avatarVariant = getAvatarVariant(student.id)

  return (
    <div className="bg-card rounded-lg shadow-md hover:shadow-lg transition-all duration-200 border border-default cursor-pointer">
      <button
        className="w-full px-6 py-4 flex items-center justify-between text-left hover:bg-card-hover rounded-lg transition-colors cursor-pointer"
        aria-controls={`student-details-${student.id}`}
      >
        <div className="flex items-center space-x-4 flex-1">
          <Avatar
            name={getAvatarName(student.first_name, student.last_name, student.id)}
            variant={avatarVariant}
            className="w-12 h-12 ring-2 ring-accent-light"
          />
          <div className="flex-1">
            <h3 className="font-semibold text-lg text-primary">
              {getFullName()}
            </h3>
            <div className="flex items-center space-x-2 text-sm text-secondary">
              {student.grade !== null && student.grade !== undefined && (
                <>
                  <GraduationCap className="w-4 h-4" />
                  <span>
                    Grade
                    {' '}
                    {student.grade}
                  </span>
                </>
              )}
              {student.dob && (
                <>
                  <span className="mx-1">•</span>
                  <span>
                    Age
                    {' '}
                    {getAge()}
                  </span>
                </>
              )}
            </div>
            {student.school_name && (
              <div className="mt-2">
                <Badge className={getSchoolColor(student.school_name)}>
                  {student.school_name}
                </Badge>
              </div>
            )}
          </div>
        </div>
      </button>
    </div>
  )
}
