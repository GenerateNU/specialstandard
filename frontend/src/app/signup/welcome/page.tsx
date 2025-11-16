'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Loader2 } from 'lucide-react'
import { useAuthContext } from '@/contexts/authContext'
import CustomAlert from '@/components/ui/CustomAlert'

export default function WelcomePage() {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [showError, setShowError] = useState(false)
  const { signup } = useAuthContext()
  
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    password: '',
  })

  // Password validation states
  const passwordChecks = {
    length: formData.password.length >= 8,
    uppercase: /[A-Z]/.test(formData.password),
    lowercase: /[a-z]/.test(formData.password),
    number: /\d/.test(formData.password),
    symbol: /[!@#$%^&*()_+\-=[\]{};':"\\|,.<>?/~`]/.test(formData.password),
  }

  const allPasswordChecksPassed = Object.values(passwordChecks).every(check => check)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!allPasswordChecksPassed) {
      setError('Please ensure your password meets all requirements')
      setShowError(true)
      return
    }

    setIsLoading(true)
    setError(null)

    try {
        // First, signup the user
        await signup({
        email: formData.email,
        password: formData.password,
        first_name: formData.firstName,
        last_name: formData.lastName,
      })

        // Save data to localStorage for next step
        localStorage.setItem('onboardingData', JSON.stringify({
          firstName: formData.firstName,
          lastName: formData.lastName,
          email: formData.email,
        }))

      
    } catch (err: any) {
      console.error('Signup error:', err)
      
      const errorData = err?.response?.data

      if (errorData?.error_code === 'user_already_exists' || errorData?.msg?.includes('already registered')) {
        setError('This email is already registered. Please try logging in instead.')
      } else if (errorData?.message) {
        const message = errorData.message
        if (typeof message === 'object' && message !== null) {
          const errorMessages = Object.values(message).filter(v => typeof v === 'string').join(', ')
          setError(errorMessages || 'Validation error occurred')
        } else if (typeof message === 'string') {
          setError(message)
        } else {
          setError('An error occurred during signup. Please try again.')
        }
      } else if (errorData?.msg) {
        setError(errorData.msg)
      } else if (err instanceof Error && err.message) {
        setError(err.message)
      } else {
        setError('An error occurred during signup. Please try again.')
      }
      setShowError(true)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen p-8">
      <div className="max-w-md w-full">
        <h1 className="text-5xl font-bold text-primary mb-8">
          Let's Get Started!
        </h1>
        
        {showError && error && (
          <div className="mb-4">
            <CustomAlert
              variant="destructive"
              title="Signup Failed"
              description={error}
              onClose={() => {
                setShowError(false)
                setError(null)
              }}
            />
          </div>
        )}
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <Input
                value={formData.firstName}
                onChange={(e) => setFormData({...formData, firstName: e.target.value})}
                placeholder="First Name"
                className="border border-accent"
                required
                disabled={isLoading}
              />
            </div>
            
            <div>
              <Input
                value={formData.lastName}
                onChange={(e) => setFormData({...formData, lastName: e.target.value})}
                placeholder="Last Name"
                className="border border-accent"
                required
                disabled={isLoading}
              />
            </div>
          </div>
          
          <div>
            <Input
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({...formData, email: e.target.value})}
              placeholder="Email"
              className="border border-accent"
              required
              disabled={isLoading}
            />
          </div>
          
          <div>
            <Input
              type="password"
              value={formData.password}
              onChange={(e) => setFormData({...formData, password: e.target.value})}
              placeholder="Password"
              className="border border-accent"
              required
              disabled={isLoading}
            />
          </div>
          
          <div className={`flex items-start gap-2 text-xs ${passwordChecks.length ? 'text-success' : 'text-secondary'}`}>
            <span>{passwordChecks.length ? '✓' : '•'}</span>
            <span>Use 8 or more characters</span>
          </div>
          
          <div className={`flex items-start gap-2 text-xs ${passwordChecks.uppercase && passwordChecks.lowercase ? 'text-success' : 'text-secondary'}`}>
            <span>{passwordChecks.uppercase && passwordChecks.lowercase ? '✓' : '•'}</span>
            <span>Use upper and lower case letters (e.g. Aa)</span>
          </div>
          
          <div className={`flex items-start gap-2 text-xs ${passwordChecks.number ? 'text-success' : 'text-secondary'}`}>
            <span>{passwordChecks.number ? '✓' : '•'}</span>
            <span>Use a number (e.g. 1234)</span>
          </div>
          
          <div className={`flex items-start gap-2 text-xs ${passwordChecks.symbol ? 'text-success' : 'text-secondary'}`}>
            <span>{passwordChecks.symbol ? '✓' : '•'}</span>
            <span>Use a symbol (e.g. !@#$)</span>
          </div>
          
          <div className="pt-4">
            <Button
              type="submit"
              size="long"
              className="w-full font-light hover:bg-accent-hover text-white"
              disabled={isLoading || !allPasswordChecksPassed}
            >
              {isLoading ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin mr-2" />
                  <span>Creating account...</span>
                </>
              ) : (
                <span>Sign up</span>
              )}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}