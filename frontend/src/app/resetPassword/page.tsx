'use client'

import { Loader2, UserPlus } from 'lucide-react'
import Image from 'next/image'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import CustomAlert from '@/components/ui/CustomAlert'
import { Input } from '@/components/ui/input'
import { useAuthContext } from '@/contexts/authContext'
import { validatePassword } from '@/lib/validatePassword'
import {createClient} from '@supabase/supabase-js'

export const dynamic = 'force-dynamic'

export default function ResetPasswordPage() {
    const [password, setPassword] = useState('')
    const [confirmPassword, setConfirmPassword] = useState('')
    const [error, setError] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(false)
    const [showError, setShowError] = useState(false)
    const [sessionReady, setSessionReady] = useState(false)
    const [sessionError, setSessionError] = useState<string | null>(null)

    const { isAuthenticated } = useAuthContext()
    const router = useRouter()
      // eslint-disable-next-line node/prefer-global/process
    const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL;
          // eslint-disable-next-line node/prefer-global/process
    const supabaseAnonKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY;

    if (!supabaseUrl) {
      // eslint-disable-next-line node/prefer-global/process
      if (process.env.NODE_ENV === "production" && typeof window === "undefined") {
        console.warn(
          "Supabase URL not found during build. Using placeholder for static generation."
        );
      } else {
        console.error(
          "Supabase URL is required. Please set NEXT_PUBLIC_SUPABASE_URL environment variable."
        );
      }
    }

    const supabase = createClient(
      supabaseUrl || "https://placeholder-for-static-build.supabase.co",
      supabaseAnonKey || "placeholder-key-for-static-build"
    );

    // Redirect if already authenticated
    useEffect(() => {
        if (!isLoading && isAuthenticated) {
            router.push('/')
        }
    }, [isAuthenticated, isLoading, router])

    useEffect(() => {
        if (error)
            setShowError(true)
    }, [error])

    useEffect(() => {
        if (typeof window === 'undefined') return

        const restoreSession = async () => {
            try {
                const hash = window.location.hash.replace('#', '')
                const params = new URLSearchParams(hash)

                const accessToken = params.get('access_token')
                const refreshToken = params.get('refresh_token')
                const type = params.get('type')

                console.warn("Extracted tokens from URL:", { 
                    hasAccessToken: !!accessToken, 
                    hasRefreshToken: !!refreshToken,
                    type 
                })

                // Check if this is a recovery/password reset type
                if (type !== 'recovery') {
                    setSessionError("Invalid reset link. This is not a password recovery link.")
                    setSessionReady(false)
                    return
                }

                if (accessToken && refreshToken) {
                    // Set the session with the tokens from the email link
                    const { data, error: setSessionError } = await supabase.auth.setSession({
                        access_token: accessToken,
                        refresh_token: refreshToken,
                    })

                    if (setSessionError) {
                        console.error("Failed to restore session:", setSessionError)
                        setSessionError("Invalid or expired password-reset link. Please request a new one.")
                        setSessionReady(false)
                        return
                    }

                    // Verify session was actually set
                    const { data: sessionData, error: getSessionError } = await supabase.auth.getSession()
                    
                    if (getSessionError || !sessionData?.session) {
                        console.error("Session not found after setting:", getSessionError)
                        if (setSessionError) {
                            setSessionError("Failed to establish session. Please try the reset link again.")
                        }
                        setSessionReady(false)
                        return
                    }

                    console.warn("Session restored successfully!")
                    setSessionReady(true)
                } else {
                    setSessionError("Invalid reset link. Missing authentication tokens.")
                    setSessionReady(false)
                }
            } catch (err) {
                console.error("Error restoring session:", err)
                setSessionError("An error occurred while processing the reset link.")
                setSessionReady(false)
            }
        }

        restoreSession()
    }, [supabase])

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

        // Check if session is ready before proceeding
        if (!sessionReady) {
            setError("Session is not ready. Please try the reset link again.")
            return
        }

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
            const { error: updateError } = await supabase.auth.updateUser({ password })

            if (updateError) {
                console.error("Password update failed:", updateError)
                setError(updateError.message || "Password unable to be reset successfully")
            } else {
                // Clear the hash from URL after successful reset
                window.history.replaceState({}, document.title, window.location.pathname)
                setError(null)
                
                // Redirect to login after 2 seconds
                setTimeout(() => {
                    // Sign out to clear the session
                    supabase.auth.signOut().then(() => {
                        router.push("/login")
                    })
                }, 2000)
            }
        }
        catch (err: unknown) {
            console.error('Reset error:', err)

            const errorData = (err as any)?.response?.data

            if (errorData?.message) {
                const message = errorData.message
                if (typeof message === 'object' && message !== null) {
                    const errorMessages = Object.values(message).filter(v => typeof v === 'string').join(', ')
                    setError(errorMessages || 'Validation error occurred')
                }
                else if (typeof message === 'string') {
                    setError(message)
                }
                else {
                    setError('An error occurred during password reset. Please try again.')
                }
            }
            else if (errorData?.msg) {
                setError(errorData.msg)
            }
            else if (err instanceof Error && err.message) {
                setError(err.message)
            }
            else {
                setError('An error occurred during password reset. Please try again.')
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
                    <h1 className="text-3xl font-bold text-primary mb-2">Reset Password</h1>
                    <p className="text-secondary">Enter your new password</p>
                </div>

                <div className="bg-card rounded-lg shadow-lg border border-default p-8 flex flex-col gap-2">
                    {sessionError && (
                        <CustomAlert
                            variant="destructive"
                            title="Reset Link Error"
                            description={sessionError}
                            onClose={() => setSessionError(null)}
                        />
                    )}

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

                    {!sessionReady && !sessionError && (
                        <div className="flex items-center justify-center py-8">
                            <Loader2 className="w-6 h-6 animate-spin text-primary" />
                            <span className="ml-3 text-secondary">Validating reset link...</span>
                        </div>
                    )}

                    {sessionReady && (
                        <form onSubmit={handleSubmit} className="space-y-6">
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
                                            <span>Resetting Password...</span>
                                        </>
                                    )
                                    : (
                                        <>
                                            <UserPlus className="w-5 h-5" />
                                            <span>Reset Password</span>
                                        </>
                                    )}
                            </Button>
                        </form>
                    )}

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