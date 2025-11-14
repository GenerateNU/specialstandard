import FlashcardGameInterface from '@/components/games/FlashcardGameInterface'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Flashcard Game',
  description: 'Practice with interactive flashcards',
}

// For testing, use hardcoded values
const sessionStudentId = 1 // Replace with real value
const sessionId = "c35ceea3-fa7d-4d14-a69d-cbed270c737f" // Optional
const studentId = "89e2d744-eec1-490e-a335-422ce79eae70" // Optional

export default function FlashcardsPage() {
  return <FlashcardGameInterface
      session_student_id={sessionStudentId}
      session_id={sessionId}
      student_id={studentId}
    />
}