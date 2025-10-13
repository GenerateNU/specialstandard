// src/lib/api/students.ts
import type {
  CreateStudentInput,
  Student,
  UpdateStudentInput,
} from '@/types/student'
import apiClient from './apiClient'

export async function fetchStudents(): Promise<Student[]> {
  try {
    // Use a high limit to get all students for now
    const response = await apiClient.get<Student[]>('/api/v1/students?limit=100')
    return response.data
  }
  catch (error) {
    console.error('Error fetching students:', error)
    throw error
  }
}

export async function fetchStudent(studentId: string): Promise<Student | null> {
  try {
    const response = await apiClient.get<Student>(`/api/v1/students/${studentId}`)
    return response.data
  }
  catch (error) {
    console.error('Error fetching student:', error)
    return null
  }
}

export async function createStudent(data: CreateStudentInput): Promise<Student | null> {
  try {
    const response = await apiClient.post<Student>('/api/v1/students', data)
    return response.data
  }
  catch (error) {
    console.error('Error creating student:', error)
    return null
  }
}

export async function updateStudent(
  studentId: string,
  data: UpdateStudentInput,
): Promise<Student | null> {
  try {
    const response = await apiClient.patch<Student>(`/api/v1/students/${studentId}`, data)
    return response.data
  }
  catch (error) {
    console.error('Error updating student:', error)
    return null
  }
}

export async function deleteStudent(studentId: string): Promise<boolean> {
  try {
    await apiClient.delete(`/api/v1/students/${studentId}`)
    return true
  }
  catch (error) {
    console.error('Error deleting student:', error)
    return false
  }
}
