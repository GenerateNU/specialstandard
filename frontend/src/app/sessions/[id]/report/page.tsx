'use client'

import { use } from 'react'
import CurriculumLayout from '@/components/curriculum/CurriculumLayout'
import { useSessionContext } from '@/contexts/sessionContext'
import { useGameResults } from '@/hooks/useGameResults'
import { useSessionStudentsForSession } from '@/hooks/useSessionStudents'
import { Avatar } from '@/components/ui/avatar'
import { getAvatarName, getAvatarVariant } from '@/lib/avatarUtils'
import { GameOverallStats } from '@/components/statistics/GamePerformanceCharts'

interface PageProps {
  params: Promise<{ id: string }>
}

function StudentReport({ 
  sessionStudent,
  sessionId 
}: { 
  sessionStudent: any
  sessionId: string 
}) {
  const { existingResults, isLoadingResults } = useGameResults({
    session_student_id: sessionStudent.session_student_id,
    session_id: sessionId,
    student_id: sessionStudent.id,
  })

  // Calculate stats from existingResults
  const totalGames = existingResults.length
  const completedGames = existingResults.filter(r => r.completed).length
  const completionRate = totalGames > 0 ? (completedGames / totalGames) * 100 : 0
  const totalTime = existingResults.reduce((sum, r) => sum + (r.time_taken_sec || 0), 0)
  const avgTimeSeconds = totalGames > 0 ? totalTime / totalGames : 0
  const totalErrors = existingResults.reduce((sum, r) => sum + (r.count_of_incorrect_attempts || 0), 0)
  const avgErrors = totalGames > 0 ? totalErrors / totalGames : 0

  const stats = {
    totalGames,
    completedGames,
    completionRate,
    avgTimeSeconds,
    totalErrors,
    avgErrors,
  }

  const avatarVariant = getAvatarVariant(sessionStudent.id)

  return (
    <div className="bg-card rounded-3xl p-6 shadow-sm border border-border/50 space-y-6">
      {/* Student Header */}
      <div className="flex items-center gap-3">
        <Avatar
          name={getAvatarName(sessionStudent.first_name, sessionStudent.last_name, sessionStudent.id)}
          variant={avatarVariant}
          className="w-12 h-12"
        />
        <div>
          <h3 className="text-lg font-semibold text-foreground">
            {sessionStudent.first_name} {sessionStudent.last_name}
          </h3>
          <p className="text-sm text-muted-foreground">Grade {sessionStudent.grade}</p>
        </div>
      </div>

      {/* Game Stats */}
      <div>
        <h4 className="text-sm font-medium text-muted-foreground mb-3">Game Performance</h4>
        {isLoadingResults ? (
          <p className="text-sm text-muted-foreground">Loading game results...</p>
        ) : totalGames === 0 ? (
          <p className="text-sm text-muted-foreground">No game data available</p>
        ) : (
          <GameOverallStats stats={stats} />
        )}
      </div>

      {/* Ratings (dawg we have no get endpoint for this i aint saving them in storage no ratings here*/}

      {/* Notes */}
      {sessionStudent.notes && (
        <div>
          <h4 className="text-sm font-medium text-muted-foreground mb-2">Notes</h4>
          <p className="text-sm text-foreground bg-background p-3 rounded-lg border border-border">
            {sessionStudent.notes}
          </p>
        </div>
      )}
    </div>
  )
}

export default function ReportPage({ params }: PageProps) {
  const { id } = use(params)
  const { session } = useSessionContext()
  const { students: sessionStudents, isLoading } = useSessionStudentsForSession(id)

  const sessionDate = session ? new Date(session.start_datetime) : new Date()
  const formattedDate = sessionDate.toLocaleDateString('en-US', {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  })

  if (isLoading) {
    return (
      <CurriculumLayout
        title="Session Report"
        subtitle={formattedDate}
        backHref={`/`}
        backLabel="Back to Home"
      >
        <div className="flex items-center justify-center min-h-[50vh]">
          <div className="text-center">
            <p className="text-muted-foreground">Loading session data...</p>
          </div>
        </div>
      </CurriculumLayout>
    )
  }

  if (sessionStudents.length === 0) {
    return (
      <CurriculumLayout
        title="Session Report"
        subtitle={formattedDate}
        backHref={`/`}
        backLabel="Back to Home"
      >
        <div className="flex items-center justify-center min-h-[50vh]">
          <div className="text-center">
            <p className="text-muted-foreground">No students in this session</p>
          </div>
        </div>
      </CurriculumLayout>
    )
  }

  return (
    <CurriculumLayout
      title="Session Report"
      subtitle={formattedDate}
      backHref={`/sessions/${id}/curriculum`}
      backLabel="Back to Curriculum"
    >
      <div className="space-y-6 max-w-6xl mx-auto py-6">
        {sessionStudents.map((sessionStudent) => (
          <StudentReport
            key={sessionStudent.session_student_id}
            sessionStudent={sessionStudent}
            sessionId={id}
          />
        ))}
      </div>
    </CurriculumLayout>
  )
}