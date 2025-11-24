'use client'

import { Suspense, useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { ArrowLeft, CheckCircle, Volume2 } from 'lucide-react'
import {
  DndContext,
  DragOverlay,
  closestCenter,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core'
import {
  SortableContext,
  verticalListSortingStrategy,
  useSortable,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import AppLayout from '@/components/AppLayout'
import { useGameContents } from '@/hooks/useGameContents'
import { useGameResults } from '@/hooks/useGameResults'

// Draggable image component
function DraggableImage({ id, url, isInSequence }) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className={`bg-card rounded-lg shadow-md p-3 cursor-grab active:cursor-grabbing ${
        isInSequence ? 'ring-2 ring-blue' : ''
      }`}
    >
      <div className="aspect-square rounded-lg overflow-hidden">
        <img 
          src={url} 
          alt="Sequence item"
          className="w-full h-full object-contain"
        />
      </div>
    </div>
  )
}

// Droppable sequence slot
function SequenceSlot({ position, imageUrl }) {
  return (
    <div className="relative bg-card-hover border-2 border-dashed border-border rounded-lg p-4 min-h-[150px] flex items-center justify-center">
      {imageUrl ? (
        <div className="w-full h-full">
          <div className="aspect-square rounded-lg overflow-hidden">
            <img 
              src={imageUrl} 
              alt={`Position ${position}`}
              className="w-full h-full object-contain"
            />
          </div>
        </div>
      ) : (
        <span className="text-secondary text-sm">Drop here ({position})</span>
      )}
    </div>
  )
}

export default function SequencingGameContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  
  const themeId = searchParams.get('themeId')
  const difficulty = searchParams.get('difficulty')
  const category = searchParams.get('category')
  const questionType = searchParams.get('questionType')
  const sessionStudentId = searchParams.get('session_student_id')
  const sessionId = searchParams.get('session_id')
  const studentId = searchParams.get('student_id')

  const { gameContents, isLoading, error } = useGameContents({
    theme_id: themeId || undefined,
    difficulty_level: difficulty ? Number.parseInt(difficulty) : undefined,
    category: category as any,
    question_type: questionType as any,
  })

  const gameResultsHook = sessionStudentId ? useGameResults({
    session_student_id: Number.parseInt(sessionStudentId),
    session_id: sessionId || undefined,
    student_id: studentId || undefined
  }) : null

  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0)
  const [score, setScore] = useState(0)
  const [gameComplete, setGameComplete] = useState(false)
  const [resultsSaved, setResultsSaved] = useState(false)
  const [cardStartTime, setCardStartTime] = useState<number | null>(null)
  const [incorrectAttempts, setIncorrectAttempts] = useState(0)
  
  // Sequencing state
  const [availableImages, setAvailableImages] = useState<string[]>([])
  const [sequence, setSequence] = useState<(string | null)[]>([])
  const [activeId, setActiveId] = useState<string | null>(null)
  const [showSuccess, setShowSuccess] = useState(false)
  const [correctAnswer, setCorrectAnswer] = useState<string[]>([])
  const [loadingAnswer, setLoadingAnswer] = useState(false)

  const currentQuestion = gameContents?.[currentQuestionIndex]

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

  // Helper to convert filename to S3 URL
  const getImageUrl = (filename: string): string => {
    // If it's already a full URL, return as-is
    if (filename.startsWith('http')) {
      return filename
    }
    // Otherwise, construct S3 URL (you may need to adjust bucket/path)
    return `https://specialstandard-bucket.s3.us-east-1.amazonaws.com/${filename}`
  }

  // Fetch correct answer from URL and convert to full image URLs
  useEffect(() => {
    const fetchAnswer = async () => {
      if (!currentQuestion?.answer) return
      
      setLoadingAnswer(true)
      try {
        // Fetch the answer URL to get the JSON array
        const response = await fetch(currentQuestion.answer)
        const text = await response.text()
        
        // Parse the JSON (might be double-encoded)
        let parsed = JSON.parse(text)
        if (typeof parsed === 'string') {
          parsed = JSON.parse(parsed)
        }
        
        // Convert filenames to full S3 URLs
        const imageUrls = parsed.map(filename => getImageUrl(filename))
        setCorrectAnswer(imageUrls)
      } catch (err) {
        console.error('Error fetching answer:', err)
        setCorrectAnswer([])
      } finally {
        setLoadingAnswer(false)
      }
    }

    fetchAnswer()
  }, [currentQuestion])

  // Initialize available images and sequence
  useEffect(() => {
    if (currentQuestion?.options && correctAnswer.length > 0) {
      // Convert option filenames to full S3 URLs
      const imageUrls = currentQuestion.options.map(filename => getImageUrl(filename))
      const shuffled = [...imageUrls].sort(() => Math.random() - 0.5)
      setAvailableImages(shuffled)
      setSequence(Array(correctAnswer.length).fill(null))
      setShowSuccess(false)
      setIncorrectAttempts(0)
    }
  }, [currentQuestionIndex, currentQuestion, correctAnswer])

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
  }, [currentQuestion])

  const handleDragStart = (event) => {
    setActiveId(event.active.id)
  }

  const handleDragEnd = (event) => {
    const { active, over } = event
    setActiveId(null)

    if (!over) return

    const draggedImage = active.id
    
    // Check if dropped on a sequence slot
    if (over.id.startsWith('slot-')) {
      const slotIndex = parseInt(over.id.split('-')[1])
      
      // Update sequence
      const newSequence = [...sequence]
      
      // If slot already has an image, return it to available
      if (newSequence[slotIndex]) {
        setAvailableImages([...availableImages, newSequence[slotIndex]!])
      }
      
      newSequence[slotIndex] = draggedImage
      setSequence(newSequence)
      
      // Remove from available
      setAvailableImages(availableImages.filter(img => img !== draggedImage))
      
      // Check if sequence is correct
      checkSequence(newSequence)
    }
  }

  const checkSequence = (currentSequence: (string | null)[]) => {
    // Only check if all slots are filled
    if (currentSequence.every(item => item !== null)) {
      const isCorrect = currentSequence.every((img, idx) => img === correctAnswer[idx])
      
      if (isCorrect) {
        setShowSuccess(true)
        setScore(score + 1)
        
        if (gameResultsHook && currentQuestion && cardStartTime) {
          const timeTaken = Math.floor((Date.now() - cardStartTime) / 1000)
          gameResultsHook.completeCard(currentQuestion.id, timeTaken, incorrectAttempts)
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
      }
    }
  }

  const handleRemoveFromSequence = (index: number) => {
    const imageToRemove = sequence[index]
    if (imageToRemove) {
      const newSequence = [...sequence]
      newSequence[index] = null
      setSequence(newSequence)
      setAvailableImages([...availableImages, imageToRemove])
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
                  {sequence.map((imageUrl, index) => (
                    <SortableContext
                      key={`slot-${index}`}
                      items={[`slot-${index}`]}
                      strategy={verticalListSortingStrategy}
                    >
                      <div 
                        id={`slot-${index}`}
                        className="relative"
                        onClick={() => handleRemoveFromSequence(index)}
                      >
                        <div className="absolute -left-8 top-1/2 -translate-y-1/2 text-lg font-bold text-secondary">
                          {index + 1}
                        </div>
                        <SequenceSlot position={index + 1} imageUrl={imageUrl} />
                      </div>
                    </SortableContext>
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
                  items={availableImages}
                  strategy={verticalListSortingStrategy}
                >
                  <div className="grid grid-cols-2 gap-4">
                    {availableImages.map((imageUrl) => (
                      <DraggableImage
                        key={imageUrl}
                        id={imageUrl}
                        url={imageUrl}
                        isInSequence={false}
                      />
                    ))}
                  </div>
                </SortableContext>
                {availableImages.length === 0 && (
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
                      src={activeId} 
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