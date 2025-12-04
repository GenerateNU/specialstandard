"use client";

import AppLayout from "@/components/AppLayout";
import { GameContentSelector } from "@/components/games/GameContentSelector";
import { useSessionContext } from "@/contexts/sessionContext";
import { useStudents } from '@/hooks/useStudents'
import { useThemes } from '@/hooks/useThemes'
import type {
  GameContent,
  GetGameContentsCategory,
  GetGameContentsQuestionType,
  Theme,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import { GameContentExerciseType } from "@/lib/api/theSpecialStandardAPI.schemas";
import { getGameContent } from "@/lib/api/game-content";
import { BadgeCheck, BookOpen, Brain, FileText, Gamepad2, Image, Loader2, Shuffle, SquareDashedMousePointer} from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import React, { Suspense } from "react";
import { useManualGameResult } from "@/hooks/useManualGameResults";
import Tooltip from '@/components/ui/tooltip'


function GamesPageContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { session, students, currentLevel, currentMonth, currentYear } = useSessionContext();
  const sessionId = searchParams.get("sessionId") ?? "00000000-0000-0000-0000-000000000000"; // could use this or the session context
  const categoryParam = searchParams.get("category") as GetGameContentsCategory | null;
  
  // Fetch theme for current month/year
  const { themes, isLoading: themesLoading } = useThemes({ 
    month: currentMonth + 1, // Convert from 0-11 to 1-12
    year: currentYear 
  });
  const currentTheme = themes && themes.length > 0 ? themes[0] : null;
  
  // Redirect to curriculum if no level is selected
  React.useEffect(() => {
    if (sessionId !== "00000000-0000-0000-0000-000000000000" && !currentLevel) {
      router.push(`/sessions/${sessionId}/curriculum`);
    }
  }, [sessionId, currentLevel, router]);
  
  const sessionStudentIds = students ? students.map((student) => student.sessionStudentId?.toString()) : [];
  const { students: allStudents } = useStudents()
  const [selectedContent, setSelectedContent] = React.useState<{
    theme: Theme;
    difficultyLevel: number;
    category: GetGameContentsCategory;
    questionType: GetGameContentsQuestionType;
  } | null>(null);

  const [pdfExercises, setPdfExercises] = React.useState<GameContent[]>([]);
  const [isLoadingPdfs, setIsLoadingPdfs] = React.useState(false);

  const [showResultModal, setShowResultModal] = React.useState(false);
  const [selectedPdf, setSelectedPdf] = React.useState<GameContent | null>(null);
  const [timeTaken, setTimeTaken] = React.useState<number>(0);
  const [completed, setCompleted] = React.useState(true);
  const [incorrectAttempts, setIncorrectAttempts] = React.useState<number>(0);

  const [selectedStudentId, setSelectedStudentId] = React.useState<string>('');
  const [showSuccessMessage, setShowSuccessMessage] = React.useState(false);


  const { submitResultAsync, isSubmitting, error } = useManualGameResult();

  // Fetch PDFs when content is selected
  React.useEffect(() => {
    const fetchPdfs = async () => {
      if (!selectedContent) return;

      setIsLoadingPdfs(true);
      try {
        const api = getGameContent();
        const response = await api.getGameContents({
          exercise_type: GameContentExerciseType.pdf,
          theme_id: selectedContent.theme.id,
          difficulty_level: selectedContent.difficultyLevel,
          category: selectedContent.category,
        });

        if (response && Array.isArray(response)) {
          const pdfItems = response.filter((item: GameContent) => item.exercise_type === GameContentExerciseType.pdf);
          const uniqueContents = Array.from(
            new Map(pdfItems.map(item => [item.id, item])).values()
          );
          setPdfExercises(uniqueContents);
        }
      } catch (err) {
        console.error('Error fetching PDFs:', err);
        setPdfExercises([]);
      } finally {
        setIsLoadingPdfs(false);
      }
    };

    fetchPdfs();
  }, [selectedContent]);

  const handleContentSelection = (selection: {
    theme: Theme;
    difficultyLevel: number;
    category: GetGameContentsCategory;
    questionType: GetGameContentsQuestionType;
  }) => {
    setSelectedContent(selection);
  };

  const handleDownloadPdf = (pdfUrl: string) => {
    window.open(pdfUrl, '_blank');
  };

  const handleOpenResultModal = (pdf: GameContent) => {

    setSelectedPdf(pdf);
    setShowResultModal(true);
    setTimeTaken(0);
    setCompleted(true);
    setIncorrectAttempts(0);
    setSelectedStudentId(sessionStudentIds[0] || ''); // Set default student

  };

const handleSubmitResult = async () => {
  if (!selectedPdf || !selectedStudentId) {
    return;
  }  
  try {
    await submitResultAsync({
      session_student_id: Number.parseInt(selectedStudentId),
      content_id: selectedPdf.id,
      time_taken_sec: timeTaken,
      completed,
      count_of_incorrect_attempts: incorrectAttempts,
    });
    
    setShowSuccessMessage(true); // Show success
    
    // Hide modal after brief delay
    setTimeout(() => {
      setShowResultModal(false);
      setSelectedPdf(null);
      setShowSuccessMessage(false);
    }, 1500);
  } catch (error) {
    console.error('Error saving result:', error);
  }
};

  // Show game selection after content is selected
  if (selectedContent) {
    return (
      <AppLayout>
        <div className="min-h-screen bg-background p-8">
          <div className="max-w-4xl mx-auto">
            <button
              onClick={() => setSelectedContent(null)}
              className="mb-6 text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
            >
              ← Back to Content
            </button>
            <h1 className="mb-8">Select an Exercise</h1>
            
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* PDF Exercises Section */}
              <div className="bg-card shadow-md rounded-4xl p-8 flex flex-col">
                <h2 className="text-2xl font-bold text-primary mb-6">PDF Exercises</h2>
                {isLoadingPdfs ? (
                  <div className="flex items-center justify-center py-12">
                    <Loader2 className="w-5 h-5 animate-spin text-accent mr-2" />
                    <span className="text-secondary">Loading PDF exercises...</span>
                  </div>
                ) : pdfExercises.length > 0 ? (
                  <div className="flex flex-col gap-4">
                    {pdfExercises.map((pdf, index) => (
                      <button
                        key={pdf.id}
                        onClick={() => handleDownloadPdf(pdf.answer)}
                        className="bg-pink cursor-pointer hover:bg-pink-hover text-white p-6 rounded-lg font-semibold transition-all hover:scale-105 text-left flex items-center justify-between gap-4 relative"
                      >
                        <div className="flex items-center gap-4">
                          <FileText className="w-6 h-6 shrink-0" />
                          <span>
                            {pdf.question_type ? `${pdf.question_type.split('_').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ')} ${index + 1}` : `PDF Exercise ${index + 1}`}
                          </span>
                        </div>
                        <Tooltip content="Submit Result">
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleOpenResultModal(pdf);
                            }}
                            className="p-2 bg-white/20 hover:bg-white/30 rounded-lg transition-colors"
                          >
                            <BadgeCheck className="w-5 h-5" />
                          </button>
                        </Tooltip>
                      </button>
                    ))}
                  </div>
                ) : (
                  <div className="flex items-center justify-center py-12 text-center text-secondary">
                    No PDF exercises available for this selection
                  </div>
                )}
              </div>

              {/* Interactive Games Section */}
              <div className="bg-card shadow-md rounded-4xl p-8 flex flex-col">
                <h2 className="text-2xl font-bold text-primary mb-6">Interactive Games</h2>
                <div className="flex flex-col gap-4">
                  <button
                    onClick={() => {
                      const params = new URLSearchParams({
                        themeId: selectedContent.theme.id,
                        difficulty: String(selectedContent.difficultyLevel),
                        category: selectedContent.category,
                        questionType: selectedContent.questionType,
                        sessionId,
                      });
                      router.push(`/games/flashcards?${params.toString()}`);
                    }}
                    className="bg-pink cursor-pointer hover:bg-pink-hover text-white p-6 rounded-lg font-semibold transition-all hover:scale-105 text-left flex items-center gap-4"
                  >
                    <BookOpen className="w-6 h-6 shrink-0" />
                    <span>Flashcards</span>
                  </button>
                  
                  <button
                    onClick={() => {
                      const params = new URLSearchParams({
                        themeId: selectedContent.theme.id,
                        difficulty: String(selectedContent.difficultyLevel),
                        category: selectedContent.category,
                        questionType: selectedContent.questionType,
                        sessionId,
                      });
                      router.push(`/games/image-matching?${params.toString()}`);
                    }}
                    className="bg-pink cursor-pointer hover:bg-pink-hover text-white p-6 rounded-lg font-semibold transition-all hover:scale-105 text-left flex items-center gap-4"
                  >
                    <Image className="w-6 h-6 shrink-0" />
                    <span>Image Matching</span>
                  </button>
                  
                  <button
                    onClick={() => {
                      const params = new URLSearchParams({
                        themeId: selectedContent.theme.id,
                        difficulty: String(selectedContent.difficultyLevel),
                        category: selectedContent.category,
                        questionType: selectedContent.questionType,
                        sessionId,
                      });
                      router.push(`/games/memorymatch?${params.toString()}`);
                    }}
                    className="bg-pink cursor-pointer hover:bg-pink-hover text-white p-6 rounded-lg font-semibold transition-all hover:scale-105 text-left flex items-center gap-4"
                  >
                    <Brain className="w-6 h-6 shrink-0" />
                    <span>Memory Match</span>
                  </button>
                  
                  <button
                    onClick={() => {
                      const params = new URLSearchParams({
                        themeId: selectedContent.theme.id,
                        difficulty: String(selectedContent.difficultyLevel),
                        category: selectedContent.category,
                        questionType: selectedContent.questionType,
                        sessionId,
                      });
                      router.push(`/games/drag-and-drop?${params.toString()}`);
                    }}
                    className="bg-pink cursor-pointer hover:bg-pink-hover text-white p-6 rounded-lg font-semibold transition-all hover:scale-105 text-left flex items-center gap-4"
                  >
                    <SquareDashedMousePointer className="w-6 h-6 shrink-0" />
                    <span>Drag and Drop</span>
                  </button>

                  <button
                    onClick={() => {
                      const params = new URLSearchParams({
                        themeId: selectedContent.theme.id,
                        difficulty: String(selectedContent.difficultyLevel),
                        category: selectedContent.category,
                        questionType: selectedContent.questionType,
                        sessionId,
                      });
                      router.push(`/games/word-image-match?${params.toString()}`);
                    }}
                    className="bg-pink cursor-pointer hover:bg-pink-hover text-white p-6 rounded-lg font-semibold transition-all hover:scale-105 text-left flex items-center gap-4"
                  >
                    <Shuffle className="w-6 h-6 shrink-0" />
                    <span>Word-Image Matching</span>
                  </button>

                  <button
                    disabled
                    className="bg-card border border-default p-6 rounded-lg opacity-50 cursor-not-allowed text-left flex items-center gap-4"
                  >
                    <Gamepad2 className="w-6 h-6 shrink-0 text-muted" />
                    <span className="text-muted font-semibold">Quiz Game (Coming Soon)</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
    {/* Result Submission Modal */}
    {showResultModal && selectedPdf && (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div className="bg-card rounded-lg shadow-xl p-8 max-w-md w-full mx-4">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-semibold">Submit Exercise Result</h2>
            <button
              onClick={() => setShowResultModal(false)}
              className="text-muted hover:text-primary transition-colors"
            >
              ✕
            </button>
          </div>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">
                Exercise: {selectedPdf.question_type?.split('_').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ')} {pdfExercises.findIndex(pdf => pdf.id === selectedPdf.id) + 1}
              </label>
            </div>

            <div>
              <label htmlFor="student" className="block text-sm font-medium mb-2">
                Select Student
              </label>
              <select
                id="student"
                value={selectedStudentId}
                onChange={(e) => setSelectedStudentId(e.target.value)}
                className="w-full px-4 py-2 border border-default rounded-lg focus:outline-none focus:ring-2 focus:ring-blue bg-background"
              >
                <option value="">Select a student...</option>
                {students?.map((student) => {
                  const studentDetails = allStudents?.find(s => s.id === student.studentId);
                  const fullName = `${studentDetails?.first_name || ''} ${studentDetails?.last_name || ''}`.trim();
                  return (
                    <option key={student.sessionStudentId} value={student.sessionStudentId?.toString()}>
                      {fullName}
                    </option>
                  );
                })}
              </select>
            </div>

            <div>
              <label htmlFor="timeTaken" className="block text-sm font-medium mb-2">
                Time Taken (minutes)
              </label>
              <input
                id="timeTaken"
                type="number"
                min="0"
                step="0.5"
                value={timeTaken / 60}
                onChange={(e) => setTimeTaken(Math.round((Number.parseFloat(e.target.value) || 0) * 60))}
                className="w-full px-4 py-2 border border-default rounded-lg focus:outline-none focus:ring-2 focus:ring-blue bg-background text-base"
                placeholder="e.g. 5 or 5.5"
              />
            </div>

            <div>
              <label htmlFor="incorrectAttempts" className="block text-sm font-medium mb-2">
                Number of Incorrect Attempts
              </label>
              <input
                id="incorrectAttempts"
                type="number"
                min="0"
                value={incorrectAttempts}
                onChange={(e) => setIncorrectAttempts(Number.parseInt(e.target.value) || 0)}
                className="w-full px-4 py-2 border border-default rounded-lg focus:outline-none focus:ring-2 focus:ring-blue bg-background"
              />
            </div>

            <div className="flex items-center gap-2">
              <input
                id="completed"
                type="checkbox"
                checked={completed}
                onChange={(e) => setCompleted(e.target.checked)}
                className="w-4 h-4"
              />
              <label htmlFor="completed" className="text-sm font-medium">
                Completed
              </label>
            </div>
          </div>

          <div className="flex gap-4 mt-6">
            <button
              onClick={() => setShowResultModal(false)}
              className="flex-1 px-4 py-2 border border-default rounded-lg hover:bg-card-hover transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleSubmitResult}
              disabled={!selectedStudentId || isSubmitting}
              className="flex-1 px-4 py-2 bg-pink text-white rounded-lg hover:bg-pink-hover transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {showSuccessMessage ? '✓ Saved!' : isSubmitting ? 'Submitting...' : 'Submit'}
            </button>
          </div>

          {error && (
            <p className="text-error text-sm mt-4 text-center">
              Failed to submit result. Please try again.
            </p>
          )}
        </div>
      </div>
    )}

      </AppLayout>
    );
  }

  // Show content selector
  // Show loading state while fetching theme
  if (themesLoading) {
    return (
      <AppLayout>
        <div className="min-h-screen bg-background flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue mx-auto mb-4"></div>
            <p className="text-muted">Loading theme...</p>
          </div>
        </div>
      </AppLayout>
    );
  }
  
  // Show error if no theme found
  if (!currentTheme) {
    return (
      <AppLayout>
        <div className="min-h-screen bg-background flex items-center justify-center">
          <div className="text-center">
            <p className="text-error mb-4">No theme found for the current month</p>
            <button
              onClick={() => router.push(`/sessions/${sessionId}/curriculum`)}
              className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover transition-colors"
            >
              Back to Curriculum
            </button>
          </div>
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <GameContentSelector
        onSelectionComplete={handleContentSelection}
        onBack={session ? () => router.push(`/sessions/${session.id}/curriculum`) : () =>router.back()}
        backLabel={session ? "Back to Curriculum" : "Back"}
        initialDifficultyLevel={currentLevel || undefined}
        initialCategory={categoryParam || undefined}
        theme={currentTheme}
      />
    </AppLayout>
  );
}

export default function GamesPage() {
  return (
    <Suspense fallback={null}>
      <GamesPageContent />
    </Suspense>
  );
}
