'use client'

import {useParams, useRouter} from "next/navigation";
import {useStudentSessions} from "@/hooks/useStudentSessions";
import moment from "moment";
import {ArrowLeft, Check, Pencil, X} from "lucide-react";
import {IndividualityBadge} from "@/components/sessions/UpcomingSession";
import AppLayout from "@/components/AppLayout";
import {useState} from "react";
import {useSessions} from "@/hooks/useSessions";

export default function SessionHistory() {
  const params = useParams()
  const studentID = params.studentID as string
  const sessionID = params.sessionID as string
  const { sessions, isLoading, error } = useStudentSessions(studentID)
  const sessionInfo = sessions?.find(item => item.session.id === sessionID)
  const session = sessionInfo?.session
  const router = useRouter()

  const { updateSession } = useSessions()
  const [isEditing, setIsEditing] = useState(false)
  const [notesValue, setNotesValue] = useState(session?.notes || "")

  if (isLoading) {
    return <div className="p-6 text-muted-foreground">Loading Session...</div>
  }

  if (error || !session) {
    return <div className="p-6 text-muted-foreground italic">Session not found.</div>
  }

  const start = moment(session.start_datetime)
  const end = moment(session.end_datetime)

  const handleSave = async () => {
    updateSession(sessionID, { notes: notesValue })

    setIsEditing(false)
  }

  const handleCancel = () => {
    setNotesValue(session.notes || "")
    setIsEditing(false)
  }

  return (
      <AppLayout>
          <div className="grow bg-background flex flex-row h-screen">
              <div className="w-[85%] p-10 flex flex-col overflow-y-scroll">
                  {/* Header */}
                  <div className="flex flex-row">
                      <ArrowLeft onClick={() => router.push(`/student/${studentID}/session-history`)}
                                 className="mt-1 mr-1 w-8 h-8 cursor-pointer"/>
                      <div className="flex items-center justify-between mb-8">
                          <h1 className="text-3xl font-bold text-primary">{session.session_name}</h1>
                      </div>
                  </div>

                  {/* Session Info Card */}
                  <div
                      key={session.id}
                      className="p-4 bg-card border-2 border-default rounded-[32px] flex flex-col
                                   justify-center min-h-20 w-full shadow-md"
                  >
                      <div className="ml-2 my-4">
                          <div className="flex flex-row">
                              <div>
                                  {/* Session title and relative time */}
                                  <div className="w-full flex justify-between items-center mb-1">
                                      <div className="font-semibold text-base">
                                          <h3>{session.session_name || `Session`}</h3>
                                      </div>
                                  </div>

                                  {/* Date & Time Range */}
                                  <div className="text-md text-muted-foreground">
                                      {start.format('dddd, MMMM D, YYYY')}
                                      {' | '}
                                      {start.format('h:mm A')}
                                      {' - '}
                                      {end.format('h:mm A')}
                                  </div>
                              </div>
                              <div className="pl-[3%] mt-[1vh] flex flex-col items-center">
                                <IndividualityBadge sessionId={session.id} />
                              </div>
                          </div>
                      </div>
                  </div>

                  {/* Notes Section */}
                  <div className="p-4 bg-card border-2 border-default rounded-[32px] flex flex-col
                                  mt-10 justify-center min-h-20 w-full shadow-md">
                      <div className="ml-2 my-4">
                          <div className="w-full flex flex-row justify-between items-center mb-3">
                              <h3>Session Notes</h3>

                              {isEditing ? (
                                <div className="flex flex-row gap-3 mt-2">
                                  <button
                                    onClick={handleSave}
                                    className="px-4 py-2 rounded-2xl bg-primary font-semibold
                                               bg-orange-disabled flex items-center gap-2"
                                  >
                                    <Check className="w-4 h-4" />
                                    Save
                                  </button>

                                  <button
                                    onClick={handleCancel}
                                    className="px-4 py-2 rounded-2xl bg-primary font-semibold
                                               bg-orange-disabled flex items-center gap-2"
                                  >
                                    <X className="w-4 h-4" />
                                    Cancel
                                  </button>
                                </div>
                              ) : (
                                <button onClick={() => setIsEditing(true)}
                                        className="flex items-center gap-1 px-3 py-1 rounded-xl
                                                   border border-default bg-orange-disabled hover:shadow-md
                                                   transition text-sm font-medium cursor-pointer">
                                  <Pencil className="w-4 h-4" />
                                  Edit
                                </button>
                              )}
                          </div>
                          {!isEditing ? (
                            <div className="text-md">
                              {session.notes || (
                                <span className="text-muted-foreground italic">No notes added.</span>
                              )}
                            </div>
                          ) : (
                            <div className="flex flex-col gap-3">
                              <textarea
                                value={notesValue}
                                onChange={(e) => setNotesValue(e.target.value)}
                                className="w-full p-4 rounded-2xl border border-default
                                           min-h-[180px] bg-white-hover text-md focus:outline-none
                                           focus:ring-2 focus:ring-primary"
                                placeholder="Add session notes here..."
                                title="Session Notes"
                              />
                            </div>
                          )}
                      </div>
                  </div>
              </div>
          </div>
      </AppLayout>
  )
}