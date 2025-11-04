'use client'

import { Loader2 } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import AppLayout from '@/components/AppLayout'
import { Button } from '@/components/ui/button'
import UpcomingSessionCard from '@/components/UpcomingSessionCard'
import { useAuthContext } from '@/contexts/authContext'

export default function Home() {
  const { isAuthenticated, isLoading } = useAuthContext()
  const router = useRouter()
  const CORNER_ROUND = 'rounded-2xl'

  // Redirect to login if not authenticated
  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login')
    }
  }, [isAuthenticated, isLoading, router])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!isAuthenticated) {
    return null
  }

  return (
    <AppLayout>
      <div className="grow bg-card-hover flex flex-row h-screen">
        {/* Main Content */}
        <div className="w-full md:w-2/3 p-10 flex flex-col gap-10 overflow-y-scroll">
          <div className="flex flex-row justify-between items-end shrink-0">
            <div className="flex flex-col items-left justify-start">
              <div className="text-4xl font-serif font-bold">Your Educator Dashboard</div>
              <div className="text-2xl">
                {new Date().toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric' })}
              </div>

            </div>
            <Button variant="outline">
              Download Newsletter
            </Button>
          </div>
          <div className={`w-full  h-96 shrink-0 bg-background p-6 gap-4 font-semibold text-xl ${CORNER_ROUND} flex flex-col`}>
            Upcoming Sessions
            <div className="flex flex-row w-full gap-3 flex-1 min-h-0">
              <div className="w-1/3 h-full gap-2 flex flex-col overflow-y-scroll pr-2">
                <UpcomingSessionCard
                  sessionName="Session Name"
                  startTime="9:15am"
                  endTime="10:00am"
                  date="11/04/2025"
                />
                <UpcomingSessionCard
                  sessionName="Session Name"
                  startTime="9:15am"
                  endTime="10:00am"
                  date="11/04/2025"
                />
                <UpcomingSessionCard
                  sessionName="Session Name"
                  startTime="9:15am"
                  endTime="10:00am"
                  date="11/04/2025"
                />
                <UpcomingSessionCard
                  sessionName="Session Name"
                  startTime="9:15am"
                  endTime="10:00am"
                  date="11/04/2025"
                />
                <UpcomingSessionCard
                  sessionName="Session Name"
                  startTime="9:15am"
                  endTime="10:00am"
                  date="11/04/2025"
                />
              </div>
              <div className={`w-2/3 h-full bg-accent p-4 bg-red-500 ${CORNER_ROUND}`}></div>
            </div>
          </div>
          <div className={`w-full h-96 shrink-0 bg-background ${CORNER_ROUND}`}></div>
        </div>
        {/* Sidebar */}
        <div className="hidden md:block w-1/3 h-screen bg-red-500 sticky top-0"></div>
      </div>
    </AppLayout>
  )
}
