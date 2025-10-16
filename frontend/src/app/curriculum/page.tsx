'use client'

import { ArrowLeft, Image as ImageIcon } from 'lucide-react'
import Image from 'next/image'
import Link from 'next/link'

export default function ImagePage() {
  return (
    <div className="min-h-screen bg-background py-8">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <header className="mb-8">
          {/* Back button */}
          <Link
            href="/"
            className="inline-flex items-center gap-2 text-secondary hover:text-primary mb-4 transition-colors group"
          >
            <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
            <span className="text-sm font-medium">Back to Home</span>
          </Link>

          {/* Centered title */}
          <div className="text-center">
            <h1 className="text-3xl font-bold text-primary mb-2">Curriculum</h1>
          </div>
        </header>

        {/* Centered image container */}
        <div className="flex items-center justify-center min-h-[60vh]">
          <div className="relative">
            <Image
              src="/curriculumMap.png" // Replace with your PNG file path
              alt="Curriculum roadmap"
              width={800} // Adjust width as needed
              height={600} // Adjust height as needed
              priority
            />
          </div>
        </div>
      </div>
    </div>
  )
}
