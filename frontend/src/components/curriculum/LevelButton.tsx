'use client'

import { cn } from '@/lib/utils'

interface LevelButtonProps {
  level: number
  onClick: () => void
  isSelected?: boolean
  className?: string
}

export default function LevelButton({ level, onClick, isSelected = false, className }: LevelButtonProps) {
  return (
    <button
      onClick={onClick}
      className={cn(
        'w-full px-8 py-6 rounded-full text-white font-bold text-xl transition-all hover:scale-105 hover:shadow-xl cursor-pointer',
        isSelected 
          ? 'bg-pink-disabled scale-105 shadow-xl ring-4 ring-pink/30'
          : 'bg-pink hover:bg-pink-hover',
        className,
      )}
    >
      Level {level}
    </button>
  )
}

