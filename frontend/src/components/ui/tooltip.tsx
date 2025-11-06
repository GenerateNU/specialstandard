'use client'

import React, { useState } from 'react'
import { createPortal } from 'react-dom'

interface TooltipProps {
  content: string
  children: React.ReactNode
  /** When false the tooltip won't be rendered and children are returned as-is */
  enabled?: boolean
}

export default function Tooltip({ content, children, enabled = true }: TooltipProps) {
  const [isVisible, setIsVisible] = useState(false)
  const [isMounted, setIsMounted] = useState(false)
  const [position, setPosition] = useState({ top: 0, left: 0 })
  const wrapperRef = React.useRef<HTMLDivElement>(null)

  if (!enabled) {
    return <>{children}</>
  }

  const handleMouseEnter = () => {
    if (wrapperRef.current) {
      const rect = wrapperRef.current.getBoundingClientRect()
      setPosition({
        top: rect.top + rect.height / 2,
        left: rect.right + 12,
      })
      setIsVisible(true)
      setTimeout(() => setIsMounted(true), 10)
    }
  }

  const handleMouseLeave = () => {
    setIsMounted(false)
    setTimeout(() => setIsVisible(false), 200)
  }

  return (
    <>
      <div
        ref={wrapperRef}
        className="relative inline-flex items-center"
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
      >
        {children}
      </div>

      {isVisible
        && typeof window !== 'undefined'
        && createPortal(
          <div
            role="tooltip"
            className={`pointer-events-none fixed whitespace-nowrap 
                rounded-lg bg-foreground text-white text-sm px-3 
                py-1.5 z-50 shadow-lg -translate-y-1/2 transition-opacity 
                duration-200 ${isMounted ? 'opacity-100' : 'opacity-0'}`}
            style={{
              top: `${position.top}px`,
              left: `${position.left}px`,
            }}
          >
            {content}
          </div>,
          document.body,
        )}
    </>
  )
}
