'use client'

import React, {useCallback, useEffect, useState} from "react";
import MatchingCard from "@/components/games/word-image-match/MatchingCard";
import { useRouter } from 'next/navigation'
import {CheckCircle, RotateCw} from "lucide-react";
import type {GetGameContentsCategory, PostGameResultInput} from "@/lib/api/theSpecialStandardAPI.schemas";
import {CATEGORIES} from "@/components/games/FlashcardGameInterface";
import {useGameContents} from "@/hooks/useGameContents";
import AppLayout from "@/components/AppLayout";
import { useSessionContext } from "@/contexts/sessionContext";
import { useStudents } from "@/hooks/useStudents";
import { getGameResult } from "@/lib/api/game-result";

export interface MatchingCardContent {
  id: string
  isImage: boolean
  value: string
  pairID: string
}

interface StudentResult {
  content_id: string
  time_taken_sec: number
  count_of_incorrect_attempts: number
  incorrect_attempts: string[]
}

interface WordImageMatchingGameInterfaceProps {
  session_student_ids: number[]
  session_id: string
  themeID: string
  themeName: string | null
  difficulty: string
  category: string
  questionType: string
}

export default function WordImageMatchingGameInterface({
  session_student_ids,
  session_id,
  themeID,
  themeName,
  difficulty,
  category,
  questionType
}: WordImageMatchingGameInterfaceProps) {
  const router = useRouter()
  const { session, students: sessionStudents } = useSessionContext()
  const { students: allStudents } = useStudents()
  
  // Prefer session context ID over prop
  const effectiveSessionId = session?.id || session_id

  // Track current student index and completed students
  const [currentStudentIndex, setCurrentStudentIndex] = useState(0)
  const [completedStudents, setCompletedStudents] = useState<Set<number>>(new Set())
  
  // Results storage for all students
  const [allResults, setAllResults] = useState<Map<number, StudentResult[]>>(new Map())
  
  // Current round tracking
  const [cardStartTimes, setCardStartTimes] = useState<Map<string, number>>(new Map())
  const [incorrectAttempts, setIncorrectAttempts] = useState<Map<string, string[]>>(new Map())

  const currentSessionStudentId = session_student_ids[currentStudentIndex]

  const { gameContents, isLoading } = useGameContents({
    theme_id: themeID || undefined,
    difficulty_level: difficulty ? Number.parseInt(difficulty) : undefined,
    category: category as any,
    question_type: questionType as any,
    applicable_game_types: ['word/image matching'],
  })

  const [shuffledCards, setShuffledCards] = useState<MatchingCardContent[]>([])
  const [selectedCards, setSelectedCards] = useState<MatchingCardContent[]>([])
  const [matchedIDs, setMatchedIDs] = useState<Set<string>>(new Set())
  const [tempWrongIDs, setTempWrongIDs] = useState<Set<string>>(new Set())
  const [gameCompleted, setGameCompleted] = useState(false)
  const [resultsSaved, setResultsSaved] = useState(false)
  const [isSaving, setIsSaving] = useState(false)

  // Get student name helper
  const getStudentName = useCallback((sessionStudentId: number) => {
    const sessionStudent = sessionStudents.find(s => s.sessionStudentId === sessionStudentId);
    if (!sessionStudent) return 'Student';
    const student = allStudents?.find(s => s.id === sessionStudent.studentId);
    return student ? `${student.first_name} ${student.last_name}` : 'Student';
  }, [sessionStudents, allStudents]);

  const imageCards = shuffledCards.filter(card => card.isImage)
  const wordCards = shuffledCards.filter(card => !card.isImage)
  const groupedCols = [
    { key: 'images', cards: imageCards },
    { key: 'words', cards: wordCards },
  ]

  // Initialize when student changes or game contents load
  useEffect(() => {
    if (gameContents.length === 0) return
    if (completedStudents.has(currentSessionStudentId)) return

    // Build cards from game contents
    const newCards: MatchingCardContent[] = gameContents.flatMap((gc) => [
      {
        id: `${gc.id}-image`,
        isImage: true,
        value: gc.answer,
        pairID: gc.id
      },
      {
        id: `${gc.id}-word`,
        isImage: false,
        value: gc.question,
        pairID: gc.id
      }
    ])

    const shuffled = [...newCards].sort(() => Math.random() - 0.5)
    setShuffledCards(shuffled)
    setSelectedCards([])
    setMatchedIDs(new Set())
    setTempWrongIDs(new Set())
    setGameCompleted(false)
    
    // Initialize start times for all content
    const startTimes = new Map<string, number>()
    gameContents.forEach(gc => {
      startTimes.set(gc.id, Date.now())
    })
    setCardStartTimes(startTimes)
    setIncorrectAttempts(new Map())
  }, [currentStudentIndex]) // Only re-run when student changes

  // Reset game function for "Play Again" button
  const resetGame = () => {
    const newCards: MatchingCardContent[] = gameContents.flatMap((gc) => [
      {
        id: `${gc.id}-image`,
        isImage: true,
        value: gc.answer,
        pairID: gc.id
      },
      {
        id: `${gc.id}-word`,
        isImage: false,
        value: gc.question,
        pairID: gc.id
      }
    ])
    const shuffled = [...newCards].sort(() => Math.random() - 0.5)
    setShuffledCards(shuffled)
    setSelectedCards([])
    setMatchedIDs(new Set())
    setTempWrongIDs(new Set())
    setGameCompleted(false)
    
    // Reset start times
    const startTimes = new Map<string, number>()
    gameContents.forEach(gc => {
      startTimes.set(gc.id, Date.now())
    })
    setCardStartTimes(startTimes)
    setIncorrectAttempts(new Map())
  }

  const canSelectCard = (card: MatchingCardContent) => {
    if (selectedCards.length === 0) return true
    const first = selectedCards[0]
    return first.isImage !== card.isImage
  }

  const handleCardClick = (card: MatchingCardContent) => {
    // De-select a selected card
    if (selectedCards.some(c => c.id === card.id)) {
      setSelectedCards(prev => prev.filter(c => c.id !== card.id))
      return
    }

    // Can't select matched cards or cards of same type
    if ((matchedIDs.has(card.id)) || (!canSelectCard(card))) return

    const newSelected = [...selectedCards, card]
    setSelectedCards(newSelected)

    if (newSelected.length === 2) {
      const [first, second] = newSelected
      if (first.pairID === second.pairID) {
        // Correct match!
        setMatchedIDs(prev => new Set(prev).add(first.id).add(second.id))
        
        // Record result for this content
        const contentId = first.pairID
        const startTime = cardStartTimes.get(contentId) || Date.now()
        const timeTaken = Math.floor((Date.now() - startTime) / 1000)
        const attempts = incorrectAttempts.get(contentId) || []
        
        const result: StudentResult = {
          content_id: contentId,
          time_taken_sec: timeTaken,
          count_of_incorrect_attempts: attempts.length,
          incorrect_attempts: attempts
        }
        
        setAllResults(prev => {
          const newMap = new Map(prev)
          const studentResults = newMap.get(currentSessionStudentId) || []
          newMap.set(currentSessionStudentId, [...studentResults, result])
          return newMap
        })
        
        setTimeout(() => setSelectedCards([]), 800)
      } else {
        // Wrong match
        const wrongSet = new Set<string>([first.id, second.id])
        setTempWrongIDs(wrongSet)

        // Track incorrect attempts
        setIncorrectAttempts(prev => {
          const newMap = new Map(prev)
          const firstAttempts = newMap.get(first.pairID) || []
          const secondAttempts = newMap.get(second.pairID) || []
          newMap.set(first.pairID, [...firstAttempts, second.value])
          newMap.set(second.pairID, [...secondAttempts, first.value])
          return newMap
        })

        setTimeout(() => {
          setTempWrongIDs(new Set())
          setSelectedCards([])
        }, 800)
      }
    }
  }

  // Check for round completion
  useEffect(() => {
    if (shuffledCards.length > 0 && matchedIDs.size === shuffledCards.length && !gameCompleted) {
      setGameCompleted(true)
      setCompletedStudents(prev => new Set(prev).add(currentSessionStudentId))
    }
  }, [matchedIDs, shuffledCards, gameCompleted, currentSessionStudentId])

  const handleSaveAllProgress = async () => {
    setIsSaving(true)
    try {
      const api = getGameResult()
      
      // Save results for all students who played
      for (const [studentId, results] of allResults.entries()) {
        for (const result of results) {
          const input: PostGameResultInput = {
            session_student_id: studentId,
            content_id: result.content_id,
            time_taken_sec: result.time_taken_sec,
            completed: true,
            count_of_incorrect_attempts: result.count_of_incorrect_attempts,
            incorrect_attempts: result.incorrect_attempts
          }
          await api.postGameResults(input)
        }
      }
      
      setResultsSaved(true)
    } catch (err) {
      console.error("Error saving results", err)
    } finally {
      setIsSaving(false)
    }
  }

  const handleNextStudent = () => {
    // Reset game state BEFORE changing student to prevent false completion
    setMatchedIDs(new Set())
    setShuffledCards([])
    setSelectedCards([])
    setGameCompleted(false)
    // Move to next student
    setCurrentStudentIndex(prev => prev + 1)
  }

  const hasMoreStudents = currentStudentIndex < session_student_ids.length - 1
  const allStudentsComplete = completedStudents.size === session_student_ids.length

  if (isLoading) {
    return (
      <AppLayout>
        <div className="min-h-screen flex items-center justify-center bg-background">
          <p className="text-secondary">Loading game content...</p>
        </div>
      </AppLayout>
    )
  }

  if (gameContents.length === 0) {
    return (
      <AppLayout>
        <div className="min-h-screen flex items-center justify-center bg-background">
          <div className="text-center bg-card p-8 rounded-sm shadow-lg max-w-md">
            <p className="text-error mb-4">No game content available for this selection.</p>
            <button
              onClick={() => router.push(effectiveSessionId ? `/sessions/${effectiveSessionId}/curriculum` : '/games')}
              className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover"
            >
              Back to Content
            </button>
          </div>
        </div>
      </AppLayout>
    )
  }

  if (gameCompleted) {
    return (
        <AppLayout>
            <div className="min-h-screen flex items-center justify-center bg-background">
                <div className="text-center bg-card p-8 rounded-sm shadow-lg max-w-md">
                    <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
                    <h2 className="mb-2">
                      {allStudentsComplete ? 'All Students Complete!' : 'Round Complete!'}
                    </h2>
                    <p className="text-secondary mb-4">
                      {getStudentName(currentSessionStudentId)} finished!
                    </p>
                    
                    {/* Show progress through students */}
                    <p className="text-sm text-muted mb-4">
                      Student {currentStudentIndex + 1} of {session_student_ids.length}
                    </p>

                    <div className="flex flex-wrap gap-3 justify-center">
                        <button
                            onClick={resetGame}
                            className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover cursor-pointer"
                        >
                            Play Again
                        </button>
                        
                        {hasMoreStudents ? (
                            <button
                                onClick={handleNextStudent}
                                className="px-6 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 cursor-pointer"
                            >
                                Next Student →
                            </button>
                        ) : (
                            <button
                                onClick={handleSaveAllProgress}
                                disabled={isSaving || resultsSaved}
                                className="px-6 py-2 bg-pink text-white rounded-lg hover:bg-pink-hover transition-colors disabled:bg-pink-disabled disabled:cursor-not-allowed"
                            >
                                {isSaving ? 'Saving...' : resultsSaved ? 'Saved!' : 'Save Progress'}
                            </button>
                        )}
                        
                        <button
                            onClick={() => router.push(effectiveSessionId ? `/sessions/${effectiveSessionId}/curriculum` : '/games')}
                            className="px-6 py-2 bg-card-hover text-primary rounded-lg hover:bg-card border border-border cursor-pointer"
                        >
                            Back to Content
                        </button>
                    </div>
                </div>
            </div>
        </AppLayout>
    )
  }

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <button onClick={() => router.push(effectiveSessionId ? `/sessions/${effectiveSessionId}/curriculum` : '/games')}
                  className="text-blue hover:text-blue-hover flex items-center gap-2 transition-colors cursor-pointer">
            ← Back to Content
          </button>
          <button onClick={resetGame}
                  className="flex items-center gap-2 text-secondary hover:text-primary transition-colors cursor-pointer">
            <RotateCw className="w-4 h-4" />
            Reset
          </button>
        </div>

        {/* Current Student Banner */}
        <div className="bg-blue text-white rounded-lg p-4 mb-6 text-center">
          <p className="text-sm opacity-90 mb-1">Current Player ({currentStudentIndex + 1} of {session_student_ids.length})</p>
          <p className="text-2xl font-bold">{getStudentName(currentSessionStudentId)}</p>
        </div>

        <h1 className="mb-4">Word-Image Match The Following!</h1>

        <div className="bg-card rounded-lg p-4 mb-6 border border-default">
          <div className="flex flex-wrap items-center gap-2 text-sm">
            <span className="text-muted">Theme: </span>
            <span className="font-medium text-primary">{themeName}</span>
            <span className="text-muted mx-2">•</span>
            <span className="text-muted">Difficulty:</span>
            <span className="font-medium text-primary">Level {difficulty}</span>
            <span className="text-muted mx-2">•</span>
            <span className="text-muted">Category:</span>
            <span className="font-medium text-primary">{CATEGORIES[category as GetGameContentsCategory]?.label || category}</span>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-8">
          {groupedCols.map(group => (
            <div key={group.key} className="flex flex-col gap-4">
              {group.cards.map(card => (
                <MatchingCard
                  key={card.id}
                  isImage={card.isImage}
                  value={card.value}
                  isSelected={selectedCards.some(c => c.id === card.id)}
                  isWrong={tempWrongIDs.has(card.id)}
                  isMatched={matchedIDs.has(card.id)}
                  onClick={() => handleCardClick(card)}
                />
              ))}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
