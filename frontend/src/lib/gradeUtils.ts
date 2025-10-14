// Utility functions to handle grade conversion between display (K-12) and storage (0-12)

export function gradeToDisplay(grade?: string | number | null): string {
  if (grade === null || grade === undefined)
    return 'not specified'
  const numGrade = typeof grade === 'string' ? Number.parseInt(grade) : grade
  if (numGrade === 0)
    return 'K'
  return numGrade.toString()
}

export function gradeToStorage(displayGrade: string): number | undefined {
  if (!displayGrade || displayGrade.trim() === '')
    return undefined
  if (displayGrade.toUpperCase() === 'K')
    return 0
  const num = Number.parseInt(displayGrade)
  return Number.isNaN(num) ? undefined : num
}

export const gradeOptions = [
  { value: 'K', label: 'Kindergarten (K)' },
  { value: '1', label: '1st Grade' },
  { value: '2', label: '2nd Grade' },
  { value: '3', label: '3rd Grade' },
  { value: '4', label: '4th Grade' },
  { value: '5', label: '5th Grade' },
  { value: '6', label: '6th Grade' },
  { value: '7', label: '7th Grade' },
  { value: '8', label: '8th Grade' },
  { value: '9', label: '9th Grade' },
  { value: '10', label: '10th Grade' },
  { value: '11', label: '11th Grade' },
  { value: '12', label: '12th Grade' },
]
