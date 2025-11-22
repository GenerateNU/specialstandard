// Client-Side Password Validation
export function validatePassword(pwd: string): string | null {
    if (pwd.length < 8) {
        return 'Password must be at least 8 characters long'
    }
    if (!/[A-Z]/.test(pwd)) {
        return 'Password must include at least one uppercase letter'
    }
    if (!/[a-z]/.test(pwd)) {
        return 'Password must include at least one lowercase letter'
    }
    if (!/\d/.test(pwd)) {
        return 'Password must include at least one digit'
    }
    if (!/[!@#$%^&*()_+\-=[\]{};':"\\|,.<>?/~`]/.test(pwd)) {
        return 'Password must include at least one special character (!@#$%^&*()_+-=[]{};:\'",.<>?/~`|)'
    }
    return null
}