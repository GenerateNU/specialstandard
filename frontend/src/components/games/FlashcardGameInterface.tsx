'use client'

import React, { useEffect, useState } from 'react'
import { Check, RotateCw, X } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useGameContents } from '@/hooks/useGameContents'
import type { 
  GetGameContentsQuestionType
} from '@/lib/api/theSpecialStandardAPI.schemas'
import { 
  GetGameContentsCategory
} from '@/lib/api/theSpecialStandardAPI.schemas'
import './flashcard.css'
import { useGameResults } from '@/hooks/useGameResults'
import { useSessionContext } from '@/contexts/sessionContext'
import { useStudents } from '@/hooks/useStudents'

export const CATEGORIES = {
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
  onMarkIncorrect?: () => void
  timeTaken?: number
  incorrectAttempts?: number
}

interface FlashcardGameInterfaceProps {
  session_student_ids: number[]
  session_id?: string
  themeId: string,
  themeName: string,
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
  onMarkIncorrect,
  timeTaken,
  incorrectAttempts
}) => {
  return (
    <div 
      className="relative w-full h-64 cursor-pointer preserve-3d transition-transform duration-700"
      style={{ transformStyle: 'preserve-3d', transform: isFlipped ? 'rotateY(180deg)' : '' }}
      onClick={onFlip}
    >
      {/* Front of card - Shows Image */}
      <div 
        className="absolute inset-0 w-full h-full bg-card rounded-xl shadow-lg p-8 flex items-center justify-center backface-hidden border-2 border-default overflow-hidden"
        
        style={{ backfaceVisibility: 'hidden' }}
      >
        <div className="text-center flex items-center justify-center w-full h-full">
          {answer.startsWith('http') || answer.startsWith('https') ? (
            <div className="w-full h-full flex items-center justify-center">
              <img 
                src={answer} 
                alt="Flashcard image"
                className="max-w-full max-h-full object-contain rounded-lg"
                style={{ maxHeight: '200px' }}
                onError={(e) => {
                  // Fallback to text if image fails to load
                  e.currentTarget.style.display = 'none'
                  const parent = e.currentTarget.parentElement
                  if (parent) {
                    const text = document.createElement('p')
                    text.className = 'text-2xl font-bold text-primary'
                    text.textContent = answer
                    parent.appendChild(text)
                  }
                }}
              />
            </div>
          ) : (
            <p className="text-2xl font-bold text-primary">{answer}</p>
          )}
        </div>
      </div>
      
      {/* Back of card - Shows Word/Question */}
      <div 
        className="absolute inset-0 w-full h-full bg-blue rounded-xl shadow-lg p-4 flex flex-col items-center justify-between backface-hidden overflow-hidden"
        style={{ backfaceVisibility: 'hidden', transform: 'rotateY(180deg)' }}
      >
        <div className="text-center flex items-center justify-center w-full flex-1">
          <div>
            <p className="text-white/70 text-sm mb-2">Answer</p>
            <p className="text-2xl font-semibold text-white">{question}</p>
          </div>
        </div>
        
        {(onMarkCorrect || onMarkIncorrect) && (
          <div className="space-y-2">
            <div className="flex items-center justify-center gap-2 text-white/70 text-sm">
              {timeTaken !== undefined && (
                <span>{timeTaken}s</span>
              )}
              {incorrectAttempts !== undefined && incorrectAttempts > 0 && (
                <span className="text-red-300">• {incorrectAttempts} incorrect</span>
              )}
            </div>
            <div className="flex items-center justify-center gap-4">
              {onMarkIncorrect && (
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    onMarkIncorrect()
                  }}
                  className="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors flex items-center gap-2"
                  title="Mark as incorrect"
                >
                  <X className="w-4 h-4" />
                  Incorrect
                </button>
              )}
              {onMarkCorrect && (
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    onMarkCorrect()
                  }}
                  className="px-4 py-2 bg-green-500 hover:bg-green-600 text-white rounded-lg transition-colors flex items-center gap-2"
                  title="Mark as correct"
                >
                  <Check className="w-4 h-4" />
                  Correct
                </button>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default function FlashcardGameInterface({ 
  session_student_ids, 
  session_id, 
  themeId,
  themeName,
  difficulty,
  category,
  questionType
}: FlashcardGameInterfaceProps) {
  const router = useRouter()
  const [currentCardIndex, setCurrentCardIndex] = useState(0)
  const [flippedCards, setFlippedCards] = useState<Set<number>>(new Set())
  const [cardStartTime, setCardStartTime] = useState<number | null>(null)
  const [timeTaken, setTimeTaken] = useState(0)
  const [resultsSaved, setResultsSaved] = useState(false)
  const [incorrectAttempts, setIncorrectAttempts] = useState<Map<number, number>>(new Map())

  const { gameContents, isLoading: contentsLoading, error: contentsError } = useGameContents({
    theme_id: themeId,
    category: category as GetGameContentsCategory,
    question_type: questionType as GetGameContentsQuestionType,
    difficulty_level: difficulty,
    question_count: 10,
  })

  // Calculate questions per student and limit total cards
  // Ensure at least 1 question per student, or use all available if fewer questions than students
  const questionsPerStudent = Math.max(1, Math.floor(gameContents.length / session_student_ids.length))
  const totalQuestionsToUse = Math.min(questionsPerStudent * session_student_ids.length, gameContents.length)
  const limitedGameContents = gameContents.length > 0 ? gameContents.slice(0, totalQuestionsToUse) : gameContents
  
  // Get current student based on card index
  const currentStudentIndex = currentCardIndex % session_student_ids.length
  const currentSessionStudentId = session_student_ids[currentStudentIndex]

  // Create game results hooks for each student
  const gameResultsHooks = session_student_ids.map(studentId => 
    useGameResults({
      session_student_id: studentId,
      session_id,
    })
  )

  const currentGameResultsHook = gameResultsHooks[currentStudentIndex]
  const startCard = currentGameResultsHook?.startCard

  // Get student names from session context
  const { students: sessionStudents, session } = useSessionContext()
  const { students: allStudents } = useStudents()
  
  // Prefer session context ID over prop
  const effectiveSessionId = session?.id || session_id
  
  const getStudentName = (sessionStudentId: number) => {
    const sessionStudent = sessionStudents.find(s => s.sessionStudentId === sessionStudentId)
    if (!sessionStudent) return 'Student'
    const student = allStudents?.find(s => s.id === sessionStudent.studentId)
    return student ? `${student.first_name} ${student.last_name}` : 'Student'
  }

  // Initialize timer when card loads
  useEffect(() => {
    if (limitedGameContents[currentCardIndex]) {
      setCardStartTime(Date.now())
      setFlippedCards(new Set())
      setTimeTaken(0)
      startCard?.(limitedGameContents[currentCardIndex])
    }
  }, [currentCardIndex, limitedGameContents.length, startCard])

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
    if (!currentGameResultsHook || !limitedGameContents[currentCardIndex]) return
    
    const finalTime = Math.floor((Date.now() - (cardStartTime || Date.now())) / 1000)
    const attempts = incorrectAttempts.get(currentCardIndex) || 0
    
    currentGameResultsHook.completeCard(limitedGameContents[currentCardIndex].id, finalTime, attempts)
    
    // Unflip the card first, then advance after a small delay for smooth transition
    setFlippedCards(prev => {
      const newSet = new Set(prev)
      newSet.delete(currentCardIndex)
      return newSet
    })
    
    setTimeout(() => {
      setCurrentCardIndex(prev => prev + 1)
      setCardStartTime(null)
    }, 300) // Wait for flip animation
  }

  const handleMarkIncorrect = () => {
    if (!limitedGameContents[currentCardIndex]) return
    
    // Increment incorrect attempts for this card
    setIncorrectAttempts(prev => {
      const newMap = new Map(prev)
      const current = newMap.get(currentCardIndex) || 0
      newMap.set(currentCardIndex, current + 1)
      return newMap
    })
    
    // Flip the card back to show the image (front side)
    setFlippedCards(prev => {
      const newSet = new Set(prev)
      newSet.delete(currentCardIndex)
      return newSet
    })
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

  if (contentsError || limitedGameContents.length === 0) {
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
          <div className="flex items-center gap-4">
            <button
              onClick={() => router.push(effectiveSessionId ? `/sessions/${effectiveSessionId}/curriculum` : '/games')}
              className="text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
            >
              ← Back to Content
            </button>
          </div>
          <button
            onClick={() => {
              setCurrentCardIndex(0)
              setFlippedCards(new Set())
              setCardStartTime(null)
              setTimeTaken(0)
              setIncorrectAttempts(new Map())
            }}
            className="flex items-center gap-2 text-secondary hover:text-primary transition-colors"
          >
            <RotateCw className="w-4 h-4" />
            Reset
          </button>
        </div>

        <h1 className="mb-4">Flashcards</h1>
        
        {/* Current Student Banner */}
        <div className="bg-blue text-white rounded-lg p-4 mb-6 text-center">
          <p className="text-sm opacity-90 mb-1">Current Player</p>
          <p className="text-2xl font-bold">{getStudentName(currentSessionStudentId)}</p>
        </div>

        {/* Display selected options */}
        <div className="bg-card rounded-lg p-4 mb-6 border border-default">
          <div className="flex flex-wrap items-center gap-2 text-sm">
            <span className="text-muted">Theme:</span>
            <span className="font-medium text-primary">{themeName}</span>
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
            <span>Card {Math.min(currentCardIndex + 1, limitedGameContents.length)} of {limitedGameContents.length}</span>
            <span>{Math.round((Math.min(currentCardIndex + 1, limitedGameContents.length) / limitedGameContents.length) * 100)}% Complete</span>
          </div>
          <div className="w-full bg-card rounded-full h-2 border border-default">
            <div 
              className="bg-blue h-full rounded-full transition-all duration-300"
              style={{ width: `${(Math.min(currentCardIndex + 1, limitedGameContents.length) / limitedGameContents.length) * 100}%` }}
            />
          </div>
        </div>

        {/* Current flashcard */}
        {limitedGameContents[currentCardIndex] && (
          <div className="mb-8">
            <Flashcard
              question={limitedGameContents[currentCardIndex].question}
              answer={limitedGameContents[currentCardIndex].answer}
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
              onMarkCorrect={currentGameResultsHook ? handleMarkCorrect : undefined}
              onMarkIncorrect={currentGameResultsHook ? handleMarkIncorrect : undefined}
              timeTaken={flippedCards.has(currentCardIndex) ? timeTaken : undefined}
              incorrectAttempts={incorrectAttempts.get(currentCardIndex) || 0}
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
            {limitedGameContents.map((_, index) => (
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
              if (currentCardIndex < limitedGameContents.length) {
                setCurrentCardIndex(prev => prev + 1)
              }
            }}
            disabled={currentCardIndex >= limitedGameContents.length}
            className={`px-6 py-2 rounded-lg font-medium transition-colors ${
              currentCardIndex >= limitedGameContents.length
                ? 'bg-blue-disabled text-white cursor-not-allowed'
                : 'bg-blue text-white hover:bg-blue-hover'
            }`}
          >
            Next
          </button>
        </div>

        {currentCardIndex >= limitedGameContents.length && (
          <div className="mt-8 p-6 bg-blue-light border border-blue rounded-lg text-center">
            <p className="text-blue font-semibold">Great job! You've completed all flashcards!</p>
            <div className="mt-4 flex gap-3 justify-center">
              <button
                onClick={() => {
                  setCurrentCardIndex(0)
                  setFlippedCards(new Set())
                  setCardStartTime(null)
                  setTimeTaken(0)
                  setResultsSaved(false)
                  setIncorrectAttempts(new Map())
                }}
                className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover transition-colors"
              >
                Start Over
              </button>
              
              <button
                onClick={async () => {
                  // Save all results for all students
                  try {
                    for (let i = 0; i < gameResultsHooks.length; i++) {
                      const hook = gameResultsHooks[i]
                      await hook.saveAllResults()
                    }
                    setResultsSaved(true)
                  } catch (error) {
                    console.error('❌ Error saving results:', error)
                  }
                }}
                disabled={gameResultsHooks.some(hook => hook.isSaving) || resultsSaved}
                className="px-6 py-2 bg-pink text-white rounded-lg hover:bg-pink-hover transition-colors disabled:bg-pink-disabled disabled:cursor-not-allowed"
              >
                {gameResultsHooks.some(hook => hook.isSaving) ? 'Saving...' : resultsSaved ? 'Saved!' : 'Save Progress'}
              </button>
              
              <button
                onClick={() => router.push(effectiveSessionId ? `/sessions/${effectiveSessionId}/curriculum` : '/games')}
                className="px-6 py-2 bg-card-hover text-primary rounded-lg hover:bg-card border border-border"
              >
                Back to Content
              </button>
            </div>
            
            {gameResultsHooks.some(hook => hook.saveError) && (
              <p className="text-error text-sm mt-2">
                Failed to save progress. Please try again.
              </p>
            )}
          </div>
        )}

        {currentGameResultsHook && limitedGameContents[currentCardIndex] && (
          <div className="text-center mt-2">
            {(() => {
              const existingResult = currentGameResultsHook.getResultForContent(limitedGameContents[currentCardIndex].id)
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