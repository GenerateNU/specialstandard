import AppLayout from "@/components/AppLayout";
import { Edit2, ExternalLink } from "lucide-react";
import React from "react";

interface ProfileData {
  name: string;
  title: string;
  avatar: string;
  schoolDistrict: string;
  schools: string[];
}

const AdminProfile: React.FC = () => {
  // I DONT FEEL LIKE PULLING YET
  const profile: ProfileData = {
    name: "Johnny Doe",
    title: "Speech-Language Pathologist",
    avatar:
      "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&h=400&fit=crop",
    schoolDistrict: "Boston Public Schools",
    schools: [
      "Brighton High School",
      "Forrest Hills Middle School",
      "Boston Latin Academy",
    ],
  };

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
              <img
                src={profile.avatar}
                alt={profile.name}
                className="w-32 h-32 rounded-full object-cover mb-4"
              />
              <h2 className="text-2xl font-semibold text-black mb-1">
                {profile.name}
              </h2>
              <p className="text-black">{profile.title}</p>
            </div>

            <div className="space-y-6">
              <div>
                <h3 className="text-sm font-medium text-black mb-3">
                  School District
                </h3>
                <span className="inline-block px-4 py-2 bg-indigo-100 text-indigo-800 rounded-full text-sm font-medium">
                  {profile.schoolDistrict}
                </span>
              </div>

              <div>
                <h3 className="text-sm font-medium text-black mb-3">Schools</h3>
                <div className="flex flex-wrap gap-2">
                  {profile.schools.map((school, index) => (
                    <span
                      key={index}
                      className="px-4 py-2 bg-gray-100 text-gray-700 rounded-full text-sm border border-gray-300"
                    >
                      {school}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </div>

          {/* Resources Card */}
          <div className="bg-white rounded-2xl shadow-sm p-6">
            <h3 className="text-lg font-semibold text-black mb-3">Resources</h3>
            <a
              href="ADD EDPLAN HERE"
              className="inline-flex items-center gap-2 text-black hover:text-grey-400 transition-colors"
            >
              <span>EdPlan</span>
              <ExternalLink size={16} />
            </a>
          </div>
        </div>
      </div>
    </AppLayout>
  );
};

export default AdminProfile;
