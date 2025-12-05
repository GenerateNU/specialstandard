'use client'

import { useEffect, useState } from 'react'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { ChevronRight } from 'lucide-react'
import Link from 'next/link'
import { Avatar } from '@/components/ui/avatar'
import { getAvatarName, getAvatarVariant } from '@/lib/avatarUtils'

import { SessionRatingCategory, SessionRatingLevel } from '@/lib/api/theSpecialStandardAPI.schemas'
import { useSessionStudent, useSessionStudents } from '@/hooks/useSessionStudents'
import { Slider } from '@/components/ui/slider'

interface RateStudentProps {
  sessionId: string
  studentId: string
  sessionStudentId: number
  firstName: string
  lastName: string
}

const RATING_CATEGORIES = [
  { id: SessionRatingCategory.visual_cue, label: 'Visual Cues' },
  { id: SessionRatingCategory.verbal_cue, label: 'Verbal Cues' },
  { id: SessionRatingCategory.gestural_cue, label: 'Gesture Cues' },
  { id: SessionRatingCategory.engagement, label: 'Engagement' },
]

const RATING_LEVELS = [
  { id: SessionRatingLevel.minimal, label: 'Minimal', value: 0 },
  { id: SessionRatingLevel.moderate, label: 'Moderate', value: 1 },
  { id: SessionRatingLevel.maximal, label: 'Maximal', value: 2 },
]

const levelToValue = (level: SessionRatingLevel | null): number => {
  if (level === null) return 0
  const found = RATING_LEVELS.find(l => l.id === level)
  return found?.value ?? 0
}

const valueToLevel = (value: number): SessionRatingLevel => {
  const found = RATING_LEVELS.find(l => l.value === value)
  return found?.id ?? SessionRatingLevel.minimal
}


export default function RateStudent({
  sessionId,
  studentId,
  sessionStudentId,
  firstName,
  lastName,
}: RateStudentProps) {
  const { sessionStudent, isLoading } = useSessionStudent(sessionStudentId, sessionId)
  const { updateSessionStudent, isUpdating } = useSessionStudents()
  const avatarVariant = getAvatarVariant(studentId)
  
  const [notes, setNotes] = useState(sessionStudent?.notes || '')
  const [ratings, setRatings] = useState<Record<SessionRatingCategory, SessionRatingLevel | null>>({
    [SessionRatingCategory.visual_cue]: null,
    [SessionRatingCategory.verbal_cue]: null,
    [SessionRatingCategory.gestural_cue]: null,
    [SessionRatingCategory.engagement]: null,
  })
  const [hasChanges, setHasChanges] = useState(false)

  useEffect(() => {
    if (sessionStudent) {
      setNotes(sessionStudent.notes || '')
      
      const ratingsMap: Record<SessionRatingCategory, SessionRatingLevel | null> = {
        [SessionRatingCategory.visual_cue]: null,
        [SessionRatingCategory.verbal_cue]: null,
        [SessionRatingCategory.gestural_cue]: null,
        [SessionRatingCategory.engagement]: null,
      }
      
      sessionStudent.ratings?.forEach((rating: { category: SessionRatingCategory; level: SessionRatingLevel }) => {
        const category = rating.category as SessionRatingCategory
        ratingsMap[category] = rating.level
      })
      
      setRatings(ratingsMap)
      setHasChanges(false)
    }
  }, [sessionStudent?.session_student_id])

  const saveChanges = async () => {
    try {
      const ratingsArray = Object.entries(ratings)
        .filter(([_, level]) => level !== null)
        .map(([category, level]) => ({
          category: category as SessionRatingCategory,
          level: level as SessionRatingLevel,
        }))

      await updateSessionStudent({
        session_id: sessionId,
        student_id: studentId,
        notes,
        ratings: ratingsArray,
      })
      
      setHasChanges(false)
    } catch (error) {
      console.error('Failed to save ratings', error)
    }
  }

  const handleRatingChange = (category: SessionRatingCategory, value: number) => {
    const currentLevel = ratings[category]
    const newLevel = valueToLevel(value)
    
    // Toggle to null if clicking the same value
    const finalLevel = currentLevel === newLevel ? null : newLevel
    
    const newRatings = { ...ratings, [category]: finalLevel }
    setRatings(newRatings)
    setHasChanges(true)
  }

  const handleNotesChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setNotes(e.target.value)
    setHasChanges(true)
  }

  if (isLoading) {
    return (
      <div className="w-full max-w-4xl mx-auto">
        <div className="bg-card rounded-3xl p-8 shadow-sm border border-border/50">
          <div className="text-center text-muted-foreground">Loading student data...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="w-full max-w-4xl mx-auto">
      <div className="bg-card rounded-3xl p-8 shadow-sm border border-border/50">
        <div className="mb-8 flex flex-col md:flex-row md:items-start md:justify-between gap-6">
          {/* Left: Avatar + Name */}
          <div className="flex flex-col items-center text-center px-11">
            <Avatar
              name={getAvatarName(firstName, lastName, studentId)}
              variant={avatarVariant}
              className="w-28 h-28 mb-3 border border-pink"
              title={`${firstName} ${lastName}`}
            />
            <h2 className="text-2xl font-bold">{`${firstName} ${lastName}`}</h2>
          </div>
          {/* Right: Notes */}
          <div className="flex-1">
            <h2 className="text-sm font-medium text-muted-foreground -mt-0.5 mb-2">
              Student Notes
            </h2>
            <Textarea
              placeholder={`Notes for ${firstName} ${lastName}...`}
              className="min-h-[120px] resize-none bg-background border-border/50 focus:border-primary/50 text-lg"
              value={notes}
              onChange={handleNotesChange}
            />
          </div>
        </div>

        <div className="space-y-8">         
          {RATING_CATEGORIES.map((category) => (
            <div key={category.id} className="space-y-4">
              <label className="text-lg font-medium text-muted-foreground ml-1">
                {category.label}
              </label>
              
              <div className="relative px-2">                
                <Slider
                  min={0}
                  max={2}
                  step={1}
                  value={[levelToValue(ratings[category.id])]}
                  onValueChange={(value) => handleRatingChange(category.id, value[0])}
                  className="relative z-10 *:data-[slot=slider-thumb]:bg-pink *:data-[slot=slider-thumb]:border-pink *:data-[slot=slider-range]:bg-pink *:data-[slot=slider-track]:bg-gray-300"
                />
                
                <div className="flex justify-between mt-2 px-2">
                  {RATING_LEVELS.map((level) => (
                    <span
                      key={level.value}
                      className={`text-md transition-colors ${
                      ratings[category.id] === level.id
                        ? 'text-pink font-semibold'
                        : 'text-muted-foreground'
                      }`}
                    >
                      {level.label}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="flex justify-between items-center mt-8">
          <Button
            onClick={saveChanges}
            disabled={!hasChanges || isUpdating}
            className="bg-pink hover:bg-pink/90"
          >
            {isUpdating ? 'Saving...' : 'Save Changes'}
          </Button>
          
          <Link 
            href={`/sessions/${sessionId}/report`}
            className="flex items-center text-sm font-medium text-muted-foreground hover:text-primary transition-colors"
          >
            View Report <ChevronRight className="w-4 h-4 ml-1" />
          </Link>
        </div>
      </div>
    </div>
  )
}