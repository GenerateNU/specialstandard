// components/AppLayout.tsx
'use client'

import type { ReactNode } from 'react'
import { Calendar, Component, GraduationCap, Home, PanelLeft, Users } from 'lucide-react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useEffect, useState } from 'react'
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
  { href: '/showcase', label: 'Components', icon: Component },
  { href: '/calendar', label: 'Calendar', icon: Calendar },
  { href: '/curriculum', label: 'Curriculum', icon: GraduationCap },
]

export default function AppLayout({ children }: AppLayoutProps) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(true)
  const pathname = usePathname()

  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen)

  // Close mobile sidebar on navigation
  useEffect(() => {
    setIsSidebarOpen(false)
  }, [pathname])

  // Close mobile sidebar on window resize to desktop
  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth >= 1024) {
        setIsSidebarOpen(true)
      }
    }

    handleResize()
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  return (
    <div className="flex h-screen bg-background">
      {/* Mobile backdrop overlay */}
      {isSidebarOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={() => setIsSidebarOpen(false)}
          aria-hidden="true"
        />
      )}

      {/* Mobile Toggle Button (visible only on mobile when sidebar is closed) */}
      {!isSidebarOpen && (
        <div className="fixed top-4 left-4 z-30 lg:hidden">
          <Button
            onClick={toggleSidebar}
            variant="secondary"
            size="icon"
            aria-label="Open sidebar"
            className="w-10 h-10 shadow-lg"
          >
            <PanelLeft className="w-5 h-5" />
          </Button>
        </div>
      )}

      {/* Actual Sidebar */}
      <aside
        className={`
          bg-card border-r border-default
          transition-all duration-300 ease-in-out
          shadow-md
          
          /* Mobile: Fixed overlay modal */
          fixed lg:sticky
          top-0
          h-full lg:h-screen
          z-50 lg:z-auto
          
          /* Mobile: slide in/out from left */
          ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
          
          /* Desktop: toggle width */
          ${isSidebarOpen ? 'w-70' : 'lg:w-14'}
        `}
      >
        <div className="flex flex-col h-full">
          <div className="flex flex-row items-center justify-between p-2">
            <div className="font-bold text-xl flex items-center overflow-hidden">
              <span className={`whitespace-nowrap ml-2 transition-opacity duration-200 ${isSidebarOpen ? 'opacity-100' : 'opacity-0 lg:opacity-0'}`}>
                The Special Standard
              </span>
            </div>
            <Button
              onClick={toggleSidebar}
              variant="secondary"
              size="icon"
              aria-label={isSidebarOpen ? 'Close sidebar' : 'Open sidebar'}
              className="w-10 h-10 flex-shrink-0 p-0"
            >
              <PanelLeft className="w-5 h-5 transition-transform duration-300" />
            </Button>
          </div>

          <Separator />

          {/* Navigation */}
          <nav className="flex-1 p-2 space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon
              const isActive = pathname === item.href

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`
                    flex items-center gap-3 rounded-lg
                    transition-all duration-200
                    h-10 flex-shrink-0 overflow-hidden justify-start
                    ${isSidebarOpen ? 'w-full px-2.5' : 'lg:w-10 w-full px-2.5 lg:px-2.5'}
                    ${
                isActive
                  ? 'bg-accent text-white font-medium'
                  : 'text-secondary hover:bg-card-hover hover:text-primary'
                }
                  `}
                >
                  <Icon className="w-5 h-5 flex-shrink-0" />
                  <span className="whitespace-nowrap">{item.label}</span>
                </Link>
              )
            })}
          </nav>

          {/* Footer */}
          {isSidebarOpen && (
            <>
              <Separator />
              <div className="p-4">
                <p className="text-xs text-muted text-center whitespace-nowrap">
                  The Special Standard © 2025
                </p>
              </div>
            </>
          )}
        </div>
      </aside>

      {/* Main content area */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Page content */}
        <main className="flex-1">
          {children}
        </main>
      </div>
    </div>
  )
}
