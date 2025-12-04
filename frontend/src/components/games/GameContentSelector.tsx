// components/games/GameContentSelector.tsx
'use client'

import React from 'react'
import { ChevronRight } from 'lucide-react'
import { useThemes } from '@/hooks/useThemes'
import { useSessionContext } from '@/contexts/sessionContext'
import { 
  GetGameContentsCategory, 
  GetGameContentsQuestionType
} from '@/lib/api/theSpecialStandardAPI.schemas'
import type { Theme } from '@/lib/api/theSpecialStandardAPI.schemas'

const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
]

const CATEGORIES = {
  [GetGameContentsCategory.receptive_language]: { 
    label: 'Receptive Language', 
    icon: '👂', 
    colorClass: 'bg-blue',
  },
  [GetGameContentsCategory.expressive_language]: { 
    label: 'Expressive Language', 
    icon: '💬', 
    colorClass: 'bg-pink',
  },
  [GetGameContentsCategory.social_pragmatic_language]: { 
    label: 'Social Pragmatic Language', 
    icon: '🤝', 
    colorClass: 'bg-orange',
  },
  [GetGameContentsCategory.speech]: { 
    label: 'Speech', 
    icon: '🗣️', 
    colorClass: 'bg-blue',
  },
}

const QUESTION_TYPES = {
  [GetGameContentsQuestionType.sequencing]: 'Sequencing',
  [GetGameContentsQuestionType.following_directions]: 'Following Directions',
  [GetGameContentsQuestionType.wh_questions]: 'WH Questions',
  [GetGameContentsQuestionType.true_false]: 'True/False',
  [GetGameContentsQuestionType.concepts_sorting]: 'Concepts & Sorting',
  [GetGameContentsQuestionType.fill_in_the_blank]: 'Fill in the Blank',
  [GetGameContentsQuestionType.categorical_language]: 'Categorical Language',
  [GetGameContentsQuestionType.emotions]: 'Emotions',
  [GetGameContentsQuestionType.teamwork_talk]: 'Teamwork Talk',
  [GetGameContentsQuestionType.express_excitement_interest]: 'Express Excitement/Interest',
  [GetGameContentsQuestionType.fluency]: 'Fluency',
  [GetGameContentsQuestionType.articulation_s]: 'Articulation - S',
  [GetGameContentsQuestionType.articulation_l]: 'Articulation - L',
}

interface GameContentSelectorProps {
  onSelectionComplete: (selection: {
    theme: Theme
    difficultyLevel: number
    category: GetGameContentsCategory
    questionType: GetGameContentsQuestionType
  }) => void
  onBack?: () => void
  backLabel?: string
  initialDifficultyLevel?: number
  initialCategory?: GetGameContentsCategory
}

export function GameContentSelector({ onSelectionComplete, onBack, backLabel, initialDifficultyLevel, initialCategory }: GameContentSelectorProps) {
  const [selectedTheme, setSelectedTheme] = React.useState<Theme | null>(null)
  const [selectedCategory, setSelectedCategory] = React.useState<GetGameContentsCategory | null>(initialCategory || null)
  const [selectedQuestionType, setSelectedQuestionType] = React.useState<GetGameContentsQuestionType | null>(null)

  const { themes, isLoading: themesLoading, error: themesError, refetch: refetchThemes } = useThemes()

  // Use the passed-in difficulty level
  const difficultyLevel = initialDifficultyLevel || 1

  React.useEffect(() => {
    if (selectedTheme && selectedCategory && selectedQuestionType) {
      onSelectionComplete({
        theme: selectedTheme,
        difficultyLevel,
        category: selectedCategory,
        questionType: selectedQuestionType,
      })
    }
  }, [selectedTheme, selectedCategory, selectedQuestionType, difficultyLevel, onSelectionComplete])

  const handleReset = () => {
    setSelectedTheme(null)
    setSelectedCategory(null)
    setSelectedQuestionType(null)
  }

  // Step 1: Theme Selection
  if (!selectedTheme) {
    if (themesError) {
      return (
        <div className="min-h-screen bg-background p-8 flex items-center justify-center">
          <div className="text-center">
            {onBack && (
              <button
                onClick={onBack}
                className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
              >
                ← {backLabel || 'Back'}
              </button>
            )}
            <p className="text-error mb-4">Failed to load themes</p>
            <button 
              onClick={() => refetchThemes()}
              className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover transition-colors"
            >
              Retry
            </button>
          </div>
        </div>
      )
    }

    return (
      <div className="min-h-screen bg-background p-8">
        <div className="max-w-6xl mx-auto">
          {onBack && (
            <button
              onClick={onBack}
              className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
            >
              ← {backLabel || 'Back'}
            </button>
          )}
          <h1 className="mb-2">Select a Theme</h1>
          <p className="text-secondary mb-8">Level: {difficultyLevel}</p>
          {themesLoading ? (
            <div className="text-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue mx-auto mb-4"></div>
              <p className="text-muted">Loading themes...</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {themes.map((theme) => (
                <button
                  key={theme.id}
                  onClick={() => setSelectedTheme(theme)}
                  className="bg-card cursor-pointer rounded-lg shadow-md p-6 hover:shadow-lg transition-all duration-200 text-left group hover:bg-card-hover border border-default hover:border-hover"
                >
                  <div className="flex items-center justify-between">
                    <h3>{theme.name}</h3>
                    <ChevronRight className="w-5 h-5 text-muted group-hover:text-primary transition-colors" />
                  </div>
                </button>
              ))}
            </div>
          )}
        </div>
      </div>
    )
  }

  // Step 2: Category Selection
  if (!selectedCategory) {
    return (
      <div className="min-h-screen bg-background p-8">
        <div className="max-w-4xl mx-auto">
          {onBack && (
            <button
              onClick={onBack}
              className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
            >
              ← {backLabel || 'Back'}
            </button>
          )}
          <h1 className="mb-2">Select a Category</h1>
          <p className="text-secondary mb-8">
            Theme: {theme.name} ({MONTHS[currentMonth]} {currentYear}) • Level {difficultyLevel}
          </p>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {Object.entries(CATEGORIES).map(([key, category]) => (
              <button
                key={key}
                onClick={() => setSelectedCategory(key as GetGameContentsCategory)}
                className="bg-card rounded-lg shadow-md p-8 hover:shadow-lg transition-all duration-200 text-left group hover:bg-card-hover border border-default hover:border-hover"
              >
                <div className="flex items-center gap-4">
                  <div className={`w-16 h-16 ${category.colorClass} rounded-full flex items-center justify-center text-2xl text-white`}>
                    {category.icon}
                  </div>
                  <div className="flex-1">
                    <h3>{category.label}</h3>
                  </div>
                  <ChevronRight className="w-5 h-5 text-muted group-hover:text-primary transition-colors" />
                </div>
              </button>
            ))}
          </div>
        </div>
      </div>
    )
  }

  // Step 2: Question Type Selection
  const category = CATEGORIES[selectedCategory]
  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-4xl mx-auto">
        <button
          onClick={() => setSelectedCategory(null)}
          className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
        >
          ← Back to Categories
        </button>
        <h1 className="mb-2">Select Question Type</h1>
        <p className="text-secondary mb-8">
          Theme: {theme.name} • Category: {category.label} • Level {difficultyLevel}
        </p>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {Object.entries(QUESTION_TYPES).map(([key, label]) => (
            <button
              key={key}
              onClick={() => setSelectedQuestionType(key as GetGameContentsQuestionType)}
              className="bg-card rounded-lg shadow-md p-4 hover:shadow-lg transition-all duration-200 text-left group flex items-center justify-between hover:bg-card-hover border border-default hover:border-hover"
            >
              <span className="font-medium text-primary">{label}</span>
              <ChevronRight className="w-5 h-5 text-muted group-hover:text-primary transition-colors" />
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}