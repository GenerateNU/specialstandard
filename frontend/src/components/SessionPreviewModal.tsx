'use client'

import type { Session } from '@/lib/api/theSpecialStandardAPI.schemas'
import { Calendar, Clock, MapPin, Maximize2, Repeat, Users, X } from 'lucide-react'
import Link from 'next/link'
import { useSessionStudentsForSession } from '@/hooks/useSessionStudents'
import { formatRecurrence } from '@/hooks/useSessions'
import { Avatar } from '@/components/ui/avatar'
import { getAvatarName, getAvatarVariant } from '@/lib/avatarUtils'

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
  const isRecurring = !!session.repetition

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
  const MODAL_WIDTH = 520
  const MODAL_HEIGHT = 320
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
        className="absolute bg-white border-4 border-pink-500 rounded-3xl shadow-2xl w-[520px] overflow-hidden"
        style={{
          left: `${adjustedX}px`,
          top: `${adjustedY}px`,
        }}
        onClick={e => e.stopPropagation()}
      >
        {/* Header with buttons */}
        <div className="absolute top-4 left-4 right-4 flex justify-between items-center z-10">
          <Link
            href={`/sessions/${session.id}`}
            className="bg-white hover:bg-gray-100 rounded-lg p-2 transition-colors"
            onClick={onClose}
            title="View full details"
          >
            <Maximize2 className="w-5 h-5 text-pink-500" />
          </Link>
          <button
            onClick={onClose}
            className="bg-white hover:bg-gray-100 rounded-full p-2 transition-colors"
            aria-label="Close"
          >
            <X className="w-5 h-5 text-pink-500" />
          </button>
        </div>

        {/* Content - Two column layout */}
        <div className="p-6 pt-12 flex gap-6">
          {/* Left Column - Session Details */}
          <div className="flex-1 space-y-4">
            {/* Session Name */}
            <div>
              <h2 className="text-lg font-bold text-gray-900">
                {session.session_name}
              </h2>
              <p className="text-sm text-gray-600">
                {session.location || 'School 3'}
              </p>
            </div>

            {/* Date */}
            <div className="flex items-start gap-2">
              <Calendar className="w-5 h-5 text-pink-500 flex-shrink-0 mt-0.5" />
              <div>
                <p className="text-sm text-gray-900 font-medium">
                  {formatDateTime(session.start_datetime)}
                </p>
              </div>
            </div>

            {/* Time */}
            <div className="flex items-start gap-2">
              <Clock className="w-5 h-5 text-pink-500 flex-shrink-0 mt-0.5" />
              <div>
                <p className="text-sm text-gray-900 font-medium">
                  {formatTimeRange()}
                </p>
              </div>
            </div>

            {/* Location */}
            {session.location && (
              <div className="flex items-start gap-2">
                <MapPin className="w-5 h-5 text-pink-500 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-sm text-gray-900 font-medium">
                    {session.location}
                  </p>
                </div>
              </div>
            )}

            {/* Recurring Info */}
            {isRecurring && session.repetition && (
              <div className="flex items-start gap-2">
                <Repeat className="w-5 h-5 text-pink-500 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-sm text-gray-900 font-medium">
                    {formatRecurrence(session.repetition)}
                  </p>
                </div>
              </div>
            )}
          </div>

          {/* Right Column - Students */}
          <div className="flex flex-col items-center">
            <p className="text-xs font-semibold text-gray-600 mb-3 uppercase tracking-wide">
              Students
            </p>
            <div className="space-y-2">
              {students.length === 0
                ? (
                    <p className="text-xs text-gray-500 text-center">
                      No students assigned
                    </p>
                  )
                : (
                    students.slice(0, 4).map((student) => (
                      <Avatar
                        key={student.id}
                        name={getAvatarName(
                          student.first_name || 'Unknown',
                          student.last_name || 'Student',
                          student.id,
                        )}
                        variant={getAvatarVariant(student.id)}
                        className="w-10 h-10 ring-2 ring-pink-200"
                      />
                    ))
                  )}
              {students.length > 4 && (
                <div className="text-xs text-gray-500 text-center mt-1">
                  +{students.length - 4} more
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}