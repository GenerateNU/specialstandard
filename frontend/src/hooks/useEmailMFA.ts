import { useState } from "react";

export interface SendVerificationResponse {
  success: boolean;
  messageId?: string;
  error?: string;
}

export interface VerifyCodeResponse {
  success: boolean;
  verified?: boolean;
  error?: string;
}

export const useEmailMFA = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [codeSent, setCodeSent] = useState(false);

  const sendVerificationCode = async (): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      // Get the temp JWT from localStorage
      const jwt =
        localStorage.getItem("temp_jwt") || localStorage.getItem("jwt");

      const response = await fetch("/api/sendVerification", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${jwt}`, // Add JWT to header
        },
      });

      const data: SendVerificationResponse = await response.json();

      if (!response.ok || !data.success) {
        throw new Error(data.error || "Failed to send verification code");
      }

      setCodeSent(true);
      return true;
    } catch (err) {
      const message = err instanceof Error ? err.message : "An error occurred";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  };

  const verifyCode = async (code: string): Promise<boolean> => {
    if (!code || code.length !== 6) {
      setError("Please enter a valid 6-digit code");
      return false;
    }

    setLoading(true);
    setError(null);

    try {
      // Get the temp JWT from localStorage
      const jwt =
        localStorage.getItem("temp_jwt") || localStorage.getItem("jwt");

      const response = await fetch("/api/verifyCode", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${jwt}`,
        },
        body: JSON.stringify({ code }),
      });

      const data: VerifyCodeResponse = await response.json();

      if (!response.ok || !data.success) {
        throw new Error(data.error || "Verification failed");
      }

      return data.verified || false;
    } catch (err) {
      const message = err instanceof Error ? err.message : "An error occurred";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  };

  const resetState = () => {
    setCodeSent(false);
    setError(null);
  };

  return {
    sendVerificationCode,
    verifyCode,
    resetState,
    loading,
    error,
    codeSent,
  };
};
