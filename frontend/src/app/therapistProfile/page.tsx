'use client'

import AppLayout from '@/components/AppLayout'
import { Button } from "@/components/ui/button";
import {useState} from "react";
import {useAuth} from "@/hooks/useAuth";
import {useRouter} from "next/navigation";
import {
    Dialog,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle
} from "@/components/ui/dialog";
import {ConfirmDialog} from "@/components/ui/confirm-dialog";
import { validatePassword } from '@/lib/validatePassword'

export default function TherapistProfile() {
    const {updatePassword, deleteAccount} = useAuth();

    const [openPassword, setOpenPassword] = useState(false)
    const [currentPassword, setCurrentPassword] = useState("")
    const [newPassword, setNewPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")
    const [passwordError, setPasswordError] = useState("")

    const [openDeleteConfirm, setOpenDeleteConfirm] = useState(false)

    const router = useRouter()

    const handlePasswordSave = () => {
        setPasswordError("")

        if (!currentPassword || !newPassword || !confirmPassword) {
            setPasswordError('All fields are required')
            return
        }
        if (newPassword !== confirmPassword) {
            setPasswordError('New passwords do not match.')
            return
        }
        if (validatePassword(newPassword)) {
            setPasswordError('Password must include at least one special character (!@#$%^&*()_+-=[]{};:\'",.<>?/~`|)')
        }

        try {
            updatePassword({password: newPassword})

            setOpenPassword(false)
            setCurrentPassword("")
            setNewPassword("")
            setConfirmPassword("")
        } catch {
            setPasswordError('Failed to update password')
        }
    }

    const handleFinalDelete = async () => {
        try {
            const userId = localStorage.getItem("userId")

            if (!userId) {
                console.error("User ID Missing")
                return;
            }

            await deleteAccount(userId)

            setOpenDeleteConfirm(false)
            setOpenDeleteConfirm(false)

            router.push("/login")
        } catch {
            const message = 'Failed to delete account'
            console.error(message)
        }
    }

    return (
        <AppLayout>
            <div className="p-10">
                <h1 className="text-3xl font-bold">Therapist Profile</h1>
            </div>

            <div className="p-10">
                <Button onClick={() => setOpenPassword(true)}>Update Password</Button>
            </div>

            <Dialog open={openPassword} onOpenChange={setOpenPassword}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Update Password</DialogTitle>
                    </DialogHeader>

                    <div className="space-y-4 mt-2">
                        <div>
                            <label className="text-sm font-medium">Current Password</label>
                            <input
                                type="password"
                                className="w-[20%] p-2 border rounded-md mt-1"
                                value={currentPassword}
                                onChange={(e) => setCurrentPassword(e.target.value)}
                            />
                        </div>

                        <div>
                            <label className="text-sm font-medium">New Password</label>
                            <input
                                type="password"
                                className="w-[20%] p-2 border rounded-md mt-1"
                                value={newPassword}
                                onChange={(e) => setNewPassword(e.target.value)}
                            />
                        </div>

                        <div>
                            <label className="text-sm font-medium">Confirm New Password</label>
                            <input
                                type="password"
                                className="w-[20%] p-2 border rounded-md mt-1"
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                            />
                        </div>

                        {passwordError && (
                            <p className="text-red-500 text-sm">{passwordError}</p>
                        )}
                    </div>

                    <DialogFooter>
                        <Button onClick={() => setOpenPassword(false)}>
                            Cancel
                        </Button>
                        <Button onClick={handlePasswordSave}>Save</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <div className="p-10">
                <Button onClick={() => setOpenDeleteConfirm(true)}>
                    Delete Account
                </Button>
            </div>

            <ConfirmDialog
                isOpen={openDeleteConfirm}
                onClose={() => setOpenDeleteConfirm(false)}
                onConfirm={handleFinalDelete}
                title="Delete Account"
                description="This action is permanent and cannot be undone."
                confirmText="Delete"
                cancelText="Cancel"
                variant="danger"
            />
        </AppLayout>
    )
}