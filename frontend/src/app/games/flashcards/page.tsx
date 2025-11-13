import FlashcardGameInterface from '@/components/games/FlashcardGameInterface'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Flashcard Game',
  description: 'Practice with interactive flashcards',
}

export default function FlashcardsPage() {
  return <FlashcardGameInterface />
}