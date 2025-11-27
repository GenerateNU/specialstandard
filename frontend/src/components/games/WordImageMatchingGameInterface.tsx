'use client'

import {useEffect, useState} from "react";
import MatchingCard from "@/components/games/word-image-match/MatchingCard";

export interface MatchingCardContent {
  id: string,
  type: 'word' | 'image'
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

export default function WordImageMatchingGameInterface({ contents }: WordImageMatchingGameInterfaceProps) {
  const [shuffledCards, setShuffledCards] = useState<MatchingCardContent[]>([])
  const [selectedCards, setSelectedCards] = useState<MatchingCardContent[]>([])
  const [matchedIDs, setMatchedIDs] = useState<Set<string>>(new Set())
  const [gameCompleted, setGameCompleted] = useState(false)

  useEffect(() => {
    const shuffled = [...contents].sort(() => Math.random() - 0.5)
      setShuffledCards(shuffled)
      setSelectedCards([])
      setMatchedIDs(new Set())
      setGameCompleted(false)
  }, [contents]);

  const handleCardClick = (card: MatchingCardContent) => {
    if (selectedCards.find(c => c.id === card.id) || matchedIDs.has(card.id)) return
    const newSelected = [...selectedCards, card]
    setSelectedCards(newSelected)

    if (newSelected.length === 2) {
      const [first, second] = newSelected
      if (first.pairID === second.pairID) {
        setMatchedIDs(prev => new Set(prev).add(first.id).add(second.id))
      }
      setTimeout(() => setSelectedCards([]), 800)
    }
  }

  useEffect(() => {
    if (shuffledCards.length > 0 && matchedIDs.size === shuffledCards.length) {
      setGameCompleted(true)
    }
  }, [matchedIDs, shuffledCards]);

  const resetGame = () => {
    const shuffled = [...contents].sort(() => Math.random() - 0.5)
    setShuffledCards(shuffled)
    setSelectedCards([])
    setMatchedIDs(new Set())
    setGameCompleted(false)
  }

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between mb-6">
          <h1 className="text-2xl font-bold">Word-Image Matching Game</h1>
          <button onClick={resetGame}
                  className="px-4 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover
                             transition-colors">
            Reset
          </button>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
          {shuffledCards.map(card => (
            <MatchingCard
              key={card.id}
              type={card.type}
              value={card.value}
              isSelected={selectedCards.some(c => c.id === card.id)}
              isMatched={matchedIDs.has(card.id)}
              onClick={() => handleCardClick(card)}
            />
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