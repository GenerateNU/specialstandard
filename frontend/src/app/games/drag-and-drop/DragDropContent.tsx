'use client'

import { useEffect, useMemo, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { ArrowLeft, CheckCircle, Volume2 } from 'lucide-react'
import {
  closestCenter,
  DndContext,
  DragOverlay,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core'
import type { DragEndEvent, DragStartEvent } from '@dnd-kit/core'
import {
  SortableContext,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import AppLayout from '@/components/AppLayout'
import { useGameContents } from '@/hooks/useGameContents'
import { useGameResults } from '@/hooks/useGameResults'
import DraggableImage from '@/components/games/drag-and-drop/DraggableImage'
import SequenceSlot from '@/components/games/drag-and-drop/SequenceSlot'

export default function SequencingGameContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const sessionStudentId = Number.parseInt(searchParams.get("sessionStudentId") ?? "0");
  const themeId = searchParams.get('themeId')
  const difficulty = searchParams.get('difficulty')
  const category = searchParams.get('category')
  const questionType = searchParams.get('questionType')
  const sessionId = searchParams.get('sessionId')
  const studentId = searchParams.get('studentId')

  const { gameContents, isLoading, error } = useGameContents({
    theme_id: themeId || undefined,
    difficulty_level: difficulty ? Number.parseInt(difficulty) : undefined,
    category: category as 'receptive_language' | 'expressive_language' | 'social_pragmatic_language' | 'speech' | undefined,
    question_type: questionType as 'sequencing' | 'following_directions' | 'wh_questions' | 'true_false' | 'concepts_sorting' | 'fill_in_the_blank' | 'categorical_language' | 'emotions' | 'teamwork_talk' | 'express_excitement_interest' | 'fluency' | 'articulation_s' | 'articulation_l' | undefined,
  })

  const gameResultsHook = useGameResults({
    session_id: sessionId || undefined,
    student_id: studentId || undefined,
    session_student_id: sessionStudentId,
  })

  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0)
  const [score, setScore] = useState(0)
  const [gameComplete, setGameComplete] = useState(false)
  const [resultsSaved, setResultsSaved] = useState(false)
  const [cardStartTime, setCardStartTime] = useState<number | null>(null)
  const [incorrectAttempts, setIncorrectAttempts] = useState(0)
  
  // Sequencing state - using FILENAMES as IDs
  const [availableFilenames, setAvailableFilenames] = useState<string[]>([])
  const [sequence, setSequence] = useState<(string | null)[]>([])
  const [activeId, setActiveId] = useState<string | null>(null)
  const [showSuccess, setShowSuccess] = useState(false)
  const [correctAnswerFilenames, setCorrectAnswerFilenames] = useState<string[]>([])
  const [loadingAnswer, setLoadingAnswer] = useState(false)

  const currentQuestion = gameContents?.[currentQuestionIndex]

  // Create a map of filename -> presigned URL
  const filenameToUrlMap = useMemo(() => {
    if (!currentQuestion?.options || !currentQuestion?.presigned_options) return {}
    
    const map: Record<string, string> = {}
    currentQuestion.options.forEach((filename, index) => {
      // Ensure we have a corresponding presigned option
      if (index < currentQuestion.presigned_options.length) {
        map[filename] = currentQuestion.presigned_options[index]
      } else {
        // Fallback or handle mismatch if needed
        map[filename] = filename // filename text fallback
      }
    })
    return map
  }, [currentQuestion])

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    })
  )

  const speakWord = (word: string): void => {
    if (typeof window !== 'undefined' && 'speechSynthesis' in window) {
      window.speechSynthesis.cancel()
      const utterance = new SpeechSynthesisUtterance(word)
      utterance.rate = 0.8
      utterance.pitch = 1.0
      utterance.volume = 1.0
      window.speechSynthesis.speak(utterance)
    }
  }

  // Parse correct answer from raw_answer field - store as FILENAMES
  useEffect(() => {
    if (!currentQuestion?.raw_answer) return
    
    setLoadingAnswer(true)
    try {      
      // Parse the raw_answer JSON string (might be double-encoded)
      const parsed = JSON.parse(currentQuestion.raw_answer)

      // Store as filenames directly
      setCorrectAnswerFilenames(parsed)
    } catch (err) {
      console.error('Error parsing raw_answer:', err)
      setCorrectAnswerFilenames([])
    } finally {
      setLoadingAnswer(false)
    }
  }, [currentQuestion])

  // Initialize available filenames and sequence
  useEffect(() => {
    if (currentQuestion?.options && correctAnswerFilenames.length > 0) {
      // Use filenames directly from options
      const shuffled = [...currentQuestion.options].sort(() => Math.random() - 0.5)
      setAvailableFilenames(shuffled)
      setSequence(Array.from({ length: correctAnswerFilenames.length }, () => null))
      setShowSuccess(false)
      setIncorrectAttempts(0)
    }
  }, [currentQuestionIndex, currentQuestion, correctAnswerFilenames])

  // Speak question on load
  useEffect(() => {
    if (currentQuestion?.question) {
      speakWord(currentQuestion.question)
    }
  }, [currentQuestionIndex, currentQuestion])
  
  // Track card start time
  useEffect(() => {
    if (currentQuestion) {
      setCardStartTime(Date.now())
      gameResultsHook?.startCard(currentQuestion)
    }
  }, [currentQuestion, gameResultsHook])

  const handleDragStart = (event: DragStartEvent) => {
    setActiveId(event.active.id as string)
  }

  const checkSequence = (currentSequence: (string | null)[]) => {
    // Only check if all slots are filled
    if (currentSequence.every(item => item !== null)) {
      const isCorrect = currentSequence.every((filename, idx) => filename === correctAnswerFilenames[idx])
      
      if (isCorrect) {
        setShowSuccess(true)
        setScore(score + 1)
        
        // Save result using the hook
        if (gameResultsHook && currentQuestion && cardStartTime) {
          const timeTaken = Math.floor((Date.now() - cardStartTime) / 1000)
          gameResultsHook.completeCard(
            currentQuestion.id, 
            timeTaken, 
            incorrectAttempts,
          )
        }

        setTimeout(() => {
          if (currentQuestionIndex < gameContents.length - 1) {
            setCurrentQuestionIndex(currentQuestionIndex + 1)
          } else {
            setGameComplete(true)
          }
        }, 2000)
      } else {
        setIncorrectAttempts(incorrectAttempts + 1)
        // Could track the incorrect attempt here if needed
      }
    }
  }

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    setActiveId(null)

    if (!over) return
    
    const draggedFilename = active.id as string
    
    // Check if dropped on a sequence slot
    if (String(over.id).startsWith('slot-')) {
      const slotIndex = Number.parseInt(String(over.id).split('-')[1])
      
      // Update sequence
      const newSequence = [...sequence]
      
      // If slot already has an image, return it to available
      if (newSequence[slotIndex]) {
        setAvailableFilenames([...availableFilenames, newSequence[slotIndex]!])
      }
      
      newSequence[slotIndex] = draggedFilename
      setSequence(newSequence)
      
      // Remove from available
      setAvailableFilenames(availableFilenames.filter(f => f !== draggedFilename))
      
      // Check if sequence is correct
      checkSequence(newSequence)
    }
  }

  const handleRemoveFromSequence = (index: number) => {
    const filenameToRemove = sequence[index]
    if (filenameToRemove) {
      const newSequence = [...sequence]
      newSequence[index] = null
      setSequence(newSequence)
      setAvailableFilenames([...availableFilenames, filenameToRemove])
    }
  }

  const handleSaveProgress = async () => {
    if (gameResultsHook) {
      await gameResultsHook.saveAllResults()
      setResultsSaved(true)
    }
  }

  const handlePlayAgain = async () => {
    if (gameResultsHook) {
      await gameResultsHook.saveAllResults()
    }
    
    setCurrentQuestionIndex(0)
    setScore(0)
    setGameComplete(false)
    setCardStartTime(null)
    setResultsSaved(false)
  }

  if (isLoading || loadingAnswer) {
    return (
      <AppLayout>
        <div className="min-h-screen flex items-center justify-center bg-background">
          <p className="text-secondary">Loading game...</p>
        </div>
      </AppLayout>
    )
  }

  if (error || !gameContents || gameContents.length === 0) {
    return (
      <AppLayout>
        <div className="min-h-screen flex items-center justify-center bg-background">
          <div className="text-center">
            <p className="text-error mb-4">No questions available for this selection.</p>
            <button
              onClick={() => router.push('/games')}
              className="px-4 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover"
            >
              Back to Games
            </button>
          </div>
        </div>
      </AppLayout>
    )
  }

  if (gameComplete) {
    return (
      <AppLayout>
        <div className="min-h-screen flex items-center justify-center bg-background">
          <div className="text-center bg-card p-8 rounded-sm shadow-lg max-w-md">
            <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
            <h2 className="mb-4">Game Complete!</h2>
            <p className="text-2xl mb-6">
              Score: {score} / {gameContents.length}
            </p>
            <div className="flex gap-4 justify-center">
              <button
                onClick={handlePlayAgain}
                className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover"
              >
                Play Again
              </button>
              {gameResultsHook && (
                <button
                  onClick={handleSaveProgress}
                  disabled={gameResultsHook.isSaving || resultsSaved}
                  className="px-6 py-2 bg-pink text-white rounded-lg hover:bg-pink-hover transition-colors disabled:bg-pink-disabled disabled:cursor-not-allowed"
                >
                  {gameResultsHook.isSaving ? 'Saving...' : resultsSaved ? 'Saved!' : 'Save Progress'}
                </button>
              )}
              <button
                onClick={() => router.push('/games')}
                className="px-6 py-2 bg-card-hover text-primary rounded-lg hover:bg-card border border-border"
              >
                Back to Games
              </button>
            </div>
            {gameResultsHook?.saveError && (
              <p className="text-error text-sm mt-2">
                Failed to save progress. Please try again.
              </p>
            )}
          </div>
        </div>
      </AppLayout>
    )
  }

  return (
    <AppLayout>
      <div className="min-h-screen bg-background p-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-between mb-8">
            <button
              onClick={() => router.push('/games')}
              className="text-blue hover:text-blue-hover flex items-center gap-2"
            >
              <ArrowLeft className="w-4 h-4" />
              Back
            </button>
            <div className="text-lg font-semibold">
              Question {currentQuestionIndex + 1} / {gameContents.length}
            </div>
            <div className="text-lg font-semibold">
              Score: {score}
            </div>
          </div>

          <div className="bg-card rounded-lg shadow-lg p-8 mb-8 text-center">
            <h1 className="text-4xl font-bold mb-4">{currentQuestion?.question}</h1>
            <button
              onClick={() => speakWord(currentQuestion?.question || '')}
              className="mx-auto flex items-center gap-2 px-6 py-3 bg-blue text-white rounded-lg hover:bg-blue-hover"
            >
              <Volume2 className="w-5 h-5" />
              Play Sound
            </button>
          </div>

          <DndContext
            sensors={sensors}
            collisionDetection={closestCenter}
            onDragStart={handleDragStart}
            onDragEnd={handleDragEnd}
          >
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* Sequence Area */}
              <div className="bg-card rounded-lg shadow-lg p-6">
                <h3 className="text-xl font-semibold mb-4">Put in order:</h3>
                <div className="space-y-4">
                  {sequence.map((filename, index) => (
                    <div 
                      key={`slot-${index}`}
                      className="relative"
                    >
                      <div className="absolute -left-8 top-1/2 -translate-y-1/2 text-lg font-bold text-secondary">
                        {index + 1}
                      </div>
                      <div onClick={() => handleRemoveFromSequence(index)}>
                        <SequenceSlot 
                          position={index + 1} 
                          filename={filename}
                          url={filename ? filenameToUrlMap[filename] : undefined}
                          slotId={`slot-${index}`}
                        />
                      </div>
                    </div>
                  ))}
                </div>
                {showSuccess && (
                  <div className="mt-4 p-4 bg-green-100 text-green-800 rounded-lg flex items-center gap-2">
                    <CheckCircle className="w-5 h-5" />
                    <span>Correct! Moving to next question...</span>
                  </div>
                )}
              </div>

              {/* Available Images Pool */}
              <div className="bg-card rounded-lg shadow-lg p-6">
                <h3 className="text-xl font-semibold mb-4">Available images:</h3>
                <SortableContext
                  items={availableFilenames}
                  strategy={verticalListSortingStrategy}
                >
                  <div className="grid grid-cols-2 gap-4">
                    {availableFilenames.map((filename) => (
                      <DraggableImage
                        key={filename}
                        id={filename}
                        filename={filename}
                        url={filenameToUrlMap[filename]}
                      />
                    ))}
                  </div>
                </SortableContext>
                {availableFilenames.length === 0 && (
                  <p className="text-secondary text-center py-8">
                    All images placed in sequence
                  </p>
                )}
              </div>
            </div>

            <DragOverlay>
              {activeId ? (
                <div className="bg-card rounded-lg shadow-xl p-3 rotate-3 cursor-grabbing">
                  <div className="aspect-square rounded-lg overflow-hidden">
                    <img 
                      src={filenameToUrlMap[activeId]} 
                      alt="Dragging"
                      className="w-full h-full object-contain"
                    />
                  </div>
                </div>
              ) : null}
            </DragOverlay>
          </DndContext>
        </div>
      </div>
    </AppLayout>
  )
}