// components/AppLayout.tsx
'use client'

import type { ReactNode } from 'react'
import { Calendar, Component, GraduationCap, Home, LogOut, PanelLeft, Users } from 'lucide-react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import Tooltip from '@/components/ui/tooltip'
import { useAuthContext } from '@/contexts/authContext'

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
  const [isSidebarOpen, setIsSidebarOpen] = useState(false)
  const pathname = usePathname()
  const { logout } = useAuthContext()

  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen)

  // Close mobile sidebar on navigation
  useEffect(() => {
    setIsSidebarOpen(false)
  }, [pathname])

  // Close mobile sidebar on window resize to desktop
  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth >= 1024) {
        setIsSidebarOpen(false)
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
          bg-black border-r border-default
          transition-all duration-300 ease-in-out
          shadow-md
          h-screen
          z-50 lg:z-auto
          
          /* Mobile: fixed positioning, slide in/out from left */
          fixed lg:sticky
          top-0
          left-0
          
          ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
          
          /* Desktop: toggle width */
          ${isSidebarOpen ? 'w-2/3 lg:w-1/6' : 'lg:w-14'}
        `}
      >
        <div className="flex flex-col h-full">
          <div className="flex flex-row items-center justify-between p-2">
            <div className="font-bold text-lg font-serif flex items-center overflow-hidden h-full">
              <span className={`whitespace-nowrap ml-2 text-white h-full transition-opacity duration-200 leading-tight ${isSidebarOpen ? 'opacity-100' : 'opacity-0 lg:opacity-0'}`}>
                The Special
                <br />
                Standard
              </span>
            </div>
            <Button
              onClick={toggleSidebar}
              variant="secondary"
              size="icon"
              aria-label={isSidebarOpen ? 'Close sidebar' : 'Open sidebar'}
              className="w-10 h-10 flex shrink-0 p-0 bg-white/10 hover:bg-white/20 border-white/20"
            >
              <PanelLeft className="w-5 h-5 transition-transform duration-300 text-white" />
            </Button>
          </div>

          <Separator className="bg-white/10" />

          {/* Navigation */}
          <nav className="flex-1 p-2 space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon
              const isActive = pathname === item.href

              return (
                <Tooltip key={item.href} content={item.label} enabled={!isSidebarOpen}>
                  <Link
                    href={item.href}
                    className={`
                    flex items-center gap-3 rounded-lg
                    transition-all duration-200
                    h-10 shrink-0 overflow-hidden justify-start
                    ${isSidebarOpen ? 'w-full px-2.5' : 'lg:w-10 w-full px-2.5 lg:px-2.5'}
                    ${
                isActive
                  ? 'bg-blue text-white font-medium'
                  : 'text-white/70 hover:bg-white/10 hover:text-white'
                }
                  `}
                  >

                    <Icon className="w-5 h-5 shrink-0" />
                    <span className="whitespace-nowrap">{item.label}</span>
                  </Link>
                </Tooltip>
              )
            })}
          </nav>

          {/* Footer */}
          <Separator className="bg-white/10" />
          <div className="p-2">
            <Tooltip content="Logout" enabled={!isSidebarOpen}>
              <Button
                onClick={() => {
                  logout()
                  window.location.href = '/login'
                }}
                variant="ghost"
                aria-label="Logout"
                className={`
                flex items-center gap-3 rounded-lg
                transition-all duration-200
                h-10 shrink-0 overflow-hidden justify-start
                text-white/70 hover:text-white hover:bg-white/10
                ${isSidebarOpen ? 'w-full px-2.5' : 'lg:w-10 w-full px-2.5 lg:px-2.5'}
              `}
              >

                <LogOut className="w-5 h-5 mr-0.5" />
                <span className="whitespace-nowrap">Logout</span>
              </Button>
            </Tooltip>
          </div>
        </div>
      </aside>

      {/* Main content area - add left padding on desktop to account for sidebar */}
      <div className="flex-1 flex flex-col min-w-0 overflow-y-auto transition-all duration-300">
        {/* Page content */}
        <main className="flex-1 w-full">
          {children}
        </main>
      </div>
    </div>
  )
}
