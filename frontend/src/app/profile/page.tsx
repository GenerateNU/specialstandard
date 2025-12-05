"use client";

import EditModal from "@/app/profile/EditModal";
import AppLayout from "@/components/AppLayout";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { useAuthContext } from "@/contexts/authContext";
import { useAuth } from "@/hooks/useAuth";
import { useTherapist, useTherapists } from "@/hooks/useTherapists";
import { validatePassword } from "@/lib/validatePassword";
import { Edit2, ExternalLink, Settings, X } from "lucide-react";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import Image from "next/image"; 

const AdminProfile: React.FC = () => {
  const { userId: therapistId } = useAuthContext();
  const { therapist, error, refetch } = useTherapist(therapistId);
  const { updateTherapist } = useTherapists({ fetchOnMount: false });
  const { updatePassword, deleteAccount } = useAuth();
  const router = useRouter();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [openDeleteConfirm, setOpenDeleteConfirm] = useState(false);

  const handleSave = async (
    therapistId: string,
    data: {
      first_name: string;
      last_name: string;
      school_ids: number[];
      district_id: number;
    }
  ) => {
    try {
      await updateTherapist(therapistId, {
        first_name: data.first_name,
        last_name: data.last_name,
        district_id: data.district_id,
        schools: data.school_ids,
      });

      // Refresh therapist data
      await refetch();
      setIsModalOpen(false);
    } catch (error) {
      console.error("Failed to update therapist:", error);
    }
  };

  const handlePasswordSave = () => {
    setPasswordError("");

    if (!currentPassword || !newPassword || !confirmPassword) {
      setPasswordError("All fields are required");
      return;
    }
    if (newPassword !== confirmPassword) {
      setPasswordError("New passwords do not match.");
      return;
    }
    if (validatePassword(newPassword)) {
      setPasswordError(
        "Password must include at least one special character (!@#$%^&*()_+-=[]{};:'\",.<>?/~`|)"
      );
      return;
    }

    try {
      updatePassword({ password: newPassword });

      setIsSettingsOpen(false);
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
      setPasswordError("");
    } catch {
      setPasswordError("Failed to update password");
    }
  };

  const handleFinalDelete = async () => {
    try {
      const userId = localStorage.getItem("userId");

      if (!userId) {
        console.error("User ID Missing");
        return;
      }

      await deleteAccount(userId);

      setOpenDeleteConfirm(false);
      setIsSettingsOpen(false);

      router.push("/login");
    } catch {
      const message = "Failed to delete account";
      console.error(message);
    }
  };

  if (error || !therapist) {
    return (
      <AppLayout>
        <div className="min-h-screen bg-gradient-to-br from-amber-50 to-orange-50 p-8">
          <div className="max-w-4xl">
            <div className="text-center py-12">
              <p className="text-black">{error || "Loading Data..."}</p>
            </div>
          </div>
        </div>
      </AppLayout>
    );
  }

  const fullName = `${therapist.first_name} ${therapist.last_name}`;
  const schoolNames = therapist.school_names ?? [];
  const districtName = therapist.district_name ?? "No district assigned";

  return (
    <AppLayout>
      <div className="min-h-screen bg-gradient-to-br from-amber-50 to-orange-50 p-8">
        <div className="max-w-4xl">
          <h1 className="text-3xl font-bold mb-8 text-black">My Profile</h1>

          {/* Profile Card */}
          <div className="bg-white rounded-2xl shadow-sm p-8 mb-6 relative">
            <div className="absolute top-6 right-6 flex items-center gap-3">
              <button
                onClick={() => setIsSettingsOpen(true)}
                className="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors cursor-pointer"
              >
                <span className="text-sm">Settings</span>
                <Settings size={18} />
              </button>
              <button
                onClick={() => setIsModalOpen(true)}
                className="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors cursor-pointer"
              >
                <span className="text-sm">Edit Profile</span>
                <Edit2 size={18} />
              </button>
            </div>

            <div className="flex flex-col items-center mb-8">
              <div className="w-32 h-32 rounded-full bg-gradient-to-br from-blue-600 to-grey-600 flex items-center justify-center mb-2">
                <span className="text-4xl font-bold text-white">
                  {therapist.first_name[0]}
                  {therapist.last_name[0]}
                </span>
              </div>
              <h2 className="text-2xl font-semibold text-black mb-1">
                {fullName}
              </h2>
              <p className="text-gray-600">Speech-Language Pathologist</p>
            </div>

            <div className="space-y-6">
              <div>
                <h3 className="text-sm font-medium text-black mb-3">
                  School District
                </h3>
                <span className="inline-block px-4 py-2 bg-indigo-100 text-gray-600 rounded-full text-sm font-medium">
                  {districtName}
                </span>
              </div>

              {schoolNames.length > 0 && (
                <div>
                  <h3 className="text-sm font-medium text-black mb-3">
                    Schools
                  </h3>
                  <div className="flex flex-wrap gap-2">
                    {schoolNames.map((school, index) => (
                      <span
                        key={index}
                        className="px-4 py-2 bg-gray-100 text-gray-700 rounded-full text-sm border border-gray-300"
                      >
                        {school}
                      </span>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Resources Card */}
          <div className="bg-white rounded-2xl shadow-sm p-6">
            <h3 className="text-lg font-semibold text-black mb-3">Resources</h3>
            <a
              href="https://www.edplan.com/blog/Account/login.aspx?ReturnURL=/blog/admin/"
              className="inline-flex items-center gap-2 text-black hover:text-grey-600 transition-colors cursor-pointer"
            >
              <span>EdPlan</span>
              <ExternalLink size={16} />
            </a>
          </div>

          {/* Doodleman Image on the side */}
          <div className="fixed bottom-4 right-5 pointer-events-none">
            <Image 
              src="/doodleman.png" 
              alt="Doodleman" 
              width={744} 
              height={541}
              className="w-auto h-[100px] sm:h-[200px] md:h-[300px] lg:h-[400px]"
            />
          </div>
        </div>
      </div>

      {/* Edit Profile Modal */}
      {therapistId && (
        <EditModal
          isOpen={isModalOpen}
          onClose={() => setIsModalOpen(false)}
          therapistId={therapistId}
          initialData={{
            first_name: therapist.first_name,
            last_name: therapist.last_name,
            school_names: therapist.school_names ?? [],
            district_id: therapist.district_id,
          }}
          onSave={handleSave}
        />
      )}

      {/* Settings Modal */}
      {isSettingsOpen && (
        <div className="fixed inset-0 backdrop-blur-sm bg-black/30 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-2xl shadow-xl max-w-md w-full p-6">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-bold text-black">Settings</h2>
              <button
                onClick={() => {
                  setIsSettingsOpen(false);
                  setPasswordError("");
                  setCurrentPassword("");
                  setNewPassword("");
                  setConfirmPassword("");
                }}
                className="text-gray-400 hover:text-gray-600 transition-colors cursor-pointer"
              >
                <X size={24} />
              </button>
            </div>

            <div className="space-y-6">
              {/* Password Update Section */}
              <div className="border-b border-gray-200 pb-6">
                <h3 className="text-lg font-semibold text-black mb-4">
                  Update Password
                </h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Current Password
                    </label>
                    <input
                      type="password"
                      value={currentPassword}
                      onChange={(e) => setCurrentPassword(e.target.value)}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-gray-500 focus:border-transparent text-gray-900"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      New Password
                    </label>
                    <input
                      type="password"
                      value={newPassword}
                      onChange={(e) => setNewPassword(e.target.value)}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-gray-500 focus:border-transparent text-gray-900"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Confirm New Password
                    </label>
                    <input
                      type="password"
                      value={confirmPassword}
                      onChange={(e) => setConfirmPassword(e.target.value)}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-gray-500 focus:border-transparent text-gray-900"
                    />
                  </div>

                  {passwordError && (
                    <p className="text-red-500 text-sm">{passwordError}</p>
                  )}

                  <button
                    onClick={handlePasswordSave}
                    className="w-full px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors cursor-pointer"
                  >
                    Update Password
                  </button>
                </div>
              </div>

              {/* Delete Account Section */}
              <div>
                <h3 className="text-lg font-semibold text-black mb-2">
                  Danger Zone
                </h3>
                <p className="text-sm text-gray-600 mb-4">
                  Deleting your account is permanent and cannot be undone.
                </p>
                <button
                  onClick={() => setOpenDeleteConfirm(true)}
                  className="w-full px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors cursor-pointer"
                >
                  Delete Account
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Delete Account Confirmation */}
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
  );
};

export default AdminProfile;
