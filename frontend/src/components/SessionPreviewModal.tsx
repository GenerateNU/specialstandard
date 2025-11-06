'use client'

import type { Session } from '@/lib/api/theSpecialStandardAPI.schemas'
import { Calendar, Clock, MapPin, Users } from 'lucide-react'
import Link from 'next/link'
import { useSessionStudentsForSession } from '@/hooks/useSessionStudents'

interface SessionPreviewModalProps {
  session: Session
  position: { x: number, y: number }
  onClose: () => void
}

export default function SessionPreviewModal({
  session,
  position,
  onClose,
}: SessionPreviewModalProps) {
  const { students } = useSessionStudentsForSession(session.id)

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric',
    })
  }

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true,
    })
  }

  const formatTimeRange = () => {
    return `${formatTime(session.start_datetime)} - ${formatTime(session.end_datetime)}`
  }

  // Smart positioning to stay on screen (only on client-side)
  const MODAL_WIDTH = 360
  const MODAL_HEIGHT = 400
  const SPACING = 10

  let adjustedX = position.x
  let adjustedY = position.y

  // Only calculate positioning if we're on the client side
  if (typeof window !== 'undefined') {
    // Check if modal would go off right edge, if so position to the left
    if (adjustedX + MODAL_WIDTH > window.innerWidth - SPACING) {
      adjustedX = position.x - MODAL_WIDTH - SPACING
    }

    // Ensure it doesn't go off left edge
    if (adjustedX < SPACING) {
      adjustedX = SPACING
    }

    // Check if modal would go off bottom, if so move it up
    if (adjustedY + MODAL_HEIGHT > window.innerHeight - SPACING) {
      adjustedY = window.innerHeight - MODAL_HEIGHT - SPACING
    }

    // Ensure it doesn't go off top
    if (adjustedY < SPACING) {
      adjustedY = SPACING
    }
  }

  return (
    <div
      className="fixed inset-0 z-50"
      onClick={onClose}
    >
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/20" />

      {/* Modal Card - positioned next to the event */}
      <div
        className="absolute bg-white border border-gray-300 rounded-xl shadow-2xl w-[360px] overflow-hidden"
        style={{
          left: `${adjustedX}px`,
          top: `${adjustedY}px`,
        }}
        onClick={e => e.stopPropagation()}
      >
        {/* Content */}
        <div className="p-5 space-y-3">
          {/* Date */}
          <div>
            <div className="text-xs text-gray-500 mb-1 flex items-center gap-1.5">
              <Calendar className="w-3.5 h-3.5" />
              <span className="font-medium">Date</span>
            </div>
            <div className="bg-gray-50 rounded-md px-3 py-2 text-sm text-gray-900">
              {formatDateTime(session.start_datetime)}
            </div>
          </div>

          {/* Time */}
          <div>
            <div className="text-xs text-gray-500 mb-1 flex items-center gap-1.5">
              <Clock className="w-3.5 h-3.5" />
              <span className="font-medium">Time</span>
            </div>
            <div className="bg-gray-50 rounded-md px-3 py-2 text-sm text-gray-900">
              {formatTimeRange()}
            </div>
          </div>

          {/* Location */}
          <div>
            <div className="text-xs text-gray-500 mb-1 flex items-center gap-1.5">
              <MapPin className="w-3.5 h-3.5" />
              <span className="font-medium">Location</span>
            </div>
            <div className="bg-gray-50 rounded-md px-3 py-2 text-sm text-gray-900">
              {session.notes || 'No location specified'}
            </div>
          </div>

          {/* Students */}
          <div>
            <div className="text-xs text-gray-500 mb-1 flex items-center gap-1.5">
              <Users className="w-3.5 h-3.5" />
              <span className="font-medium">Students</span>
            </div>
            <div className="bg-gray-50 rounded-md px-3 py-2 text-sm text-gray-900">
              {students.length === 0
                ? (
                    <span className="text-gray-500">No students assigned</span>
                  )
                : (
                    <>
                      {students.slice(0, 3).map((student, index) => (
                        <div key={student.id || index}>
                          {student.first_name}
                          {' '}
                          {student.last_name}
                        </div>
                      ))}
                      {students.length > 3 && (
                        <div className="text-gray-500 mt-1">
                          and
                          {' '}
                          {students.length - 3}
                          {' '}
                          more
                          {' '}
                          {students.length - 3 === 1 ? 'student' : 'students'}
                        </div>
                      )}
                    </>
                  )}
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="border-t border-gray-200 p-3 bg-gray-50 flex gap-2">
          <Link
            href={`/sessions/${session.id}`}
            className="flex-1 bg-blue hover:bg-blue-hover text-white text-center py-2 px-3 rounded-md text-sm font-medium transition-colors"
            onClick={onClose}
          >
            View Details
          </Link>
          <button
            onClick={onClose}
            className="px-4 py-2 bg-white hover:bg-gray-100 text-gray-700 rounded-md text-sm font-medium transition-colors border border-gray-300"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  )
}
