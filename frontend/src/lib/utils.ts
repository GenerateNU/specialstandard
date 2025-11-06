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

// Deterministically pick a DiceBear avatar variant from an ID string
export function getAvatarVariant(id: string): 'avataaars' | 'lorelei' | 'micah' | 'miniavs' | 'big-smile' | 'personas' {
  const variants = ['avataaars', 'lorelei', 'micah', 'miniavs', 'big-smile', 'personas'] as const

  let hash = 0
  for (let i = 0; i < id.length; i++) {
    hash = ((hash << 5) - hash) + id.charCodeAt(i)
    // Convert to 32-bit integer
    hash = hash & hash
  }

  return variants[Math.abs(hash) % variants.length]
}
