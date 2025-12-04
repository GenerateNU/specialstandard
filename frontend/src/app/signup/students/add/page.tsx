"use client";

import { Button } from "@/components/ui/button";
import CustomAlert from "@/components/ui/CustomAlert";
import { Dropdown } from "@/components/ui/dropdown";
import { Input } from "@/components/ui/input";
import { useSchools } from "@/hooks/useSchools";
import { useStudents } from "@/hooks/useStudents";
import { ArrowLeft, Loader2, Trash2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

// Grade options
const gradeOptions = [
  { label: "Kindergarten", value: "0" },
  { label: "Grade 1", value: "1" },
  { label: "Grade 2", value: "2" },
  { label: "Grade 3", value: "3" },
  { label: "Grade 4", value: "4" },
  { label: "Grade 5", value: "5" },
  { label: "Grade 6", value: "6" },
  { label: "Grade 7", value: "7" },
  { label: "Grade 8", value: "8" },
  { label: "Grade 9", value: "9" },
  { label: "Grade 10", value: "10" },
  { label: "Grade 11", value: "11" },
  { label: "Grade 12", value: "12" },
  { label: "Graduated", value: "-1" },
];

export default function AddStudentsPage() {
  const router = useRouter();
  const { addStudent } = useStudents();
  const { schools } = useSchools();

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showError, setShowError] = useState(false);

  // Form fields
  const [formData, setFormData] = useState({
    firstName: "",
    lastName: "",
    school: "",
    grade: "",
    dob: "",
  });

  // IEP goals
  const [iepGoals, setIepGoals] = useState<string[]>([""]);

  // Load therapist data from localStorage
  const [therapistId, setTherapistId] = useState<string>("");
  const [therapistProfile, setTherapistProfile] = useState<any>(null);

  useEffect(() => {
    const userId =
      localStorage.getItem("temp_userId") || localStorage.getItem("userId");
    const profileData = localStorage.getItem("therapistProfile");

    if (userId) {
      setTherapistId(userId);
    }

    if (profileData) {
      const profile = JSON.parse(profileData);
      setTherapistProfile(profile);

      // Pre-select school if therapist only has one school
      if (profile.schoolIds && profile.schoolIds.length === 1) {
        setFormData((prev) => ({ ...prev, school: profile.schoolIds[0] }));
      }
    }

    if (!userId) {
      router.push("/signup/welcome");
    }
  }, [router]);

  // Filter schools based on therapist's schools
  const availableSchools = therapistProfile?.schoolIds
    ? schools.filter((school) =>
        therapistProfile.schoolIds.includes(school.id!.toString())
      )
    : schools;

  const schoolOptions = availableSchools.map((school) => ({
    label: school.name,
    value: school.id!.toString(),
  }));

  const handleBack = () => {
    router.push("/signup/link");
  };

  const handleAddGoal = () => {
    setIepGoals([...iepGoals, ""]);
  };

  const handleRemoveGoal = (index: number) => {
    setIepGoals(iepGoals.filter((_, i) => i !== index));
  };

  const handleGoalChange = (index: number, value: string) => {
    const updatedGoals = [...iepGoals];
    updatedGoals[index] = value;
    setIepGoals(updatedGoals);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.firstName || !formData.lastName) {
      setError("Please enter student name");
      setShowError(true);
      return;
    }

    if (!formData.school) {
      setError("Please select a school");
      setShowError(true);
      return;
    }

    if (!therapistId) {
      setError("Therapist ID not found. Please restart the signup process.");
      setShowError(true);
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      // Filter out empty IEP goals
      const validGoals = iepGoals.filter((goal) => goal.trim() !== "");

      await addStudent({
        first_name: formData.firstName,
        last_name: formData.lastName,
        school_id: Number(formData.school),
        therapist_id: therapistId,
        grade: formData.grade ? Number(formData.grade) : null,
        dob: formData.dob || null,
        iep: validGoals.length > 0 ? validGoals : null,
      });

      // Save student data to localStorage for summary
      const studentsData = localStorage.getItem("onboardingStudents");
      const existingStudents = studentsData ? JSON.parse(studentsData) : [];
      existingStudents.push({
        name: `${formData.firstName} ${formData.lastName}`,
        school: availableSchools.find(
          (s) => s.id!.toString() === formData.school
        )?.name,
        grade: formData.grade,
        dob: formData.dob,
      });
      localStorage.setItem(
        "onboardingStudents",
        JSON.stringify(existingStudents)
      );
      router.push("/signup/students");
    } catch (err: any) {
      console.error("Add student error:", err);
      setError(err?.message || "Failed to add student. Please try again.");
      setShowError(true);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen p-8">
      <div className="max-w-md w-full">
        <button
          onClick={handleBack}
          className="mb-6 flex items-center text-secondary hover:text-primary cursor-pointer transition-colors"
        >
          <ArrowLeft className="w-4 h-4 mr-1" />
          Back
        </button>

        <h1 className="text-3xl font-bold text-primary mb-8">Add Students</h1>

        {showError && error && (
          <div className="mb-4">
            <CustomAlert
              variant="destructive"
              title="Error"
              description={error}
              onClose={() => {
                setShowError(false);
                setError(null);
              }}
            />
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Student Name */}
          <div className="flex gap-4">
            <div className="flex-1">
              <p className="text-xs text-secondary mb-1">
                First Name (Initial Only)
              </p>
              <Input
                value={formData.firstName}
                onChange={(e) =>
                  setFormData({ ...formData, firstName: e.target.value })
                }
                placeholder="Enter student's first initial"
                className="border border-gray-300"
                required
                disabled={isSubmitting}
              />
            </div>

            <div className="flex-1">
              <p className="text-xs text-secondary mb-1">
                Last Name (Initial Only)
              </p>
              <Input
                value={formData.lastName}
                onChange={(e) =>
                  setFormData({ ...formData, lastName: e.target.value })
                }
                placeholder="Enter student's last initial"
                className="border border-gray-300"
                required
                disabled={isSubmitting}
              />
            </div>
          </div>
          {/* School Selection */}
          <div>
            <p className="text-xs text-secondary mb-1">Student's School</p>
            <Dropdown
              items={schoolOptions.map((opt) => ({
                label: opt.label!,
                value: opt.value,
                onClick: () => setFormData({ ...formData, school: opt.value }),
              }))}
              value={formData.school}
              placeholder="Select school"
              align="left"
              className="w-md border border-gray-300"
            />
          </div>

          {/* Grade Selection */}
          <div>
            <Dropdown
              items={gradeOptions.map((opt) => ({
                label: opt.label,
                value: opt.value,
                onClick: () => setFormData({ ...formData, grade: opt.value }),
              }))}
              value={formData.grade}
              placeholder="Select grade"
              className="w-md border-gray-300"
            />
          </div>

          {/* Date of Birth */}
          <div>
            <p className="text-xs text-secondary mt-1">
              Date of birth (optional)
            </p>
            <Input
              type="date"
              value={formData.dob}
              onChange={(e) =>
                setFormData({ ...formData, dob: e.target.value })
              }
              className="border border-gray-300"
              disabled={isSubmitting}
            />
          </div>

          {/* IEP Goals */}
          <div>
            <label className="block text-sm font-medium text-primary mb-2">
              IEP Goals (optional)
            </label>
            {iepGoals.map((goal, index) => (
              <div key={index} className="flex gap-2 mb-2">
                <Input
                  value={goal}
                  onChange={(e) => handleGoalChange(index, e.target.value)}
                  placeholder={`Enter student's goals`}
                  className="border border-gray-300"
                  disabled={isSubmitting}
                />
                {iepGoals.length > 1 && (
                  <button
                    type="button"
                    onClick={() => handleRemoveGoal(index)}
                    className="text-red-500 hover:text-red-700"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                )}
              </div>
            ))}
            <button
              type="button"
              onClick={handleAddGoal}
              className="text-xs hover:text-accent-hover underline"
            >
              + Add goal
            </button>
          </div>

          <div className="pt-4">
            <Button
              type="submit"
              size="long"
              className="w-full text-white"
              disabled={isSubmitting}
            >
              {isSubmitting ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin mr-2" />
                  <span>Adding student...</span>
                </>
              ) : (
                <span>Save</span>
              )}
            </Button>
          </div>

          <div className="text-center">
            <button
              type="button"
              onClick={() => router.push("/signup/sessions/add")}
              className="text-sm clock cursor-pointer text-secondary hover:text-primary underline"
            >
              Skip this step
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
