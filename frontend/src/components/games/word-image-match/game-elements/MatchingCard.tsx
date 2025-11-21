'use client'

import React from 'react'
import clsx from 'clsx'

interface MatchingCardProps {
  type: 'word' | 'image'
  value: string
  isSelected: boolean
  isMatched: boolean
  onClick: () => void
  disabled?: boolean
}

export default function MatchingCard({
  type,
  value,
  isSelected,
  isMatched,
  onClick,
  disabled = false,
 }: MatchingCardProps) {
  return (
    <button onClick={onClick}
            disabled={disabled || isMatched}
            className={clsx("relative flex items-center justify-center rounded-lg border p-4 bg-card transition-all",
              "hover:shadow-md active:scale-[0.97]",
              isSelected && "border-blue shadow-lg bg-card-hover",
              isMatched && "border-green bg-green/15 cursor-not-allowed"
            )}>
      {type === 'word' ? (
        <span className="text-lg font-medium text-foreground">{value}</span>
      ) : (
        <img src={value} alt="Image" className="w-20 h-20 object-contain"/>
      )}
    </button>
  )
}