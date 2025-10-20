"use client"
import { useAuthContext } from "@/contexts/authContext";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { AuthMFA } from "./AuthMFA";
import { EnrollMFA, supabase } from "./EnrollMFA";

export function AppWithMFA({ children }: { children: React.ReactNode }) {
  const [readyToShow, setReadyToShow] = useState(false);
  const [showMFAScreen, setShowMFAScreen] = useState(false);
  const { showMFAEnroll, setShowMFAEnroll } = useAuthContext();
  const router = useRouter();

  useEffect(() => {
    (async () => {
      try {
        const { data, error } =
          await supabase.auth.mfa.getAuthenticatorAssuranceLevel();
        if (error) {
          throw error;
        }

        console.log(data);

        if (data.nextLevel === "aal2" && data.nextLevel !== data.currentLevel) {
          setShowMFAScreen(true);
        }
      } finally {
        setReadyToShow(true);
      }
    })();
  }, []);

  if (showMFAEnroll) {
    return (
      <EnrollMFA
        onEnrolled={() => {
          setShowMFAEnroll(false);
          router.push('/students');
        }}
        onCancelled={() => {
          setShowMFAEnroll(false);
          router.push('/students');
        }}
      />
    );
  }

  if (readyToShow) {
    if (showMFAScreen) {
      return <AuthMFA />;
    }
  }
  return <>{children}</>;
}

export default AppWithMFA;
