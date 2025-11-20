'use client'

import React, { Suspense, useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import CustomAlert from '@/components/ui/CustomAlert'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { useAuth } from '@/hooks/useAuth'

const passwordSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters long'),
  confirmPassword: z.string().min(6, 'Password must be at least 6 characters long'),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
})

const otpSchema = z.object({
  otp: z.string().min(6, 'OTP must be 6 digits').max(6, 'OTP must be 6 digits'),
})

function ResetPasswordContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { updatePassword } = useAuth()
  const [error, setError] = useState<string | null>(null)
  const [showError, setShowError] = useState(false)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)
  const [step, setStep] = useState<'password' | 'otp'>('password')
  const [userEmail, setUserEmail] = useState('')

  const passwordForm = useForm<z.infer<typeof passwordSchema>>({
    resolver: zodResolver(passwordSchema),
    defaultValues: {
      email: '',
      password: '',
      confirmPassword: '',
    },
  })

  const otpForm = useForm<z.infer<typeof otpSchema>>({
    resolver: zodResolver(otpSchema),
    defaultValues: {
      otp: '',
    },
  })

  const isPasswordLoading = passwordForm.formState.isSubmitting
  const isOtpLoading = otpForm.formState.isSubmitting

  useEffect(() => {
    if (error) {
      setShowError(true)
    }
  }, [error])

  const onPasswordSubmit = async (values: z.infer<typeof passwordSchema>) => {
    try {
      setError(null)
      
      const token = searchParams.get('token')
      
      if (!token) {
        setError('Invalid reset link. Please request a new password reset.')
        return
      }

      // Send OTP to email
      const response = await fetch('/api/v1/auth/send-reset-otp', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: values.email,
          token,
        }),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        setError(errorData.message || 'Failed to send OTP')
        return
      }

      setUserEmail(values.email)
      passwordForm.setValue('password', values.password)
      passwordForm.setValue('confirmPassword', values.confirmPassword)
      setStep('otp')
    } catch (err: unknown) {
      console.error('Error sending OTP:', err)
      setError('An error occurred. Please try again.')
    }
  }

  const onOtpSubmit = async (values: z.infer<typeof otpSchema>) => {
    try {
      setError(null)

      const password = passwordForm.getValues('password')

      // Verify OTP and reset password
      await updatePassword({
        email: userEmail,
        password,
        otp: values.otp,
      })

      setSuccessMessage('Password reset successful! Redirecting to login...')
      
      // Clear auth state
      localStorage.removeItem('jwt')
      localStorage.removeItem('userId')
      
      setTimeout(() => {
        router.push('/login')
      }, 2000)
    } catch (err: unknown) {
      console.error('Password reset error:', err)
      
      const errorData = (err as any)?.response?.data
      
      if (errorData?.message) {
        const message = errorData.message
        if (typeof message === 'object' && message !== null) {
          const errorMessages = Object.values(message)
            .filter((v) => typeof v === 'string')
            .join(', ')
          setError(errorMessages || 'Validation error occurred')
        } else if (typeof message === 'string') {
          setError(message)
        } else {
          setError('Invalid OTP. Please try again.')
        }
      } else if (errorData?.msg) {
        setError(errorData.msg)
      } else if (err instanceof Error && err.message) {
        setError(err.message)
      } else {
        setError('Invalid OTP. Please try again.')
      }
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-background px-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-primary mb-2">Reset Password</h1>
          <p className="text-secondary">
            {step === 'password' ? 'Enter your new password' : 'Enter the OTP sent to your email'}
          </p>
        </div>

        <div className="bg-card rounded-lg shadow-lg border border-default p-8 flex flex-col gap-2">
          {showError && error && (
            <CustomAlert
              variant="destructive"
              title="Reset Failed"
              description={error}
              onClose={() => {
                setShowError(false)
                setError(null)
              }}
            />
          )}

          {successMessage && (
            <CustomAlert
              variant="default"
              title="Success"
              description={successMessage}
              onClose={() => setSuccessMessage(null)}
            />
          )}

          {step === 'password' ? (
            <Form {...passwordForm}>
              <form onSubmit={passwordForm.handleSubmit(onPasswordSubmit)} className="space-y-6">
                <FormField
                  control={passwordForm.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email *</FormLabel>
                      <FormControl>
                        <Input
                          type="email"
                          placeholder="your@email.com"
                          disabled={isPasswordLoading}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={passwordForm.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>New Password *</FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder="Enter new password"
                          disabled={isPasswordLoading}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={passwordForm.control}
                  name="confirmPassword"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Confirm Password *</FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder="Confirm new password"
                          disabled={isPasswordLoading}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <Button type="submit" disabled={isPasswordLoading} size="long" className="w-full">
                  {isPasswordLoading ? (
                    <>
                      <Loader2 className="w-5 h-5 animate-spin" />
                      <span>Sending OTP...</span>
                    </>
                  ) : (
                    <span>Continue</span>
                  )}
                </Button>
              </form>
            </Form>
          ) : (
            <Form {...otpForm}>
              <form onSubmit={otpForm.handleSubmit(onOtpSubmit)} className="space-y-6">
                <FormField
                  control={otpForm.control}
                  name="otp"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>OTP *</FormLabel>
                      <FormControl>
                        <Input
                          type="text"
                          placeholder="000000"
                          disabled={isOtpLoading}
                          maxLength={6}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <Button type="submit" disabled={isOtpLoading} size="long" className="w-full">
                  {isOtpLoading ? (
                    <>
                      <Loader2 className="w-5 h-5 animate-spin" />
                      <span>Resetting Password...</span>
                    </>
                  ) : (
                    <span>Reset Password</span>
                  )}
                </Button>

                <Button 
                  type="button" 
                  variant="ghost" 
                  size="sm" 
                  onClick={() => setStep('password')}
                  disabled={isOtpLoading}
                  className="w-full"
                >
                  Back
                </Button>
              </form>
            </Form>
          )}

          <div className="mt-6 text-center">
            <p className="text-sm text-secondary">
              Remember your password?{' '}
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

export default function ResetPasswordPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center bg-background">
          <Loader2 className="w-8 h-8 animate-spin text-primary" />
        </div>
      }
    >
      <ResetPasswordContent />
    </Suspense>
  )
}