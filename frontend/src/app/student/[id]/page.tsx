"use client";

import AppLayout from "@/components/AppLayout";
import RecentSession from "@/components/attendance/RecentSession";
import StudentSchedule from "@/components/schedule/StudentSchedule";
import SessionNotes from "@/components/sessions/sessionNotes";
import { Avatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { useRecentlyViewedStudents } from "@/hooks/useRecentlyViewedStudents";
import { useStudents } from "@/hooks/useStudents";
import { getAvatarVariant } from "@/lib/utils";
import {
  ChevronLeft,
  CirclePlus,
  PencilLine,
  Save,
  Trash2,
  X,
} from "lucide-react";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";

function StudentPage() {
  const params = useParams();
  const studentId = params?.id as string;

  if (!params || !studentId) {
    return (
      <div className="min-h-screen h-screen flex items-center justify-center bg-background">
        <div className="text-error">Invalid student ID</div>
      </div>
    );
  }

  const { students, isLoading, updateStudent } = useStudents();
  const student = students.find((s) => s.id === studentId);
  const { addRecentStudent } = useRecentlyViewedStudents();

  // Track this student as recently viewed - only when studentId changes
  useEffect(() => {
    if (student && studentId) {
      addRecentStudent(student);
    }
  }, [studentId]); // Only track when navigating to a new student

  const [edit, setEdit] = useState(false);
  const [iepGoals, setIepGoals] = useState<string[]>([]);
  const [isSaving, setIsSaving] = useState(false);

  // Initialize IEP goals from student data
  useEffect(() => {
    if (student?.iep) {
      try {
        // Try to parse as JSON array
        const parsed = JSON.parse(student.iep);
        if (Array.isArray(parsed)) {
          setIepGoals(parsed);
        } else {
          // If it's a string, split by newlines or treat as single goal
          setIepGoals([student.iep]);
        }
      } catch {
        // If parsing fails, treat as plain text (split by newlines)
        setIepGoals(student.iep.split("\n").filter((g) => g.trim()));
      }
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
      // Save IEP goals as JSON string
      const iepString = JSON.stringify(filteredGoals);
      await updateStudent(student.id, { iep: iepString });
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
    if (student?.iep) {
      try {
        const parsed = JSON.parse(student.iep);
        setIepGoals(Array.isArray(parsed) ? parsed : [student.iep]);
      } catch {
        setIepGoals(student.iep.split("\n").filter((g) => g.trim()));
      }
    } else {
      setIepGoals([]);
    }
    setEdit(false);
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

  return (
    <AppLayout>
      <div className="w-full h-screen bg-background">
        <div
          className={`w-full h-full grid grid-rows-2 gap-8 ${PADDING} relative`}
        >
          {/* Edit toggle button */}
          <div className="absolute top-1/2 right-5 z-20 flex gap-2">
            {edit ? (
              <>
                <Button
                  onClick={handleSave}
                  disabled={isSaving}
                  className="w-12 h-12 p-0 bg-green-600 hover:bg-green-700"
                  size="icon"
                >
                  <Save size={20} />
                </Button>
                <Button
                  onClick={handleCancel}
                  disabled={isSaving}
                  className="w-12 h-12 p-0 bg-red-600 hover:bg-red-700"
                  size="icon"
                >
                  <X size={20} />
                </Button>
              </>
            ) : (
              <Button
                onClick={() => setEdit(!edit)}
                className="w-12 h-12 p-0"
                variant="secondary"
                size="icon"
              >
                <PencilLine size={20} />
              </Button>
            )}
          </div>

          <div className="flex gap-8">
            {/* pfp and initials */}
            <div className="flex flex-col items-left justify-between gap-2 w-1/6">
              <Button
                variant="outline"
                className={`w-1/2 pl-4 flex flex-row items-left justify-start ${CORNER_ROUND}`}
                onClick={() => window.history.back()}
              >
                <ChevronLeft />
                Back
              </Button>
              <div className="w-full aspect-square border-2 border-accent rounded-full">
                <Avatar
                  name={fullName + student.id}
                  variant={avatarVariant}
                  className="w-full h-full"
                />
              </div>
              <div
                className={`w-full h-1/5 text-3xl font-bold flex items-center 
                justify-center bg-background border-border border-2 ${PADDING}
                 ${CORNER_ROUND}`}
              >
                {initials}
              </div>
            </div>
            {/* student schedule */}
            <div
              className={`flex-[3] ${CORNER_ROUND} overflow-hidden bg-blue flex flex-col justify-between ${PADDING}`}
            >
              <StudentSchedule studentId={studentId} className="h-3/4" />
              <Button
                className="h-1/5 rounded-2xl text-lg font-bold "
                variant="secondary"
              >
                View Student Schedule
              </Button>
            </div>
            <div
              className={`bg-pink flex-2 flex flex-col items-center justify-between ${CORNER_ROUND} ${PADDING}`}
            >
              <div className="w-full h-3/4 text-3xl font-bold flex items-center rounded-2xl">
                <RecentSession studentId={studentId} />
              </div>
              <Button
                className="w-full h-1/5 rounded-2xl text-lg font-bold "
                variant="secondary"
              >
                View Student Attendance
              </Button>
            </div>
          </div>

          {/* Goals and Session Notes */}
          <div className="grid grid-cols-2 gap-8 overflow-hidden">
            <div className="gap-2 flex flex-col overflow-hidden">
              <div className="w-full text-2xl text-primary flex items-baseline font-semibold">
                IEP Goals
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
                    rounded-2xl transition bg-background select-none border-2 border-border ${PADDING} ${!edit && "hover:scale-99"}`}
                    >
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
                          />
                          <Button
                            onClick={() => deleteGoal(index)}
                            variant="ghost"
                            size="icon"
                            className="text-red-600 hover:text-red-700 hover:bg-red-100 flex-shrink-0"
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

            <div className="gap-2 flex flex-col overflow-hidden">
              <div className="gap-2 flex flex-col overflow-hidden">
                <div className="w-full text-2xl text-primary flex items-baseline font-semibold">
                  Session Notes
                </div>
                <div className="flex-1 overflow-y-auto">
                  <SessionNotes studentId={studentId} />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AppLayout>
  );
}

export default StudentPage;
