"use client";

import { EmailMFAVerification } from "@/components/ui/MFA";
import { useAuthContext } from "@/contexts/authContext";
import { Loader2 } from "lucide-react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function CompletePage() {
  const router = useRouter();
  const [userName, setUserName] = useState("");
  const [needsMFA, setNeedsMFA] = useState(true);
  const [tempUserId, setTempUserId] = useState<string | null>(null);

  const { completeMFALogin } = useAuthContext();

  useEffect(() => {
    // Get user name from localStorage
    const onboardingData = localStorage.getItem("onboardingData");
    if (onboardingData) {
      const data = JSON.parse(onboardingData);
      setUserName(data.firstName || "");
    }

    // Get temp userId from localStorage (stored during signup)
    const storedTempUserId = localStorage.getItem("temp_userId");
    if (storedTempUserId) {
      setTempUserId(storedTempUserId);
    } else {
      // If no temp userId, something went wrong - redirect to signup
      console.error("No temp userId found");
      router.push("/signup");
    }
  }, [router]);

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

  // Show loading if we don't have userId yet
  if (!tempUserId) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  // Show MFA verification first
  if (needsMFA) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background px-4">
        <EmailMFAVerification
          onVerified={handleMFAVerified}
          userId={tempUserId}
        />
      </div>
    );
  }

  // Success screen (brief loading state while redirecting)
  return (
    <div className="flex items-center justify-center min-h-screen p-8 bg-background">
      <div className="max-w-lg w-full bg-background rounded-2xl shadow-lg p-12 text-center">
        <h1 className="text-3xl font-bold text-primary mb-4">
          {userName ? `${userName}, you're all set!` : "You're all set!"}
        </h1>

        <p className="text-secondary mb-8">
          Your account is ready to go. You can head straight to your dashboard
          and view your sessions on your schedule!
        </p>

        <div className="mb-12 bg-transparent rounded-md p-1 w-fit mx-auto">
          <Image
            src="/littlemegaphone.png"
            alt="The Special Standard"
            width={180}
            height={50}
            priority
          />
        </div>

        <Loader2 className="w-6 h-6 animate-spin text-primary mx-auto" />
      </div>
    </div>
  );
}
