// components/CreateSessionDialog.tsx

import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { MultiSelect } from "../ui/multiselect";
import { RepetitionSelector } from "./repetition-selector";
import {
  AlertCircle,
  Calendar,
  CheckCircle,
  Clock,
  FileText,
  MapPin,
  User,
} from "lucide-react";
import type { StudentBody } from "@/hooks/useStudents";
import type {
  PostSessionsBody,
  Session,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import type { RepetitionConfig } from "./repetition-selector";
import { formatRecurrence, useSessions } from "@/hooks/useSessions";

interface CreateSessionDialogProps {
  open: boolean;
  setOpen: (open: boolean) => void;
  therapistId: string;
  students?: Array<StudentBody>;
  onSubmit?: (data: PostSessionsBody) => Promise<void>;
  initialDateTime?: { start: Date; end: Date };
}

export function CreateSessionDialog({
  open,
  setOpen,
  therapistId,
  students = [],
  onSubmit,
  initialDateTime,
}: CreateSessionDialogProps) {
  const form = useForm<{
    session_name: string;
    student_ids: string[];
    sessionDate: string;
    startTime: string;
    endTime: string;
    location: string;
    notes?: string;
  }>({
    defaultValues: {
      session_name: "",
      student_ids: [],
      sessionDate: new Date().toISOString().split("T")[0],
      startTime: "09:00",
      endTime: "10:00",
      location: "",
      notes: "",
    },
  });

  const [repetitionConfig, setRepetitionConfig] = React.useState<
    RepetitionConfig | undefined
  >();
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [submitError, setSubmitError] = React.useState<string | null>(null);
  const [createdSessions, setCreatedSessions] = React.useState<Session[] | null>(
    null
  );
  const { refetch } = useSessions();

  // Reset form when dialog opens/closes
  useEffect(() => {
    if (open && initialDateTime) {
      const startDate = new Date(initialDateTime.start);
      const endDate = new Date(initialDateTime.end);

      const startHours = startDate.getHours().toString().padStart(2, "0");
      const startMinutes = startDate.getMinutes().toString().padStart(2, "0");
      const endHours = endDate.getHours().toString().padStart(2, "0");
      const endMinutes = endDate.getMinutes().toString().padStart(2, "0");

      const sessionDate = startDate.toISOString().split("T")[0];

      form.reset({
        session_name: "",
        student_ids: [],
        sessionDate,
        startTime: `${startHours}:${startMinutes}`,
        endTime: `${endHours}:${endMinutes}`,
        location: "",
        notes: "",
      });
      setRepetitionConfig(undefined);
      setCreatedSessions(null);
      setSubmitError(null);
    } else if (open && !initialDateTime) {
      form.reset({
        session_name: "",
        student_ids: [],
        sessionDate: new Date().toISOString().split("T")[0],
        startTime: "09:00",
        endTime: "10:00",
        location: "",
        notes: "",
      });
      setRepetitionConfig(undefined);
      setCreatedSessions(null);
      setSubmitError(null);
    }
  }, [open, initialDateTime, form]);

  const handleSubmit = async (data: any) => {
    try {
      setSubmitError(null);
      setIsSubmitting(true);
      setCreatedSessions(null);

      // Validation
      if (!data.session_name.trim()) {
        form.setError("session_name", {
          type: "manual",
          message: "Session name is required",
        });
        setIsSubmitting(false);
        return;
      }

      if (data.student_ids.length === 0) {
        form.setError("student_ids", {
          type: "manual",
          message: "Please select at least one student",
        });
        setIsSubmitting(false);
        return;
      }

      const [year, month, day] = data.sessionDate.split('-').map(Number);
      const sessionDate = new Date(year, month - 1, day);

      const [startHour, startMin] = data.startTime.split(":").map(Number);
      const [endHour, endMin] = data.endTime.split(":").map(Number);

      const startDateTime = new Date(sessionDate);
      startDateTime.setHours(startHour, startMin, 0, 0);

      const endDateTime = new Date(sessionDate);
      endDateTime.setHours(endHour, endMin, 0, 0);

      if (endDateTime <= startDateTime) {
        setSubmitError("End time must be after start time");
        setIsSubmitting(false);
        return;
      }

      const postBody: PostSessionsBody = {
        session_name: data.session_name,
        start_datetime: startDateTime.toISOString(),
        end_datetime: endDateTime.toISOString(),
        therapist_id: therapistId,
        notes: data.notes || undefined,
        location: data.location || undefined,
        student_ids: data.student_ids,
        repetition: repetitionConfig,
      };

      if (onSubmit) {
        await onSubmit(postBody);
        form.reset();
        setRepetitionConfig(undefined);
        refetch();
        setOpen(false);
      }
    } catch (error) {
      setSubmitError(
        error instanceof Error ? error.message : "Failed to create session"
      );
      console.error("Error creating session:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    form.reset();
    setRepetitionConfig(undefined);
    setCreatedSessions(null);
    setSubmitError(null);
    setOpen(false);
  };

  // Success screen after session creation
  if (createdSessions && createdSessions.length > 0) {
    const firstSession = createdSessions[0];
    const isRecurring = createdSessions.length > 1;

    return (
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="w-full max-w-2xl">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <CheckCircle className="w-6 h-6 text-green-600" />
              Sessions Created Successfully
            </DialogTitle>
          </DialogHeader>

          <div className="space-y-4">
            {/* Session Summary */}
            <div className="border rounded-lg p-4 bg-slate-50 space-y-3">
              <div>
                <p className="font-semibold text-lg">
                  {firstSession.session_name}
                </p>
                <p className="text-sm text-gray-600">
                  {new Date(firstSession.start_datetime).toLocaleString()}
                </p>
              </div>

              {isRecurring && (
                <div className="bg-blue-50 border border-blue-200 rounded p-3 space-y-2">
                  <p className="font-medium text-sm text-blue-900">
                    Recurring Sessions Created
                  </p>
                  <p className="text-sm text-blue-800">
                    <span className="font-semibold">{createdSessions.length}</span>{" "}
                    session{createdSessions.length !== 1 ? "s" : ""} created
                  </p>
                  {firstSession.repitition && (
                    <p className="text-sm text-blue-800">
                      {formatRecurrence(firstSession.repitition)}
                    </p>
                  )}
                </div>
              )}

              <div className="grid grid-cols-2 gap-3 text-sm">
                <div>
                  <p className="text-gray-600">Students</p>
                  <p className="font-medium">
                    {students
                      .filter((s) => form.getValues("student_ids").includes(s.id))
                      .map((s) => `${s.first_name} ${s.last_name}`)
                      .join(", ")}
                  </p>
                </div>
                {firstSession.location && (
                  <div>
                    <p className="text-gray-600">Location</p>
                    <p className="font-medium">{firstSession.location}</p>
                  </div>
                )}
              </div>
            </div>

            {/* Sessions List (if multiple) */}
            {isRecurring && createdSessions.length > 0 && (
              <details className="border rounded-lg p-3">
                <summary className="cursor-pointer font-medium text-sm hover:text-blue-600">
                  View all {createdSessions.length} sessions
                </summary>
                <div className="mt-3 max-h-64 overflow-y-auto space-y-2">
                  {createdSessions.map((session, idx) => (
                    <div
                      key={session.id}
                      className="text-sm p-2 bg-gray-50 rounded flex justify-between items-center"
                    >
                      <span>
                        {idx + 1}. {new Date(session.start_datetime).toLocaleDateString()}{" "}
                        at{" "}
                        {new Date(session.start_datetime).toLocaleTimeString([], {
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                      </span>
                    </div>
                  ))}
                </div>
              </details>
            )}
          </div>

          <DialogFooter className="gap-2 sm:gap-0 pt-4 border-t">
            <Button variant="outline" onClick={handleClose}>
              Close
            </Button>
            <Button
              onClick={() => {
                setCreatedSessions(null);
                form.reset();
                setRepetitionConfig(undefined);
                setSubmitError(null);
              }}
            >
              Create Another
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    );
  }

  // Main form screen
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Create New Session</DialogTitle>
          <DialogDescription>
            Schedule a therapy session with optional recurring schedule
          </DialogDescription>
        </DialogHeader>

        {submitError && (
          <div className="flex gap-2 p-3 bg-red-50 border border-red-200 rounded text-red-700 text-sm">
            <AlertCircle className="w-4 h-4 flex-shrink-0 mt-0.5" />
            <p>{submitError}</p>
          </div>
        )}

        <Form {...form}>
          <div className="space-y-4">
            {/* Basic Info Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Session Name */}
              <FormField
                control={form.control}
                name="session_name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="flex items-center gap-1 text-sm">
                      <FileText className="w-4 h-4" />
                      Session Name *
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder="e.g., Speech Therapy"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Date */}
              <FormField
                control={form.control}
                name="sessionDate"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="flex items-center gap-1 text-sm">
                      <Calendar className="w-4 h-4" />
                      Date *
                    </FormLabel>
                    <FormControl>
                      <Input type="date" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Start Time */}
              <FormField
                control={form.control}
                name="startTime"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="flex items-center gap-1 text-sm">
                      <Clock className="w-4 h-4" />
                      Start Time *
                    </FormLabel>
                    <FormControl>
                      <Input type="time" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* End Time */}
              <FormField
                control={form.control}
                name="endTime"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="flex items-center gap-1 text-sm">
                      <Clock className="w-4 h-4" />
                      End Time *
                    </FormLabel>
                    <FormControl>
                      <Input type="time" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Location */}
              <FormField
                control={form.control}
                name="location"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="flex items-center gap-1 text-sm">
                      <MapPin className="w-4 h-4" />
                      Location
                    </FormLabel>
                    <FormControl>
                      <Input placeholder="e.g., Room 234" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* Students - Full Width */}
            <FormField
              control={form.control}
              name="student_ids"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="flex items-center gap-1 text-sm">
                    <User className="w-4 h-4" />
                    Students *
                  </FormLabel>
                  <FormControl>
                    <MultiSelect
                      options={students.map((student) => ({
                        label: `${student.first_name} ${student.last_name}`,
                        value: student.id,
                      }))}
                      value={field.value}
                      onValueChange={field.onChange}
                      placeholder="Select students for this session"
                      showTags={true}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Notes - Full Width */}
            <FormField
              control={form.control}
              name="notes"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="flex items-center gap-1 text-sm">
                    <FileText className="w-4 h-4" />
                    Session Notes
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Goals, activities, or special considerations..."
                      rows={3}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Repetition Selector - Full Width */}
            <RepetitionSelector
              value={repetitionConfig}
              onChange={setRepetitionConfig}
              sessionDate={form.watch("sessionDate")}
              sessionTime={form.watch("startTime")}
            />

            <DialogFooter className="gap-2 sm:gap-0 pt-4 border-t">
              <Button type="button" variant="outline" onClick={handleClose}>
                Cancel
              </Button>
              <Button
                type="button"
                disabled={isSubmitting}
                onClick={() => handleSubmit(form.getValues())}
              >
                {isSubmitting ? "Creating..." : "Create Session"}
              </Button>
            </DialogFooter>
          </div>
        </Form>
      </DialogContent>
    </Dialog>
  );
}