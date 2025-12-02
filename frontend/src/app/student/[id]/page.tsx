"use client";
import {
  ArrowLeft,
  ChevronLeft,
  CirclePlus,
  PencilLine,
  RotateCw,
  Save,
  Trash,
  Trash2,
  Trophy,
  User,
  X } from "lucide-react";

import AppLayout from "@/components/AppLayout";
import { PageHeader } from "@/components/PageHeader";
import UpcomingSession from "@/components/sessions/UpcomingSession";
import Link from "next/link";
import {useParams, useRouter} from "next/navigation";
import React, { useEffect, useState } from "react";

import SchoolTag from "@/components/school/schoolTag";
import { Avatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { useRecentlyViewedStudents } from "@/hooks/useRecentlyViewedStudents";
import { useStudents } from "@/hooks/useStudents";


import CustomPieChart from "@/components/statistics/PieChart";
import { useStudentAttendance } from "@/hooks/useStudentAttendance";
import { useStudentRatings } from "@/hooks/useStudentRatings";
import { getAvatarVariant } from "@/lib/avatarUtils";
import SessionRatingsChart from "@/components/statistics/SessionRatingsChart";



function mapLevelToNumber(level: string): number {
  switch (level) {
    case "minimal": return 1;
    case "low": return 2;
    case "moderate": return 3;
    case "high": return 4;
    case "maximal": return 5;
    default: return 0;
  }
}

function StudentPage() {
  const params = useParams();
  const studentId = params.id as string;

  const { students, isLoading, updateStudent, deleteStudent } = useStudents();
  const student = students.find((s) => s.id === studentId);
  const { addRecentStudent } = useRecentlyViewedStudents();

  const { attendance, isLoading: attendanceLoading } = useStudentAttendance({
    studentId: studentId || "",
  });

  // Session ratings hook
  const { ratings } = useStudentRatings({ studentId });

  // Track this student as recently viewed - only when studentId changes
  useEffect(() => {
    if (student && studentId) {
      addRecentStudent(student);
    }
  }, [studentId]); // Only track when navigating to a new student

  const [edit, setEdit] = useState(false);
  const [iepGoals, setIepGoals] = useState<string[]>([]);
  const [isSaving, setIsSaving] = useState(false);
  const router = useRouter()

  // Initialize IEP goals from student data
  useEffect(() => {
    if (student?.iep && Array.isArray(student.iep)) {
      setIepGoals(student.iep);
    } else {
      setIepGoals([]);
    }
  }, [student?.iep]);

  const CORNER_ROUND = "rounded-4xl";
  const PADDING = "p-5";

  const fullName = student ? `${student.first_name} ${student.last_name}` : "";
  const initials = student
    ? `${student.first_name[0]}.${student.last_name[0]}.`
    : "";
  const avatarVariant = student ? getAvatarVariant(student.id) : "lorelei";

  const handleSave = async () => {
    if (!student) return;

    // Filter out empty goals before saving
    const filteredGoals = iepGoals.filter((goal) => goal.trim() !== "");

    setIsSaving(true);
    try {
      // Save IEP goals as array
      await updateStudent(student.id, { iep: filteredGoals });
      setIepGoals(filteredGoals);
      setEdit(false);
    } catch (error) {
      console.error("Failed to save IEP goals:", error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleCancel = () => {
    // Reset to original data
    if (student?.iep && Array.isArray(student.iep)) {
      setIepGoals(student.iep);
    } else {
      setIepGoals([]);
    }
    setEdit(false);
  };

  const handleDelete = async () => {
    if (
      // eslint-disable-next-line no-alert
      window.confirm(
        `Are you sure you want to delete the student: ${student?.first_name}? This action cannot be undone.`
      )
    ) {
      try {
        await deleteStudent(studentId);
        window.history.back();
      } catch (error) {
        console.error("Failed to delete student:", error);
      }
    }
  };

  const addGoal = () => {
    setIepGoals([...iepGoals, ""]);
  };

  const updateGoal = (index: number, value: string) => {
    const newGoals = [...iepGoals];
    newGoals[index] = value;
    setIepGoals(newGoals);
  };

  const deleteGoal = (index: number) => {
    setIepGoals(iepGoals.filter((_, i) => i !== index));
  };

  if (isLoading) {
    return (
      <div className="min-h-screen h-screen flex items-center justify-center bg-background">
        <div className="text-primary">Loading student...</div>
      </div>
    );
  }

  if (!student) {
    return (
      <div className="min-h-screen h-screen flex items-center justify-center bg-background">
        <div className="text-error">Student not found</div>
      </div>
    );
  }

  // Prepare chart data from ratings
  const chartData = ratings
    .map((entry) => {
      // Separate ratings by category
      const visualRatings = entry.ratings.filter(r => r.category === 'visual_cue');
      const verbalRatings = entry.ratings.filter(r => r.category === 'verbal_cue');
      const gesturalRatings = entry.ratings.filter(r => r.category === 'gestural_cue');

      const calcAverage = (ratings: any[]) => {
        if (ratings.length === 0) return null;
        const sum = ratings.map(r => mapLevelToNumber(r.level)).reduce((a, b) => a + b, 0);
        return Math.round(sum / ratings.length);
      };

      // Use session_date or fallback to session_id
      return {
        session: entry.session_date || entry.session_id,
        visual_cue: calcAverage(visualRatings),
        verbal_cue: calcAverage(verbalRatings),
        gestural_cue: calcAverage(gesturalRatings),
      };
    });

    const engagementData = ratings
    .map((entry) => {
      const engagementRatings = entry.ratings.filter(r => r.category === 'engagement');

      const calcAverage = (ratings: any[]) => {
        if (ratings.length === 0) return null;
        const sum = ratings.map(r => mapLevelToNumber(r.level)).reduce((a, b) => a + b, 0);
        return Math.round(sum / ratings.length);
      };

      return {
        session: entry.session_date || entry.session_id,
        engagement: calcAverage(engagementRatings),
      };
    });

  return (
    <AppLayout>
      <div className="w-full h-screen bg-background">
        <div className="w-full h-full flex flex-col gap-6 p-10 relative overflow-y-auto">
          <div className="flex flex-col gap-3">
            <Link
              href="/students"
              className="inline-flex items-center gap-2 text-secondary hover:text-primary transition-colors group w-fit"
            >
              <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
              <span className="text-sm font-medium">Back to Students</span>
            </Link>
            <PageHeader
              title="Student Profile"
              icon={User}
              className="mb-0!"
              actions={
                <Button
                  variant="outline"
                  className={`w-fit p-4 flex flex-row items-center gap-2 ${CORNER_ROUND} shrink-0`}
                  onClick={handleDelete}
                >
                  <Trash />
                  Delete
                </Button>
              }
            />
          </div>

          <div className="flex flex-col gap-6 shrink-0">
            <div className="grid grid-cols-2 gap-6 h-60">
              {/* Student Profile */}
              <div
                className={`flex-1 bg-card border-2 border-default ${CORNER_ROUND} overflow-hidden flex flex-col relative`}
              >
                {/* Edit Profile Button - Separate Section */}
                <div className="flex justify-end p-3 flex-shrink-0 relative z-10">
                  <Button
                    onClick={() => {
                      /* Navigate to edit page */
                    }}
                    variant="secondary"
                    className="flex items-center gap-2 hover:bg-accent h-8"
                    size="sm"
                  >
                    <span className="text-base font-medium">Edit Profile</span>
                    <PencilLine size={18} />
                  </Button>
                </div>

                {/* Content - positioned absolutely to center in entire card */}
                <div className="absolute inset-0 flex items-center justify-center gap-8 p-5 pointer-events-none">
                  <div className="max-h-full aspect-square border-2 border-default rounded-full flex-shrink-0 pointer-events-auto">
                    <Avatar
                      name={fullName + student.id}
                      variant={avatarVariant}
                      className="w-full h-full"
                    />
                  </div>

                  <div className="flex flex-col gap-3 flex-1 pointer-events-auto">
                    <div className="text-4xl font-bold text-primary">
                      {initials}
                    </div>

                    <div className="flex items-center gap-3 flex-wrap">
                      <span className="text-xl font-medium text-primary">
                        Grade {student.grade}
                      </span>
                      {student.school_name && (
                        <SchoolTag schoolName={student.school_name} />
                      )}
                    </div>
                  </div>
                </div>
              </div>

              {/* Upcoming Sessions */}
              <div
                className={`flex-1 bg-card border-2 border-default ${CORNER_ROUND} ${PADDING} flex flex-col gap-4`}
              >
                <h2>
                  Upcoming Sessions
                </h2>
                  <UpcomingSession studentId={studentId} />
              </div>
            </div>
          </div>
          {/* Attendance */}
          <div
            className={`bg-card border-2 flex flex-col gap-4 border-default ${CORNER_ROUND} ${PADDING} h-[30vh] min-h-[220px]`}
          >
          <h2>
            Goal Progress
          </h2>
            {attendanceLoading ? (
              <div className="flex items-center justify-center h-full">
                <div className="text-sm text-muted-foreground">
                  Loading attendance...
                </div>
              </div>
            ) : (
              <div
                className={
                  attendance?.total_count === 0
                    ? "opacity-50 pointer-events-none"
                    : ""
                }
              >
                <CustomPieChart
                  percentage={
                    attendance && attendance.total_count > 0
                      ? Math.round(
                          (attendance.present_count / attendance.total_count) *
                            100
                        )
                      : 0
                  }
                  title="Attendance"
                  showPlaceholder={!attendance || attendance.total_count === 0}
                />
              </div>
            )}
          </div>
          {/* Goals, Session Notes, and Ratings */}
          <div className="grid grid-cols-2 gap-6 h-[25vh] min-h-[300px]">
            {/* IEP Goals */}
            <div className="bg-card border-2 border-default rounded-4xl p-5 gap-6 flex flex-col overflow-hidden relative h-full">
              <h2>
                IEP Goals
              </h2>
              {/* Edit toggle button */}
              <div className="absolute top-5 right-5 z-10 flex gap-2">
                {edit ? (
                  <>
                    <Button
                      onClick={handleSave}
                      disabled={isSaving}
                      className="w-10 h-10 p-0 bg-green-600 hover:bg-green-700"
                      size="icon"
                    >
                      <Save size={20} />
                    </Button>
                    <Button
                      onClick={handleCancel}
                      disabled={isSaving}
                      className="w-10 h-10 p-0 bg-red-600 hover:bg-red-700"
                      size="icon"
                    >
                      <X size={20} />
                    </Button>
                  </>
                ) : (
                  <Button
                    onClick={() => setEdit(!edit)}
                    className="w-10 h-10 p-0"
                    variant="secondary"
                    size="icon"
                  >
                    <PencilLine size={20} />
                  </Button>
                )}
              </div>
              <div className="flex-1 overflow-y-auto flex flex-col gap-2">
                {iepGoals.length === 0 && !edit ? (
                  <div className="text-muted-foreground italic">
                    No IEP goals set
                  </div>
                ) : (
                  iepGoals.map((goal, index) => (
                    <div
                      key={index}
                      className={`w-full text-lg flex items-center gap-2
                    rounded-2xl transition bg-card cursor-pointer select-none border-2 border-border ${PADDING} ${!edit && "hover:scale-99"}`}
                    onClick={() => setEdit(true)}
                    >
                      <Trophy size={20} className="flex-shrink-0" />
                      {edit ? (
                        <>

                          <input
                            value={goal}
                            onChange={(e) => updateGoal(index, e.target.value)}
                            onBlur={() =>
                              goal.trim() === "" && deleteGoal(index)
                            }
                            className="flex-1 bg-transparent outline-none py-1 leading-normal"
                            placeholder="Enter IEP goal..."
                            autoFocus
                          />
                          <Button
                            onClick={() => deleteGoal(index)}
                            variant="ghost"
                            size="icon"
                            className="text-red-600 hover:text-red-700 hover:bg-red-100 flex-shrink-0 w-8 h-8"
                          >
                            <Trash2 size={18} />
                          </Button>
                        </>
                      ) : (
                        <span className="py-1 leading-normal">{goal}</span>
                      )}
                    </div>
                  ))
                )}
              </div>
              {edit && (
                <Button
                  onClick={addGoal}
                  variant="outline"
                  className="w-full border-2 rounded-full gap-2"
                >
                  <CirclePlus size={20} />
                  Add IEP Goal
                </Button>
              )}
            </div>

            {/* Session Notes */}
            <div className="bg-card border-2 border-default rounded-4xl p-5 gap-2 flex flex-col overflow-hidden h-full">
              <h2 >
                Session Notes
              </h2>
              <div className="flex-1 overflow-y-auto">
                <SessionNotes studentId={studentId} />
              </div>
            </div>
          </div>
          <div className="w-full flex flex-col gap-8">
            <h1>Progress History</h1>
            <SessionRatingsChart
              chartData={chartData}
              title="Ratings"
              categories={[
                { key: "visual_cue", label: "Visual Cue", color: "var(--color-blue)" },
                { key: "verbal_cue", label: "Verbal Cue", color: "var(--color-pink)" },
                { key: "gestural_cue", label: "Gestural Cue", color: "var(--color-orange)" }
              ]}
            />
            <SessionRatingsChart
              title="Engagement"
              chartData={engagementData}
              categories={[
                { key: "engagement", label: "Engagement", color: "var(--color-blue)" }
              ]}
            />
          </div>
        </div>
      </div>
    </AppLayout>
  );
}

export default StudentPage;
