import type {
  PostVerificationSendCodeBody,
  PostVerificationVerifyBody,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import { getVerification } from "@/lib/api/verification";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";

export const useEmailMFA = (userId: string) => {
  const [codeSent, setCodeSent] = useState(false);
  const api = getVerification();

  const sendCodeMutation = useMutation({
    mutationFn: (body: PostVerificationSendCodeBody) =>
      api.postVerificationSendCode(body),
    onSuccess: () => {
      setCodeSent(true);
    },
  });

  const verifyCodeMutation = useMutation({
    mutationFn: (body: PostVerificationVerifyBody) =>
      api.postVerificationVerify(body),
  });

  const sendVerificationCode = async (): Promise<boolean> => {
    try {
      await sendCodeMutation.mutateAsync({ user_id: userId });
      return true;
    } catch (err) {
      console.error("Send code error:", err);
      return false;
    }
  };

  const verifyCode = async (code: string): Promise<boolean> => {
    if (!code || code.length !== 6) {
      return false;
    }

    try {
      const response = await verifyCodeMutation.mutateAsync({
        code,
        user_id: userId,
      });
      return response.verified || false;
    } catch (err) {
      console.error("Verify code error:", err);
      return false;
    }
  };

  const resetState = () => {
    setCodeSent(false);
    sendCodeMutation.reset();
    verifyCodeMutation.reset();
  };

  return {
    sendVerificationCode,
    verifyCode,
    resetState,
    loading: sendCodeMutation.isPending || verifyCodeMutation.isPending,
    error:
      sendCodeMutation.error?.message ||
      verifyCodeMutation.error?.message ||
      null,
    codeSent,
  };
};
