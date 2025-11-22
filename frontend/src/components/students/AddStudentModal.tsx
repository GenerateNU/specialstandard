"use client";

import type { CreateStudentInput } from "@/lib/api/theSpecialStandardAPI.schemas";

import { Calendar, FileText, GraduationCap, User } from "lucide-react";
import { useState } from "react";

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
import { useStudents } from "@/hooks/useStudents";
import { gradeOptions, gradeToStorage } from "@/lib/gradeUtils";
import { useForm } from "react-hook-form";
import { useAuthContext } from "@/contexts/authContext";

interface AddStudentModalProps {
  trigger?: React.ReactNode;
}

export default function AddStudentModal({ trigger }: AddStudentModalProps) {
  const [open, setOpen] = useState(false);
  const { addStudent } = useStudents();
  const { userId } = useAuthContext();

  type CreateStudentFormInput = Omit<CreateStudentInput, 'grade' | 'iep'> & {
    grade?: string
    iep?: string // Keep as string for textarea input
  }

  const form = useForm<CreateStudentFormInput>({
    defaultValues: {
      first_name: "",
      last_name: "",
      dob: "",
      therapist_id: userId ?? undefined, // use auth hook to get ID?, change hard codedd\
      grade: "",
      iep: "",
    },
  });

  const onSubmit = async (data: CreateStudentFormInput) => {
    try {
      // Convert frontend data format to backend-expected format
      const backendData = {
        first_name: data.first_name,
        last_name: data.last_name,
        dob: data.dob || undefined,
        therapist_id: userId ?? undefined, // Use the proper UUID
        grade: gradeToStorage(data.grade ?? ""), // Convert K to 0, numbers to numbers
        // Convert IEP string to array (split by newlines for multiple goals)
        iep: data.iep ? data.iep.split('\n').map((goal: string) => goal.trim()).filter((goal: string) => goal) : undefined,
        school_id: 1, // TODO: Get this from context or props
      };

      // Add the student using the hook with converted data
      addStudent(backendData as any);

      // Reset form and close modal
      form.reset();
      setOpen(false);
    } catch (error) {
      console.error("Error adding student:", error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger || (
          <Button className="flex items-center gap-2">
            <User className="w-4 h-4" />
            Add Student
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <User className="w-5 h-5 text-accent" />
            Add New Student
          </DialogTitle>
          <DialogDescription>
            Add a new student to your roster. All fields marked with * are
            required.
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
                      onValueChange={(value) => field.onChange(value)} // no conversion here
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
              <Button type="submit">Add Student</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
