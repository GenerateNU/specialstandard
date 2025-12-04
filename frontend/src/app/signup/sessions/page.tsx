"use client";

import UpcomingSessionCard from "@/components/dashboard/UpcomingSessionCard";
import { Avatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { useSessions } from "@/hooks/useSessions";
import { useSessionStudentsForSession } from "@/hooks/useSessionStudents";
import { getAvatarName, getAvatarVariant } from "@/lib/avatarUtils";
import { ArrowLeft, Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function SessionsPage() {
  const router = useRouter();
  const { sessions, isLoading: loadingSessions, refetch } = useSessions();
  const [selectedSessionId, setSelectedSessionId] = useState<string>("");

  // Use the hook to get students for the selected session
  const { students: sessionStudents, isLoading: loadingStudents } =
    useSessionStudentsForSession(selectedSessionId);

  // Get the first session's ID when sessions load
  useEffect(() => {
    if (sessions.length > 0 && !selectedSessionId) {
      setSelectedSessionId(sessions[0].id);
    }
  }, [sessions, selectedSessionId]);

  useEffect(() => {
    const userId =
      localStorage.getItem("temp_userId") || localStorage.getItem("userId"); // Check temp_userId first

    if (!userId) {
      router.push("/signup/welcome");
    }

    // Refetch sessions when page loads
    refetch();
  }, [router, refetch]);

  const handleBack = () => {
    router.push("/signup/students");
  };

  const handleAddSession = () => {
    router.push("/signup/sessions/add");
  };

  const handleFinish = () => {
    router.push("/signup/complete");
  };

  // Helper functions for formatting
  const formatTime = (datetime: string) => {
    const date = new Date(datetime);
    return date.toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
      hour12: true,
    });
  };

  const formatDateString = (datetime: string) => {
    const date = new Date(datetime);
    return date.toLocaleDateString("en-US", {
      weekday: "long",
      month: "long",
      day: "numeric",
    });
  };

  if (loadingSessions) {
    return (
      <div className="flex items-center justify-center min-h-screen p-8">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center min-h-screen p-8">
      <div className="max-w-2xl w-full">
        <div className="flex justify-between items-center mb-6">
          <button
            onClick={handleBack}
            className="flex items-center cursor-pointer text-secondary hover:text-primary transform hover:scale-125 transition-colors"
          >
            <ArrowLeft className="w-4 h-4 mr-1" />
            Back
          </button>
        </div>

        <h1 className="text-3xl font-bold text-primary mb-8">Your Sessions</h1>

        {sessions.length === 0 ? (
          <>
            <div className="bg-card rounded-lg border border-default p-12 text-center mb-6">
              <p className="text-secondary mb-6">
                No sessions scheduled yet. Add your first session to get
                started.
              </p>
            </div>

            <div className="flex gap-3">
              <Button
                onClick={handleAddSession}
                variant="outline"
                className="flex-1"
              >
                Add Session
              </Button>

              <Button
                onClick={handleFinish}
                className="flex-1 hover:bg-accent-hover text-white"
              >
                Finish
              </Button>
            </div>
          </>
        ) : (
          <>
            {/* Session Card Display */}
            <div className="bg-card rounded-lg border border-default p-6 mb-6">
              {sessions.length > 0 && (
                <div className="space-y-4">
                  <div className="text-sm text-secondary mb-2">
                    <div>Session Name</div>
                    <div className="font-semibold text-primary">
                      {formatTime(sessions[0].start_datetime)} -{" "}
                      {formatTime(sessions[0].end_datetime)}, Does not repeat
                    </div>
                    {sessions[0].location && (
                      <div className="mt-1">{sessions[0].location}</div>
                    )}
                  </div>

                  <div>
                    <div className="text-sm text-secondary mb-2">Students</div>
                    <div className="space-y-2">
                      {loadingStudents ? (
                        <div className="text-sm text-secondary">
                          Loading students...
                        </div>
                      ) : sessionStudents.length > 0 ? (
                        sessionStudents.map((student: any) => (
                          <div
                            key={student.id}
                            className="flex items-center gap-3 p-2 bg-background rounded-lg border border-default"
                          >
                            <Avatar
                              name={getAvatarName(
                                student.first_name,
                                student.last_name,
                                student.id
                              )}
                              variant={getAvatarVariant(student.id)}
                              className="w-8 h-8"
                            />
                            <span className="text-sm text-primary">
                              {student.first_name} {student.last_name}
                            </span>
                            <span className="text-xs text-secondary ml-auto">
                              ID
                            </span>
                          </div>
                        ))
                      ) : (
                        <div className="text-sm text-secondary">
                          No students in this session
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* Show other sessions as cards */}
            {sessions.length > 1 && (
              <div className="space-y-2 mb-6">
                <p className="text-sm text-secondary mb-2">
                  Other scheduled sessions:
                </p>
                {sessions.slice(1).map((session) => (
                  <UpcomingSessionCard
                    key={session.id}
                    sessionName={session.session_name}
                    startTime={formatTime(session.start_datetime)}
                    endTime={formatTime(session.end_datetime)}
                    date={formatDateString(session.start_datetime)}
                    className="text-white"
                  />
                ))}
              </div>
            )}

            <div className="flex gap-3">
              <Button
                onClick={handleAddSession}
                variant="outline"
                className="flex-1"
              >
                Add Session
              </Button>

              <Button
                onClick={handleFinish}
                className="flex-1 hover:bg-accent-hover text-white"
              >
                Finish
              </Button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
