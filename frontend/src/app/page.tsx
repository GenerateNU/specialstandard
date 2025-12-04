"use client";

import AppLayout from "@/components/AppLayout";
import MiniCalendar from "@/components/dashboard/MiniCalendar";
import RecentStudentCard from "@/components/dashboard/RecentStudentCard";
import UpcomingSessionCard from "@/components/dashboard/UpcomingSessionCard";
import StudentSchedule from "@/components/schedule/StudentSchedule";
import { Button } from "@/components/ui/button";
import { useAuthContext } from "@/contexts/authContext";
import { useNewsletter } from "@/hooks/useNewsletter";
import { useRecentlyViewedStudents } from "@/hooks/useRecentlyViewedStudents";
import { useSessions } from "@/hooks/useSessions";
import { useSessionStudentsForSession } from "@/hooks/useSessionStudents";
import { useTherapists } from "@/hooks/useTherapists";
import type { Session } from "@/lib/api/theSpecialStandardAPI.schemas";
import { formatDateString, formatTime, getTherapistName } from "@/lib/utils";
import { ChevronDown, Loader2, UserCircle } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function Home() {
  const { isAuthenticated, isLoading, userId } = useAuthContext();
  const router = useRouter();
  const CORNER_ROUND = "rounded-2xl";

  const {
    downloadNewsletter,
    isLoading: newsletterLoading,
    error: newsletterError,
  } = useNewsletter();

  // Recently viewed students
  const { recentStudents } = useRecentlyViewedStudents();

  // Fetch therapists data
  const { therapists } = useTherapists();
  const currentTherapist =
    therapists.find((t) => t.id === userId) || therapists[0];

  // Fetch all sessions for today
  const { sessions: allSessions, isLoading: sessionsLoading } = useSessions({
    startdate: new Date(new Date().setHours(0, 0, 0, 0)).toISOString(),
  });

  const sessions = allSessions.slice(0, 5);

  const [selectedSession, setSelectedSession] = useState<Session | null>(null);
  const [openSection, setOpenSection] = useState<
    "students" | "curriculum" | null
  >("students");

  // Fetch students for the selected session
  const { students: sessionStudents, isLoading: studentsLoading } =
    useSessionStudentsForSession(selectedSession?.id || "");

  // Set first session as selected by default
  useEffect(() => {
    if (sessions.length > 0 && !selectedSession) {
      setSelectedSession(sessions[0]);
    }
  }, [sessions, selectedSession]);

  // Redirect to login if not authenticated
  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  const handleDownloadNewsletter = async () => {
    await downloadNewsletter();
  };

  return (
    <AppLayout>
      <div className="grow bg-background flex flex-row h-screen">
        {/* Main Content */}
        <div className="w-full md:w-7/10 p-10 flex flex-col gap-8 overflow-y-scroll">
          {/* Header */}
          <div className="flex flex-row justify-between items-end shrink-0">
            <div className="flex flex-col items-left justify-start">
              <h1 className="text-4xl font-bold">Your Educator Dashboard</h1>
              <p className="text-2xl text-secondary">
                {new Date().toLocaleDateString("en-US", {
                  weekday: "long",
                  month: "long",
                  day: "numeric",
                })}
              </p>
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                onClick={handleDownloadNewsletter}
                disabled={newsletterLoading}
                className="text-pink border-pink hover:bg-pink hover:text-white"
              >
                {newsletterLoading ? (
                  <>
                    <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                    Downloading...
                  </>
                ) : (
                  "This Month's Newsletter ↓"
                )}
              </Button>
              {newsletterError && (
                <span className="text-sm text-error">{newsletterError}</span>
              )}
            </div>
          </div>

          {/* Upcoming Sessions */}
          <div
            className={`w-full bg-card p-6 gap-4 ${CORNER_ROUND} flex flex-col transition`}
          >
            <div className="w-full flex items-center justify-between">
              <h3 className="text-xl font-semibold">Upcoming Sessions</h3>
              <Button
                size="sm"
                variant="default"
                onClick={() => router.push("/calendar?view=card")}
                className="bg-pink hover:bg-pink-hover text-white"
              >
                View All Sessions
              </Button>
            </div>

            <div className="grid grid-cols-[1fr_2fr] w-full gap-4 items-start">
              {/* Sessions List */}
              <div className="gap-2 flex flex-col">
                {sessionsLoading ? (
                  <div className="flex items-center justify-center py-8">
                    <Loader2 className="w-6 h-6 animate-spin text-primary" />
                  </div>
                ) : sessions.length > 0 ? (
                  sessions.map((session, index) => (
                    <div
                      key={session.id}
                      onClick={() => setSelectedSession(session)}
                      className="cursor-pointer animate-in fade-in slide-in-from-left-3 duration-300"
                      style={{ animationDelay: `${index * 75}ms` }}
                    >
                      <UpcomingSessionCard
                        className={`transition-all duration-200 ${
                          selectedSession?.id === session.id
                            ? "bg-orange-disabled ring-2 ring-offset-1 ring-blue-disabled scale-[1.02]"
                            : "bg-white-disabled hover:scale-[1.01]"
                        }`}
                        sessionID={session.id}
                        sessionName={session.session_name}
                        startTime={formatTime(session.start_datetime)}
                        endTime={formatTime(session.end_datetime)}
                        date={formatDateString(session.start_datetime)}
                      />
                    </div>
                  ))
                ) : (
                  <div className="flex items-center justify-center py-8 text-muted">
                    No upcoming sessions
                  </div>
                )}
              </div>


              {/* Session Details Panel */}
              <div
                className={`p-4 text-sm font-normal ${CORNER_ROUND} justify-start flex flex-col gap-4 self-start transition-all duration-300 bg-orange-disabled`}
              >
                {selectedSession ? (
                  <div className="flex flex-col gap-6 animate-in fade-in duration-300">
                    <div className="flex flex-row flex-1">
                      <div className="flex flex-col flex-1">
                        <strong className="text-primary">{selectedSession.session_name}</strong>
                        <strong className="text-secondary">
                          {getTherapistName(
                            selectedSession.therapist_id,
                            therapists
                          )}
                        </strong>
                      </div>
                      <div className="flex flex-col flex-1 text-right">
                        <span className="text-secondary">
                          {formatTime(selectedSession.start_datetime)} –{" "}
                          {formatTime(selectedSession.end_datetime)}
                        </span>
                        <span className="text-secondary">
                          {formatDateString(selectedSession.start_datetime)}
                        </span>
                        {selectedSession.location && (
                          <div className="text-sm text-secondary">
                            <strong>Location:</strong> {selectedSession.location}
                          </div>
                        )}
                      </div>
                    </div>
                    {selectedSession.notes && (
                      <div className="text-sm text-secondary">
                        <strong>Notes:</strong> {selectedSession.notes}
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="flex items-center justify-center py-8 text-muted">
                    <span>Select a session to view details</span>
                  </div>
                )}

                {selectedSession && (
                  <div className="w-full gap-3 flex flex-col mt-2">
                    {/* Students Section */}
                    <div className="flex flex-col w-full border-b-2 border-border overflow-hidden">
                      <button
                        className="flex items-center justify-between px-2 py-2 cursor-pointer group hover:bg-orange/20 transition-colors"
                        onClick={() =>
                          setOpenSection(
                            openSection === "students" ? null : "students"
                          )
                        }
                      >
                        <h4 className="font-semibold text-primary">Students</h4>
                        <ChevronDown
                          size={18}
                          className={`transition-transform duration-300 ease-in-out text-primary ${
                            openSection === "students" ? "rotate-180" : ""
                          }`}
                        />
                      </button>
                      <div
                        className={`grid transition-all duration-300 ease-in-out ${
                          openSection === "students"
                            ? "grid-rows-[1fr] opacity-100"
                            : "grid-rows-[0fr] opacity-0"
                        }`}
                      >
                        <div className="overflow-hidden">
                          <div className="mt-1 mb-3 text-sm font-normal space-y-1 px-2">
                            {studentsLoading ? (
                              <div className="flex items-center gap-2 text-muted">
                                <Loader2 className="w-4 h-4 animate-spin" />
                                Loading students...
                              </div>
                            ) : sessionStudents.length > 0 ? (
                              <div className="flex flex-wrap gap-x-2 gap-y-0.5">
                                {sessionStudents.map((student, index) => (
                                  <span
                                    key={student.id}
                                    className="text-sm text-secondary animate-in fade-in duration-300"
                                    style={{
                                      animationDelay: `${index * 50}ms`,
                                    }}
                                  >
                                    {student.first_name} {student.last_name}
                                    {index < sessionStudents.length - 1 && ","}
                                  </span>
                                ))}
                              </div>
                            ) : (
                              <div className="text-muted">
                                No students assigned
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Curriculum Section */}
                    <div className="flex flex-col w-full border-b-2 border-border overflow-hidden">
                      <button
                        className="flex items-center justify-between px-2 py-2 cursor-pointer group hover:bg-orange/20 transition-colors"
                        onClick={() =>
                          setOpenSection(
                            openSection === "curriculum" ? null : "curriculum"
                          )
                        }
                      >
                        <h4 className="font-semibold text-primary">Curriculum</h4>
                        <ChevronDown
                          size={18}
                          className={`transition-transform duration-300 ease-in-out text-primary ${
                            openSection === "curriculum" ? "rotate-180" : ""
                          }`}
                        />
                      </button>
                      <div
                        className={`grid transition-all duration-300 ease-in-out ${
                          openSection === "curriculum"
                            ? "grid-rows-[1fr] opacity-100"
                            : "grid-rows-[0fr] opacity-0"
                        }`}
                      >
                        <div className="overflow-hidden">
                          <div className="mt-1 mb-3 text-sm font-normal px-2">
                            <div className="text-muted">
                              Curriculum content coming soon...
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Schedule */}
          <div
            className={`w-full shrink-0 p-6 bg-card flex items-start flex-col gap-4 ${CORNER_ROUND}`}
          >
            <div className="w-full flex items-center justify-between">
              <h3 className="text-xl font-semibold">Your Schedule</h3>
              <Button
                size="sm"
                variant="default"
                onClick={() => router.push("/calendar")}
                className="bg-pink hover:bg-pink-hover text-white"
              >
                View Full Schedule
              </Button>
            </div>
            <div className="w-full h-120">
              <StudentSchedule initialView="work_week" className="h-full" />
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="flex flex-col p-10 w-3/10 h-screen bg-white-hover sticky top-0 space-y-6 overflow-y-auto">
          <Link href="/profile" className="block">
            <div className="flex flex-row justify-between items-center hover:shadow-md hover:-translate-y-1 cursor-pointer text-black transition p-4 rounded-2xl hover:bg-orange-disabled shrink-0">
              <div>
                <h4 className="text-black font-semibold">
                  {currentTherapist
                    ? `${currentTherapist.first_name} ${currentTherapist.last_name}`
                    : "Loading..."}
                </h4>
                <p className="text-sm text-black/70">My Profile</p>
              </div>
              <UserCircle size={36} strokeWidth={1} className="text-black" />
            </div>
          </Link>

          <div className="flex flex-col w-full p-6 !bg-white rounded-xl gap-4 min-h-0 flex-1 overflow-hidden">
            <div>
              <h3 className="text-lg font-semibold text-primary">Students</h3>
              <p className="text-sm text-muted">Recently Viewed</p>
            </div>
            <div className="w-full flex-1 flex flex-col gap-2 overflow-y-auto min-h-0">
              {recentStudents.length > 0 ? (
                recentStudents.map((student, index) => (
                  <div
                    key={student.id}
                    className="animate-in fade-in duration-300"
                    style={{ animationDelay: `${index * 50}ms` }}
                  >
                    <RecentStudentCard
                      id={student.id}
                      firstName={student.first_name}
                      lastName={student.last_name}
                      grade={student.grade}
                    />
                  </div>
                ))
              ) : (
                <div className="flex-1 flex items-center justify-center text-sm text-muted">
                  No recently viewed students
                </div>
              )}
            </div>
            <Button
              className="w-full shrink-0 bg-pink hover:bg-pink-hover text-white"
              onClick={() => router.push("/students")}
            >
              View All Students
            </Button>
          </div>


          <div className="w-full bg-card rounded-xl shrink-0">
            <MiniCalendar />
          </div>
        </div>
      </div>
    </AppLayout>
  );
}