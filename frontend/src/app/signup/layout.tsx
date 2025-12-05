"use client";
import { Check } from "lucide-react";
import Image from "next/image";
import { usePathname } from "next/navigation";

const steps = [
  { id: "welcome", path: "/signup/welcome" },
  { id: "link", path: "/signup/link" },
  { id: "students", path: "/signup/students" },
  { id: "sessions", path: "/signup/sessions" },
];

export default function OnboardingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();

  // Determine current step index
  const getCurrentStepIndex = () => {
    // Handle sub-routes like students-add, students-edit
    const basePath = pathname.split("/").slice(0, 3).join("/");
    const stepIndex = steps.findIndex((step) => step.path === basePath);

    // If exact match not found, check if it's a sub-route
    if (stepIndex === -1) {
      const pathSegment = pathname.split("/")[2];
      if (pathSegment?.startsWith("students")) return 3;
      if (pathSegment?.startsWith("sessions")) return 4;
      if (pathSegment === "complete") return steps.length;
    }

    return stepIndex;
  };

  const currentStepIndex = getCurrentStepIndex();

  return (
    <div className="min-h-screen flex bg-background">
      {/* Left Sidebar Progress */}
      <div
        className="fixed left-0 top-0 w-64 border-default p-8 flex flex-col bg-black border-r border-default 
        transition-all duration-300 ease-in-out shadow-md h-screen z-50 lg:z-auto text-center items-center"
      >
        <div className="mb-12 bg-transparent rounded-md p-1 w-fit">
          <Image
            src="/littleguy.png"
            alt="The Special Standard"
            width={140}
            height={30}
            priority
          />
        </div>

        <div className="flex-1">
          <div className="space-y-1">
            {steps.map((step, index) => {
              const isCompleted = index < currentStepIndex;
              const isCurrent = index === currentStepIndex;
              const isUpcoming = index > currentStepIndex;

              return (
                <div key={step.id} className="flex items-start">
                  <div className="flex flex-col items-center">
                    {/* Step indicator */}
                    <div
                      className={`w-12 h-12 rounded-full flex items-center justify-center text-sm font-medium
                      ${isCompleted ? "bg-accent border-2 border-gray-200 text-white" : ""}
                      ${isCurrent ? "bg-accent text-amber-400 border-2 border-amber-400" : ""}
                      ${isUpcoming ? "border-gray-200 border-2 text-white" : ""}
                      ${!isCompleted && !isCurrent && !isUpcoming ? "text-primary" : ""}
                    `}
                    >
                      {isCompleted ? (
                        <Check className="w-4 h-4" />
                      ) : (
                        <span>{index + 1}</span>
                      )}
                    </div>

                    {/* Connector line */}
                    {index < steps.length - 1 && (
                      <div
                        className={`
                        w-0.5 h-16 mt-2
                        ${isCompleted ? "bg-accent" : "bg-gray-200"}
                      `}
                      />
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 ml-64 overflow-y-auto min-h-screen">
        {children}
      </div>
    </div>
  );
}
