'use client'

import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { useAuthContext } from '@/contexts/authContext'

export default function SignupPage() {
  const {isAuthenticated, isLoading} = useAuthContext()
  const router = useRouter()

  // Redirect if already authenticated
  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.push('/')
    } else {
      router.push('/signup/welcome')
    }
  }, [isAuthenticated, isLoading, router])
}