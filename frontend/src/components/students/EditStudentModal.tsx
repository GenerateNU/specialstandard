"use client";

import type { UpdateStudentInput } from "@/lib/api/theSpecialStandardAPI.schemas";

import { Building2, Calendar, FileText, GraduationCap, PencilLine } from "lucide-react";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Dropdown } from "@/components/ui/dropdown";
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
import { useStudents, type StudentBody } from "@/hooks/useStudents";
import { useSchools } from "@/hooks/useSchools";
import { gradeOptions } from "@/lib/gradeUtils";
import { useForm } from "react-hook-form";

interface EditStudentModalProps {
  student: StudentBody;
  trigger?: React.ReactNode;
  onSuccess?: () => void;
}

interface EditStudentFormInput {
  first_name: string;
  last_name: string;
  dob?: string;
  grade?: string;
  iep?: string;
  school_id?: string;
}

// Format date string to YYYY-MM-DD for HTML date input
function formatDateForInput(dateString?: string | null): string {
  if (!dateString) return "";
  try {
    const date = new Date(dateString);
    if (Number.isNaN(date.getTime())) return "";
    return date.toISOString().split("T")[0];
  } catch {
    return "";
  }
}

export default function EditStudentModal({ student, trigger, onSuccess }: EditStudentModalProps) {
  const [open, setOpen] = useState(false);
  const { updateStudent } = useStudents();
  const { schools, isLoading: isLoadingSchools, error: schoolsError } = useSchools();

  const form = useForm<EditStudentFormInput>({
    defaultValues: {
      first_name: student.first_name,
      last_name: student.last_name,
      dob: formatDateForInput(student.dob),
      grade: student.grade || "",
      iep: student.iep?.join("\n") || "",
      school_id: student.school_id?.toString() || "",
    },
  });

  // Reset form when student changes or modal opens
  useEffect(() => {
    if (open) {
      form.reset({
        first_name: student.first_name,
        last_name: student.last_name,
        dob: formatDateForInput(student.dob),
        grade: student.grade || "",
        iep: student.iep?.join("\n") || "",
        school_id: student.school_id?.toString() || "",
      });
    }
  }, [open, student, form]);

  const onSubmit = async (data: EditStudentFormInput) => {
    try {
      // Convert form data to update format
      const updateData: UpdateStudentInput = {
        first_name: data.first_name,
        last_name: data.last_name,
        dob: data.dob || null,
        // Convert grade display value to storage format
        grade: data.grade
          ? data.grade.toUpperCase() === "K"
            ? 0
            : Number.parseInt(data.grade)
          : null,
        // Convert IEP string to array
        iep: data.iep
          ? data.iep.split("\n").map((goal: string) => goal.trim()).filter((goal: string) => goal)
          : null,
        school_id: data.school_id ? Number.parseInt(data.school_id) : undefined,
      };

      // Update the student
      updateStudent(student.id, updateData);

      setOpen(false);
      onSuccess?.();
    } catch (error) {
      console.error("Error updating student:", error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger || (
          <Button variant="secondary" className="flex items-center gap-2 hover:bg-accent h-8" size="sm">
            <span className="text-base font-medium">Edit Profile</span>
            <PencilLine size={18} />
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="max-w-md max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <PencilLine className="w-5 h-5 text-accent" />
            Edit Student Profile
          </DialogTitle>
          <DialogDescription>
            Update {student.first_name}'s profile information.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="first_name"
                rules={{ required: "First name is required" }}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>First Name *</FormLabel>
                    <FormControl>
                      <Input placeholder="John" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="last_name"
                rules={{ required: "Last name is required" }}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Last Name *</FormLabel>
                    <FormControl>
                      <Input placeholder="Doe" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="dob"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="flex items-center gap-1">
                    <Calendar className="w-4 h-4" />
                    Date of Birth
                  </FormLabel>
                  <FormControl>
                    <Input type="date" {...field} value={field.value ?? ""} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="school_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="flex items-center gap-1">
                    <Building2 className="w-4 h-4" />
                    School
                  </FormLabel>
                  <FormControl>
                    <Dropdown
                      value={field.value || ""}
                      onValueChange={(value) => field.onChange(value)}
                      placeholder={
                        isLoadingSchools
                          ? "Loading schools..."
                          : schoolsError
                            ? "Error loading schools"
                            : schools.length === 0
                              ? "No schools available"
                              : "Select school..."
                      }
                      items={schools.map((school) => ({
                        label: school.name,
                        value: school.id.toString(),
                      }))}
                      className="w-full justify-between"
                    />
                  </FormControl>
                  {schoolsError && (
                    <p className="text-sm text-error mt-1">
                      Unable to load schools. You may need to contact support.
                    </p>
                  )}
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="grade"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="flex items-center gap-1">
                    <GraduationCap className="w-4 h-4" />
                    Grade Level
                  </FormLabel>
                  <FormControl>
                    <Dropdown
                      value={field.value || ""}
                      onValueChange={(value) => field.onChange(value)}
                      placeholder="Select grade level..."
                      items={gradeOptions}
                      className="w-full justify-between"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="iep"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="flex items-center gap-1">
                    <FileText className="w-4 h-4" />
                    IEP Goals
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Enter IEP goals (one per line)&#10;Example:&#10;Improve articulation of /r/ sound&#10;Increase expressive vocabulary by 20 words"
                      rows={3}
                      {...field}
                      value={field.value ?? ""}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter className="gap-2 sm:gap-0">
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  form.reset();
                  setOpen(false);
                }}
              >
                Cancel
              </Button>
              <Button type="submit">Save Changes</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

