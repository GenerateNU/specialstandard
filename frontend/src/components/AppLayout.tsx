// components/AppLayout.tsx
'use client'

import type { ReactNode } from 'react'
import { Grid3x3, Home, Menu, Users, X } from 'lucide-react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'

interface AppLayoutProps {
  children: ReactNode
}

interface NavItem {
  href: string
  label: string
  icon: React.ComponentType<{ className?: string }>
}

const navItems: NavItem[] = [
  { href: '/', label: 'Home', icon: Home },
  { href: '/students', label: 'Students', icon: Users },
  { href: '/showcase', label: 'Components', icon: Grid3x3 },
]

export default function AppLayout({ children }: AppLayoutProps) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false)
  const pathname = usePathname()

  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen)

  return (
    <div className="flex min-h-screen bg-background">
      {/* Mobile overlay */}
      {isSidebarOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={toggleSidebar}
        />
      )}

      {/* Sidebar */}
      <aside
        className={`
          fixed lg:static inset-y-0 left-0 z-50
          w-64 bg-card border-r border-default
          transform transition-transform duration-200 ease-in-out
          ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
        `}
      >
        <div className="flex flex-col h-full">
          {/* Header */}
          <div className="flex items-center justify-between p-6">
            <h2 className="text-xl font-bold text-primary">The Special Standard</h2>
            <Button
              onClick={toggleSidebar}
              variant="ghost"
              size="icon"
              className="lg:hidden"
              aria-label="Close sidebar"
            >
              <X className="w-5 h-5" />
            </Button>
          </div>

          <Separator />

          {/* Navigation */}
          <nav className="flex-1 p-4 space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon
              const isActive = pathname === item.href

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={() => setIsSidebarOpen(false)}
                  className={`
                    flex items-center gap-3 px-4 py-3 rounded-lg
                    transition-all duration-200
                    ${
                isActive
                  ? 'bg-accent text-white font-medium'
                  : 'text-secondary hover:bg-card-hover hover:text-primary'
                }
                  `}
                >
                  <Icon className="w-5 h-5" />
                  <span>{item.label}</span>
                </Link>
              )
            })}
          </nav>

          {/* Footer */}
          <Separator className="" />
          <div className="p-4">
            <p className="text-xs text-muted text-center">
              The Special Standard © 2025
            </p>
          </div>
        </div>
      </aside>

      {/* Main content area */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Mobile header */}
        <header className="lg:hidden sticky top-0 z-30 bg-card px-4 py-3">
          <Button
            onClick={toggleSidebar}
            variant="ghost"
            size="icon"
            aria-label="Open sidebar"
          >
            <Menu className="w-6 h-6" />
          </Button>
          <Separator />
        </header>

        {/* Page content */}
        <main className="flex-1">
          {children}
        </main>
      </div>
    </div>
  )
}
