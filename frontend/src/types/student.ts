// src/types/student.ts
export interface Student {
  id: string
  first_name: string
  last_name: string
  dob?: string
  therapist_id: string
  grade?: string
  iep?: string
  created_at: string
  updated_at: string
}

export interface CreateStudentInput {
  first_name: string
  last_name: string
  dob?: string
  therapist_id: string
  grade?: string
  iep?: string
}

export interface UpdateStudentInput {
  first_name?: string
  last_name?: string
  dob?: string
  therapist_id?: string
  grade?: string
  iep?: string
}
