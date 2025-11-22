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
import {validatePassword} from "@/app/therapistProfile/page"
import {createClient} from '@supabase/supabase-js'

export default function ResetPasswordPage() {
    const [password, setPassword] = useState('')
    const [confirmPassword, setConfirmPassword] = useState('')
    const [error, setError] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(false)
    const [showError, setShowError] = useState(false)

    const { isAuthenticated } = useAuthContext()
    const router = useRouter()

    const supabase = createClient(
        process.env.NEXT_PUBLIC_SUPABASE_URL!,
        process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
        {
            auth: {
                persistSession: true,
                autoRefreshToken: true,
            },
        }
    )

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
            const hash = window.location.hash.replace('#', '')
            const params = new URLSearchParams(hash)

            const accessToken = params.get('access_token')
            const refreshToken = params.get('refresh_token')

            if (accessToken && refreshToken) {
                const { error } = await supabase.auth.setSession({
                    access_token: accessToken,
                    refresh_token: refreshToken,
                })

                if (error) {
                    console.warn("Failed to restore session:", error)
                    setError("Invalid or expired password-reset link.")
                } else {
                    console.warn("Session restored successfully!")
                }
            }

            await supabase.auth.getSession()
        }

        restoreSession()
    }, [])

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
            const { error } = await supabase.auth.updateUser({ password })

            if (error) {
                console.error(error)
                setError("Password unable to be reset successfully")
            } else {
                setTimeout(() => router.push("/login"), 2000)
            }
        }
        catch (err: unknown) {
            console.error('Reset error:', err)

            // Type guard for axios error
            const errorData = (err as any)?.response?.data

            // Handle various error response formats
            if (errorData?.message) {
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
                    setError('An error occurred during reset password. Please try again.')
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
                    <h1 className="text-3xl font-bold text-primary mb-2">Reset Password</h1>
                    <p className="text-secondary">Type in your new and re-confirm your new password</p>
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