'use client'

import type { ReactNode } from 'react'
import { ArrowLeft } from 'lucide-react'
import Link from 'next/link'

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
  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="bg-blue text-white px-8 py-6">
        <div className="flex items-center justify-between mb-4">
          <Link
            href={backHref}
            className="inline-flex items-center gap-2 text-white hover:text-white/80 transition-colors group"
          >
            <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
            <span className="text-sm font-medium">{backLabel}</span>
          </Link>

          {headerContent}
        </div>

        <div>
          <h1 className="text-5xl font-bold mb-2">{title}</h1>
          {subtitle && <p className="text-xl opacity-90">{subtitle}</p>}
        </div>
      </div>

      {/* Content */}
      <div className="p-8">
        {children}
      </div>
    </div>
  )
}

