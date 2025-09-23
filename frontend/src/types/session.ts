// src/types/session.ts
export interface Session {
  id: string
  start_datetime: string
  end_datetime: string
  therapist_id: string
  notes?: string | null
  created_at?: string | null
  updated_at?: string | null
}

export interface CreateSessionInput {
  start_datetime: string
  end_datetime: string
  therapist_id: string
  notes?: string | null
}

export interface UpdateSessionInput {
  start_datetime?: string
  end_datetime?: string
  therapist_id?: string
  notes?: string | null
}
