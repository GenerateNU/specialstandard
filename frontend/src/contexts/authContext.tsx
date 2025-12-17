"use client";

import { useAuth } from "@/hooks/useAuth";
import { useTherapist } from "@/hooks/useTherapists";
import type { PostAuthLoginBody, PostAuthSignupBody, Therapist } from "@/lib/api/theSpecialStandardAPI.schemas";
import { useRouter } from "next/navigation";
import { createContext, useContext, useEffect, useState } from "react";

interface AuthContextType {
  userId: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  therapistProfile: Therapist | null;
  isProfileLoading: boolean;
  login: (
    credentials: PostAuthLoginBody
  ) => Promise<{ requiresMFA: boolean; userId?: string | null }>;
  completeMFALogin: () => void;
  signup: (credentials: PostAuthSignupBody) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [userId, setUserId] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [pendingMFAAuth, setPendingMFAAuth] = useState<{
    jwt: string;
    userId: string;
  } | null>(null);
  const router = useRouter();
  const { userLogin, userLogout, userSignup } = useAuth();
  
  const { 
    therapist: therapistProfile, 
    isLoading: isProfileLoading, 
  } = useTherapist(userId);

  // Check if user is authenticated on mount
  useEffect(() => {
    const checkAuth = () => {
      const storedUserId = localStorage.getItem("userId");
      const storedJwt = localStorage.getItem("jwt");

      if (storedJwt && storedUserId) {
        setUserId(storedUserId);
      } else {
        setUserId(null);
      }
      setIsLoading(false);
    };
    checkAuth();
  }, []);

  const login = async (credentials: PostAuthLoginBody) => {
    try {
      const response = await userLogin(credentials);
      console.warn("Login response:", response);

      if (!response.access_token || !response.user?.id) {
        throw new Error("Invalid login response");
      }

      // Check if backend indicates MFA is required for this login
      // When EMAIL_VERIFICATION_ENABLED=false, backend returns requires_mfa: false
      const requiresMFA = response.needs_mfa ?? true; // Default to true for safety
      
      if (requiresMFA) {
        // Store temporarily and require MFA
        localStorage.setItem("temp_jwt", response.access_token);
        localStorage.setItem("temp_userId", response.user.id);

        setPendingMFAAuth({
          jwt: response.access_token,
          userId: response.user.id,
        });

        return { requiresMFA: true, userId: response.user.id };
      } else {
        // MFA not required (verification disabled) - authenticate immediately
        localStorage.setItem("jwt", response.access_token);
        localStorage.setItem("userId", response.user.id);
        setUserId(response.user.id);
        
        return { requiresMFA: false, userId: response.user.id };
      }
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  };

  // Call this after MFA verification succeeds
  const completeMFALogin = () => {
    let finalUserId: string | null = null;
    
    if (pendingMFAAuth) {
      localStorage.setItem("jwt", pendingMFAAuth.jwt);
      localStorage.setItem("userId", pendingMFAAuth.userId);

      // Clean up temp storage
      localStorage.removeItem("temp_jwt");
      localStorage.removeItem("temp_userId");

      finalUserId = pendingMFAAuth.userId;
      setPendingMFAAuth(null);
    } else {
      // Fallback: if pendingMFAAuth is null but we have temp credentials
      const tempJwt = localStorage.getItem("temp_jwt");
      const tempUserId = localStorage.getItem("temp_userId");

      if (tempJwt && tempUserId) {
        localStorage.setItem("jwt", tempJwt);
        localStorage.setItem("userId", tempUserId);

        // Clean up temp storage
        localStorage.removeItem("temp_jwt");
        localStorage.removeItem("temp_userId");

        finalUserId = tempUserId;
      }
    }
    
    if (finalUserId) {
      setUserId(finalUserId);
    }
  };

  const signup = async (credentials: PostAuthSignupBody) => {
    try {
      const response = await userSignup(credentials);

      if (!response.access_token || !response.user?.id) {
        throw new Error("Invalid signup response");
      }

      // Check if MFA is required
      const requiresMFA = response.needs_mfa ?? true; // Default to true for safety

      // Store temp credentials
      localStorage.setItem("temp_jwt", response.access_token);
      localStorage.setItem("temp_userId", response.user.id);
      localStorage.setItem("signup_requires_mfa", String(requiresMFA)); // Store MFA requirement

      setPendingMFAAuth({
        jwt: response.access_token,
        userId: response.user.id,
      });

      router.push("/signup/link");
    } catch (error) {
      console.error("Signup failed:", error);
      throw error;
    }
  };

  const logout = () => {
    // Clear localStorage
    localStorage.removeItem("jwt");
    localStorage.removeItem("userId");
    localStorage.removeItem("temp_jwt");
    localStorage.removeItem("temp_userId");
    localStorage.removeItem("recentlyViewedStudents");
    localStorage.removeItem("signup_requires_mfa"); // Clean up signup MFA flag

    setPendingMFAAuth(null);

    // Clear any remaining cookies 
    document.cookie.split(";").forEach((cookie) => {
      const eqPos = cookie.indexOf("=");
      const name =
        eqPos > -1 ? cookie.substring(0, eqPos).trim() : cookie.trim();
      document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`;
    });

    userLogout();
    setUserId(null);
    router.push("/login");
  };

  return (
    <AuthContext.Provider
      value={{
        userId,
        isAuthenticated: !!userId,
        isLoading,
        therapistProfile,
        isProfileLoading,
        login,
        completeMFALogin,
        signup,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuthContext() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}