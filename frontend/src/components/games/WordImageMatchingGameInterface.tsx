'use client'

import React, {useEffect, useState} from "react";
import MatchingCard from "@/components/games/word-image-match/MatchingCard";
import { useRouter } from 'next/navigation'
import {RotateCw} from "lucide-react";
import type {GetGameContentsCategory} from "@/lib/api/theSpecialStandardAPI.schemas";
import {CATEGORIES} from "@/components/games/FlashcardGameInterface";
import {useGameContents} from "@/hooks/useGameContents";

export interface MatchingCardContent {
  id: string
  isImage: boolean
  value: string
  pairID: string
}

interface WordImageMatchingGameInterfaceProps {
  contents: MatchingCardContent[]
  session_student_id: number
  session_id: string
  student_id: string
  themeID: string
  themeName: string | null
  difficulty: string
  category: string
  questionType: string
}

export default function WordImageMatchingGameInterface({
  contents,
  session_student_id,
  session_id,
  student_id,
  themeID,
  themeName,
  difficulty,
  category,
  questionType
}: WordImageMatchingGameInterfaceProps) {
  const [shuffledCards, setShuffledCards] = useState<MatchingCardContent[]>([])
  const [selectedCards, setSelectedCards] = useState<MatchingCardContent[]>([])
  const [matchedIDs, setMatchedIDs] = useState<Set<string>>(new Set())
  const [gameCompleted, setGameCompleted] = useState(false)

  const imageCards = shuffledCards.filter(card => card.isImage)
  const wordCards = shuffledCards.filter(card => !card.isImage)
  const groupedCols = [
    { key: 'images', cards: imageCards },
    { key: 'words', cards: wordCards },
  ]

  const router = useRouter()

  const [tempWrongIDs, setTempWrongIDs] = useState<Set<string>>(new Set())

  const { gameContents, isLoading, error } = useGameContents({
    theme_id: themeID || undefined,
    difficulty_level: difficulty ? Number.parseInt(difficulty) : undefined,
    category: category as any,
    question_type: questionType as any,
  })
  const cards: MatchingCardContent[] = gameContents.flatMap((gc) => [
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

  useEffect(() => {
    const shuffled = [...cards].sort(() => Math.random() - 0.5)
    setShuffledCards(shuffled)
    setSelectedCards([])
    setMatchedIDs(new Set())
    setTempWrongIDs(new Set())
    setGameCompleted(false)
  }, [gameContents]);

  const canSelectCard = (card: MatchingCardContent) => {
    if (selectedCards.length === 0) return true
    const first = selectedCards[0]
    return first.isImage !== card.isImage
  }

  const handleCardClick = (card: MatchingCardContent) => {
    // De-select a selected card.
    if (selectedCards.some(c => c.id === card.id)) {
      setSelectedCards(prev => prev.filter(c => c.id !== card.id))
      return
    }

    // Can't select matched cards or cards of same type (two images, etc.)
    if ((matchedIDs.has(card.id)) || (!canSelectCard(card))) {
      return
    }

    const newSelected = [...selectedCards, card]
    setSelectedCards(newSelected)

    if (newSelected.length === 2) {
      const [first, second] = newSelected
      if (first.pairID === second.pairID) {
        setMatchedIDs(prev => new Set(prev).add(first.id).add(second.id))
        setTimeout(() => setSelectedCards([]), 800)
      } else {
        const wrongSet = new Set<string>([first.id, second.id])
        setTempWrongIDs(wrongSet)

        setTimeout(() => {
          setTempWrongIDs(new Set())
          setSelectedCards([])
        }, 800)
      }
    }
  }

  useEffect(() => {
    if (shuffledCards.length > 0 && matchedIDs.size === shuffledCards.length) {
      setGameCompleted(true)
    }
  }, [matchedIDs, shuffledCards]);

  const resetGame = () => {
    const shuffled = [...cards].sort(() => Math.random() - 0.5)
    setShuffledCards(shuffled)
    setSelectedCards([])
    setMatchedIDs(new Set())
    setGameCompleted(false)
  }

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <button onClick={() => router.back()}
                  className="text-blue hover:text-blue-hover flex items-center gap-2 transition-colors cursor-pointer">
            ← Back
          </button>
          <button onClick={resetGame}
                  className="flex items-center gap-2 text-secondary hover:text-primary transition-colors cursor-pointer">
            <RotateCw className="w-4 h-4" />
            Reset
          </button>
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
                  id={card.id}
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

        {gameCompleted && (
          <div className="mt-8 p-6 bg-green/20 border border-green rounded-lg text-center">
            <p className="text-green-500 font-semibold text-xl">You matched all the cards!</p>
            <button onClick={resetGame}
                    className="mt-4 px-6 bg-green-300 text-white rounded-lg hover:bg-green-hover
                               transition-colors">
              Play Again!
            </button>
          </div>
        )}
      </div>
    </div>
  )
}