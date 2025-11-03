/**
 * Avatar utility functions for consistent avatar generation across the app
 */

export type AvatarVariant = 'avataaars' | 'lorelei' | 'micah' | 'miniavs' | 'big-smile' | 'personas'

const AVATAR_VARIANTS: AvatarVariant[] = ['avataaars', 'lorelei', 'micah', 'miniavs', 'big-smile', 'personas']

/**
 * Get a deterministic avatar variant based on a student ID
 * Uses a simple hash function to ensure the same ID always gets the same avatar
 * @param id - Student ID (optional)
 * @returns Avatar variant string
 */
export function getAvatarVariant(id?: string): AvatarVariant {
  // Default to first variant if no ID
  if (!id) {
    return AVATAR_VARIANTS[0]
  }

  // Simple hash function to get consistent index
  let hash = 0
  for (let i = 0; i < id.length; i++) {
    hash = ((hash << 5) - hash) + id.charCodeAt(i)
    hash = hash & hash // Convert to 32-bit integer
  }

  return AVATAR_VARIANTS[Math.abs(hash) % AVATAR_VARIANTS.length]
}

/**
 * Get student initials from first and last name
 * @param firstName - Student's first name (optional)
 * @param lastName - Student's last name (optional)
 * @returns Two-letter initials, defaults to "??" if names are missing
 */
export function getStudentInitials(firstName?: string, lastName?: string): string {
  const first = firstName?.charAt(0) || '?'
  const last = lastName?.charAt(0) || '?'
  return `${first}${last}`.toUpperCase()
}

/**
 * Generate a consistent avatar name for DiceBear
 * Format: "FirstName LastName{id}" for deterministic avatar generation
 * @param firstName - Student's first name
 * @param lastName - Student's last name
 * @param id - Student's ID
 * @returns Formatted name string for avatar
 */
export function getAvatarName(firstName: string, lastName: string, id: string): string {
  return `${firstName} ${lastName}${id}`
}
