'use client'

import React, { useEffect, useState } from 'react'
import { Check, RotateCw } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useGameContents } from '@/hooks/useGameContents'
import type { 
  GetGameContentsQuestionType,
  Theme
} from '@/lib/api/theSpecialStandardAPI.schemas'
import { 
  GetGameContentsCategory
} from '@/lib/api/theSpecialStandardAPI.schemas'
import './flashcard.css'
import { useGameResults } from '@/hooks/useGameResults'

const CATEGORIES = {
  [GetGameContentsCategory.receptive_language]: { 
    label: 'Receptive Language', 
    icon: '👂', 
    colorClass: 'bg-blue',
  },
  [GetGameContentsCategory.expressive_language]: { 
    label: 'Expressive Language', 
    icon: '💬', 
    colorClass: 'bg-pink',
  },
  [GetGameContentsCategory.social_pragmatic_language]: { 
    label: 'Social Pragmatic Language', 
    icon: '🤝', 
    colorClass: 'bg-orange',
  },
  [GetGameContentsCategory.speech]: { 
    label: 'Speech', 
    icon: '🗣️', 
    colorClass: 'bg-blue',
  },
}

interface FlashcardProps {
  question: string
  answer: string
  isFlipped: boolean
  onFlip: () => void
  onMarkCorrect?: () => void
  timeTaken?: number
}

interface FlashcardGameInterfaceProps {
  session_student_id?: number
  session_id?: string
  student_id?: string
  themeId: string
  difficulty: number
  category: string
  questionType: string
}

const Flashcard: React.FC<FlashcardProps> = ({ 
  question, 
  answer, 
  isFlipped, 
  onFlip,
  onMarkCorrect,
  timeTaken
}) => {
  return (
    <div 
      className="relative w-full h-64 cursor-pointer preserve-3d transition-transform duration-700"
      style={{ transformStyle: 'preserve-3d', transform: isFlipped ? 'rotateY(180deg)' : '' }}
      onClick={onFlip}
    >
      {/* Front of card */}
      <div 
        className="absolute inset-0 w-full h-full bg-card rounded-xl shadow-lg p-8 flex items-center justify-center backface-hidden border-2 border-default"
        style={{ backfaceVisibility: 'hidden' }}
      >
        <div className="text-center">
          <p className="text-muted text-sm mb-2">Question</p>
          <p className="text-xl font-semibold text-primary">{question}</p>
        </div>
      </div>
      
      {/* Back of card */}
      <div 
        className="absolute inset-0 w-full h-full bg-blue rounded-xl shadow-lg p-8 flex flex-col items-center justify-center backface-hidden"
        style={{ backfaceVisibility: 'hidden', transform: 'rotateY(180deg)' }}
      >
        <div className="text-center flex-1 flex items-center">
          <p className="text-2xl font-bold text-white">{answer}</p>
        </div>
        
        {onMarkCorrect && (
          <div className="mt-4 flex items-center gap-3 text-white/70 text-sm">
            {timeTaken !== undefined && (
              <span>{timeTaken}s</span>
            )}
            <button
              onClick={(e) => {
                e.stopPropagation()
                onMarkCorrect()
              }}
              className="p-2 hover:bg-white/20 rounded-full transition-colors"
              title="Mark as correct"
            >
              <Check className="w-5 h-5" />
            </button>
          </div>
        )}
      </div>
    </div>
  )
}

export default function FlashcardGameInterface({ 
  session_student_id, 
  session_id, 
  student_id,
  themeId,
  difficulty,
  category,
  questionType
}: FlashcardGameInterfaceProps) {
  const router = useRouter()
  const [currentCardIndex, setCurrentCardIndex] = useState(0)
  const [flippedCards, setFlippedCards] = useState<Set<number>>(new Set())
  const [cardStartTime, setCardStartTime] = useState<number | null>(null)
  const [timeTaken, setTimeTaken] = useState(0)
  //const [setThemeName] = useState<string>('')

  const { gameContents, isLoading: contentsLoading, error: contentsError } = useGameContents({
    theme_id: themeId,
    category: category as GetGameContentsCategory,
    question_type: questionType as GetGameContentsQuestionType,
    difficulty_level: difficulty,
    question_count: 10,
  })

  const gameResultsHook = session_student_id ? useGameResults({
    session_student_id,
    session_id,
    student_id
  }) : null

  const startCard = gameResultsHook?.startCard

  // Initialize timer when card loads
  useEffect(() => {
    if (gameContents[currentCardIndex]) {
      setCardStartTime(Date.now())
      setFlippedCards(new Set())
      setTimeTaken(0)
      startCard?.(gameContents[currentCardIndex])
    }
  }, [currentCardIndex, gameContents, startCard])

  // Update timer display
  useEffect(() => {
    if (cardStartTime === null) return
    
    const interval = setInterval(() => {
      setTimeTaken(Math.floor((Date.now() - cardStartTime) / 1000))
    }, 100)
    
    return () => clearInterval(interval)
  }, [cardStartTime])

  // Fetch theme name once
  useEffect(() => {
    // In a real app, you'd fetch this from the API
    // For now, we'll set a default
    //setThemeName('Theme')
  }, [themeId])

  const handleMarkCorrect = () => {
    if (!gameResultsHook || !gameContents[currentCardIndex]) return
    
    const finalTime = Math.floor((Date.now() - (cardStartTime || Date.now())) / 1000)
    gameResultsHook.completeCard(gameContents[currentCardIndex].id, finalTime)
    
    setCurrentCardIndex(prev => prev + 1)
    setCardStartTime(null)
  }

  if (contentsLoading) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue mx-auto mb-4"></div>
          <p className="text-secondary">Loading flashcards...</p>
        </div>
      </div>
    )
  }

  if (contentsError || gameContents.length === 0) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <p className="text-error mb-4">
            {contentsError ? 'Failed to load flashcards' : 'No flashcards available'}
          </p>
          <button 
            onClick={() => router.back()}
            className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover transition-colors"
          >
            Go Back
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-8">
          <button
            onClick={() => router.back()}
            className="text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
          >
            ← Back
          </button>
          <button
            onClick={() => {
              setCurrentCardIndex(0)
              setFlippedCards(new Set())
              setCardStartTime(null)
              setTimeTaken(0)
            }}
            className="flex items-center gap-2 text-secondary hover:text-primary transition-colors"
          >
            <RotateCw className="w-4 h-4" />
            Reset
          </button>
        </div>

        <h1 className="mb-4">Flashcards</h1>
        
        {/* Display selected options */}
        <div className="bg-card rounded-lg p-4 mb-6 border border-default">
          <div className="flex flex-wrap items-center gap-2 text-sm">
            <span className="text-muted">Difficulty:</span>
            <span className="font-medium text-primary">Level {difficulty}</span>
            <span className="text-muted mx-2">•</span>
            <span className="text-muted">Category:</span>
            <span className="font-medium text-primary">{CATEGORIES[category as GetGameContentsCategory]?.label || category}</span>
          </div>
        </div>
        
        {/* Progress indicator */}
        <div className="mb-8">
          <div className="flex items-center justify-between text-sm text-secondary mb-2">
            <span>Card {currentCardIndex + 1} of {gameContents.length}</span>
            <span>{Math.round((currentCardIndex / gameContents.length) * 100)}% Complete</span>
          </div>
          <div className="w-full bg-card rounded-full h-2 border border-default">
            <div 
              className="bg-blue h-full rounded-full transition-all duration-300"
              style={{ width: `${(currentCardIndex / gameContents.length) * 100}%` }}
            />
          </div>
        </div>

        {/* Current flashcard */}
        {gameContents[currentCardIndex] && (
          <div className="mb-8">
            <Flashcard
              question={gameContents[currentCardIndex].question}
              answer={gameContents[currentCardIndex].answer}
              isFlipped={flippedCards.has(currentCardIndex)}
              onFlip={() => {
                setFlippedCards(prev => {
                  const newSet = new Set(prev)
                  if (newSet.has(currentCardIndex)) {
                    newSet.delete(currentCardIndex)
                  } else {
                    newSet.add(currentCardIndex)
                  }
                  return newSet
                })
              }}
              onMarkCorrect={gameResultsHook ? handleMarkCorrect : undefined}
              timeTaken={flippedCards.has(currentCardIndex) ? timeTaken : undefined}
            />
            <p className="text-center text-muted mt-4 text-sm">Click card to flip</p>
          </div>
        )}

        {/* Navigation buttons */}
        <div className="flex items-center justify-between">
          <button
            onClick={() => setCurrentCardIndex(prev => Math.max(0, prev - 1))}
            disabled={currentCardIndex === 0}
            className={`px-6 py-2 rounded-lg font-medium transition-colors ${
              currentCardIndex === 0
                ? 'bg-card text-disabled cursor-not-allowed border border-default'
                : 'bg-card border border-default text-primary hover:bg-card-hover hover:border-hover'
            }`}
          >
            Previous
          </button>

          <div className="flex gap-2">
            {gameContents.map((_, index) => (
              <button
                key={index}
                onClick={() => setCurrentCardIndex(index)}
                className={`w-2 h-2 rounded-full transition-colors ${
                  index === currentCardIndex ? 'bg-blue' : 'bg-card border border-default'
                }`}
              />
            ))}
          </div>

          <button
            onClick={() => {
              if (currentCardIndex < gameContents.length - 1) {
                setCurrentCardIndex(prev => prev + 1)
              }
            }}
            disabled={currentCardIndex === gameContents.length - 1}
            className={`px-6 py-2 rounded-lg font-medium transition-colors ${
              currentCardIndex === gameContents.length - 1
                ? 'bg-blue-disabled text-white cursor-not-allowed'
                : 'bg-blue text-white hover:bg-blue-hover'
            }`}
          >
            Next
          </button>
        </div>

        {currentCardIndex >= gameContents.length && (
          <div className="mt-8 p-6 bg-blue-light border border-blue rounded-lg text-center">
            <p className="text-blue font-semibold">Great job! You've completed all flashcards!</p>
            <div className="mt-4 flex gap-3 justify-center">
              <button
                onClick={() => {
                  setCurrentCardIndex(0)
                  setFlippedCards(new Set())
                  setCardStartTime(null)
                  setTimeTaken(0)
                }}
                className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover transition-colors"
              >
                Start Over
              </button>
              
              {gameResultsHook && (
                <button
                  onClick={() => gameResultsHook.saveAllResults()}
                  disabled={gameResultsHook.isSaving}
                  className="px-6 py-2 bg-pink text-white rounded-lg hover:bg-pink-hover transition-colors disabled:bg-pink-disabled"
                >
                  {gameResultsHook.isSaving ? 'Saving...' : 'Save Progress'}
                </button>
              )}
            </div>
            
            {gameResultsHook?.saveError && (
              <p className="text-error text-sm mt-2">
                Failed to save progress. Please try again.
              </p>
            )}
          </div>
        )}

        {gameResultsHook && gameContents[currentCardIndex] && (
          <div className="text-center mt-2">
            {(() => {
              const existingResult = gameResultsHook.getResultForContent(gameContents[currentCardIndex].id)
              if (existingResult?.completed) {
                return (
                  <span className="text-success text-sm">
                    ✓ Previously completed in {existingResult.time_taken_sec}s
                  </span>
                )
              }
              return null
            })()}
          </div>
        )}
      </div>
    </div>
  )
}