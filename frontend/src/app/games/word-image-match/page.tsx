import type {Metadata} from 'next';
import { LoadingSpinner } from '@/app/games/flashcards/page';

export const metadata: Metadata = {
  title: 'Word-Image Matching Game',
  description: 'Match words with their correct images',
}

export default function WordImageMatchingPage() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <ImageMatchingContent />
    </Suspense>
  )
}