'use client'

import {useParams, useRouter} from "next/navigation";
import {useStudentSessions} from "@/hooks/useStudentSessions";
import moment from "moment";
import {ArrowLeft, Check, Pencil, X} from "lucide-react";
import {IndividualityBadge} from "@/components/sessions/UpcomingSession";
import AppLayout from "@/components/AppLayout";
import {useEffect, useState} from "react";
import {useSessions} from "@/hooks/useSessions";
import {useSessionStudents} from "@/hooks/useSessionStudents";

export default function SessionHistory() {
  const params = useParams()
  const studentID = params.studentID as string
  const sessionID = params.sessionID as string
  const { sessions, isLoading, error } = useStudentSessions(studentID)
  const sessionInfo = sessions?.find(item => item.session.id === sessionID)
  const session = sessionInfo?.session
  const router = useRouter()

  const { updateSession } = useSessions()
  const [isStudentSessionNotesEditing, setIsStudentSessionNotesEditing] = useState(false)
  const [isSessionNotesEditing, setIsSessionNotesEditing] = useState(false)
  const [notesValue, setNotesValue] = useState("")
  const [studentSessionNotesValue, setStudentSessionNotesValue] = useState("")

  const { updateSessionStudent } = useSessionStudents()

  // Sync state with data when it loads
  useEffect(() => {
    if (session?.notes) {
      setNotesValue(session.notes)
    }
    if (sessionInfo?.notes) {
      setStudentSessionNotesValue(sessionInfo.notes)
    }
  }, [session, sessionInfo])

  if (isLoading) {
    return <AppLayout><div className="p-6 text-muted-foreground">Loading Session...</div></AppLayout>
  }

  if (error || !session) {
    return <AppLayout><div className="p-6 text-muted-foreground italic">Session not found.</div></AppLayout>
  }

  const start = moment(session.start_datetime)
  const end = moment(session.end_datetime)

  const sessionNotesHandleSave = async () => {
    await updateSession(sessionID, { notes: notesValue })
    setIsSessionNotesEditing(false)
  }

  const sessionNotesHandleCancel = () => {
    setNotesValue(session.notes || "")
    setIsSessionNotesEditing(false)
  }

  const studentSessionNotesHandleSave = async () => {
    await updateSessionStudent({
      session_id: session.id,
      student_id: sessionInfo?.student_id,
      notes: studentSessionNotesValue 
    })
    setIsStudentSessionNotesEditing(false)
  }

  const studentSessionNotesHandleCancel = () => {
    setStudentSessionNotesValue(sessionInfo?.notes || "")
    setIsStudentSessionNotesEditing(false)
  }

  return (
    <AppLayout>
      <div className="grow bg-background flex h-screen">
        <div className="flex-1 p-10 flex flex-col overflow-y-scroll">
          {/* Header */}
          <div className="flex items-center gap-3 mb-8">
            <ArrowLeft 
              onClick={() => router.push(`/student/${studentID}/session-history`)}
              className="w-8 h-8 cursor-pointer hover:opacity-70 transition" 
            />
            <h1 className="text-4xl font-bold text-primary">{session.session_name}</h1>
          </div>

          {/* Session Info Card */}
          <div className="p-6 bg-card border-2 border-default rounded-2xl shadow-md mb-8">
            <div className="flex items-start justify-between">
              <div>
                <h2 className="font-semibold text-lg mb-2">{session.session_name || "Session"}</h2>
                <p className="text-muted-foreground">
                  {start.format('dddd, MMMM D, YYYY')} • {start.format('h:mm A')} - {end.format('h:mm A')}
                </p>
              </div>
              <IndividualityBadge sessionId={session.id} />
            </div>
          </div>

          {/* Notes Container */}
          <div className="bg-card border-2 border-default rounded-2xl shadow-md flex flex-col">
            <div className="p-8 space-y-8">
              {/* Student Specific Notes Section */}
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-semibold">Student Specific Notes</h3>
                  {!isStudentSessionNotesEditing && (
                    <button 
                      onClick={() => setIsStudentSessionNotesEditing(true)}
                      className="flex items-center gap-2 px-3 py-2 rounded-lg border border-default 
                        bg-orange-disabled hover:shadow-md transition text-sm font-medium cursor-pointer"
                    >
                      <Pencil className="w-4 h-4" />
                      Edit
                    </button>
                  )}
                </div>

                {!isStudentSessionNotesEditing ? (
                  <div className="text-base text-foreground">
                    {studentSessionNotesValue ? (
                      <p className="whitespace-pre-wrap">{studentSessionNotesValue}</p>
                    ) : (
                      <span className="text-muted-foreground italic">No notes added.</span>
                    )}
                  </div>
                ) : (
                  <div className="flex flex-col gap-3">
                    <textarea
                      value={studentSessionNotesValue}
                      onChange={(e) => setStudentSessionNotesValue(e.target.value)}
                      className="w-full p-4 rounded-xl border border-default
                        min-h-[200px] bg-white text-base focus:outline-none
                        focus:ring-2 focus:ring-primary resize-none"
                      placeholder="Add student-specific notes here..."
                      title="Student-Specific Notes"
                    />
                    <div className="flex gap-3 justify-end">
                      <button
                        onClick={studentSessionNotesHandleCancel}
                        className="px-4 py-2 rounded-lg border border-default bg-white hover:bg-muted
                          font-semibold flex items-center gap-2 transition"
                      >
                        <X className="w-4 h-4" />
                        Cancel
                      </button>
                      <button
                        onClick={studentSessionNotesHandleSave}
                        className="px-4 py-2 rounded-lg bg-primary text-white hover:bg-primary/90
                          font-semibold flex items-center gap-2 transition"
                      >
                        <Check className="w-4 h-4" />
                        Save
                      </button>
                    </div>
                  </div>
                )}
              </div>

              {/* Divider */}
              <div className="border-t border-default"></div>

              {/* Session Notes Section */}
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-semibold">Session Notes</h3>
                  {!isSessionNotesEditing && (
                    <button 
                      onClick={() => setIsSessionNotesEditing(true)}
                      className="flex items-center gap-2 px-3 py-2 rounded-lg border border-default 
                        bg-orange-disabled hover:shadow-md transition text-sm font-medium cursor-pointer"
                    >
                      <Pencil className="w-4 h-4" />
                      Edit
                    </button>
                  )}
                </div>

                {!isSessionNotesEditing ? (
                  <div className="text-base text-foreground">
                    {notesValue ? (
                      <p className="whitespace-pre-wrap">{notesValue}</p>
                    ) : (
                      <span className="text-muted-foreground italic">No notes added.</span>
                    )}
                  </div>
                ) : (
                  <div className="flex flex-col gap-3">
                    <textarea
                      value={notesValue}
                      onChange={(e) => setNotesValue(e.target.value)}
                      className="w-full p-4 rounded-xl border border-default
                        min-h-[200px] bg-white text-base focus:outline-none
                        focus:ring-2 focus:ring-primary resize-none"
                      placeholder="Add session notes here..."
                      title="Session Notes"
                    />
                    <div className="flex gap-3 justify-end">
                      <button
                        onClick={sessionNotesHandleCancel}
                        className="px-4 py-2 rounded-lg border border-default bg-white hover:bg-muted
                          font-semibold flex items-center gap-2 transition"
                      >
                        <X className="w-4 h-4" />
                        Cancel
                      </button>
                      <button
                        onClick={sessionNotesHandleSave}
                        className="px-4 py-2 rounded-lg bg-primary text-white hover:bg-primary/90
                          font-semibold flex items-center gap-2 transition"
                      >
                        <Check className="w-4 h-4" />
                        Save
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </AppLayout>
  )
}