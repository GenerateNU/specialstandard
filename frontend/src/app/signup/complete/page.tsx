"use client";

import { EmailMFAVerification } from "@/components/ui/MFA";
import { useAuthContext } from "@/contexts/authContext";
import { CheckCircle, Loader2 } from "lucide-react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function CompletePage() {
  const router = useRouter();
  const [userName, setUserName] = useState("");
  const [needsMFA, setNeedsMFA] = useState(true);

  const { completeMFALogin } = useAuthContext();

  useEffect(() => {
    // Get user name from localStorage
    const onboardingData = localStorage.getItem("onboardingData");
    if (onboardingData) {
      const data = JSON.parse(onboardingData);
      setUserName(data.firstName || "");
    }
  }, []);

  const handleMFAVerified = () => {
    completeMFALogin();
    setNeedsMFA(false);

    // Clean up onboarding data
    localStorage.removeItem("onboardingData");
    localStorage.removeItem("therapistProfile");
    localStorage.removeItem("onboardingStudents");
    localStorage.removeItem("onboardingSessions");

    // Redirect immediately
    router.push("/");
  };

  // Show MFA verification first
  if (needsMFA) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background px-4">
        <EmailMFAVerification onVerified={handleMFAVerified} />
      </div>
    );
  }

  // Success screen (brief loading state while redirecting)
  return (
    <div className="flex items-center justify-center min-h-screen p-8 bg-background">
      <div className="max-w-lg w-full bg-background rounded-2xl shadow-lg p-12 text-center">
        <div className="mb-8 flex justify-center">
          <div className="w-24 h-24 bg-accent-light rounded-full flex items-center justify-center">
            <CheckCircle className="w-12 h-12 text-accent" />
          </div>
        </div>

        <h1 className="text-3xl font-bold text-primary mb-4">
          {userName ? `${userName}, you're all set!` : "You're all set!"}
        </h1>

        <p className="text-secondary mb-8">
          Your account is ready to go. You can head straight to your dashboard
          and view your sessions on your schedule!
        </p>

        <div className="mb-12 bg-white rounded-md p-1 w-fit mx-auto">
          <Image
            src="/tss.png"
            alt="The Special Standard"
            width={140}
            height={30}
            priority
          />
        </div>

        <Loader2 className="w-6 h-6 animate-spin text-primary mx-auto" />
      </div>
    </div>
  );
}
