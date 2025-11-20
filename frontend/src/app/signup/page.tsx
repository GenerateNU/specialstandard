'use client'

import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import CustomAlert from '@/components/ui/CustomAlert'
import { Input } from '@/components/ui/input'
import { useAuthContext } from '@/contexts/authContext'
import { validatePassword } from '@/lib/validatePassword'

export default function SignupPage() {

  const { isAuthenticated, isLoading } = useAuthContext()
  const router = useRouter()

  // Redirect if already authenticated
  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.push('/')
    }
    else {
      router.push('/signup/welcome')
    }
  }, [isAuthenticated, isLoading, router])

  useEffect(() => {
    if (error)
      setShowError(true)
  }, [error])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    )
  }

  if (isAuthenticated) {
    return null
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    // Validate passwords match
    if (password !== confirmPassword) {
      setError('Passwords do not match')
      return
    }

    // Validate password strength
    const passwordError = validatePassword(password)
    if (passwordError) {
      setError(passwordError)
      return
    }

    setIsLoading(true)

    try {
      await signup({
        email,
        password,
        ...(firstName && { first_name: firstName }),
        ...(lastName && { last_name: lastName }),
      })
      // AuthContext will handle redirect to /
    }
    catch (err: unknown) {
      console.error('Signup error:', err)

      // Type guard for axios error
      const errorData = (err as any)?.response?.data

      // Check for "user already exists" error
      if (errorData?.error_code === 'user_already_exists' || errorData?.msg?.includes('already registered')) {
        setError('This email is already registered. Please try logging in instead.')
      }
      // Handle various error response formats
      else if (errorData?.message) {
        const message = errorData.message
        // If message is an object (validation errors), extract a meaningful error
        if (typeof message === 'object' && message !== null) {
          const errorMessages = Object.values(message).filter(v => typeof v === 'string').join(', ')
          setError(errorMessages || 'Validation error occurred')
        }
        else if (typeof message === 'string') {
          setError(message)
        }
        else {
          setError('An error occurred during signup. Please try again.')
        }
      }
      else if (errorData?.msg) {
        setError(errorData.msg)
      }
      else if (err instanceof Error && err.message) {
        setError(err.message)
      }
      else {
        setError('An error occurred during signup. Please try again.')
      }
    }
    finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-background px-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <Image
            src="/tss.png"
            alt="The Special Standard logo"
            width={180}
            height={38}
            className="mx-auto mb-6"
            priority
          />
          <h1 className="text-3xl font-bold text-primary mb-2">Create Account</h1>
          <p className="text-secondary">Sign up to get started</p>
        </div>

        <div className="bg-card rounded-lg shadow-lg border border-default p-8 flex flex-col gap-2">
          {showError && error && (
            <CustomAlert
              variant="destructive"
              title="Signup Failed"
              description={error}
              onClose={() => {
                setShowError(false)
                setError(null)
              }}
            />
          )}

          <form onSubmit={handleSubmit} className="space-y-6 bg">
            <div>
              <label
                htmlFor="email"
                className="block text-sm font-medium text-primary mb-2"
              >
                Email *
              </label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={e => setEmail(e.target.value)}
                required
                disabled={isLoading}
                placeholder="therapist@example.com"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label
                  htmlFor="firstName"
                  className="block text-sm font-medium text-primary mb-2"
                >
                  First Name
                </label>
                <Input
                  id="firstName"
                  type="text"
                  value={firstName}
                  onChange={e => setFirstName(e.target.value)}
                  disabled={isLoading}
                  placeholder="John"
                />
              </div>

              <div>
                <label
                  htmlFor="lastName"
                  className="block text-sm font-medium text-primary mb-2"
                >
                  Last Name
                </label>
                <Input
                  id="lastName"
                  type="text"
                  value={lastName}
                  onChange={e => setLastName(e.target.value)}
                  disabled={isLoading}
                  placeholder="Doe"
                />
              </div>
            </div>

            <div>
              <label
                htmlFor="password"
                className="block text-sm font-medium text-primary mb-2"
              >
                Password *
              </label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                required
                disabled={isLoading}
                placeholder="••••••••"
              />
              <p className="text-xs text-secondary mt-1">
                Must be 8+ characters with uppercase, lowercase, number, and special character
              </p>
            </div>

            <div>
              <label
                htmlFor="confirmPassword"
                className="block text-sm font-medium text-primary mb-2"
              >
                Confirm Password *
              </label>
              <Input
                id="confirmPassword"
                type="password"
                value={confirmPassword}
                onChange={e => setConfirmPassword(e.target.value)}
                required
                disabled={isLoading}
                placeholder="••••••••"
              />
            </div>

            <Button type="submit" disabled={isLoading} size="long">
              {isLoading
                ? (
                    <>
                      <Loader2 className="w-5 h-5 animate-spin" />
                      <span>Creating account...</span>
                    </>
                  )
                : (
                    <>
                      <UserPlus className="w-5 h-5" />
                      <span>Sign Up</span>
                    </>
                  )}
            </Button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-secondary">
              Already have an account?
              {' '}
              <Link href="/login">
                <Button variant="link" size="sm">
                  Sign in
                </Button>
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
