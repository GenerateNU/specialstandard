'use client'

import { ChevronDown } from 'lucide-react'
import * as React from 'react'
import { cn } from '@/lib/utils'

export interface DropdownItem {
  label: string
  value: string
  onClick?: () => void
  icon?: React.ReactNode
  disabled?: boolean
}

export interface DropdownProps {
  trigger?: React.ReactNode
  items: DropdownItem[]
  className?: string
  align?: 'left' | 'right' | 'center'
  value?: string
  onValueChange?: (value: string) => void
  placeholder?: string
}

const Dropdown = React.forwardRef<HTMLDivElement, DropdownProps>(
  ({ trigger, items, className, align = 'left', value, onValueChange, placeholder = 'Select...' }, ref) => {
    const [isOpen, setIsOpen] = React.useState(false)
    const [selectedValue, setSelectedValue] = React.useState(value)
    const dropdownRef = React.useRef<HTMLDivElement>(null)

    // Update selected value when prop changes
    React.useEffect(() => {
      setSelectedValue(value)
    }, [value])

    // Close dropdown when clicking outside
    React.useEffect(() => {
      const handleClickOutside = (event: MouseEvent) => {
        if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
          setIsOpen(false)
        }
      }

      document.addEventListener('mousedown', handleClickOutside)
      return () => {
        document.removeEventListener('mousedown', handleClickOutside)
      }
    }, [])

    const handleSelect = (item: DropdownItem) => {
      setSelectedValue(item.value)
      onValueChange?.(item.value)
      item.onClick?.()
      setIsOpen(false)
    }

    const selectedItem = items.find(item => item.value === selectedValue)
    const displayContent = trigger || (
      <span className="flex items-center gap-2">
        {selectedItem?.icon}
        {selectedItem?.label || placeholder}
      </span>
    )

    const alignmentClasses = {
      left: 'left-0',
      right: 'right-0',
      center: 'left-1/2 -translate-x-1/2',
    }

    return (
      <div ref={dropdownRef} className="relative inline-block">
        <button
          type="button"
          onClick={() => setIsOpen(!isOpen)}
          className={cn(
            'inline-flex items-center justify-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors',
            'border border-border bg-background hover:bg-card-hover',
            'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-accent w-full',
            className,
          )}
        >
          {displayContent}
          <ChevronDown className={cn(
            'h-4 w-4 transition-transform',
            isOpen && 'rotate-180',
          )}
          />
        </button>

        {isOpen && (
          <div
            ref={ref}
            className={cn(
              'absolute z-50 mt-2 overflow-hidden rounded-md border border-border bg-background shadow-lg',
              alignmentClasses[align],
            )}
          >
            <div className="p-1">
              {items.map((item, index) => (
                <button
                  key={`${item.value}-${index}`}
                  type="button"
                  disabled={item.disabled}
                  onClick={() => handleSelect(item)}
                  className={cn(
                    'flex w-full gap-2 rounded-sm px-2 py-1.5 text-sm transition-colors',
                    'hover:bg-card-hover hover:text-primary',
                    'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-accent',
                    'disabled:pointer-events-none disabled:opacity-50',
                    '[&>svg]:h-4 [&>svg]:w-4',
                    selectedValue === item.value && 'bg-accent text-foreground', className,
                  )}
                >
                  {item.icon}
                  {item.label}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>
    )
  },
)
Dropdown.displayName = 'Dropdown'

export { Dropdown }
