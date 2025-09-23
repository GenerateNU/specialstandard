// src/components/SessionModal.tsx
'use client'

import type { Session } from '@/types/session'
import { X } from 'lucide-react'
import { useState } from 'react'

interface SessionModalProps {
  session: Session | null
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
}

export default function SessionModal({
  session,
  isOpen,
  onClose,
}: SessionModalProps) {
  const [formData, setFormData] = useState({
    start_datetime: session ? new Date(session.start_datetime).toISOString().slice(0, 16) : '',
    end_datetime: session ? new Date(session.end_datetime).toISOString().slice(0, 16) : '',
    therapist_id: session?.therapist_id || '',
    notes: session?.notes || '',
  })

  if (!isOpen)
    return null

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({ ...prev, [name]: value }))
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-card rounded-lg p-6 w-full max-w-md">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold text-primary">
            {session ? 'Edit Session' : 'New Session'}
          </h2>
          <button
            onClick={onClose}
            className="text-secondary hover:text-primary transition-colors"
            title="Close"
            aria-label="Close"
          >
            <X className="w-6 h-6" />
          </button>
        </div>

        <form className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-secondary mb-1">
              Start Date & Time *
            </label>
            <input
              type="datetime-local"
              name="start_datetime"
              value={formData.start_datetime}
              onChange={handleChange}
              required
              placeholder="Select start date and time"
              title="Start Date & Time"
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-secondary mb-1">
              End Date & Time *
            </label>
            <input
              type="datetime-local"
              name="end_datetime"
              value={formData.end_datetime}
              onChange={handleChange}
              required
              title="End Date & Time"
              placeholder="Select end date and time"
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-secondary mb-1">
              Therapist ID *
            </label>
            <input
              type="text"
              name="therapist_id"
              value={formData.therapist_id}
              onChange={handleChange}
              required
              placeholder="e.g., 123e4567-e89b-12d3-a456-426614174000"
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-secondary mb-1">
              Notes
            </label>
            <textarea
              name="notes"
              value={formData.notes}
              onChange={handleChange}
              rows={3}
              placeholder="Session notes..."
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
            />
          </div>

        </form>
      </div>
    </div>
  )
}
