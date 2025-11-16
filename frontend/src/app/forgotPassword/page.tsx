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
import {useAuth} from "@/hooks/useAuth";

export default function ForgotPasswordPage() {
    const [email, setEmail] = useState('')
    const [error, setError] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(false)
    const [showError, setShowError] = useState(false)

    const { isAuthenticated } = useAuthContext()
    const { forgotPassword } = useAuth()
    const router = useRouter()

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
        setIsLoading(true)

        try {
            await forgotPassword({ email })
            // AuthContext will handle redirect to /

            router.push("/login")
        }
        catch (err: unknown) {
            console.error('Forgot Password error:', err)

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
                    setError('An error occurred during resetting password. Please try again.')
                }
            }
            else if (errorData?.msg) {
                setError(errorData.msg)
            }
            else if (err instanceof Error && err.message) {
                setError(err.message)
            }
            else {
                setError('An error occurred during resetting password. Please try again.')
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
                    <p className="text-secondary">Provide your email to get going</p>
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
                                        <span>Reset Password</span>
                                    </>
                                )}
                        </Button>
                    </form>

                    <div className="mt-6 text-center">
                        <p className="text-sm text-secondary">
                            Remembered It?
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
