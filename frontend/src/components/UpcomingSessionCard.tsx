'use client'

import { ArrowRight } from 'lucide-react'

interface UpcomingSessionCardProps {
  sessionName: string
  startTime: string
  endTime: string
  date: string
  onClick?: () => void
}

export default function UpcomingSessionCard({
  sessionName,
  startTime,
  endTime,
  date,
  onClick,
}: UpcomingSessionCardProps) {
  return (
    <button
      onClick={onClick}
      className="h-16 bg-accent hover:bg-accent-hover transition w-full rounded-lg p-2 flex flex-row items-center justify-between text-background cursor-pointer group"
    >
      <div className="flex flex-col text-xs font-normal justify-start text-left">
        <span>{sessionName}</span>
        <span>
          {startTime}
          {' '}
          –
          {' '}
          {endTime}
        </span>
        <span>{date}</span>
      </div>
      <ArrowRight size={18} className="shrink-0 transition-transform duration-200 ease-out group-hover:translate-x-1 will-change-transform" />
    </button>
  )
}
