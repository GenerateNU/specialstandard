'use client'

import { useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";
import WordImageMatchingGameInterface from "@/components/games/WordImageMatchingGameInterface";
import { StudentSelector } from "@/components/games/StudentSelector";
import { useSessionContext } from "@/contexts/sessionContext";

export default function WordImageMatchingContent() {
    const searchParams = useSearchParams()
    const router = useRouter()
    const { session } = useSessionContext()

    const themeId = searchParams.get('themeId')
    const themeName = searchParams.get('themeName')
    const themeWeek = searchParams.get('themeWeek')
    const difficulty = searchParams.get('difficulty')
    const category = searchParams.get('category')
    const questionType = searchParams.get('questionType')
    const sessionId = searchParams.get('sessionId') || '00000000-0000-0000-0000-000000000000'
    const sessionStudentIdsParam = searchParams.get('sessionStudentIds')
    
    const effectiveSessionId = session?.id || sessionId

    const [selectedStudentIds, setSelectedStudentIds] = useState<string[]>(
        sessionStudentIdsParam ? sessionStudentIdsParam.split(',') : []
    )

    if (!themeId || !difficulty || !category || !questionType) {
        return (
            <div className="min-h-screen bg-background p-8 flex items-center justify-center">
                <div className="text-center">
                    <p className="text-error mb-4">
                        Missing game parameters. Please select content first.
                    </p>
                    <a href="/games" className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover
                                      transition-colors inline-block">
                        Go Back
                    </a>
                </div>
            </div>
        )
    }

    // Show student selector if no students selected yet
    if (selectedStudentIds.length === 0) {
        return (
            <StudentSelector
                gameTitle="Word-Image Matching"
                onStudentsSelected={(studentIds) => {
                    setSelectedStudentIds(studentIds)
                    // Update URL with selected students
                    const params = new URLSearchParams(searchParams.toString())
                    params.set('sessionStudentIds', studentIds.join(','))
                    router.replace(`/games/word-image-match?${params.toString()}`)
                }}
                onBack={() => router.push(effectiveSessionId ? `/sessions/${effectiveSessionId}/curriculum` : '/games')}
            />
        )
    }

    return (
        <WordImageMatchingGameInterface
            session_student_ids={selectedStudentIds.map(id => Number.parseInt(id))}
            session_id={sessionId}
            themeID={themeId}
            themeWeek={themeWeek ? Number.parseInt(themeWeek) : 1}
            themeName={themeName}
            difficulty={difficulty}
            category={category}
            questionType={questionType}
        />
    )
}
