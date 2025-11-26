'use client'

import {useSearchParams} from "next/navigation";

const sessionStudentId = 1
const sessionId = "c35ceea3-fa7d-4d14-a69d-cbed270c737f"
const studentId = "89e2d744-eec1-490e-a335-422ce79eae70"

export default function WordImageMatchingContent() {
  const searchParams = useSearchParams()

  const themeId = searchParams.get('themeId')
  const themeName = searchParams.get('themeName')
  const difficulty = searchParams.get('difficulty')
  const category = searchParams.get('category')
  const questionType = searchParams.get('questionType')

  if (!themeId || !difficulty || !category || !questionType) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <p className="text-error mb-4">
            Missing game parameters. Please select content first.
          </p>
          <a href="/games" className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover
                                      transition-colors inline-block">
            Go Back
          </a>
        </div>
      </div>
    )
  }

  return (
    <ImageMatchingGameInterface
      session_student_id={sessionStudentId}
      session_id={sessionId}
      student_id={studentId}
      themeId={themeId}
      themeName={themeName}
      difficulty={difficulty}
      category={category}
      questionType={questionType}
    />
  )
}