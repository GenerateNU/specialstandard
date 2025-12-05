"use client";

import StudentCard from "@/components/students/studentCard";
import { Button } from "@/components/ui/button";
import { useStudents } from "@/hooks/useStudents";
import { ArrowLeft, Loader2, Plus } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function StudentsPage() {
  const router = useRouter();
  const { students, isLoading, refetch } = useStudents();
  const [therapistId, setTherapistId] = useState<string>("");

  useEffect(() => {
    const userId =
      localStorage.getItem("temp_userId") || localStorage.getItem("userId"); // Check temp_userId first
    if (userId) {
      setTherapistId(userId);
    } else {
      router.push("/signup/welcome");
    }

    refetch();
  }, [router, refetch]);

  const handleBack = () => {
    router.back();
  };

  const handleAddStudent = () => {
    router.push("/signup/students/add");
  };

  const handleContinue = () => {
    router.push("/signup/sessions/add");
  };

  // Filter students for this therapist
  const therapistStudents = students.filter(
    (student) => student.therapist_id === therapistId
  );

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen p-8">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center min-h-screen p-8">
      <div className="max-w-2xl w-full">
        <button
          onClick={handleBack}
          className="mb-6 flex items-center cursor-pointer text-secondary hover:text-primary transform hover:scale-125 transition-colors"
        >
          <ArrowLeft className="w-4 h-4 mr-1" />
          Back
        </button>

        <div className="flex items-center justify-between mb-8">
          <h1 className="text-3xl font-bold text-primary">Your Students</h1>
        </div>

        {therapistStudents.length === 0 ? (
          <div className="bg-card rounded-lg border border-default p-12 text-center mb-6">
            <p className="text-secondary mb-6">
              No students added yet. Add your first student to get started.
            </p>
            <Button
              onClick={handleAddStudent}
              className="text-white"
            >
              <Plus className="w-4 h-4 mr-2" />
              Add Your First Student
            </Button>
          </div>
        ) : (
          <>
            <div className="space-y-3 mb-6">
              {therapistStudents.map((student) => (
                <StudentCard key={student.id} student={student} />
              ))}
            </div>

            <div className="bg-accent-light rounded-lg p-4 mb-6">
              <p className="text-sm text-primary">
                <strong>{therapistStudents.length}</strong> student
                {therapistStudents.length !== 1 ? "s" : ""} added
              </p>
            </div>
          </>
        )}

        <div className="flex gap-3">
          <Button
            onClick={handleAddStudent}
            variant="outline"
            className="flex-1"
          >
            Add Student
          </Button>

          <Button onClick={handleContinue} className="flex-1 text-white">
            Continue
          </Button>
        </div>

        {therapistStudents.length === 0 && (
          <div className="text-center mt-4">
            <button
              onClick={handleContinue}
              className="text-sm text-secondary cursor-pointer hover:text-primary underline"
            >
              Skip this step
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
