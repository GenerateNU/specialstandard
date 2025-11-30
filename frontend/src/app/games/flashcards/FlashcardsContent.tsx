'use client'

import { useSearchParams } from 'next/navigation'
import FlashcardGameInterface from '@/components/games/FlashcardGameInterface'


export function FlashcardsContent() {
  const searchParams = useSearchParams()
  const sessionStudentId = Number.parseInt(searchParams.get("sessionStudentId") ?? "0");
  const themeId = searchParams.get('themeId')
  const themeName = searchParams.get('themeName')
  const difficulty = searchParams.get('difficulty')
  const category = searchParams.get('category')
  const questionType = searchParams.get('questionType')
  const sessionId = searchParams.get('sessionId') || '00000000-0000-0000-0000-000000000000'

  if (!themeId || !difficulty || !category || !questionType) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <p className="text-error mb-4">Missing game parameters. Please select content first.</p>
          <a href="/games" className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover transition-colors inline-block">
            Go Back
          </a>
        </div>
      </div>
    )
  }

  return (
    <FlashcardGameInterface
      session_student_id={sessionStudentId}
      session_id={sessionId}
      themeId={themeId}
      themeName={themeName || 'Theme'}
      difficulty={Number.parseInt(difficulty)}
      category={category}
      questionType={questionType}
    />
  )
}