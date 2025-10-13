"use client";

import {
  ArrowRight,
  Edit,
  LogOut,
  Plus,
  Settings,
  Trash,
  User,
  Users,
} from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { useState } from "react";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { useForm } from "react-hook-form";
// Import all components
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
// New components
import { Checkbox } from "@/components/ui/checkbox";
import CustomAlert from "@/components/ui/CustomAlert";
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
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";

export default function ComponentShowcase() {
  const [showAlert, setShowAlert] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState("");
  const [selectedAction, setSelectedAction] = useState("");
  const [checkboxChecked, setCheckboxChecked] = useState(false);

  // Form setup
  const form = useForm({
    defaultValues: {
      username: "",
      email: "",
      message: "",
      terms: false,
    },
  });

  const onSubmit = (data: any) => {
    console.warn("Form submitted:", data);
  };

  return (
    <div className="font-sans min-h-screen p-8 pb-20 sm:p-20">
      {/* Header */}
      <header className="mb-12">
        <Image
          src="/tss.png"
          alt="The Special Standard logo"
          width={180}
          height={38}
          priority
          className="mb-8"
        />
        <h1 className="text-4xl font-bold mb-4 tracking-tight text-primary">
          Component Showcase
        </h1>
        <p className="text-secondary">
          A comprehensive display of all UI components with the custom styling
          system
        </p>
      </header>

      <main className="space-y-16 max-w-6xl">
        {/* Buttons Section */}
        <section>
          <h2 className="text-2xl font-semibold mb-6 text-primary">Buttons</h2>

          <div className="space-y-6">
            <div>
              <h3 className="text-lg font-medium mb-3 text-secondary">
                Variants
              </h3>
              <div className="flex flex-wrap gap-3">
                <Button>Default</Button>
                <Button variant="destructive">Destructive</Button>
                <Button variant="outline">Outline</Button>
                <Button variant="secondary">Secondary</Button>
                <Button variant="ghost">Ghost</Button>
                <Button variant="link">Link</Button>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-medium mb-3 text-secondary">Sizes</h3>
              <div className="flex flex-wrap items-center gap-3">
                <Button size="sm">Small</Button>
                <Button>Default</Button>
                <Button size="lg">Large</Button>
                <Button size="icon">
                  <Plus className="h-4 w-4" />
                </Button>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-medium mb-3 text-secondary">
                States
              </h3>
              <div className="flex flex-wrap gap-3">
                <Button disabled>Disabled</Button>
                <Button variant="outline" disabled>
                  Disabled Outline
                </Button>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-medium mb-3 text-secondary">
                Tab Buttons
              </h3>
              <div className="flex gap-0">
                <Button variant="tab" active={true}>
                  Active Tab
                </Button>
                <Button variant="tab" active={false}>
                  Inactive Tab
                </Button>
                <Button variant="tab" active={false}>
                  Another Tab
                </Button>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-medium mb-3 text-secondary">
                Full Width
              </h3>
              <Button size="long">Full Width Button</Button>
            </div>
          </div>
        </section>

        {/* Form Controls Section */}
        <section>
          <h2 className="text-2xl font-semibold mb-6 text-primary">
            Form Controls
          </h2>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="space-y-4">
              <div>
                <Label htmlFor="input-example">Input Field</Label>
                <Input id="input-example" placeholder="Enter text here..." />
              </div>

              <div>
                <Label htmlFor="email-example">Email Input</Label>
                <Input
                  id="email-example"
                  type="email"
                  placeholder="email@example.com"
                />
              </div>

              <div>
                <Label htmlFor="disabled-input">Disabled Input</Label>
                <Input
                  id="disabled-input"
                  placeholder="Cannot edit this"
                  disabled
                />
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <Label htmlFor="textarea-example">Textarea</Label>
                <Textarea
                  id="textarea-example"
                  placeholder="Enter your message here..."
                  rows={5}
                />
              </div>
            </div>
          </div>

          <Separator className="my-8" />

          {/* New Form Controls */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="space-y-6">
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="terms"
                  checked={checkboxChecked}
                  onCheckedChange={setCheckboxChecked}
                />
                <Label htmlFor="terms">Accept terms and conditions</Label>
              </div>
            </div>
          </div>
        </section>

        {/* React Hook Form Section */}
        <section>
          <h2 className="text-2xl font-semibold mb-6 text-primary">
            Form with Validation
          </h2>

          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmit)}
              className="space-y-6 max-w-md"
            >
              <FormField
                control={form.control}
                name="username"
                rules={{ required: "Username is required" }}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Username</FormLabel>
                    <FormControl>
                      <Input placeholder="johndoe" {...field} />
                    </FormControl>
                    <FormDescription>
                      This is your public display name.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="email"
                rules={{
                  required: "Email is required",
                  pattern: {
                    value: /^\S[^\s@]*@\S+$/,
                    message: "Invalid email address",
                  },
                }}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input
                        type="email"
                        placeholder="john@example.com"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit">Submit Form</Button>
            </form>
          </Form>
        </section>

        {/* Alerts Section */}
        <section>
          <h2 className="text-2xl font-semibold mb-6 text-primary">Alerts</h2>

          <div className="space-y-4">
            <Alert>
              <AlertTitle>Default Alert</AlertTitle>
              <AlertDescription>
                This is a default alert with some informational content.
              </AlertDescription>
            </Alert>

            {showAlert && (
              <CustomAlert
                variant="success"
                title="Success!"
                description="Your operation completed successfully."
                onClose={() => setShowAlert(false)}
              />
            )}

            <CustomAlert
              variant="warning"
              title="Warning"
              description="Please review your input before proceeding."
            />

            <CustomAlert
              variant="destructive"
              title="Error"
              description="Something went wrong. Please try again."
            />
          </div>
        </section>

        {/* Dropdown Section - Updated */}
        <section>
          <h2 className="text-2xl font-semibold mb-6 text-primary">
            Dropdown (with value tracking)
          </h2>

          <div className="space-y-6">
            <div className="flex gap-4">
              <Dropdown
                value={selectedUser}
                onValueChange={setSelectedUser}
                placeholder="Select user..."
                align="left"
                items={[
                  {
                    label: "My Profile",
                    value: "profile",
                    icon: <User />,
                  },
                  {
                    label: "Settings",
                    value: "settings",
                    icon: <Settings />,
                  },
                  {
                    label: "Logout",
                    value: "logout",
                    icon: <LogOut />,
                  },
                ]}
              />

              <Dropdown
                value={selectedAction}
                onValueChange={setSelectedAction}
                placeholder="Choose action..."
                align="right"
                items={[
                  {
                    label: "Edit",
                    value: "edit",
                    icon: <Edit />,
                  },
                  {
                    label: "Delete",
                    value: "delete",
                    icon: <Trash />,
                    disabled: true,
                  },
                ]}
              />
            </div>

            {(selectedUser || selectedAction) && (
              <div className="text-sm text-muted">
                Selected values: User = "{selectedUser || "none"}
                ", Action = "{selectedAction || "none"}"
              </div>
            )}
          </div>
        </section>

        {/* Badges Section */}
        <section>
          <h2 className="text-2xl font-semibold mb-6 text-primary">Badges</h2>

          <div className="flex flex-wrap gap-2">
            <Badge>Default</Badge>
            <Badge variant="secondary">Secondary</Badge>
            <Badge variant="destructive">Destructive</Badge>
            <Badge variant="outline">Outline</Badge>
            <Badge variant="success">Success</Badge>
            <Badge variant="warning">Warning</Badge>
          </div>
        </section>

        {/* Dialog Section */}
        <section>
          <h2 className="text-2xl font-semibold mb-6 text-primary">Dialog</h2>

          <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
            <DialogTrigger asChild>
              <Button>Open Dialog</Button>
            </DialogTrigger>
            <DialogContent className="max-w-md">
              <DialogHeader>
                <DialogTitle>Edit Profile</DialogTitle>
                <DialogDescription>
                  Make changes to your profile here. Click save when you're
                  done.
                </DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="name" className="text-right">
                    Name
                  </Label>
                  <Input
                    id="name"
                    defaultValue="Pedro Duarte"
                    className="col-span-3"
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="username" className="text-right">
                    Username
                  </Label>
                  <Input
                    id="username"
                    defaultValue="@peduarte"
                    className="col-span-3"
                  />
                </div>
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setDialogOpen(false)}>
                  Cancel
                </Button>
                <Button onClick={() => setDialogOpen(false)}>
                  Save changes
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </section>

        {/* Card Examples Section */}
        <section>
          <h2 className="text-2xl font-semibold mb-6 text-primary">Cards</h2>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 max-w-4xl">
            {/* Custom card from original design */}
            <Link
              href="/students"
              className="group p-6 bg-card rounded-xl shadow-lg hover:shadow-xl transition-all duration-200 border border-border hover:bg-card-hover hover:border-hover"
            >
              <div className="flex items-center justify-between mb-3">
                <Users className="w-10 h-10 text-accent" />
                <ArrowRight className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform group-hover:text-accent" />
              </div>
              <h3 className="text-xl font-semibold text-primary mb-2">
                View Students
              </h3>
              <p className="text-secondary text-sm">
                Browse and manage all student records in the system
              </p>
            </Link>

            {/* Card component example */}
            <Card>
              <CardHeader>
                <CardTitle>Team Member</CardTitle>
                <CardDescription>Active team member profile</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex items-center space-x-4">
                  <Avatar>
                    <AvatarImage src="https://github.com/shadcn.png" />
                    <AvatarFallback>JD</AvatarFallback>
                  </Avatar>
                  <div className="space-y-1">
                    <p className="text-sm font-medium leading-none">John Doe</p>
                    <p className="text-sm text-muted">john@example.com</p>
                  </div>
                </div>
              </CardContent>
              <CardFooter className="flex gap-2">
                <Button variant="outline" size="sm">
                  View Profile
                </Button>
                <Button size="sm">Message</Button>
              </CardFooter>
            </Card>
          </div>
        </section>
      </main>
    </div>
  );
}
