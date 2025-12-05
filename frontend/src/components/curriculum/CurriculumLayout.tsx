'use client'

import type { ReactNode } from 'react'
import { ArrowLeft } from 'lucide-react'
import Link from 'next/link'
import { useSessionContext } from '@/contexts/sessionContext'

interface CurriculumLayoutProps {
  children: ReactNode
  title: string
  subtitle?: string
  backHref: string
  backLabel?: string
  headerContent?: ReactNode
}

export default function CurriculumLayout({
  children,
  title,
  subtitle,
  backHref,
  backLabel = 'Back',
  headerContent,
}: CurriculumLayoutProps) {

  const { clearSession } = useSessionContext()

  return (
    <div className="min-h-screen">
      {/* Header */}
      <div className="px-8 py-6">
        <div className="flex items-center justify-between mb-4">
          <Link
            href={backHref}
            onClick={backHref === '/' ? () => clearSession() : undefined}
            className="inline-flex items-center gap-2"
          >
            <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
            <span className="text-xl font-semibold cursor-pointer">{backLabel}</span>
          </Link>

          {headerContent}
        </div>

        <div>
          <h1 className="text-5xl font-bold mb-2">{title}</h1>
          {subtitle && <p className="text-xl opacity-90">{subtitle}</p>}
        </div>
      </div>

      {/* Content */}
      <div className="">
        {children}
      </div>
    </div>
  )
}

