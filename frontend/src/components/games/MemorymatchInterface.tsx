"use client";

import { useGameContents } from "@/hooks/useGameContents";
import { useGameResults } from "@/hooks/useGameResults";
import type { GetGameContentsQuestionType } from "@/lib/api/theSpecialStandardAPI.schemas";
import { GetGameContentsCategory } from "@/lib/api/theSpecialStandardAPI.schemas";
import { RotateCw } from "lucide-react";
import { useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import "./flashcard.css";

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

interface MemorymatchGameInterfaceProps {
  session_student_id?: number;
  session_id?: string;
  student_id?: string;
  themeId: string;
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
      {availableWords.length > 1 ? (
        // ----- normal wheel -----
        <div
          className="w-full h-full rounded-full relative overflow-hidden transition-transform duration-[4000ms] ease-out"
          style={{
            transform: `rotate(${rotation}deg)`,
          }}
        >
          <svg viewBox="0 0 100 100" className="w-full h-full">
            {availableWords.map((word, index) => (
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
                  fill="white"
                  fontSize="3"
                  fontWeight="600"
                  textAnchor="middle"
                  transform={`rotate(${
                    index * segmentAngle + segmentAngle / 2
                  } 50 50)`}
                >
                  {word.question}
                </text>
              </g>
            ))}
          </svg>
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-8 h-8 bg-gray-800 rounded-full border-2 border-gray-700"></div>
        </div>
      ) : (
        // ----- clean placeholder -----
        <div className="w-full h-full flex items-center justify-center">
          <div className="text-xl font-semibold text-primary">
            Only one word left!
          </div>
        </div>
      )}

      {/* Spin button */}
      <button
        onClick={onSpin}
        disabled={isSpinning || availableWords.length === 0}
        className={`absolute left-1/2 -translate-x-1/2 rounded-lg font-semibold transition-all ${
          isSpinning || availableWords.length === 0
            ? "bg-gray-400 text-gray-200 cursor-not-allowed"
            : "bg-zinc-300 text-white hover:scale-105 hover:shadow-lg"
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

export default function MemorymatchGameInterface({
  session_student_id,
  session_id,
  student_id,
  themeId,
  themeName,
  difficulty,
  category,
  questionType,
}: MemorymatchGameInterfaceProps) {
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

  const {
    gameContents,
    isLoading: contentsLoading,
    error: contentsError,
  } = useGameContents({
    theme_id: themeId,
    category: category as GetGameContentsCategory,
    question_type: questionType as GetGameContentsQuestionType,
    difficulty_level: difficulty,
    question_count: 10,
  });

  const gameResultsHook = session_student_id
    ? useGameResults({
        session_student_id,
        session_id,
        student_id,
      })
    : null;

  const startCard = gameResultsHook?.startCard;

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
    // Get available words
    const availableWords = gameContents.filter(
      (content) => !completedIds.has(content.id)
    );

    // Special handling for single word
    if (availableWords.length === 1) {
      const finalWord = availableWords[0];

      setSelectedWordIndex(0);
      setIsSpinning(false);
      setShowActionButtons(true);

      if (gameResultsHook) {
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
      if (gameResultsHook) {
        startCard?.(availableWords[selectedIndex]);
        setCardStartTime(Date.now());
      }
    }, 4000);
  };

  const handleMarkResult = (isCorrect: boolean) => {
    if (selectedWordIndex === null) return;

    // Get available words to find the actual content
    const availableWords = gameContents.filter(
      (content) => !completedIds.has(content.id)
    );
    const selectedContent = availableWords[selectedWordIndex];

    if (!selectedContent) return;

    setShowActionButtons(false);

    if (isCorrect) {
      // Mark as completed by adding ID
      setCompletedIds((prev) => new Set([...prev, selectedContent.id]));

      if (gameResultsHook) {
        const finalTime = Math.floor(
          (Date.now() - (cardStartTime || Date.now())) / 1000
        );
        gameResultsHook.completeCard(selectedContent.id, finalTime);
      }
    }

    // Reset immediately for snappier feel
    setSelectedWordIndex(null);
    setCardStartTime(null);
    setTimeTaken(0);
  };

  const handleReset = () => {
    setSelectedWordIndex(null);
    setShowActionButtons(false);
    setSpinRotation(0);
    setCompletedIds(new Set());
    setCardStartTime(null);
    setTimeTaken(0);
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

  if (contentsError || gameContents.length === 0) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <p className="text-error mb-4">
            {contentsError
              ? "Failed to load Memory Match"
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
          <button
            onClick={() => router.back()}
            className="text-blue hover:text-blue-hover flex items-center gap-2 transition-colors"
          >
            ← Back
          </button>
          <button
            onClick={handleReset}
            className="flex items-center gap-2 text-secondary hover:text-primary transition-colors"
          >
            <RotateCw className="w-4 h-4" />
            Reset
          </button>
        </div>

        <h1 className="mb-4">Flashcards</h1>

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

        {/* Progress indicator */}
        <div className="mb-8">
          <div className="flex items-center justify-between text-sm text-secondary mb-2">
            <span>
              Card {completedIds.size + 1} of {gameContents.length}
            </span>
            <span>
              {Math.round((completedIds.size / gameContents.length) * 100)}%
              Complete
            </span>
          </div>
          <div className="w-full bg-card rounded-full h-2 border border-default">
            <div
              className="bg-blue h-full rounded-full transition-all duration-300"
              style={{
                width: `${(completedIds.size / gameContents.length) * 100}%`,
              }}
            />
          </div>
        </div>

        {/* Spinner Game */}
        {gameContents.length > 0 && (
          <div className="mb-8">
            {/* Container for wheel and side buttons */}
            <div className="relative flex items-center justify-center gap-4">
              {/* Incorrect button - left side */}
              {selectedWordIndex !== null && showActionButtons && (
                <button
                  onClick={() => handleMarkResult(false)}
                  className="absolute left-0 px-4 py-2 bg-red-300 text-white rounded-lg hover:bg-red-400 transition-colors animate-fade-in"
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
                words={gameContents}
                onSpin={handleSpin}
                isSpinning={isSpinning}
                selectedWord={
                  selectedWordIndex !== null
                    ? gameContents.filter((c) => !completedIds.has(c.id))[
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
                  className="absolute right-0 px-4 py-2 bg-green-300 text-white rounded-lg hover:bg-green-400 transition-colors animate-fade-in"
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
                      gameContents.filter((c) => !completedIds.has(c.id))[
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
        {completedIds.size === gameContents.length && (
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

              {gameResultsHook && (
                <button
                  onClick={() => gameResultsHook.saveAllResults()}
                  disabled={gameResultsHook.isSaving}
                  className="px-6 py-2 bg-pink text-white rounded-lg hover:bg-pink-hover transition-colors disabled:bg-pink-disabled"
                >
                  {gameResultsHook.isSaving ? "Saving..." : "Save Progress"}
                </button>
              )}
            </div>

            {gameResultsHook?.saveError && (
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
