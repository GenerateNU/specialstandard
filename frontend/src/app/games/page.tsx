"use client";

import AppLayout from "@/components/AppLayout";
import { GameContentSelector } from "@/components/games/GameContentSelector";
import type {
  GetGameContentsCategory,
  GetGameContentsQuestionType,
  Theme,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import {BookOpen, Brain, Gamepad2, Image, Shuffle} from "lucide-react";
import { useRouter } from "next/navigation";
import React from "react";

export default function GamesPage() {
  const router = useRouter();
  const [selectedContent, setSelectedContent] = React.useState<{
    theme: Theme;
    difficultyLevel: number;
    category: GetGameContentsCategory;
    questionType: GetGameContentsQuestionType;
  } | null>(null);

  const handleContentSelection = (selection: {
    theme: Theme;
    difficultyLevel: number;
    category: GetGameContentsCategory;
    questionType: GetGameContentsQuestionType;
  }) => {
    setSelectedContent(selection);
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
            <h1 className="mb-8">Select a Game</h1>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <button
                onClick={() => {
                  const params = new URLSearchParams({
                    themeId: selectedContent.theme.id,
                    difficulty: String(selectedContent.difficultyLevel),
                    category: selectedContent.category,
                    questionType: selectedContent.questionType,
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
      </AppLayout>
    );
  }

  // Show content selector
  return (
    <AppLayout>
      <GameContentSelector onSelectionComplete={handleContentSelection} />
    </AppLayout>
  );
}
