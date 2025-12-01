import { useEmailMFA } from "@/hooks/useEmailMFA";
import React, { useState } from "react";

interface EmailMFAVerificationProps {
  onVerified: () => void;
  userId: string;
}

export const EmailMFAVerification: React.FC<EmailMFAVerificationProps> = ({
  onVerified,
  userId,
}) => {
  const [code, setCode] = useState("");
  const { sendVerificationCode, verifyCode, loading, error, codeSent } =
    useEmailMFA(userId);

  const handleSendCode = async () => {
    const success = await sendVerificationCode();
    if (!success) {
      // we do this because hook will throw error
    }
  };

  const handleVerifyCode = async (e: React.FormEvent) => {
    e.preventDefault();
    const verified = await verifyCode(code);
    if (verified) {
      onVerified();
    }
  };

  const handleCodeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/\D/g, "").slice(0, 6);
    setCode(value);
  };

  if (!codeSent) {
    return (
      <div className="p-6 max-w-md mx-auto bg-white rounded-lg shadow-md">
        <h2 className="text-2xl font-bold mb-4 text-black">
          Email Verification Required
        </h2>
        <p className="mb-6 text-gray-600">
          To continue, we need to verify your identity. Click below to receive a
          verification code via email.
        </p>

        {error && (
          <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        <div className="flex gap-3">
          <button
            onClick={handleSendCode}
            disabled={loading}
            className="flex-1 bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? "Sending..." : "Send Verification Code"}
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 max-w-md mx-auto bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4 text-black">
        Enter Verification Code
      </h2>
      <p className="mb-6 text-gray-600">
        We've sent a 6-digit code to your email. Please enter it below.
      </p>

      {error && (
        <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
          {error}
        </div>
      )}

      <form onSubmit={handleVerifyCode}>
        <input
          type="text"
          value={code}
          onChange={handleCodeChange}
          placeholder="000000"
          className="w-full text-center text-3xl text-gray-600 tracking-widest p-4 border-2 border-gray-300 rounded mb-4 focus:outline-none focus:border-blue-600"
          maxLength={6}
          autoComplete="off"
          autoFocus
        />

        <div className="flex gap-3">
          <button
            type="submit"
            disabled={loading || code.length !== 6}
            className="flex-1 bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? "Verifying..." : "Verify Code"}
          </button>

          <button
            type="button"
            onClick={handleSendCode}
            disabled={loading}
            className="px-4 py-2 border border-gray-300 text-gray-600 rounded hover:bg-gray-50"
          >
            Resend Code
          </button>
        </div>
      </form>
    </div>
  );
};
