"use client";

import type { CreateStudentInput } from "@/lib/api/theSpecialStandardAPI.schemas";

import { Building2, Calendar, FileText, GraduationCap, User } from "lucide-react";
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
import { useSchools } from "@/hooks/useSchools";
import { gradeOptions, gradeToStorage } from "@/lib/gradeUtils";
import { useForm } from "react-hook-form";
import { useAuthContext } from "@/contexts/authContext";

interface AddStudentModalProps {
  trigger?: React.ReactNode;
}

export default function AddStudentModal({ trigger }: AddStudentModalProps) {
  const [open, setOpen] = useState(false);
  const { addStudent } = useStudents();
  const { userId, therapistProfile } = useAuthContext();
  const { schools, isLoading: isLoadingSchools, error: schoolsError } = useSchools();

  // 1. Get the list of school IDs the therapist is associated with
  const therapistSchoolIds = therapistProfile?.schools || [];

  // 2. Filter the schools to only include those the therapist is allowed to see
  const availableSchools = schools.filter(school => 
    therapistSchoolIds.length === 0 || therapistSchoolIds.includes(Number(school.id))
  );
  
  // Update the dropdown items to use the filtered list
  const schoolOptions = availableSchools.map(school => ({
    label: school.name,
    value: school.id.toString(),
  }));

  type CreateStudentFormInput = Omit<CreateStudentInput, 'grade' | 'iep' | 'school_id'> & {
    grade?: string
    iep?: string
    school_id?: string // Keep as string for form
  }

  const form = useForm<CreateStudentFormInput>({
    defaultValues: {
      first_name: "",
      last_name: "",
      dob: "",
      therapist_id: userId ?? undefined,
      grade: "",
      iep: "",
      school_id: "",
    },
  });
  
  // --- START FIX: Centralize the closing logic ---
  const handleOpenChange = (newOpenState: boolean) => {
    setOpen(newOpenState);
    
    // When the dialog is closing (newOpenState is false), reset the form errors and values.
    if (!newOpenState) {
      form.reset();
    }
  };
  // --- END FIX ---

  const onSubmit = async (data: CreateStudentFormInput) => {
    try {
      // Convert frontend data format to backend-expected format
      const backendData = {
        first_name: data.first_name,
        last_name: data.last_name,
        dob: data.dob || undefined,
        therapist_id: userId ?? undefined,
        // Only convert/include grade if a value exists
        grade: data.grade ? gradeToStorage(data.grade) : undefined,
        // Convert IEP string to array...
        iep: data.iep ? data.iep.split('\n').map((goal: string) => goal.trim()).filter((goal: string) => goal) : undefined,
        school_id: data.school_id ? Number.parseInt(data.school_id) : undefined,
      };

      // Add the student using the hook with converted data
      addStudent(backendData as any);

      // Reset form and close modal using the shared handler
      handleOpenChange(false);
    } catch (error) {
      console.error("Error adding student:", error);
    }
  };

  // Get today's date in 'YYYY-MM-DD' format for validation
  const today = new Date().toISOString().split('T')[0];

  return (
    // Pass the centralized handler to onOpenChange
    <Dialog open={open} onOpenChange={handleOpenChange}> 
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
              rules={{
                max: {
                  value: today,
                  message: "Date of Birth cannot be in the future.",
                },
              }}
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="flex items-center gap-1">
                    <Calendar className="w-4 h-4" />
                    Date of Birth
                  </FormLabel>
                  <FormControl>
                    <Input 
                      type="date" 
                      max={today} // Added max attribute for UI validation
                      {...field} 
                      value={field.value ?? ""} 
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="school_id"
              rules={{ required: schools.length > 0 ? "School is required" : false }}
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
                            : availableSchools.length === 0 // Use filtered list here
                              ? "No schools available in your district" 
                              : "Select school..."
                      }
                      // Use the filtered schoolOptions here
                      items={schoolOptions} 
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
                  // Use the shared handler here too
                  handleOpenChange(false);
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