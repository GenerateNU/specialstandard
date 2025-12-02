// components/games/GameContentSelector.tsx
'use client'

import React from 'react'
import { ChevronRight, FileText, Loader2} from 'lucide-react'
import { useThemes } from '@/hooks/useThemes'
import { 
  GetGameContentsCategory, 
  GetGameContentsQuestionType
} from '@/lib/api/theSpecialStandardAPI.schemas'
import type { Theme, GameContent } from '@/lib/api/theSpecialStandardAPI.schemas'
import { getGameContent } from '@/lib/api/game-content'


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
}

export function GameContentSelector({ onSelectionComplete, onBack, backLabel }: GameContentSelectorProps) {
  const [selectedTheme, setSelectedTheme] = React.useState<Theme | null>(null)
  const [selectedDifficulty, setSelectedDifficulty] = React.useState<number | null>(null)
  const [selectedCategory, setSelectedCategory] = React.useState<GetGameContentsCategory | null>(null)
  const [selectedQuestionType, setSelectedQuestionType] = React.useState<GetGameContentsQuestionType | null>(null)
  const [pdfExercises, setPdfExercises] = React.useState<GameContent[]>([])
  const [isLoadingPdfs, setIsLoadingPdfs] = React.useState(false)

  const { themes, isLoading: themesLoading, error: themesError, refetch: refetchThemes } = useThemes()

  React.useEffect(() => {
    if (selectedTheme && selectedDifficulty && selectedCategory && selectedQuestionType) {
      onSelectionComplete({
        theme: selectedTheme,
        difficultyLevel: selectedDifficulty,
        category: selectedCategory,
        questionType: selectedQuestionType,

      })
    }
  }, [selectedTheme, selectedDifficulty, selectedCategory, selectedQuestionType, onSelectionComplete])

  React.useEffect(() => {
      const fetchPdfs = async () => {
        if (!selectedTheme || !selectedDifficulty || !selectedCategory) return

        setIsLoadingPdfs(true)
        
        try {
          const api = getGameContent()
          
          const response = await api.getGameContents({
            exercise_type: 'pdf',
            theme_id: selectedTheme.id,
            difficulty_level: selectedDifficulty,
            category: selectedCategory
          })
          
          if (response && Array.isArray(response)) {
            const pdfItems = response.filter((item: GameContent) => item.exercise_type === 'pdf')
            
            // Remove duplicates by id
            const uniqueContents = Array.from(
              new Map(pdfItems.map(item => [item.id, item])).values()
            )
            
            setPdfExercises(uniqueContents)
          }
        } catch (err) {
          console.error('💥 Error fetching pdf exercises:', err)
        } finally {
          setIsLoadingPdfs(false)
        }
      }
      fetchPdfs()
    }, [selectedTheme, selectedDifficulty, selectedCategory])
  
  

  const handleReset = () => {
    setSelectedDifficulty(null)
    setSelectedTheme(null)
    setSelectedCategory(null)
    setSelectedQuestionType(null)
    setPdfExercises([])

  }

  const handleDownloadPdf = (pdfUrl: string) => { // TODO: also put the answers??
    window.open(pdfUrl, '_blank')
  }

  // Step 0: Difficulty Selection
  if (!selectedDifficulty) {
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
          <h1 className="mb-8">Select Difficulty Level</h1>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {[1, 2, 3].map((level) => (
              <button
                key={level}
                onClick={() => setSelectedDifficulty(level)}
                className="bg-card rounded-lg shadow-md p-8 hover:shadow-lg transition-all duration-200 group hover:bg-card-hover border border-default hover:border-hover text-left"
              >
                <h3 className="font-bold text-lg mb-2">Level {level}</h3>
                <p className="text-secondary text-sm">
                  {level === 1 && 'Beginner - Perfect for starting out'}
                  {level === 2 && 'Intermediate - Build your skills'}
                  {level === 3 && 'Advanced - Challenge yourself'}
                </p>
              </button>
            ))}
          </div>
        </div>
      </div>
    )
  }

  // Step 1: Theme Selection
  if (!selectedTheme) {
    if (themesError) {
      return (
        <div className="min-h-screen bg-background p-8 flex items-center justify-center">
          <div className="text-center">
            <button
              onClick={() => setSelectedDifficulty(null)}
              className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
            >
              ← Back to Difficulty
            </button>
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
          <button
            onClick={() => setSelectedDifficulty(null)}
            className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
          >
            ← Back to Difficulty
          </button>
          <h1 className="mb-2">Select a Theme</h1>
          <p className="text-secondary mb-8">Difficulty Level: {selectedDifficulty}</p>
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
                  className="bg-card rounded-lg shadow-md p-6 hover:shadow-lg transition-all duration-200 text-left group hover:bg-card-hover border border-default hover:border-hover"
                >
                  <div className="flex items-center justify-between">
                    <h3>{theme.name}</h3>
                    <ChevronRight className="w-5 h-5 text-muted group-hover:text-primary transition-colors" />
                  </div>
                  <p className="text-secondary mt-2">
                    {`Explore ${theme.name} themed content`}
                  </p>
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
          <button
            onClick={() => {
              setSelectedTheme(null)
              handleReset()
            }}
            className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
          >
            ← Back to Themes
          </button>
          <h1 className="mb-2">Select a Category</h1>
          <p className="text-secondary mb-8">Theme: {selectedTheme?.name}</p>
          
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

  // Step 3: Question Type Selection
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
<h1 className="mb-2">Select Exercise Type</h1>
        <p className="text-secondary mb-8">
          Theme: {selectedTheme?.name} / Category: {category.label}
        </p>
        
        {/* PDF Exercises Section */}
        <div className="mb-12">
          <h2 className="text-lg font-semibold mb-4 text-primary">PDF Exercises</h2>
          {isLoadingPdfs ? (
            <div className="bg-card rounded-lg shadow-md p-6 flex items-center justify-center">
              <Loader2 className="w-5 h-5 animate-spin text-accent mr-2" />
              <span className="text-secondary">Loading PDF exercises...</span>
            </div>
          ) : pdfExercises.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {pdfExercises.map((pdf) => (
                <button
                  key={pdf.id}
                  onClick={() => handleDownloadPdf(pdf.answer)}
                  className="bg-card rounded-lg shadow-md p-4 hover:shadow-lg transition-all duration-200 text-left group flex items-center justify-between hover:bg-card-hover border border-default hover:border-hover"
                >
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-purple rounded-full flex items-center justify-center text-white">
                      <FileText className="w-5 h-5" />
                    </div>
                    <span className="font-medium text-primary">
                      {pdf.question_type ? pdf.question_type.split('_').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ') : 'PDF Exercise'}
                    </span>
                  </div>
                  <Download className="w-5 h-5 text-muted group-hover:text-primary transition-colors" />
                </button>
              ))}
            </div>
          ) : (
            <div className="bg-card rounded-lg shadow-md p-6 text-center text-secondary">
              No PDF exercises available for this selection
            </div>
          )}
        </div>

        {/* Interactive Games Section */}
        <div>
          <h2 className="text-lg font-semibold mb-4 text-primary">Interactive Games</h2>
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
    </div>
  )
}
