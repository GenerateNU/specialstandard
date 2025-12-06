"use client";

import { useGameContents } from "@/hooks/useGameContents";
import { useGameResults } from "@/hooks/useGameResults";
import type { GetGameContentsQuestionType } from "@/lib/api/theSpecialStandardAPI.schemas";
import { GetGameContentsCategory } from "@/lib/api/theSpecialStandardAPI.schemas";
import { RotateCw } from "lucide-react";
import { useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import "./flashcard.css";
import { useSessionContext } from "@/contexts/sessionContext";
import { useStudents } from "@/hooks/useStudents";

const CATEGORIES = {
  [GetGameContentsCategory.receptive_language]: {
    label: "Receptive Language",
    icon: "👂",
    colorClass: "bg-blue",
  },
  [GetGameContentsCategory.expressive_language]: {
    label: "Expressive Language",
    icon: "💬",
    colorClass: "bg-pink",
  },
  [GetGameContentsCategory.social_pragmatic_language]: {
    label: "Social Pragmatic Language",
    icon: "🤝",
    colorClass: "bg-orange",
  },
  [GetGameContentsCategory.speech]: {
    label: "Speech",
    icon: "🗣️",
    colorClass: "bg-blue",
  },
};

interface SpinnerGameInterfaceProps {
  session_student_ids: number[];
  session_id?: string;
  themeId: string;
  themeWeek: number;
  themeName: string;
  difficulty: number;
  category: string;
  questionType: string;
}

// Spinner component
const SpinnerWheel: React.FC<{
  words: any[];
  onSpin: () => void;
  isSpinning: boolean;
  selectedWord: string | null;
  rotation: number;
  completedIds: Set<string>;
}> = ({ words, onSpin, isSpinning, rotation, completedIds }) => {
  // Only show segments for words that are not yet completed
  const availableWords = words.filter((w) => !completedIds.has(w.id));
  const segmentAngle = availableWords.length
    ? 360 / availableWords.length
    : 360;
  const colors = [
    "#FF6B6B",
    "#4ECDC4",
    "#45B7D1",
    "#96CEB4",
    "#FFEAA7",
    "#DDA0DD",
    "#98D8C8",
    "#FFB6C1",
    "#87CEEB",
    "#FFD700",
  ];

  // Function to create SVG path for pie segment
  const createPieSlice = (index: number) => {
    const startAngle = index * segmentAngle;
    const endAngle = startAngle + segmentAngle;

    // Convert to radians
    const startRad = (startAngle * Math.PI) / 180;
    const endRad = (endAngle * Math.PI) / 180;

    // Calculate points
    const x1 = 50 + 50 * Math.cos(startRad);
    const y1 = 50 + 50 * Math.sin(startRad);
    const x2 = 50 + 50 * Math.cos(endRad);
    const y2 = 50 + 50 * Math.sin(endRad);

    // Large arc flag for segments > 180°
    const largeArc = segmentAngle > 180 ? 1 : 0;

    return `M 50 50 L ${x1} ${y1} A 50 50 0 ${largeArc} 1 ${x2} ${y2} Z`;
  };

  return (
    <div
      className="relative mx-auto mb-6"
      style={{
        width: "30vw",
        height: "30vw",
        maxWidth: "400px",
        maxHeight: "400px",
        minWidth: "250px",
        minHeight: "250px",
      }}
    >
      {/* Pointer */}
      <div
        className="absolute top-0 left-1/2 -translate-x-1/2 z-20"
        style={{
          width: 0,
          height: 0,
          borderLeft: "12px solid transparent",
          borderRight: "12px solid transparent",
          borderBottom: "20px solid #333",
        }}
      ></div>
      {/* Wheel - SVG based for clean pie segments */}
      <div
        className="w-full h-full rounded-full relative overflow-hidden"
        style={{
          transform: `rotate(${rotation}deg)`,
          transition: 'transform 4000ms ease-out'
        }}
      >
        <svg viewBox="0 0 100 100" className="w-full h-full">
          {availableWords.length === 1 ? (
            // Single word - full circle
            <g key={availableWords[0].id}>
              <circle
                cx="50"
                cy="50"
                r="50"
                fill={colors[0]}
                stroke="#333"
                strokeWidth="0.5"
              />
              <text
                x="50"
                y="30"
                fill="black"
                fontSize="6"
                fontWeight="600"
                textAnchor="middle"
              >
                {availableWords[0].question}
              </text>
            </g>
          ) : (
            // Multiple words - pie segments
            availableWords.map((word, index) => (
              <g key={word.id}>
                <path
                  d={createPieSlice(index)}
                  fill={colors[index % colors.length]}
                  stroke="#333"
                  strokeWidth="0.5"
                />
                <text
                  x="70"
                  y="50"
                  fill="black"
                  fontSize="5"
                  fontWeight="600"
                  textAnchor="middle"
                  transform={`rotate(${
                    index * segmentAngle + segmentAngle / 2
                  } 50 50)`}
                >
                  {word.question}
                </text>
              </g>
            ))
          )}
        </svg>
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-8 h-8 bg-gray-800 rounded-full border-2 border-gray-700"></div>
      </div>

      {/* Spin button */}
      <button
        onClick={onSpin}
        disabled={isSpinning || availableWords.length === 0}
        className={`absolute left-1/2 -translate-x-1/2 rounded-lg font-semibold transition-all ${
          isSpinning || availableWords.length === 0
            ? "bg-gray-400 text-gray-200 cursor-not-allowed"
            : "bg-zinc-300 text-black hover:scale-105 hover:shadow-lg"
        }`}
        style={{
          bottom: "-25%",
          padding: "clamp(0.5rem, 1.5vw, 1rem) clamp(1rem, 3vw, 2rem)",
          fontSize: "clamp(0.875rem, 2vw, 1.125rem)",
        }}
      >
        {availableWords.length === 0
          ? "All Done!"
          : availableWords.length === 1
            ? "Choose Word"
            : isSpinning
              ? "Spinning..."
              : "Spin the Wheel!"}
      </button>
    </div>
  );
};

export default function SpinnerGameInterface({
  session_student_ids,
  session_id,
  themeId,
  themeWeek,
  themeName,
  difficulty,
  category,
  questionType,
}: SpinnerGameInterfaceProps) {
  const router = useRouter();
  const [cardStartTime, setCardStartTime] = useState<number | null>(null);
  const [timeTaken, setTimeTaken] = useState(0);

  // Spinner-specific state
  const [isSpinning, setIsSpinning] = useState(false);
  const [selectedWordIndex, setSelectedWordIndex] = useState<number | null>(
    null
  );
  const [showActionButtons, setShowActionButtons] = useState(false);
  const [spinRotation, setSpinRotation] = useState(0);
  const [completedIds, setCompletedIds] = useState<Set<string>>(new Set());
  const [questionCount, setQuestionCount] = useState(0);
  const [resultsSaved, setResultsSaved] = useState(false);
  const [incorrectAttempts, setIncorrectAttempts] = useState<Map<string, number>>(new Map());

  const {
    gameContents,
    isLoading: contentsLoading,
    error: contentsError,
  } = useGameContents({
    theme_id: themeId,
    theme_week: themeWeek ?? undefined,
    category: category as GetGameContentsCategory,
    question_type: questionType as GetGameContentsQuestionType,
    difficulty_level: difficulty,
    question_count: 10,
    applicable_game_types: ['spinner'],
  });

  // Calculate questions per student and limit total cards
  // Ensure at least 1 question per student, or use all available if fewer questions than students
  const questionsPerStudent = Math.max(1, Math.floor(gameContents.length / session_student_ids.length));
  const totalQuestionsToUse = Math.min(questionsPerStudent * session_student_ids.length, gameContents.length);
  const limitedGameContents = gameContents.length > 0 ? gameContents.slice(0, totalQuestionsToUse) : gameContents;
  
  // Get current student based on question count
  const currentStudentIndex = questionCount % session_student_ids.length;
  const currentSessionStudentId = session_student_ids[currentStudentIndex];

  // Create game results hooks for each student
  const gameResultsHooks = session_student_ids.map(studentId => 
    useGameResults({
      session_student_id: studentId,
      session_id,
    })
  );

  const currentGameResultsHook = gameResultsHooks[currentStudentIndex];
  const startCard = currentGameResultsHook?.startCard;

  // Get student names from session context
  const { students: sessionStudents, session } = useSessionContext();
  const { students: allStudents } = useStudents();
  
  // Prefer session context ID over prop
  const effectiveSessionId = session?.id || session_id;
  
  const getStudentName = (sessionStudentId: number) => {
    const sessionStudent = sessionStudents.find(s => s.sessionStudentId === sessionStudentId);
    if (!sessionStudent) return 'Student';
    const student = allStudents?.find(s => s.id === sessionStudent.studentId);
    return student ? `${student.first_name} ${student.last_name}` : 'Student';
  };

  // Update timer display
  useEffect(() => {
    if (cardStartTime === null) return;

    const interval = setInterval(() => {
      setTimeTaken(Math.floor((Date.now() - cardStartTime) / 1000));
    }, 1000); // Changed from 100ms to 1000ms for less frequent updates

    return () => clearInterval(interval);
  }, [cardStartTime]);

  // Handler functions
  const handleSpin = () => {
    // Get available words from limited game contents
    const availableWords = limitedGameContents.filter(
      (content) => !completedIds.has(content.id)
    );

    // Special handling for single word
    if (availableWords.length === 1) {
      const finalWord = availableWords[0];

      setSelectedWordIndex(0);
      setIsSpinning(false);
      setShowActionButtons(true);

      if (currentGameResultsHook) {
        startCard?.(finalWord);
        setCardStartTime(Date.now());
      }
      return;
    }

    if (isSpinning || availableWords.length === 0) return;

    setIsSpinning(true);
    setShowActionButtons(false);

    // Random spin calculation
    const spins = 5 + Math.random() * 3; // 5-8 full rotations
    const randomAngle = Math.random() * 360;
    const totalRotation = spinRotation + (spins * 360 + randomAngle);

    setSpinRotation(totalRotation);

    const segmentAngle = 360 / availableWords.length;

    // The pointer is at the top. When wheel rotates clockwise by X degrees,
    // the segment that WAS at position (360 - X) is now at the top
    const normalizedRotation = totalRotation % 360;

    // Which segment is now at the top after rotating?
    let segmentAtTop = (360 - normalizedRotation - 90) % 360;
    if (segmentAtTop < 0) segmentAtTop += 360;

    // Now find which segment index this angle belongs to
    const selectedIndex =
      Math.floor(segmentAtTop / segmentAngle) % availableWords.length;

    setTimeout(() => {
      setSelectedWordIndex(selectedIndex);
      setIsSpinning(false);
      setShowActionButtons(true);

      // Start tracking time for this word
      if (currentGameResultsHook) {
        startCard?.(availableWords[selectedIndex]);
        setCardStartTime(Date.now());
      }
    }, 4000);
  };

  const handleMarkResult = (isCorrect: boolean) => {
    if (selectedWordIndex === null) return;

    // Get available words to find the actual content
    const availableWords = limitedGameContents.filter(
      (content) => !completedIds.has(content.id)
    );
    const selectedContent = availableWords[selectedWordIndex];

    if (!selectedContent) return;

    if (isCorrect) {
      // Mark as completed by adding ID
      setCompletedIds((prev) => new Set([...prev, selectedContent.id]));

      if (currentGameResultsHook) {
        const finalTime = Math.floor(
          (Date.now() - (cardStartTime || Date.now())) / 1000
        );
        const attempts = incorrectAttempts.get(selectedContent.id) || 0;
        currentGameResultsHook.completeCard(selectedContent.id, finalTime, attempts);
      }
      
      // Increment question count to move to next student
      setQuestionCount(prev => prev + 1);
      
      // Reset for next word
      setShowActionButtons(false);
      setSelectedWordIndex(null);
      setCardStartTime(null);
      setTimeTaken(0);
    } else {
      // Track incorrect attempt
      setIncorrectAttempts(prev => {
        const newMap = new Map(prev);
        const current = newMap.get(selectedContent.id) || 0;
        newMap.set(selectedContent.id, current + 1);
        return newMap;
      });
      
      // Hide buttons and reset to let them spin again
      setShowActionButtons(false);
      setSelectedWordIndex(null);
      setCardStartTime(null);
      setTimeTaken(0);
    }
  };

  const handleReset = () => {
    setSelectedWordIndex(null);
    setShowActionButtons(false);
    setSpinRotation(0);
    setCompletedIds(new Set());
    setCardStartTime(null);
    setTimeTaken(0);
    setQuestionCount(0);
    setResultsSaved(false);
    setIncorrectAttempts(new Map());
  };

  if (contentsLoading) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue mx-auto mb-4"></div>
          <p className="text-secondary">Loading Game...</p>
        </div>
      </div>
    );
  }

  if (contentsError || limitedGameContents.length === 0) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <p className="text-error mb-4">
            {contentsError
              ? "Failed to load Spinner"
              : "No words available"}
          </p>
          <button
            onClick={() => router.back()}
            className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover transition-colors"
          >
            Go Back
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-4">
            <button
              onClick={() => router.push(effectiveSessionId ? `/sessions/${effectiveSessionId}/curriculum` : '/games')}
              className="text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
            >
              ← Back to Content
            </button>
          </div>
          <button
            onClick={handleReset}
            className="flex items-center gap-2 text-secondary hover:text-primary transition-colors"
          >
            <RotateCw className="w-4 h-4" />
            Reset
          </button>
        </div>

        <h1 className="mb-4">Spinner</h1>

        {/* Current Student Banner */}
        <div className="bg-blue text-white rounded-lg p-4 mb-6 text-center">
          <p className="text-sm opacity-90 mb-1">Current Player</p>
          <p className="text-2xl font-bold">{getStudentName(currentSessionStudentId)}</p>
        </div>

        {/* Display selected options */}
        <div className="bg-card rounded-lg p-4 mb-6 border border-default">
          <div className="flex flex-wrap items-center gap-2 text-sm">
            <span className="text-muted">Theme:</span>
            <span className="font-medium text-primary">{themeName}</span>
            <span className="text-muted">Difficulty:</span>
            <span className="font-medium text-primary">Level {difficulty}</span>
            <span className="text-muted mx-2">•</span>
            <span className="text-muted">Category:</span>
            <span className="font-medium text-primary">
              {CATEGORIES[category as GetGameContentsCategory]?.label ||
                category}
            </span>
          </div>
        </div>

        {/* Progress indicator - hide when all complete */}
        {completedIds.size < limitedGameContents.length && (
        <div className="mb-8">
          <div className="flex items-center justify-between text-sm text-secondary mb-2">
            <span>
              Card {completedIds.size + 1} of {limitedGameContents.length}
            </span>
            <span>
              {Math.round((completedIds.size / limitedGameContents.length) * 100)}%
              Complete
            </span>
          </div>
          <div className="w-full bg-card rounded-full h-2 border border-default">
            <div
              className="bg-blue h-full rounded-full transition-all duration-300"
              style={{
                width: `${(completedIds.size / limitedGameContents.length) * 100}%`,
              }}
            />
          </div>
        </div>
        )}

        {/* Spinner Game - hide when all complete */}
        {limitedGameContents.length > 0 && completedIds.size < limitedGameContents.length && (
          <div className="mb-8">
            {/* Container for wheel and side buttons */}
            <div className="relative flex items-center justify-center gap-4">
              {/* Incorrect button - left side */}
              {selectedWordIndex !== null && showActionButtons && (
                <button
                  onClick={() => handleMarkResult(false)}
                  className="absolute left-0 px-4 py-2 bg-red-300 text-black rounded-lg hover:bg-red-400 transition-colors animate-fade-in"
                  style={{
                    left: "5%",
                    top: "50%",
                    transform: "translateY(-50%)",
                  }}
                >
                  Incorrect
                </button>
              )}

              {/* Spinner Wheel */}
              <SpinnerWheel
                words={limitedGameContents}
                onSpin={handleSpin}
                isSpinning={isSpinning}
                selectedWord={
                  selectedWordIndex !== null
                    ? limitedGameContents.filter((c) => !completedIds.has(c.id))[
                        selectedWordIndex
                      ]?.question
                    : null
                }
                rotation={spinRotation}
                completedIds={completedIds}
              />

              {/* Correct button - right side */}
              {selectedWordIndex !== null && showActionButtons && (
                <button
                  onClick={() => handleMarkResult(true)}
                  className="absolute right-0 px-4 py-2 bg-green-300 text-black rounded-lg hover:bg-green-400 transition-colors animate-fade-in"
                  style={{
                    right: "5%",
                    top: "50%",
                    transform: "translateY(-50%)",
                  }}
                >
                  Correct
                </button>
              )}
            </div>

            {/* Result display - selected word */}
            <div className="text-center mt-20 min-h-[60px]">
              {selectedWordIndex !== null && (
                <div>
                  <p className="text-3xl font-bold text-primary">
                    {
                      limitedGameContents.filter((c) => !completedIds.has(c.id))[
                        selectedWordIndex
                      ]?.question
                    }
                  </p>
                </div>
              )}
            </div>

            {/* Timer display */}
            <div className="text-center text-muted text-sm">
              {selectedWordIndex !== null && timeTaken > 0 && (
                <span>Time: {timeTaken}s</span>
              )}
            </div>
          </div>
        )}

        {/* Completion message */}
        {completedIds.size === limitedGameContents.length && (
          <div className="mt-8 p-6 bg-blue-light border border-blue rounded-lg text-center">
            <p className="text-blue font-semibold">
              Great job! You've completed all words!
            </p>
            <div className="mt-4 flex gap-3 justify-center">
              <button
                onClick={handleReset}
                className="px-6 py-2 bg-blue text-white rounded-lg hover:bg-blue-hover transition-colors"
              >
                Start Over
              </button>

              <button
                onClick={async () => {
                  // Save all results for all students
                  for (const hook of gameResultsHooks) {
                    await hook.saveAllResults();
                  }
                  setResultsSaved(true);
                }}
                disabled={gameResultsHooks.some(hook => hook.isSaving) || resultsSaved}
                className="px-6 py-2 bg-pink text-white rounded-lg hover:bg-pink-hover transition-colors disabled:bg-pink-disabled disabled:cursor-not-allowed"
              >
                {gameResultsHooks.some(hook => hook.isSaving) ? "Saving..." : resultsSaved ? "Saved!" : "Save Progress"}
              </button>
              
              <button
                onClick={() => router.push(effectiveSessionId ? `/sessions/${effectiveSessionId}/curriculum` : '/games')}
                className="px-6 py-2 bg-card-hover text-primary rounded-lg hover:bg-card border border-border"
              >
                Back to Content
              </button>
            </div>

            {gameResultsHooks.some(hook => hook.saveError) && (
              <p className="text-error text-sm mt-2">
                Failed to save progress. Please try again.
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
