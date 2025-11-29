'use client'

import { Suspense, useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { ArrowLeft, CheckCircle, Volume2, XCircle } from 'lucide-react'
import AppLayout from '@/components/AppLayout'
import { useGameContents } from '@/hooks/useGameContents'
import { useGameResults } from '@/hooks/useGameResults'

function ImageMatchingGameContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const themeId = searchParams.get('themeId')
  const difficulty = searchParams.get('difficulty')
  const category = searchParams.get('category')
  const questionType = searchParams.get('questionType')
  const sessionStudentId = Number.parseInt(searchParams.get('sessionStudentId') || '0')
  const sessionId = searchParams.get('sessionId') || '00000000-0000-0000-0000-000000000000'

  const { gameContents, isLoading, error } = useGameContents({
    theme_id: themeId || undefined,
    difficulty_level: difficulty ? Number.parseInt(difficulty) : undefined,
    category: category as any,
    question_type: questionType as any,
  })

  const gameResultsHook = useGameResults({
    session_student_id: sessionStudentId,
    session_id: sessionId || undefined,
  })

  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0)
  const [score, setScore] = useState(0)
  const [showFeedback, setShowFeedback] = useState<'correct' | 'incorrect' | null>(null)
  const [gameComplete, setGameComplete] = useState(false)
  const [wrongAnswerUrl, setWrongAnswerUrl] = useState<string>('')
  const [imageOptions, setImageOptions] = useState<string[]>([])
  const [cardStartTime, setCardStartTime] = useState<number | null>(null)
  const [resultsSaved, setResultsSaved] = useState(false)
  const currentQuestion = gameContents?.[currentQuestionIndex]

  type SpeakWordFn = (word: string) => void

  const speakWord: SpeakWordFn = (word: string): void => {
    if (typeof window !== 'undefined' && 'speechSynthesis' in window) {
      window.speechSynthesis.cancel()

      const utterance = new SpeechSynthesisUtterance(word)
      utterance.rate = 0.8
      utterance.pitch = 1.0
      utterance.volume = 1.0

      window.speechSynthesis.speak(utterance)
    }
  }

  const handleImageClick = (selectedAnswer: string) => {
    if (showFeedback) return
    const isCorrect = selectedAnswer === currentQuestion?.answer

    if (isCorrect) {
      setScore(score + 1)
      setShowFeedback('correct')
      
      if (gameResultsHook && currentQuestion && cardStartTime) {
        const timeTaken = Math.floor((Date.now() - cardStartTime) / 1000)
        gameResultsHook.completeCard(currentQuestion.id, timeTaken)
      }
    } else {
      setShowFeedback('incorrect')
    }

    setTimeout(() => {
      setShowFeedback(null)
      if (currentQuestionIndex < gameContents.length - 1) {
        setCurrentQuestionIndex(currentQuestionIndex + 1)
      } else {
        setGameComplete(true)
      }
    }, 1500)
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

  useEffect(() => {
    if (currentQuestion && wrongAnswerUrl) {
      const shuffled = [currentQuestion.answer, wrongAnswerUrl].sort(() => Math.random() - 0.5)
      setImageOptions(shuffled)
    }
  }, [currentQuestion, wrongAnswerUrl])

  useEffect(() => {
    if (gameContents && currentQuestion) {
      const otherQuestions = gameContents.filter((_, idx) => idx !== currentQuestionIndex)
      if (otherQuestions.length > 0) {
        const randomQuestion = otherQuestions[Math.floor(Math.random() * otherQuestions.length)]
        setWrongAnswerUrl(randomQuestion.answer)
      }
    }
  }, [currentQuestionIndex, gameContents])

  useEffect(() => {
    if (currentQuestion?.question && !showFeedback) {
      speakWord(currentQuestion.question)
    }
  }, [currentQuestionIndex, currentQuestion])

  useEffect(() => {
    if (currentQuestion) {
      setCardStartTime(Date.now())
      gameResultsHook?.startCard(currentQuestion)
    }
  }, [currentQuestion])

  if (isLoading) {
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
        <div className="max-w-4xl mx-auto">
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

          <div className="grid grid-cols-2 gap-8">
            {imageOptions.map((answerUrl, index) => (
              <button
                key={index}
                onClick={() => handleImageClick(answerUrl)}
                disabled={showFeedback !== null}
                className={`relative bg-card rounded-lg shadow-lg p-4 hover:shadow-xl transition-all duration-200 border-4 ${
                  showFeedback === 'correct' && answerUrl === currentQuestion?.answer
                    ? 'border-green-500'
                    : showFeedback === 'incorrect' && answerUrl !== currentQuestion?.answer
                    ? 'border-red-500'
                    : 'border-transparent hover:border-blue'
                } ${showFeedback ? 'pointer-events-none' : ''}`}
              >
                <div className="aspect-square rounded-lg overflow-hidden">
                  <img 
                    src={answerUrl} 
                    alt="Option"
                    className="w-full h-full object-contain"
                  />
                </div>

                {showFeedback === 'correct' && answerUrl === currentQuestion?.answer && (
                  <div className="absolute inset-0 flex items-center justify-center bg-green-500 opacity-20">
                    <CheckCircle className="w-16 h-16 text-green-500" />
                  </div>
                )}
                {showFeedback === 'incorrect' && answerUrl !== currentQuestion?.answer && (
                  <div className="absolute inset-0 flex items-center justify-center bg-red-500 opacity-20">
                    <XCircle className="w-16 h-16 text-red-500" />
                  </div>
                )}
              </button>
            ))}
          </div>
        </div>
      </div>
    </AppLayout>
  )
}

export default function ImageMatchingGame() {
  return (
    <Suspense fallback={
      <AppLayout>
        <div className="min-h-screen flex items-center justify-center bg-background">
          <p className="text-secondary">Loading game...</p>
        </div>
      </AppLayout>
    }>
      <ImageMatchingGameContent />
    </Suspense>
  )
}