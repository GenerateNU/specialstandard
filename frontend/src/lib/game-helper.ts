// File: lib/game-helpers.ts

import { 
  GetGameContentsCategory, 
  GetGameContentsQuestionType 
} from '@/lib/api/theSpecialStandardAPI.schemas'

// Define which question types are available for each category
export const CATEGORY_QUESTION_MAPPING: Record<GetGameContentsCategory, GetGameContentsQuestionType[]> = {
  [GetGameContentsCategory.receptive_language]: [
    GetGameContentsQuestionType.following_directions,
    GetGameContentsQuestionType.wh_questions,
    GetGameContentsQuestionType.true_false,
    GetGameContentsQuestionType.concepts_sorting,
    GetGameContentsQuestionType.categorical_language,
  ],
  [GetGameContentsCategory.expressive_language]: [
    GetGameContentsQuestionType.sequencing,
    GetGameContentsQuestionType.fill_in_the_blank,
    GetGameContentsQuestionType.wh_questions,
    GetGameContentsQuestionType.categorical_language,
  ],
  [GetGameContentsCategory.social_pragmatic_language]: [
    GetGameContentsQuestionType.emotions,
    GetGameContentsQuestionType.teamwork_talk,
    GetGameContentsQuestionType.express_excitement_interest,
  ],
  [GetGameContentsCategory.speech]: [
    GetGameContentsQuestionType.fluency,
    GetGameContentsQuestionType.articulation_s,
    GetGameContentsQuestionType.articulation_l,
  ],
}

// Helper function to get available question types for a category
export function getQuestionTypesForCategory(category: GetGameContentsCategory): GetGameContentsQuestionType[] {
  return CATEGORY_QUESTION_MAPPING[category] || []
}

// Helper to check if a question type is available for a category
export function isQuestionTypeAvailableForCategory(
  category: GetGameContentsCategory,
  questionType: GetGameContentsQuestionType
): boolean {
  const availableTypes = getQuestionTypesForCategory(category)
  return availableTypes.includes(questionType)
}

// Get a display-friendly category name
export function getCategoryDisplayName(category: GetGameContentsCategory): string {
  const displayNames: Record<GetGameContentsCategory, string> = {
    [GetGameContentsCategory.receptive_language]: 'Receptive Language',
    [GetGameContentsCategory.expressive_language]: 'Expressive Language',
    [GetGameContentsCategory.social_pragmatic_language]: 'Social Pragmatic Language',
    [GetGameContentsCategory.speech]: 'Speech',
  }
  return displayNames[category] || category
}

// Get a display-friendly question type name
export function getQuestionTypeDisplayName(questionType: GetGameContentsQuestionType): string {
  const displayNames: Record<GetGameContentsQuestionType, string> = {
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
  return displayNames[questionType] || questionType
}

// Helper to shuffle an array (useful for randomizing flashcards)
export function shuffleArray<T>(array: T[]): T[] {
  const shuffled = [...array]
  for (let i = shuffled.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]]
  }
  return shuffled
}

// Helper to save/load game progress (optional)
export const GameProgress = {
  save: (gameId: string, progress: any) => {
    if (typeof window !== 'undefined') {
      localStorage.setItem(`game_progress_${gameId}`, JSON.stringify(progress))
    }
  },
  
  load: (gameId: string) => {
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem(`game_progress_${gameId}`)
      return saved ? JSON.parse(saved) : null
    }
    return null
  },
  
  clear: (gameId: string) => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem(`game_progress_${gameId}`)
    }
  }
}