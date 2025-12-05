'use client'
import React from 'react'
import clsx from 'clsx'

interface MatchingCardProps {
  isImage: boolean
  value: string
  isSelected: boolean
  isWrong: boolean
  isMatched: boolean
  onClick: () => void
  disabled?: boolean
}

export default function MatchingCard({
  isImage,
  value,
  isSelected = false,
  isWrong = false,
  isMatched = false,
  onClick,
  disabled = false,
}: MatchingCardProps) {
  return (
    <button 
      onClick={onClick}
      disabled={disabled || isMatched}
      className={clsx(
        "relative flex items-center justify-center rounded-lg transition-all",
        "hover:shadow-md active:scale-[0.97]",
        "w-full h-[200px]",

        isImage ? "p-6" : "p-6",

        // BACKGROUND STATES (no borders)
        isMatched && "bg-green-200 cursor-not-allowed",
        isWrong && "bg-red-200",
        isSelected && !isWrong && !isMatched && "bg-blue-200",

        // Default background
        !isMatched && !isWrong && !isSelected && "bg-card"
      )}
    >
      {!isImage ? (
        <span className="text-3xl font-semibold text-foreground text-center leading-relaxed break-words max-w-full px-2">
          {value}
        </span>
      ) : (
        <img 
          src={value} 
          alt="Matching card image" 
          className="w-full h-full object-contain"
        />
      )}
    </button>
  )
}
