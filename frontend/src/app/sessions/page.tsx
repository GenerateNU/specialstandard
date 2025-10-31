'use client'

import Link from 'next/link'
import { useSessions } from '@/hooks/useSessions'

export default function SessionsListPage() {
  const { sessions, isLoading, error } = useSessions()

  if (isLoading) {
    return <div className="p-8">Loading sessions...</div>
  }

  if (error) {
    return (
      <div className="p-8 text-error">
        Error:
        {' '}
        {error}
      </div>
    )
  }

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6">All Sessions</h1>

      {sessions.length === 0
        ? (
            <p>No sessions found</p>
          )
        : (
            <div className="space-y-4">
              {sessions.map((session, index) => (
                <Link
                  key={session.id || `session-${index}`}
                  href={`/sessions/${session.id}`}
                  className="block p-4 border border-default rounded-lg hover:bg-card-hover transition-colors"
                >
                  <div className="font-semibold">
                    Session ID:
                    {' '}
                    {session.id}
                  </div>
                  <div className="text-sm text-secondary">
                    Start:
                    {' '}
                    {new Date(session.start_datetime).toLocaleString()}
                  </div>
                  <div className="text-sm text-secondary">
                    End:
                    {' '}
                    {new Date(session.end_datetime).toLocaleString()}
                  </div>
                  {session.notes && (
                    <div className="text-sm text-secondary mt-2">
                      Notes:
                      {' '}
                      {session.notes}
                    </div>
                  )}
                </Link>
              ))}
            </div>
          )}
    </div>
  )
}
