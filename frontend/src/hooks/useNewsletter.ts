import { useState } from 'react'
import { getNewsletter } from '@/lib/api/newsletter'

export function useNewsletter() {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const downloadNewsletter = async (date?: Date) => {
    setIsLoading(true)
    setError(null)

    try {
      // Use provided date or default to today
      const targetDate = date || new Date()
      const yyyy = targetDate.getFullYear()
      const mm = String(targetDate.getMonth() + 1).padStart(2, '0')
      const dd = String(targetDate.getDate()).padStart(2, '0')
      const dateStr = `${yyyy}-${mm}-${dd}`

      // Use orval-generated API call
      const { getNewsletterByDate } = getNewsletter()
      const response = await getNewsletterByDate({ date: dateStr })

      if (!response.s3_url) {
        throw new Error('No presigned URL returned from server')
      }

      // Open the presigned URL in a new tab to download
      window.open(response.s3_url, '_blank')
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to download newsletter'
      setError(errorMessage)
      console.error('Newsletter download error:', err)
      throw err
    } finally {
      setIsLoading(false)
    }
  }

  return {
    downloadNewsletter,
    isLoading,
    error,
  }
}