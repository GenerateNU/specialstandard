// components/games/GameContentSelector.tsx
'use client'

import React from 'react'
import { ChevronRight } from 'lucide-react'
import { 
  GetGameContentsCategory, 
  GetGameContentsQuestionType
} from '@/lib/api/theSpecialStandardAPI.schemas'
import type { Theme } from '@/lib/api/theSpecialStandardAPI.schemas'

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
  theme: Theme
}

export function GameContentSelector({ onSelectionComplete, onBack, backLabel, initialDifficultyLevel, initialCategory, theme }: GameContentSelectorProps) {
  const [selectedCategory, setSelectedCategory] = React.useState<GetGameContentsCategory | null>(initialCategory || null)
  const [selectedQuestionType, setSelectedQuestionType] = React.useState<GetGameContentsQuestionType | null>(null)

  // Use the passed-in difficulty level
  const difficultyLevel = initialDifficultyLevel || 1

  React.useEffect(() => {
    if (theme && selectedCategory && selectedQuestionType) {
      onSelectionComplete({
        theme,
        difficultyLevel,
        category: selectedCategory,
        questionType: selectedQuestionType,
      })
    }
  }, [theme, selectedCategory, selectedQuestionType, difficultyLevel, onSelectionComplete])

  // Step 1: Category Selection
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
            Theme: {theme.name} • Level {difficultyLevel}
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
        
        <div className="bg-card border-2 border-default rounded-4xl p-8">
          <div className='pb-2'>
          Game – Question Type
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {Object.entries(QUESTION_TYPES).map(([key, label]) => (
              <button
                key={key}
                onClick={() => setSelectedQuestionType(key as GetGameContentsQuestionType)}
                className="bg-pink hover:bg-pink-hover text-white p-6 rounded-xl font-semibold transition-all hover:scale-102 cursor-pointer text-left"
              >
                {label}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}