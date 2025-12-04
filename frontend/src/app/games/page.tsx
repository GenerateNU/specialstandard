"use client";

import AppLayout from "@/components/AppLayout";
import { GameContentSelector } from "@/components/games/GameContentSelector";
import { useSessionContext } from "@/contexts/sessionContext";
import { useStudents } from '@/hooks/useStudents'
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
  const { session, students } = useSessionContext();
  const sessionId = searchParams.get("sessionId") ?? "00000000-0000-0000-0000-000000000000"; // could use this or the session context
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
            {/* PDF Exercises Section */}
            <div className="mb-12">
              <h2 className="text-lg font-semibold mb-4 text-primary">PDF Exercises</h2>
              {isLoadingPdfs ? (
                <div className="bg-card rounded-lg shadow-md p-6 flex items-center justify-center">
                  <Loader2 className="w-5 h-5 animate-spin text-accent mr-2" />
                  <span className="text-secondary">Loading PDF exercises...</span>
                </div>
              ) : pdfExercises.length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {pdfExercises.map((pdf, index) => (
                  <div
                    key={pdf.id}
                    onClick={() => handleDownloadPdf(pdf.answer)}
                    className="bg-card rounded-lg shadow-md p-8 hover:shadow-lg transition-all duration-200 group hover:bg-card-hover border border-default hover:border-hover text-center relative cursor-pointer"
                  >
                    <Tooltip content="Submit Result">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleOpenResultModal(pdf);
                        }}
                        className="absolute top-4 right-4 p-1.5 bg-pink text-white rounded-lg hover:bg-pink-hover transition-colors"
                      >
                        <BadgeCheck className="w-4 h-4" />
                      </button>
                    </Tooltip>
                    
                    <FileText className="w-12 h-12 text-blue mb-4 mx-auto" />
                    <h3 className="mb-2">
                      {pdf.question_type ? `${pdf.question_type.split('_').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ')} ${index + 1}` : `PDF Exercise ${index + 1}`}
                    </h3>
                    <p className="text-secondary text-sm">Click to download</p>
                  </div>
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
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
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
                  className="bg-card rounded-lg shadow-md p-8 hover:shadow-lg transition-all duration-200 group hover:bg-card-hover border border-default hover:border-hover"
                >
                  <BookOpen className="w-12 h-12 text-blue mb-4 mx-auto" />
                  <h3 className="mb-2">Flashcards</h3>
                  <p className="text-secondary text-sm">
                    Practice with interactive flashcards
                  </p>
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
                  className="cursor-pointer bg-card rounded-lg shadow-md p-8 hover:shadow-lg transition-all duration-200 group hover:bg-card-hover border border-default hover:border-hover"
                >
                  <Image className="w-12 h-12 text-blue mb-4 mx-auto" />
                  <h3 className="mb-2">Image Matching</h3>
                  <p className="text-secondary text-sm">
                    Match words with images
                  </p>
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
                  className="bg-card rounded-lg shadow-md p-8 hover:shadow-lg transition-all duration-200 group hover:bg-card-hover border border-default hover:border-hover"
                >
                  <Brain className="w-12 h-12 text-blue mb-4 mx-auto" />
                  <h3 className="mb-2">Memory Match</h3>
                  <p className="text-secondary text-sm">
                    Spin a wheel to test your skills!
                  </p>
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
                  className="bg-card rounded-lg shadow-md p-8 hover:shadow-lg transition-all duration-200 group hover:bg-card-hover border border-default hover:border-hover"
                >
                  <SquareDashedMousePointer className="w-12 h-12 text-blue mb-4 mx-auto" />
                  <h3 className="mb-2">Drag and Drop</h3>
                  <p className="text-secondary text-sm">
                    Drag and drop the story in order!
                  </p>
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
                  className="bg-card rounded-lg shadow-md p-8 hover:shadow-lg transition-all duration-200 group hover:bg-card-hover border border-default hover:border-hover"
                >
                  <Shuffle className="w-12 h-12 text-blue mb-4 mx-auto" />
                  <h3 className="mb-2">Word-Image Matching</h3>
                  <p className="text-secondary text-sm">
                    Match many words to many images in this game!
                  </p>
                </button>

                <button
                  disabled
                  className="bg-card rounded-lg shadow-md p-8 opacity-50 cursor-not-allowed border border-default"
                >
                  <Gamepad2 className="w-12 h-12 text-muted mb-4 mx-auto" />
                  <h3 className="mb-2 text-muted">Quiz Game</h3>
                  <p className="text-disabled text-sm">Coming soon</p>
                </button>
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
  return (
    <AppLayout>
      <GameContentSelector
        onSelectionComplete={handleContentSelection}
        onBack={session ? () => router.push(`/sessions/${session.id}/curriculum`) : () =>router.back()}
        backLabel={session ? "Back to Curriculum" : "Back"}
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
