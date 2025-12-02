"use client";

import AppLayout from "@/components/AppLayout";
import ComputerGirl from "@/components/ui/computer-girl";
import { useAuthContext } from "@/contexts/authContext";
import { useTherapist } from "@/hooks/useTherapists";
import { Edit2, ExternalLink } from "lucide-react";
import React from "react";

const AdminProfile: React.FC = () => {
  const { userId: therapistId } = useAuthContext();
  const { therapist, error } = useTherapist(therapistId);

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
            <button className="absolute top-6 right-6 flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors">
              <span className="text-sm">Edit Profile</span>
              <Edit2 size={18} />
            </button>

            <div className="flex flex-col items-center mb-8">
              <div className="w-32 h-32 rounded-full bg-gradient-to-br from-blue-600 to-grey-600 flex items-center justify-center mb-2">
                {/*I AM NOT SURE HOW WE ARE GOING TO HAVE IMAGES ON FILE FOR THE THERAPIST, SO I USED THIS HACK*/}
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
                <span className="inline-block px-4 py-2 bg-indigo-100 text-indigo-800 rounded-full text-sm font-medium">
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
              href="ADD EDPLAN HERE"
              className="inline-flex items-center gap-2 text-black hover:text-grey-600 transition-colors"
            >
              <span>EdPlan</span>
              <ExternalLink size={16} />
            </a>
          </div>

          {/* Computer Girl on the side */}
          <div className="fixed bottom-4 right-5 pointer-events-none">
            <ComputerGirl />
          </div>
        </div>
      </div>
    </AppLayout>
  );
};

export default AdminProfile;
