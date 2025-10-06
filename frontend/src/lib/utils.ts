import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function capitalizeFirstLetter(string: string) {
  return string.charAt(0).toUpperCase() + string.slice(1)
}

export function formatDate(date: Date, format: 'shortMonth' | 'year') {
  if (!date)
    return ''
  date = new Date(date)

  return format === 'shortMonth'
    ? date.toLocaleDateString('en-US', { month: 'short', timeZone: 'UTC' })
    : date.getFullYear()
}

export function formatISOString(date: string) {
  const dateStr = date.split('T')[0]
  const [, month, day] = dateStr.split('-')
  return `${month}/${day}`
}

export function formatNumber(number?: number, appendString: string = '') {
  if (!number)
    return appendString
  return Number(number?.toFixed()).toLocaleString()
}
