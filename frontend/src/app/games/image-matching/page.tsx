'use client'

import { Suspense, useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { ArrowLeft, CheckCircle, Volume2, XCircle } from 'lucide-react'
import AppLayout from '@/components/AppLayout'
import { useGameContents } from '@/hooks/useGameContents'
import { useGameResults } from '@/hooks/useGameResults'
import { StudentSelector } from '@/components/games/StudentSelector'
import { useSessionContext } from '@/contexts/sessionContext'
import { useStudents } from '@/hooks/useStudents'

function ImageMatchingGameContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const sessionStudentIdsParam = searchParams.get('sessionStudentIds')

  // Show student selector if no students selected yet
  if (!sessionStudentIdsParam) {
    return (
      <StudentSelector
        gameTitle="Image Matching"
        onBack={() => router.back()}
        onStudentsSelected={(studentIds) => {
          // Update URL with selected students
          const params = new URLSearchParams(searchParams.toString())
          params.set('sessionStudentIds', studentIds.join(','))
          router.replace(`/games/image-matching?${params.toString()}`)
        }}
      />
    )
  }

  // Students are selected, show the actual game
  return <ImageMatchingGame />
}

function ImageMatchingGame() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const themeId = searchParams.get('themeId')
  const difficulty = searchParams.get('difficulty')
  const category = searchParams.get('category')
  const questionType = searchParams.get('questionType')
  const sessionStudentIdsParam = searchParams.get('sessionStudentIds')!
  const sessionId = searchParams.get('sessionId') || '00000000-0000-0000-0000-000000000000'

  const selectedStudentIds = sessionStudentIdsParam.split(',')
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0)
  const [score, setScore] = useState(0)
  const [showFeedback, setShowFeedback] = useState<'correct' | 'incorrect' | null>(null)
  const [gameComplete, setGameComplete] = useState(false)
  const [wrongAnswerUrl, setWrongAnswerUrl] = useState<string>('')
  const [imageOptions, setImageOptions] = useState<string[]>([])
  const [cardStartTime, setCardStartTime] = useState<number | null>(null)
  const [resultsSaved, setResultsSaved] = useState(false)

  const { gameContents, isLoading, error } = useGameContents({
    theme_id: themeId || undefined,
    difficulty_level: difficulty ? Number.parseInt(difficulty) : undefined,
    category: category as any,
    question_type: questionType as any,
  })

  // Get student names from session context
  const { students: sessionStudents } = useSessionContext()
  const { students: allStudents } = useStudents()

  // Calculate questions per student and limit total questions
  const questionsPerStudent = Math.floor(gameContents.length / selectedStudentIds.length)
  const totalQuestionsToUse = questionsPerStudent * selectedStudentIds.length;
  const limitedGameContents = gameContents.slice(0, totalQuestionsToUse);
  
  // Get current student based on question index
  const currentStudentIndex = currentQuestionIndex % selectedStudentIds.length;
  const currentSessionStudentId = Number.parseInt(selectedStudentIds[currentStudentIndex]);

  // Create game results hooks for each student
  const gameResultsHooks = selectedStudentIds.map(studentId => 
    useGameResults({
      session_student_id: Number.parseInt(studentId),
      session_id: sessionId || undefined,
    })
  );

  const gameResultsHook = gameResultsHooks[currentStudentIndex];
  
  const getStudentName = (sessionStudentId: number) => {
    const sessionStudent = sessionStudents.find(s => s.sessionStudentId === sessionStudentId);
    if (!sessionStudent) return 'Student';
    const student = allStudents?.find(s => s.id === sessionStudent.studentId);
    return student ? `${student.first_name} ${student.last_name}` : 'Student';
  };

  const currentQuestion = limitedGameContents?.[currentQuestionIndex]

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
      if (currentQuestionIndex < limitedGameContents.length - 1) {
        setCurrentQuestionIndex(currentQuestionIndex + 1)
      } else {
        setGameComplete(true)
      }
    }, 1500)
  }

  const handleSaveProgress = async () => {
    // Save all results for all students
    for (const hook of gameResultsHooks) {
      await hook.saveAllResults()
    }
    setResultsSaved(true)
  }

  const handlePlayAgain = async () => {
    // Save all results for all students
    for (const hook of gameResultsHooks) {
      await hook.saveAllResults()
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
    if (limitedGameContents && currentQuestion) {
      const otherQuestions = limitedGameContents.filter((_, idx) => idx !== currentQuestionIndex)
      if (otherQuestions.length > 0) {
        const randomQuestion = otherQuestions[Math.floor(Math.random() * otherQuestions.length)]
        setWrongAnswerUrl(randomQuestion.answer)
      }
    }
  }, [currentQuestionIndex, limitedGameContents])

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
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentQuestionIndex])

  if (isLoading) {
    return (
      <AppLayout>
        <div className="min-h-screen flex items-center justify-center bg-background">
          <p className="text-secondary">Loading game...</p>
        </div>
      </AppLayout>
    )
  }

  if (error || !limitedGameContents || limitedGameContents.length === 0) {
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
              Score: {score} / {limitedGameContents.length}
            </p>
            <div className="flex gap-4 justify-center">
              <button
                onClick={handlePlayAgain}
                className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover"
              >
                Play Again
              </button>
              <button
                onClick={handleSaveProgress}
                disabled={gameResultsHooks.some(hook => hook.isSaving) || resultsSaved}
                className="px-6 py-2 bg-pink text-white rounded-lg hover:bg-pink-hover transition-colors disabled:bg-pink-disabled disabled:cursor-not-allowed"
              >
                {gameResultsHooks.some(hook => hook.isSaving) ? 'Saving...' : resultsSaved ? 'Saved!' : 'Save Progress'}
              </button>
              <button
                onClick={() => router.push('/games')}
                className="px-6 py-2 bg-card-hover text-primary rounded-lg hover:bg-card border border-border"
              >
                Back to Games
              </button>
            </div>
            {gameResultsHooks.some(hook => hook.saveError) && (
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
              Question {currentQuestionIndex + 1} / {limitedGameContents.length}
            </div>
            <div className="text-lg font-semibold">
              Score: {score}
            </div>
          </div>

          {/* Current Student Banner */}
          <div className="bg-blue text-white rounded-lg p-4 mb-6 text-center">
            <p className="text-sm opacity-90 mb-1">Current Player</p>
            <p className="text-2xl font-bold">{getStudentName(currentSessionStudentId)}</p>
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

export default function ImageMatchingPage() {
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