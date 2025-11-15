'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { CheckCircle } from 'lucide-react'
import Image from 'next/image'


export default function CompletePage() {
  const router = useRouter()
  const [userName, setUserName] = useState('')
  
  useEffect(() => {
    // Get user name from localStorage
    const onboardingData = localStorage.getItem('onboardingData')
    if (onboardingData) {
      const data = JSON.parse(onboardingData)
      setUserName(data.firstName || '')
    }
    
    // Clean up onboarding data from localStorage (optional)
    localStorage.removeItem('onboardingData')
    localStorage.removeItem('therapistProfile')
    localStorage.removeItem('onboardingStudents')
    localStorage.removeItem('onboardingSessions')
  }, [])
  
  const handleGoToDashboard = () => {
    router.push('/')
  }

  return (
    <div className="flex items-center justify-center min-h-screen p-8 bg-background">
      <div className="max-w-lg w-full bg-white rounded-2xl shadow-lg p-12 text-center">
        <div className="mb-8 flex justify-center">
          <div className="w-24 h-24 bg-accent-light rounded-full flex items-center justify-center">
            <CheckCircle className="w-12 h-12 text-accent" />
          </div>
        </div>
        
        <h1 className="text-3xl font-bold text-primary mb-4">
          {userName ? `${userName}, you're all set!` : "_, you're all set!"}
        </h1>
        
        <p className="text-secondary mb-8">
          Your account is ready to go. You can head straight to your dashboard and view your sessions on your schedule!
        </p>
        
        <div className="mb-12 bg-white rounded-md p-1 w-fit mx-auto">
          <Image
            src="/tss.png"
            alt="The Special Standard"
            width={140}
            height={30}
            priority
          />
        </div>
        <Button
          onClick={handleGoToDashboard}
          className="bg-white text-primary border-2 border-primary hover:bg-gray-50 px-8 py-2"
        >
          Dashboard
        </Button>
      </div>
    </div>
  )
}