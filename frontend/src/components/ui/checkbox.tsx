'use client'

import { Check } from 'lucide-react'
import * as React from 'react'
import { cn } from '@/lib/utils'

export interface CheckboxProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  checked?: boolean | 'indeterminate'
  onCheckedChange?: (checked: boolean) => void
}

const Checkbox = React.forwardRef<HTMLButtonElement, CheckboxProps>(
  ({ className, checked = false, onCheckedChange, ...props }, ref) => {
    return (
      <button
        type="button"
        role="checkbox"
        aria-checked={
          checked === 'indeterminate'
            ? 'mixed'
            : checked === true
              ? 'true'
              : 'false'
        }
        ref={ref}
        onClick={() => onCheckedChange?.(!checked)}
        className={cn(
          'peer h-4 w-4 shrink-0 rounded-sm border border-border shadow',
          'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-accent',
          'disabled:cursor-not-allowed disabled:opacity-50',
          checked && 'bg-accent text-white border-accent',
          className,
        )}
        {...props}
      >
        {checked && <Check className="h-3 w-3 mx-auto" />}
      </button>
    )
  },
)
Checkbox.displayName = 'Checkbox'

export { Checkbox }
