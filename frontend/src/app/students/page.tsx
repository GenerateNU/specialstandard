'use client'

import {AlertCircle, ChevronsUpDown, ListFilter, Loader2, Plus, Users, X} from 'lucide-react'
import Link from 'next/link'
import { useMemo, useState } from 'react'
import AppLayout from '@/components/AppLayout'
import { PageHeader } from '@/components/PageHeader'
import AddStudentModal from '@/components/students/AddStudentModal'
import StudentCard from '@/components/students/studentCard'
import { Badge } from '@/components/ui/badge'
import { useStudents } from '@/hooks/useStudents'
import { getSchoolColor } from '@/lib/utils'
import {gradeToDisplay} from "@/lib/gradeUtils";

export default function StudentsPage() {
  const { students, isLoading, error, refetch } = useStudents()
  const [selectedSchools, setSelectedSchools] = useState<string[]>([])
  const [selectedGrades, setSelectedGrades] = useState<string[]>([])

  const [filterOpen, setFilterOpen] = useState(false)
  const [sortAZ, setSortAZ] = useState(true)

  // Get unique schools from students
  const uniqueSchools = useMemo(() => {
    const schools = students
      .map(s => s.school_name)
      .filter((name): name is string => !!name)
    return Array.from(new Set(schools)).sort()
  }, [students])

  const gradeColor = useMemo(() => {
    const freq: Record<string, number> = {}
    students.forEach(s => {
      if (s.school_name) freq[s.school_name] = (freq[s.school_name] || 0) + 1
    })
    const mostFreqSchool = Object.entries(freq).sort((a, b) => b[1] - a[1])[0]?.[0]
    return mostFreqSchool ? getSchoolColor(mostFreqSchool) : 'bg-orange'
  }, [students])

  const uniqueGrades = useMemo(() => {
    const grades = students.map(s => s.grade).filter((g): g is string => !!g)
    return Array.from(new Set(grades)).sort()
  }, [students])

  // Filter students based on selected schools
  const filteredStudents = useMemo(() => {
    let results = students

    if (selectedSchools.length > 0 || selectedGrades.length > 0) {
      results = students.filter(student =>
        (student.school_name && selectedSchools.includes(student.school_name))
        || (student.grade && selectedGrades.includes(student.grade))
      )
    }

    return results.sort((stuA, stuB) => {
      const firstCheck = sortAZ
        ? stuA.first_name.toLowerCase().localeCompare(stuB.first_name.toLowerCase())
        : stuB.first_name.toLowerCase().localeCompare(stuA.first_name.toLowerCase())
      const secondCheck = sortAZ
        ? stuA.last_name.toLowerCase().localeCompare(stuB.last_name.toLowerCase())
        : stuB.last_name.toLowerCase().localeCompare(stuA.last_name.toLowerCase())
      return (firstCheck !== 0) ? firstCheck : secondCheck
    })
  }, [students, selectedSchools, selectedGrades])

  // Toggle school filter
  const toggleSchoolFilter = (schoolName: string) => {
    setSelectedSchools(prev => 
      prev.includes(schoolName)
        ? prev.filter(s => s !== schoolName)
        : [...prev, schoolName]
    )
  }

  const toggleGradeFilter = (grade: string) => {
    setSelectedGrades(prev => prev.includes(grade)
                                        ? prev.filter(g => g !== grade)
                                        : [...prev, grade])
  }

  // Clear all filters
  const clearFilters = () => {
    setSelectedSchools([])
    setSelectedGrades([])
  }

  const numFilters = selectedSchools.length + selectedGrades.length

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
      <div className="min-h-screen bg-background">
        <div className="w-full p-10 relative">
          <PageHeader
            title="Students"
            icon={Users}
            description="View and manage all student information"
            actions={<AddStudentModal />}
          />

          <div className="mb-3 relative flex flex-wrap gap-y-2">
            <Badge className={`border border-pink text-primary h-full text-md px-5 py-2 
                               ${filterOpen ? "bg-pink text-white" : "bg-background text-primary"}`}
                   onClick={() => setFilterOpen(!filterOpen)}>
              <ListFilter className={`w-4 h-4 mr-2 ${filterOpen ? "text-white" : ""}`} />
              Filter
              {(numFilters > 0) && (
                <div className={`ml-2 flex justify-center items-center w-6 h-6 text-sm rounded-full
                                 ${filterOpen ? "bg-background text-pink" : "bg-pink text-white"}`}>
                  {numFilters}
                </div>
              )}
            </Badge>

              {filterOpen && (
                  <div className="absolute left-0 top-full mt-2 bg-white rounded-xl shadow-lg p-2 max-w-[50%] z-50">
                      <div className="w-full h-full bg-white/55 p-4 w-50 h-50">
                          <h3 className="mb-4">School</h3>
                          {uniqueSchools.length > 0 && (
                              <div className="mb-4 flex flex-wrap items-center gap-2">
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
                                              <Badge className={`${getSchoolColor(schoolName)} text-md h-full px-4 py-0.5`}>
                                                  {schoolName}
                                                  {isSelected && <X className="w-3 h-3 ml-2 inline" />}
                                              </Badge>
                                          </button>
                                      )
                                  })}
                              </div>
                          )}
                          <h3 className="mb-4">Grade Level</h3>
                          {uniqueGrades.length > 0 && (
                              <div className="flex flex-wrap items-center gap-2">
                                  {uniqueGrades.map(grade => {
                                      const isSelected = selectedGrades.includes(grade)
                                      return (
                                          <button key={grade}
                                                  onClick={() => toggleGradeFilter(grade)}
                                                  className={`transition-all ${
                                                      isSelected ? 'opacity-100' : 'opacity-40 hover:opacity-70'
                                                  }`}>
                                              <Badge className={`${gradeColor} ${isSelected ? 'ring-2 ring-accent ring-offset-2' : '' }
                                                         text-md h-full px-4 py-0.5`}>
                                                  Grade {grade}
                                                  {isSelected && <X className="w-3 h-3 ml-1 inline" />}
                                              </Badge>
                                          </button>
                                      )
                                  })}
                              </div>
                          )}
                      </div>
                  </div>
              )}
              <Badge className="mx-4 border border-pink text-primary h-full text-md px-5 py-2"
                     onClick={() => setSortAZ(!sortAZ)}>
                  <ChevronsUpDown className="w-4 h-4 mr-2 text-primary" />
                  Sort {sortAZ ? 'A-Z' : 'Z-A'}
              </Badge>

              {selectedSchools.length > 0 && selectedSchools.map(schoolName => {
                          return (
                              <button
                                  key={schoolName}
                                  onClick={() => toggleSchoolFilter(schoolName)}
                                  className="transition-all opacity-100"
                              >
                                  <Badge className={`${getSchoolColor(schoolName)}
                                                    text-md h-full mx-1.5 text-primary text-md px-5 py-2`}>
                                      {schoolName}
                                      <X className="w-3 h-3 ml-2 inline" />
                                  </Badge>
                              </button>
                          )
              })}
              {selectedGrades.length > 0 && selectedGrades.map(grade => {
                return (
                  <button
                    key={grade}
                    onClick={() => toggleGradeFilter(grade)}
                    className="transition-all opacity-100"
                  >
                      <Badge className={`${gradeColor} text-md h-full mx-1.5 text-primary text-md px-5 py-2`}>
                        Grade {gradeToDisplay(grade)}
                        <X className="w-3 h-3 ml-2 inline" />
                      </Badge>
                  </button>
                )
              })}
          </div>

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
