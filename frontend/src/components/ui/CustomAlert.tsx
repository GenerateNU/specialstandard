import { AlertCircle, AlertTriangle, CheckCircle, X } from 'lucide-react'
import React from 'react'
import { Alert, AlertDescription, AlertTitle } from './alert'

export interface CustomAlertProps {
  variant?: 'default' | 'destructive' | 'warning' | 'success'
  title: string
  description?: string
  onClose?: () => void
}

const iconMap = {
  default: AlertCircle,
  destructive: AlertCircle,
  warning: AlertTriangle,
  success: CheckCircle,
}

function CustomAlert({
  variant = 'default',
  title,
  description,
  onClose,
}: CustomAlertProps) {
  const Icon = iconMap[variant]

  return (
    <Alert variant={variant}>
      <div className="flex items-start space-x-4">
        <Icon className="h-4 w-4 mt-0.5" />
        <div className="flex-1">
          <AlertTitle>{title}</AlertTitle>
          {description && <AlertDescription>{description}</AlertDescription>}
        </div>
        {onClose && (
          <button
            onClick={onClose}
            title="Close alert"
            className="rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-accent focus:ring-offset-2"
          >
            <X className="h-4 w-4" />
          </button>
        )}
      </div>
    </Alert>
  )
}

export default CustomAlert
