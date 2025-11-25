'use client'

import { AlertCircle, ArrowLeft, Loader2, Plus, Users, X } from 'lucide-react'
import Link from 'next/link'
import { useMemo, useState } from 'react'
import AppLayout from '@/components/AppLayout'
import AddStudentModal from '@/components/students/AddStudentModal'
import StudentCard from '@/components/students/studentCard'
import { Badge } from '@/components/ui/badge'
import { useStudents } from '@/hooks/useStudents'
import { getSchoolColor } from '@/lib/utils'

export default function StudentsPage() {
  const { students, isLoading, error, refetch } = useStudents()
  const [selectedSchools, setSelectedSchools] = useState<string[]>([])

  // Get unique schools from students
  const uniqueSchools = useMemo(() => {
    const schools = students
      .map(s => s.school_name)
      .filter((name): name is string => !!name)
    return Array.from(new Set(schools)).sort()
  }, [students])

  // Filter students based on selected schools
  const filteredStudents = useMemo(() => {
    if (selectedSchools.length === 0) {
      return students
    }
    return students.filter(student => 
      student.school_name && selectedSchools.includes(student.school_name)
    )
  }, [students, selectedSchools])

  // Toggle school filter
  const toggleSchoolFilter = (schoolName: string) => {
    setSelectedSchools(prev => 
      prev.includes(schoolName)
        ? prev.filter(s => s !== schoolName)
        : [...prev, schoolName]
    )
  }

  // Clear all filters
  const clearFilters = () => {
    setSelectedSchools([])
  }

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <Loader2 className="w-8 h-8 animate-spin text-accent mx-auto mb-4" />
          <p className="text-secondary">Loading students...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <AppLayout>
        <div className="min-h-screen flex items-center justify-center bg-background">
          <div className="text-center max-w-md">
            <AlertCircle className="w-12 h-12 text-error mx-auto mb-4" />
            <h2 className="text-xl font-semibold text-primary mb-2">
              Error Loading Students
            </h2>
            <p className="text-secondary mb-4">{error}</p>
            <button
              onClick={() => refetch()}
              className="px-4 py-2 bg-accent text-white rounded-lg hover:bg-accent-hover transition-colors"
            >
              Try Again
            </button>
          </div>
        </div>
      </AppLayout>
    )
  }

  return (
    <AppLayout>
      <div className="min-h-screen bg-background py-8">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <header className="mb-8">
            {/* Back button */}
            <Link
              href="/"
              className="inline-flex items-center gap-2 text-secondary hover:text-primary mb-4 transition-colors group"
            >
              <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
              <span className="text-sm font-medium">Back to Home</span>
            </Link>

            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center space-x-3">
                <Users className="w-8 h-8 text-accent" />
                <h1>Students</h1>
              </div>
              <AddStudentModal />
            </div>
            <p className="text-secondary">
              View and manage all student information
            </p>
          </header>

          {/* School Filters */}
          {uniqueSchools.length > 0 && (
            <div className="mb-6 flex flex-wrap items-center gap-2">
              <span className="text-sm font-medium text-secondary">Filter:</span>
              {uniqueSchools.map(schoolName => {
                const isSelected = selectedSchools.includes(schoolName)
                return (
                  <button
                    key={schoolName}
                    onClick={() => toggleSchoolFilter(schoolName)}
                    className={`transition-all ${
                      isSelected 
                        ? 'opacity-100' 
                        : 'opacity-40 hover:opacity-70'
                    }`}
                  >
                    <Badge className={`${getSchoolColor(schoolName)} ${isSelected ? 'ring-2 ring-accent ring-offset-2' : ''}`}>
                      {schoolName}
                      {isSelected && <X className="w-3 h-3 ml-1 inline" />}
                    </Badge>
                  </button>
                )
              })}
            </div>
          )}

          {students.length === 0
            ? (
                <div className="text-center py-12">
                  <Users className="w-16 h-16 text-muted mx-auto mb-4 opacity-30" />
                  <h2 className="text-xl font-semibold text-primary mb-2">
                    No Students Found
                  </h2>
                  <p className="text-secondary mb-6">
                    There are no students in the system yet. Get started by adding
                    your first student.
                  </p>
                  <AddStudentModal
                    trigger={(
                      <button className="px-6 py-3 bg-accent text-white rounded-lg hover:bg-accent-hover transition-colors flex items-center gap-2 mx-auto">
                        <Plus className="w-4 h-4" />
                        Add Your First Student
                      </button>
                    )}
                  />
                </div>
              )
            : filteredStudents.length === 0
              ? (
                  <div className="text-center py-12">
                    <Users className="w-16 h-16 text-muted mx-auto mb-4 opacity-30" />
                    <h2 className="text-xl font-semibold text-primary mb-2">
                      No Students Match Filter
                    </h2>
                    <p className="text-secondary mb-6">
                      Try selecting different schools or clear the filters to see all students.
                    </p>
                    <button
                      onClick={clearFilters}
                      className="px-6 py-3 bg-accent text-white rounded-lg hover:bg-accent-hover transition-colors flex items-center gap-2 mx-auto"
                    >
                      Clear Filters
                    </button>
                  </div>
                )
              : (
                  <div className="space-y-4">
                    <div className="flex justify-between items-center mb-4">
                      <p className="text-sm text-secondary">
                        Showing
                        {' '}
                        {filteredStudents.length}
                        {' '}
                        of
                        {' '}
                        {students.length}
                        {' '}
                        student
                        {students.length !== 1 ? 's' : ''}
                      </p>
                    </div>
                    <div className="grid grid-cols-3 gap-4">
                      {filteredStudents.map(student => (
                        <Link key={student.id} href={`/student/${student.id}`}>
                          <StudentCard student={student} />
                        </Link>
                      ))}
                    </div>
                  </div>
                )}
        </div>
      </div>
    </AppLayout>
  )
}
