import type { Metadata } from "next";
import { Suspense } from "react";
import { MemorymatchContent } from "./MemorymatchContent";

export const metadata: Metadata = {
  title: "Flashcard Game",
  description: "Practice with interactive flashcards",
};

function LoadingSpinner() {
  return (
    <div className="min-h-screen bg-background p-8 flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue mx-auto mb-4"></div>
        <p className="text-secondary">Loading...</p>
      </div>
    </div>
  );
}

export default function MemorymatchPage() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <MemorymatchContent />
    </Suspense>
  );
}
