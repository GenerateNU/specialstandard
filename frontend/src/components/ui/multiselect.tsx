'use client'

import { Check, X } from 'lucide-react'
import * as React from 'react'
import { cn } from '@/lib/utils'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './select'

export interface MultiSelectOption {
  label: string
  value: string
  icon?: React.ReactNode
  disabled?: boolean
}

export interface MultiSelectProps {
  options: MultiSelectOption[]
  value?: string[]
  onValueChange?: (value: string[]) => void
  placeholder?: string
  className?: string
  disabled?: boolean
  maxCount?: number
  showCount?: boolean
  showTags?: boolean
}

const MultiSelect = React.forwardRef<HTMLDivElement, MultiSelectProps>(
  (
    {
      options,
      value = [],
      onValueChange,
      placeholder = 'Select items...',
      className,
      disabled,
      maxCount,
      showCount = true,
      showTags = true,
    },
    ref,
  ) => {
    const [internalValue, setInternalValue] = React.useState<string[]>(value)

    React.useEffect(() => {
      setInternalValue(value)
    }, [value])

    const handleSelect = (selectedValue: string) => {
      let newValue: string[]

      if (internalValue.includes(selectedValue)) {
        // Remove if already selected
        newValue = internalValue.filter(v => v !== selectedValue)
      }
      else {
        // Add if not selected and under max count
        if (maxCount && internalValue.length >= maxCount) {
          return
        }
        newValue = [...internalValue, selectedValue]
      }

      setInternalValue(newValue)
      onValueChange?.(newValue)
    }

    const handleRemove = (valueToRemove: string) => {
      const newValue = internalValue.filter(v => v !== valueToRemove)
      setInternalValue(newValue)
      onValueChange?.(newValue)
    }

    const getDisplayText = () => {
      if (internalValue.length === 0) {
        return placeholder
      }

      if (showCount && internalValue.length > 0) {
        const itemText = internalValue.length === 1 ? 'item' : 'items'
        return `${internalValue.length} ${itemText} selected`
      }

      // Show first few selected items
      const selectedLabels = internalValue
        .slice(0, 2)
        .map(v => options.find(opt => opt.value === v)?.label)
        .filter(Boolean)
        .join(', ')

      if (internalValue.length > 2) {
        return `${selectedLabels}, +${internalValue.length - 2} more`
      }

      return selectedLabels
    }

    return (
      <div ref={ref} className="space-y-2">
        <Select
          value={internalValue[internalValue.length - 1] || ''}
          onValueChange={handleSelect}
          disabled={disabled}
        >
          <SelectTrigger className={cn('w-full', className)}>
            <SelectValue>
              {getDisplayText()}
            </SelectValue>
          </SelectTrigger>
          <SelectContent className="max-h-44 bg-background overflow-y-auto">
            {options.map((option) => {
              const isSelected = internalValue.includes(option.value)
              const isDisabled = !!(option.disabled
                || (maxCount && internalValue.length >= maxCount && !isSelected))

              return (
                <SelectItem
                  key={option.value}
                  value={option.value}
                  disabled={isDisabled}
                  className="cursor-pointer bg-secondary hover:bg-accent/50"
                >
                  <div className="flex items-center gap-2 w-full">
                    <div className={cn(
                      'flex h-4 w-4 items-center justify-center rounded-sm border',
                      isSelected
                        ? 'bg-primary border-primary'
                        : 'border-muted-foreground',
                    )}
                    >
                      {isSelected && (
                        <Check className="h-3 w-3 text-primary-foreground" />
                      )}
                    </div>
                    {option.icon && (
                      <span className="flex-shrink-0">{option.icon}</span>
                    )}
                    <span className="flex-1">{option.label}</span>
                  </div>
                </SelectItem>
              )
            })}
          </SelectContent>
        </Select>

        {showTags && internalValue.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {internalValue.map((val) => {
              const option = options.find(opt => opt.value === val)
              if (!option)
                return null

              return (
                <div
                  key={val}
                  className="flex items-center gap-1 bg-secondary text-secondary-foreground px-2 py-1 rounded-md text-sm"
                >
                  {option.icon && (
                    <span className="flex-shrink-0">{option.icon}</span>
                  )}
                  <span>{option.label}</span>
                  <button
                    type="button"
                    onClick={() => handleRemove(val)}
                    disabled={disabled}
                    className="ml-1 hover:text-destructive focus:outline-none disabled:opacity-50"
                  >
                    <X className="h-3 w-3" />
                  </button>
                </div>
              )
            })}
          </div>
        )}
      </div>
    )
  },
)

MultiSelect.displayName = 'MultiSelect'

export { MultiSelect }
