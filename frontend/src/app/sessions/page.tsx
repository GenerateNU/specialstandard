// app/sessions/page.tsx (App Router) or pages/sessions.tsx (Pages Router)
'use client'

import type { Session } from '@/types/session'
import { ArrowLeft, Calendar, Plus } from 'lucide-react'
import Link from 'next/link'
import { useState } from 'react'
import SessionCalendar from '@/components/sessionCalendar'
import SessionModal from '@/components/sessionModal'
import { useSessions } from '@/hooks/useSessions'

export default function SessionsPage() {
  const { sessions, isLoading, error, refetch } = useSessions()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [selectedSession, setSelectedSession] = useState<Session | null>(null)

  const handleNewSession = () => {
    setSelectedSession(null)
    setIsModalOpen(true)
  }

  const handleEditSession = (session: Session) => {
    setSelectedSession(session)
    setIsModalOpen(true)
  }

  return (
    <div className="min-h-screen bg-gray-50 p-4">
      <div className="max-w-7xl mx-auto">
        <div className="bg-card rounded-xl shadow-lg p-6">
          {/* Header */}
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-4">
              <Link
                href="/"
                className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <ArrowLeft className="w-5 h-5" />
              </Link>
              <h1 className="text-3xl font-bold text-primary flex items-center gap-2">
                <Calendar className="w-8 h-8 text-accent" />
                Therapy Sessions
              </h1>
            </div>
            <button
              onClick={handleNewSession}
              className="flex items-center gap-2 bg-accent hover:bg-accent-dark text-white px-4 py-2 rounded-lg transition-colors"
            >
              <Plus className="w-5 h-5" />
              New Session
            </button>
          </div>

          {/* Error Display */}
          {error && (
            <div className="mb-4 p-4 bg-red-100 border border-red-300 text-red-700 rounded-lg">
              {error}
            </div>
          )}

          {/* Loading State */}
          {isLoading
            ? (
                <div className="flex justify-center items-center h-64">
                  <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-accent"></div>
                </div>
              )
            : (
                <SessionCalendar
                  sessions={sessions}
                  onSessionClick={handleEditSession}
                  onDateClick={handleNewSession}
                />
              )}
        </div>
      </div>

      {/* Session Modal */}
      {isModalOpen && (
        <SessionModal
          session={selectedSession}
          isOpen={isModalOpen}
          onClose={() => {
            setIsModalOpen(false)
            setSelectedSession(null)
          }}
          onSuccess={refetch}
        />
      )}
    </div>
  )
}
