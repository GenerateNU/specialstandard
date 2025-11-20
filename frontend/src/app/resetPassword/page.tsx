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

const formSchema = z.object({
  password: z.string().min(6, 'Password must be at least 6 characters long'),
  confirmPassword: z.string().min(6, 'Password must be at least 6 characters long'),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
})

function ResetPasswordContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { updatePassword } = useAuth()
  const [error, setError] = useState<string | null>(null)
  const [showError, setShowError] = useState(false)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      password: '',
      confirmPassword: '',
    },
  })

  const isLoading = form.formState.isSubmitting

  useEffect(() => {
    if (error) {
      setShowError(true)
    }
  }, [error])

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      setError(null)
      
      // Get the reset token from URL params (set by Supabase email)
      const token = searchParams.get('token')

      console.warn('Reset token:', token)
      
      if (!token) {
        setError('Invalid reset link. Please request a new password reset.')
        return
      }

      // Call your backend API to update password with reset token
      await updatePassword({
        password: values.password,
        token, // Pass the reset token from the email link
      })

      setSuccessMessage('Password reset successful! Redirecting to login...')
      
      // Clear auth state and redirect
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
          setError('An error occurred during password reset. Please try again.')
        }
      } else if (errorData?.msg) {
        setError(errorData.msg)
      } else if (err instanceof Error && err.message) {
        setError(err.message)
      } else {
        setError('An error occurred during password reset. Please try again.')
      }
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-background px-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-primary mb-2">Reset Password</h1>
          <p className="text-secondary">Enter your new password below</p>
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

          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>New Password *</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="Enter new password"
                        disabled={isLoading}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="confirmPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Confirm Password *</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="Confirm new password"
                        disabled={isLoading}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" disabled={isLoading} size="long" className="w-full">
                {isLoading ? (
                  <>
                    <Loader2 className="w-5 h-5 animate-spin" />
                    <span>Resetting Password...</span>
                  </>
                ) : (
                  <span>Reset Password</span>
                )}
              </Button>
            </form>
          </Form>

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