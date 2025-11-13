import Link from 'next/link'
import { BookOpen, Brain, Gamepad2 } from 'lucide-react'

export default function GamesPage() {
  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="mb-8">Select a Game</h1>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Link
            href="/games/flashcards"
            className="bg-card rounded-lg shadow-md p-8 hover:shadow-lg transition-all duration-200 group hover:bg-card-hover border border-default hover:border-hover block"
          >
            <BookOpen className="w-12 h-12 text-blue mb-4 mx-auto" />
            <h3 className="mb-2 text-center">Flashcards</h3>
            <p className="text-secondary text-sm text-center">Practice with interactive flashcards</p>
          </Link>
          
          <div className="bg-card rounded-lg shadow-md p-8 opacity-50 cursor-not-allowed border border-default">
            <Brain className="w-12 h-12 text-muted mb-4 mx-auto" />
            <h3 className="mb-2 text-muted text-center">Memory Match</h3>
            <p className="text-disabled text-sm text-center">Coming soon</p>
          </div>
          
          <div className="bg-card rounded-lg shadow-md p-8 opacity-50 cursor-not-allowed border border-default">
            <Gamepad2 className="w-12 h-12 text-muted mb-4 mx-auto" />
            <h3 className="mb-2 text-muted text-center">Quiz Game</h3>
            <p className="text-disabled text-sm text-center">Coming soon</p>
          </div>
        </div>
      </div>
    </div>
  )
}