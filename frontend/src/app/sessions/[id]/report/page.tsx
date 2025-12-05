'use client'

import { use } from 'react'
import CurriculumLayout from '@/components/curriculum/CurriculumLayout'

interface PageProps {
  params: Promise<{ id: string }>
}

export default function ReportPage({ params }: PageProps) {
  const { id } = use(params)

  return (
    <CurriculumLayout
      title="Session Report"
      backHref={`/sessions/${id}/curriculum`}
      backLabel="Back to Curriculum"
    >
      <div className="flex items-center justify-center min-h-[50vh]">
        <div className="text-center">
          <h2 className="text-2xl font-bold mb-4">Session Report</h2>
          <p className="text-muted-foreground">Report view for session {id} is coming soon.</p>
        </div>
      </div>
    </CurriculumLayout>
  )
}
