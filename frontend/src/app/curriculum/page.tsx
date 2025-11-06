'use client'

import { AlertCircle, ArrowLeft, FileText, Loader2, RefreshCcw } from 'lucide-react'
import Link from 'next/link'
import AppLayout from '@/components/AppLayout'
import { useResources } from '@/hooks/useResources'

export default function Curriculum() {
  const { resources, isLoading, error, refetch } = useResources()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <Loader2 className="w-8 h-8 animate-spin text-accent mx-auto mb-4" />
          <p className="text-secondary">Loading resources...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <AppLayout>
        <div className="min-h-screen flex items-center justify-center bg-background">
          <div className="text-center max-w-md">
            <AlertCircle className="w-12 h-12 text-error mx-auto mb-4" />
            <h2 className="text-xl font-semibold text-primary mb-2">Error Loading Resources</h2>
            <p className="text-secondary mb-4">{error}</p>
            <button
              onClick={() => refetch()}
              className="px-4 py-2 bg-accent text-white rounded-lg hover:bg-accent-hover transition-colors flex items-center gap-2 mx-auto"
            >
              <RefreshCcw className="w-4 h-4" />
              Try Again
            </button>
          </div>
        </div>
      </AppLayout>
    )
  }

  return (
    <AppLayout>
      <div className="min-h-screen bg-background py-8">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Header */}
          <header className="mb-8">
            <Link
              href="/"
              className="inline-flex items-center gap-2 text-secondary hover:text-primary mb-4 transition-colors group"
            >
              <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
              <span className="text-sm font-medium">Back to Home</span>
            </Link>

            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center space-x-3">
                <FileText className="w-8 h-8 text-accent" />
                <h1 className="text-3xl font-bold text-primary">Curriculum Resources</h1>
              </div>
            </div>
            <p className="text-secondary">
              View and access all available learning materials.
            </p>
          </header>

          {/* Resource list */}
          {resources.length === 0
            ? (
                <div className="text-center py-16">
                  <FileText className="w-16 h-16 text-muted mx-auto mb-4 opacity-30" />
                  <h2 className="text-xl font-semibold text-primary mb-2">No Resources Available</h2>
                  <p className="text-secondary mb-6">
                    There are currently no learning materials uploaded. Check back later!
                  </p>
                </div>
              )
            : (
                <div className="grid gap-6">
                  {resources.filter(resource => resource.id === 'bd751100-042c-4091-8e36-28e0f3d2fd35').map(resource => (
                    <div
                      key={resource.id}
                      className="bg-card p-6 rounded-2xl shadow-sm border border-border hover:shadow-md transition-shadow"
                    >
                      <div className="flex justify-between items-center mb-3">
                        <h2 className="text-xl font-semibold text-primary">
                          {resource.title || 'Untitled Resource'}
                        </h2>
                      </div>
                      <p className="text-secondary mb-4">Meet Molly The Mink!</p>

                      {resource.presigned_url
                        ? (
                            <iframe
                              src={resource.presigned_url}
                              className="w-full h-[80vh] rounded-lg border"
                            />
                          )
                        : (
                            <p className="text-muted italic">No file available for this resource.</p>
                          )}
                    </div>
                  ))}
                </div>
              )}
        </div>
      </div>
    </AppLayout>
  )
}
