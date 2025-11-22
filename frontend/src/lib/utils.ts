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

// Format ISO datetime string to time (e.g., "2:30 PM")
export function formatTime(datetime: string) {
  return new Date(datetime).toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
  })
}

// Format ISO datetime string to date (e.g., "11/06/2025")
export function formatDateString(datetime: string) {
  return new Date(datetime).toLocaleDateString('en-US', {
    month: '2-digit',
    day: '2-digit',
    year: 'numeric',
  })
}

// Get therapist full name from therapist list by ID
export function getTherapistName(
  therapistId: string,
  therapists: Array<{ id: string, first_name: string, last_name: string }>,
) {
  const therapist = therapists.find(t => t.id === therapistId)
  return therapist ? `${therapist.first_name} ${therapist.last_name}` : 'Unknown Therapist'
}
