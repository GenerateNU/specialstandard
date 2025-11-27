import type {Metadata} from 'next';
import { LoadingSpinner } from '@/app/games/flashcards/page';
import WordImageMatchingContent from "@/app/games/word-image-match/WordImageMatchingContent";
import {Suspense} from "react";

export const metadata: Metadata = {
  title: 'Word-Image Matching Game',
  description: 'Match words with their correct images',
}

export default function WordImageMatchingPage() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <WordImageMatchingContent />
    </Suspense>
  )
}