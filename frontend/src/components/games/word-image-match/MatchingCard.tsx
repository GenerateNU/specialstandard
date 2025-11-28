'use client'

import React from 'react'
import clsx from 'clsx'

interface MatchingCardProps {
  type: 'word' | 'image'
  value: string
  isSelected: boolean
  isWrong: boolean
  isMatched: boolean
  onClick: () => void
  disabled?: boolean
}

export default function MatchingCard({
  type,
  value,
  isSelected = false,
  isWrong = false,
  isMatched = false,
  onClick,
  disabled = false,
 }: MatchingCardProps) {
  return (
    <button onClick={onClick}
            disabled={disabled || isMatched}
            className={clsx("relative flex items-center justify-center rounded-lg border p-4 bg-card transition-all",
              "hover:shadow-md active:scale-[0.97]",
              isMatched && "border-green-500 border-2 bg-green/15 cursor-not-allowed",
              isWrong && "border-red-500 border-2 shadow-red-500/40",
              isSelected && !isWrong && !isMatched && "border-blue shadow-lg bg-card-hover"
            )}>
      {type === 'word' ? (
        <span className="text-xl font-semibold text-foreground">{value}</span>
      ) : (
        // TODO: Change it back for image display!
        <span className="text-lg font-medium text-foreground">{value}</span>
        // <img src={value} alt="Image" className="w-20 h-20 object-contain"/>
      )}
    </button>
  )
}