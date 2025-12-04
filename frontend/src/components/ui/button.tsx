import type { VariantProps } from 'class-variance-authority'
import { Slot } from '@radix-ui/react-slot'
import { cva } from 'class-variance-authority'
import * as React from 'react'
import { cn } from '@/lib/utils'

const buttonVariants = cva(
  'inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-accent disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 cursor-pointer',
  {
    variants: {
      variant: {
        default: 'bg-pink hover:bg-pink-hover text-background shadow',
        destructive: 'bg-error text-white shadow-sm hover:bg-error/90',
        outline:
          'border border-border hover:ring-1 transition-all shadow-sm hover:text-primary',
        secondary: 'text-primary hover:-translate-y-1 hover:shadow-md transition-all',
        ghost: 'text-secondary hover:text-primary',
        link: 'text-accent underline-offset-4 hover:underline',
        tab: 'px-8 py-2 text-secondary hover:text-primary transition-colors rounded-none',
      },
      size: {
        default: 'h-9 pr-4 py-2 px-2',
        sm: 'h-8 rounded-md px-3 text-xs',
        lg: 'h-10 rounded-md px-8',
        long: 'w-full h-10 rounded-md px-8',
        icon: 'h-9 w-9',
        dropdown: 'px-4 w-full',
      },
      active: {
        true: 'border-b-2 border-accent text-primary',
        false: 'text-secondary',
      },
    },
    compoundVariants: [
      {
        variant: 'tab',
        active: true,
        className: 'border-b-2 border-accent text-primary',
      },
    ],
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
  VariantProps<typeof buttonVariants> {
  asChild?: boolean
  active?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, active, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : 'button'
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, active, className }))}
        ref={ref}
        {...props}
      />
    )
  },
)
Button.displayName = 'Button'

export { Button, buttonVariants }
