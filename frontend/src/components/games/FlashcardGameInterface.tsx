'use client'

import React, { useEffect, useState } from 'react'
import { BookOpen, Brain, Check, ChevronRight, Gamepad2, RotateCw } from 'lucide-react'
import { useThemes } from '@/hooks/useThemes'
import { useGameContents } from '@/hooks/useGameContents'
import { 
  GetGameContentsCategory, 
  GetGameContentsQuestionType
  
} from '@/lib/api/theSpecialStandardAPI.schemas'
import type {Theme} from '@/lib/api/theSpecialStandardAPI.schemas';
import './flashcard.css'
import { useGameResults } from '@/hooks/useGameResults'

// Category configuration with your custom colors
const CATEGORIES = {
  [GetGameContentsCategory.receptive_language]: { 
    label: 'Receptive Language', 
    icon: '👂', 
    colorClass: 'bg-blue',
    hoverClass: 'hover:bg-blue-hover',
    lightClass: 'bg-blue-light'
  },
  [GetGameContentsCategory.expressive_language]: { 
    label: 'Expressive Language', 
    icon: '💬', 
    colorClass: 'bg-pink',
    hoverClass: 'hover:bg-pink-hover',
    lightClass: 'bg-pink-light'
  },
  [GetGameContentsCategory.social_pragmatic_language]: { 
    label: 'Social Pragmatic Language', 
    icon: '🤝', 
    colorClass: 'bg-orange',
    hoverClass: 'hover:bg-orange-hover',
    lightClass: 'bg-orange-light'
  },
  [GetGameContentsCategory.speech]: { 
    label: 'Speech', 
    icon: '🗣️', 
    colorClass: 'bg-blue',
    hoverClass: 'hover:bg-blue-hover',
    lightClass: 'bg-blue-light'
  },
}

const QUESTION_TYPES = {
  [GetGameContentsQuestionType.sequencing]: 'Sequencing',
  [GetGameContentsQuestionType.following_directions]: 'Following Directions',
  [GetGameContentsQuestionType.wh_questions]: 'WH Questions',
  [GetGameContentsQuestionType.true_false]: 'True/False',
  [GetGameContentsQuestionType.concepts_sorting]: 'Concepts & Sorting',
  [GetGameContentsQuestionType.fill_in_the_blank]: 'Fill in the Blank',
  [GetGameContentsQuestionType.categorical_language]: 'Categorical Language',
  [GetGameContentsQuestionType.emotions]: 'Emotions',
  [GetGameContentsQuestionType.teamwork_talk]: 'Teamwork Talk',
  [GetGameContentsQuestionType.express_excitement_interest]: 'Express Excitement/Interest',
  [GetGameContentsQuestionType.fluency]: 'Fluency',
  [GetGameContentsQuestionType.articulation_s]: 'Articulation - S',
  [GetGameContentsQuestionType.articulation_l]: 'Articulation - L',
}

// Flashcard component
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
  student_id 
}: FlashcardGameInterfaceProps) {
  const [selectedTheme, setSelectedTheme] = useState<Theme | null>(null)
  const [selectedCategory, setSelectedCategory] = useState<GetGameContentsCategory | null>(null)
  const [selectedQuestionType, setSelectedQuestionType] = useState<GetGameContentsQuestionType | null>(null)
  const [selectedGame, setSelectedGame] = useState<'flashcards' | 'memory' | 'quiz' | null>(null)
  const [currentCardIndex, setCurrentCardIndex] = useState(0)
  const [flippedCards, setFlippedCards] = useState<Set<number>>(new Set())
  const [cardStartTime, setCardStartTime] = useState<number | null>(null)
  const [timeTaken, setTimeTaken] = useState(0)

  const { themes, isLoading: themesLoading, error: themesError, refetch: refetchThemes } = useThemes()
  const { gameContents, isLoading: contentsLoading, error: contentsError } = useGameContents({
    theme_id: selectedTheme?.id,
    category: selectedCategory || undefined,
    question_type: selectedQuestionType || undefined,
    question_count: 10,
  })

  const gameResultsHook = session_student_id ? useGameResults({
    session_student_id,
    session_id,
    student_id
  }) : null

  const startCard = gameResultsHook?.startCard;

  // Initialize timer when card loads
  useEffect(() => {
    if (selectedGame === 'flashcards' && gameContents[currentCardIndex]) {
      setCardStartTime(Date.now())
      setFlippedCards(new Set())
      setTimeTaken(0)
      startCard?.(gameContents[currentCardIndex])
    }
  }, [currentCardIndex, selectedGame, gameContents])

  // Update timer display
  useEffect(() => {
    if (cardStartTime === null) return
    
    const interval = setInterval(() => {
      setTimeTaken(Math.floor((Date.now() - cardStartTime) / 1000))
    }, 100)
    
    return () => clearInterval(interval)
  }, [cardStartTime])

  const handleMarkCorrect = () => {
    if (!gameResultsHook || !gameContents[currentCardIndex]) return
    
    const finalTime = Math.floor((Date.now() - (cardStartTime || Date.now())) / 1000)
    gameResultsHook.completeCard(gameContents[currentCardIndex].id, finalTime)
    
    // Move to next card
    setCurrentCardIndex(prev => prev + 1)
    
    setCardStartTime(null)
  }

  const resetSelection = () => {
    setSelectedTheme(null)
    setSelectedCategory(null)
    setSelectedQuestionType(null)
    setSelectedGame(null)
    setCurrentCardIndex(0)
    setFlippedCards(new Set())
    setCardStartTime(null)
    setTimeTaken(0)
  }

  const handleCardFlip = (index: number) => {
    setFlippedCards(prev => {
      const newSet = new Set(prev)
      if (newSet.has(index)) {
        newSet.delete(index)
      } else {
        newSet.add(index)
      }
      return newSet
    })
  }

  // Error state
  if (themesError && !selectedTheme) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <p className="text-error mb-4">Failed to load themes</p>
          <button 
            onClick={() => refetchThemes()}
            className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover transition-colors"
          >
            Retry
          </button>
        </div>
      </div>
    )
  }

  // Step 1: Theme Selection
  if (!selectedTheme) {
    return (
      <div className="min-h-screen bg-background p-8">
        <div className="max-w-6xl mx-auto">
          <h1 className="mb-8">Select a Theme</h1>
          {themesLoading ? (
            <div className="text-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue mx-auto mb-4"></div>
              <p className="text-muted">Loading themes...</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {themes.map((theme) => (
                <button
                  key={theme.id}
                  onClick={() => setSelectedTheme(theme)}
                  className="bg-card rounded-lg shadow-md p-6 hover:shadow-lg transition-all duration-200 text-left group hover:bg-card-hover border border-default hover:border-hover"
                >
                  <div className="flex items-center justify-between">
                    <h3>{theme.name}</h3>
                    <ChevronRight className="w-5 h-5 text-muted group-hover:text-primary transition-colors" />
                  </div>
                  <p className="text-secondary mt-2">
                    {`Explore ${theme.name} themed content`}
                  </p>
                </button>
              ))}
            </div>
          )}
        </div>
      </div>
    )
  }

  // Step 2: Category Selection
  if (!selectedCategory) {
    return (
      <div className="min-h-screen bg-background p-8">
        <div className="max-w-4xl mx-auto">
          <button
            onClick={resetSelection}
            className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
          >
            ← Back to Themes
          </button>
          <h1 className="mb-2">Select a Category</h1>
          <p className="text-secondary mb-8">Theme: {selectedTheme.name}</p>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {Object.entries(CATEGORIES).map(([key, category]) => (
              <button
                key={key}
                onClick={() => setSelectedCategory(key as GetGameContentsCategory)}
                className="bg-card rounded-lg shadow-md p-8 hover:shadow-lg transition-all duration-200 text-left group hover:bg-card-hover border border-default hover:border-hover"
              >
                <div className="flex items-center gap-4">
                  <div className={`w-16 h-16 ${category.colorClass} rounded-full flex items-center justify-center text-2xl text-white`}>
                    {category.icon}
                  </div>
                  <div className="flex-1">
                    <h3>{category.label}</h3>
                  </div>
                  <ChevronRight className="w-5 h-5 text-muted group-hover:text-primary transition-colors" />
                </div>
              </button>
            ))}
          </div>
        </div>
      </div>
    )
  }

  // Step 3: Question Type Selection
  if (!selectedQuestionType) {
    const category = CATEGORIES[selectedCategory]
    return (
      <div className="min-h-screen bg-background p-8">
        <div className="max-w-4xl mx-auto">
          <button
            onClick={() => setSelectedCategory(null)}
            className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
          >
            ← Back to Categories
          </button>
          <h1 className="mb-2">Select Question Type</h1>
          <p className="text-secondary mb-8">
            Theme: {selectedTheme.name} / Category: {category.label}
          </p>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {Object.entries(QUESTION_TYPES).map(([key, label]) => (
              <button
                key={key}
                onClick={() => setSelectedQuestionType(key as GetGameContentsQuestionType)}
                className="bg-card rounded-lg shadow-md p-4 hover:shadow-lg transition-all duration-200 text-left group flex items-center justify-between hover:bg-card-hover border border-default hover:border-hover"
              >
                <span className="font-medium text-primary">{label}</span>
                <ChevronRight className="w-5 h-5 text-muted group-hover:text-primary transition-colors" />
              </button>
            ))}
          </div>
        </div>
      </div>
    )
  }

  // Step 4: Game Selection
  if (!selectedGame) {
    return (
      <div className="min-h-screen bg-background p-8">
        <div className="max-w-4xl mx-auto">
          <button
            onClick={() => setSelectedQuestionType(null)}
            className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
          >
            ← Back to Question Types
          </button>
          <h1 className="mb-8">Select a Game</h1>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <button
              onClick={() => setSelectedGame('flashcards')}
              className="bg-card rounded-lg shadow-md p-8 hover:shadow-lg transition-all duration-200 group hover:bg-card-hover border border-default hover:border-hover"
            >
              <BookOpen className="w-12 h-12 text-blue mb-4 mx-auto" />
              <h3 className="mb-2">Flashcards</h3>
              <p className="text-secondary text-sm">Practice with interactive flashcards</p>
            </button>
            
            <button
              disabled
              className="bg-card rounded-lg shadow-md p-8 opacity-50 cursor-not-allowed border border-default"
            >
              <Brain className="w-12 h-12 text-muted mb-4 mx-auto" />
              <h3 className="mb-2 text-muted">Memory Match</h3>
              <p className="text-disabled text-sm">Coming soon</p>
            </button>
            
            <button
              disabled
              className="bg-card rounded-lg shadow-md p-8 opacity-50 cursor-not-allowed border border-default"
            >
              <Gamepad2 className="w-12 h-12 text-muted mb-4 mx-auto" />
              <h3 className="mb-2 text-muted">Quiz Game</h3>
              <p className="text-disabled text-sm">Coming soon</p>
            </button>
          </div>
        </div>
      </div>
    )
  }

  // Step 5: Flashcard Game
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
            onClick={() => setSelectedGame(null)}
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
            onClick={() => setSelectedGame(null)}
            className="text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
          >
            ← Back to Games
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
            <span className="text-muted">Theme:</span>
            <span className="font-medium text-primary">{selectedTheme.name}</span>
            <span className="text-muted mx-2">•</span>
            <span className="text-muted">Category:</span>
            <span className="font-medium text-primary">{CATEGORIES[selectedCategory].label}</span>
            <span className="text-muted mx-2">•</span>
            <span className="text-muted">Type:</span>
            <span className="font-medium text-primary">{QUESTION_TYPES[selectedQuestionType]}</span>
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
              onFlip={() => handleCardFlip(currentCardIndex)}
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

        {gameResultsHook && selectedGame === 'flashcards' && gameContents[currentCardIndex] && (
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