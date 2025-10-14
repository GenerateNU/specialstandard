'use client'

import type { StudentBody } from '@/hooks/useStudents'
import {
  Calendar,
  ChevronDown,
  ChevronUp,
  FileText,
  GraduationCap,
  User,
} from 'lucide-react'

import { useState } from 'react'

import { gradeToDisplay } from '@/lib/gradeUtils'

interface StudentCardProps {
  student: StudentBody
}

export default function StudentCard({ student }: StudentCardProps) {
  const [isExpanded, setIsExpanded] = useState(false)

  const formatDate = (dateString?: string) => {
    if (!dateString)
      return 'N/A'
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }

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

  return (
    <div className="bg-card rounded-lg shadow-md hover:shadow-lg transition-all duration-200 border border-default">
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full px-6 py-4 flex items-center justify-between text-left hover:bg-card-hover rounded-lg transition-colors"
        aria-expanded={isExpanded ? 'true' : 'false'}
        aria-controls={`student-details-${student.id}`}
      >
        <div className="flex items-center space-x-4">
          <div className="w-12 h-12 bg-accent-light rounded-full flex items-center justify-center">
            <User className="w-6 h-6 text-accent" />
          </div>
          <div>
            <h3 className="font-semibold text-lg text-primary">
              {getFullName()}
            </h3>
            <div className="flex items-center space-x-2 text-sm text-secondary">
              {student.grade !== null && student.grade !== undefined && (
                <>
                  <GraduationCap className="w-4 h-4" />
                  <span>
                    Grade
                    {student.grade}
                  </span>
                </>
              )}
              {student.dob && (
                <>
                  <span className="mx-1">•</span>
                  <span>
                    Age
                    {getAge()}
                  </span>
                </>
              )}
            </div>
          </div>
        </div>
        <div className="ml-4">
          {isExpanded
            ? (
                <ChevronUp className="w-5 h-5 text-muted" />
              )
            : (
                <ChevronDown className="w-5 h-5 text-muted" />
              )}
        </div>
      </button>

      {isExpanded && (
        <div
          id={`student-details-${student.id}`}
          className="px-6 pb-4 space-y-4 border-t border-default"
        >
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 pt-4">
            <div className="space-y-3">
              <div className="flex items-start space-x-3">
                <Calendar className="w-5 h-5 text-accent mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-primary">
                    Date of Birth
                  </p>
                  <p className="text-sm text-secondary">
                    {formatDate(student.dob)}
                  </p>
                </div>
              </div>

              <div className="flex items-start space-x-3">
                <User className="w-5 h-5 text-accent mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-primary">
                    Therapist ID
                  </p>
                  <p className="text-sm text-secondary font-mono text-xs">
                    {student.therapist_id}
                  </p>
                </div>
              </div>
            </div>

            <div className="space-y-3">
              <div className="flex items-start space-x-3">
                <GraduationCap className="w-5 h-5 text-accent mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-primary">
                    Grade Level
                  </p>
                  <p className="text-sm text-secondary">{student.grade}</p>
                </div>
              </div>

              {student.iep && (
                <div className="flex items-start space-x-3">
                  <FileText className="w-5 h-5 text-accent mt-0.5" />
                  <div>
                    <p className="text-sm font-medium text-primary">IEP</p>
                    <p className="text-sm text-secondary">{student.iep}</p>
                  </div>
                </div>
              )}
            </div>
          </div>

          <div className="pt-3 border-t border-default">
            <div className="flex justify-between text-xs text-muted">
              <span>
                Created:
                {formatDate(student.created_at)}
              </span>
              <span>
                Updated:
                {formatDate(student.updated_at)}
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
