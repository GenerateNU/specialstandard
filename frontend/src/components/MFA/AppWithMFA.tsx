import { useEffect, useState } from "react";
import { AuthMFA } from "./AuthMFA";
import { supabase } from "./EnrollMFA";

export function AppWithMFA({ children }: { children: React.ReactNode }) {
  const [readyToShow, setReadyToShow] = useState(false);
  const [showMFAScreen, setShowMFAScreen] = useState(false);

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

  if (readyToShow) {
    if (showMFAScreen) {
      return <AuthMFA />;
    }
  }
  return <>{children}</>;
}

export default AppWithMFA;
