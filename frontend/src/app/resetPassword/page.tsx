'use client'

import { Loader2, LockKeyhole } from 'lucide-react'
import Image from 'next/image'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import CustomAlert from '@/components/ui/CustomAlert'
import { Input } from '@/components/ui/input'
import { useAuthContext } from '@/contexts/authContext'
import { useAuth } from '@/hooks/useAuth'
import { validatePassword } from '@/lib/validatePassword'

export const dynamic = 'force-dynamic'

export default function ResetPasswordPage() {
    const [password, setPassword] = useState('')
    const [confirmPassword, setConfirmPassword] = useState('')
    const [error, setError] = useState<string | null>(null)
    const [success, setSuccess] = useState<string | null>(null)
    const [showError, setShowError] = useState(false)
    const [showSuccess, setShowSuccess] = useState(false)
    const [token, setToken] = useState<string | null>(null)
    const [tokenError, setTokenError] = useState<string | null>(null)
    const [isValidating, setIsValidating] = useState(true)

    const { isAuthenticated } = useAuthContext()
    const { updatePassword, updatePasswordMutation } = useAuth()
    const router = useRouter()

    // Extract token from URL hash on component mount
    useEffect(() => {
        if (typeof window === 'undefined') return

        setIsValidating(true)

        const hash = window.location.hash.replace('#', '')
        const params = new URLSearchParams(hash)

        const accessToken = params.get('access_token')
        const type = params.get('type')

        console.warn('Extracted from hash:', {
            hasAccessToken: !!accessToken,
            type,
        })

        // Validate token and type
        if (type !== 'recovery') {
            setTokenError(
                'Invalid reset link. This does not appear to be a password recovery link.'
            )
            setIsValidating(false)
            return
        }

        if (!accessToken) {
            setTokenError(
                'Invalid reset link. Missing authentication token. Please request a new password reset link.'
            )
            setIsValidating(false)
            return
        }

        // Token is valid
        setToken(accessToken)
        setIsValidating(false)
    }, [])

    useEffect(() => {
        if (!updatePasswordMutation.isPending && isAuthenticated) {
            router.push('/')
        }
    }, [isAuthenticated, updatePasswordMutation.isPending, router])

    useEffect(() => {
        if (error) {
            setShowError(true)
        }
    }, [error])

    useEffect(() => {
        if (success) {
            setShowSuccess(true)
        }
    }, [success])

    if (isValidating) {
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
        setSuccess(null)

        if (!token) {
            setError('Reset token is missing. Please request a new password reset link.')
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

        try {
            // Call the API with both password and token
            await updatePassword({
                password,
                token,
            })

            setSuccess(
                'Password has been reset successfully! You will be redirected to login shortly.'
            )
            setPassword('')
            setConfirmPassword('')

            // Redirect to login after 2 seconds
            setTimeout(() => {
                router.push('/login')
            }, 2000)
        } catch (err: unknown) {
            console.error('Reset error:', err)

            const errorData = (err as any)?.response?.data

            if (errorData?.message) {
                const message = errorData.message
                if (typeof message === 'object' && message !== null) {
                    const errorMessages = Object.values(message)
                        .filter(v => typeof v === 'string')
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
                    <Image
                        src="/tss.png"
                        alt="The Special Standard logo"
                        width={180}
                        height={38}
                        className="mx-auto mb-6"
                        priority
                    />
                    <h1 className="text-3xl font-bold text-primary mb-2">Reset Password</h1>
                    <p className="text-secondary">Enter your new password below</p>
                </div>

                <div className="bg-card rounded-lg shadow-lg border border-default p-8 flex flex-col gap-2">
                    {tokenError && (
                        <CustomAlert
                            variant="destructive"
                            title="Invalid Reset Link"
                            description={tokenError}
                            onClose={() => setTokenError(null)}
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

                    {showSuccess && success && (
                        <CustomAlert
                            variant="default"
                            title="Success"
                            description={success}
                            onClose={() => {
                                setShowSuccess(false)
                                setSuccess(null)
                            }}
                        />
                    )}

                    {token && !tokenError && (
                        <form onSubmit={handleSubmit} className="space-y-6">
                            <div>
                                <label
                                    htmlFor="password"
                                    className="block text-sm font-medium text-primary mb-2"
                                >
                                    New Password *
                                </label>
                                <Input
                                    id="password"
                                    type="password"
                                    value={password}
                                    onChange={e => setPassword(e.target.value)}
                                    required
                                    disabled={updatePasswordMutation.isPending}
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
                                    disabled={updatePasswordMutation.isPending}
                                    placeholder="••••••••"
                                />
                            </div>

                            <Button
                                type="submit"
                                disabled={updatePasswordMutation.isPending}
                                size="long"
                            >
                                {updatePasswordMutation.isPending ? (
                                    <>
                                        <Loader2 className="w-5 h-5 animate-spin" />
                                        <span>Resetting Password...</span>
                                    </>
                                ) : (
                                    <>
                                        <LockKeyhole className="w-5 h-5" />
                                        <span>Reset Password</span>
                                    </>
                                )}
                            </Button>
                        </form>
                    )}

                    <div className="mt-6 text-center">
                        <p className="text-sm text-secondary">
                            Remember your password?
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