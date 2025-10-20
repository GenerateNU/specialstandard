import { ArrowRight, Users } from 'lucide-react'
import Image from 'next/image'

import Link from 'next/link'
import AppLayout from '@/components/AppLayout'

export default function Home() {
  return (
    <AppLayout>
      <div className="font-sans grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
        <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
          <Image
            src="/tss.png"
            alt="The Special Standard logo"
            width={180}
            height={38}
            priority
          />
          <div className="text-center sm:text-left">
            <h1 className="text-4xl font-bold mb-4 tracking-tight text-primary">
              Welcome to The Special Standard!
            </h1>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 w-full max-w-2xl">
            <Link
              href="/students"
              className="group p-6 bg-card rounded-xl shadow-lg hover:shadow-xl transition-all duration-200 border border-default hover:bg-card-hover hover:border-hover"
            >
              <div className="flex items-center justify-between mb-3">
                <Users className="w-10 h-10 text-accent" />
                <ArrowRight className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform group-hover:text-accent" />
              </div>
              <h2 className="text-xl font-semibold text-primary mb-2">
                View Students
              </h2>
              <p className="text-secondary text-sm">
                Browse and manage all student records in the system
              </p>
            </Link>

            <Link
              href="/calendar"
              className="group p-6 bg-card rounded-xl shadow-lg hover:shadow-xl transition-all duration-200 border border-default hover:bg-card-hover hover:border-hover"
            >
              <div className="flex items-center justify-between mb-3">
                <Users className="w-10 h-10 text-accent" />
                <ArrowRight className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform group-hover:text-accent" />
              </div>
              <h2 className="text-xl font-semibold text-primary mb-2">
                View Calendar
              </h2>
              <p className="text-secondary text-sm">
                View and manage events and schedules
              </p>
            </Link>

            <Link
              href="/curriculum"
              className="group p-6 bg-card rounded-xl shadow-lg hover:shadow-xl transition-all duration-200 border border-default hover:bg-card-hover hover:border-hover"
            >
              <div className="flex items-center justify-between mb-3">
                <Users className="w-10 h-10 text-accent" />
                <ArrowRight className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform group-hover:text-accent" />
              </div>
              <h2 className="text-xl font-semibold text-primary mb-2">
                View Curriculum
              </h2>
              <p className="text-secondary text-sm">
                View and manage curriculum resources
              </p>
            </Link>

            <Link
              href="/showcase"
              className="group p-6 bg-card rounded-xl shadow-lg hover:shadow-xl transition-all duration-200 border border-default hover:bg-card-hover hover:border-hover"
            >
              <div className="flex items-center justify-between mb-3">
                <Users className="w-10 h-10 text-accent" />
                <ArrowRight className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform group-hover:text-accent" />
              </div>
              <h2 className="text-xl font-semibold text-primary mb-2">
                View Components
              </h2>
              <p className="text-secondary text-sm">
                Browse our UI component library
              </p>
            </Link>

            <div className="p-6 bg-accent-light rounded-xl border border-default opacity-60 cursor-not-allowed">
              <div className="flex items-center justify-between mb-3">
                <div className="w-10 h-10 bg-card-hover rounded-lg"></div>
                <ArrowRight className="w-5 h-5 text-muted" />
              </div>
              <h2 className="text-xl font-semibold text-primary mb-2">
                More Features
              </h2>
              <p className="text-secondary text-sm">Coming soon...</p>
            </div>
          </div>
        </main>
      </div>
    </AppLayout>
  )
}
