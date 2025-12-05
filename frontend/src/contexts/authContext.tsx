"use client";

import { useAuth } from "@/hooks/useAuth";
import { useTherapist } from "@/hooks/useTherapists"; // Import the hook to fetch profile
import type { PostAuthLoginBody, PostAuthSignupBody, Therapist } from "@/lib/api/theSpecialStandardAPI.schemas"; // Combined type import
import { useRouter } from "next/navigation";
import { createContext, useContext, useEffect, useState } from "react";

// 1. Update the AuthContextType to include the profile and its loading state
interface AuthContextType {
  userId: string | null;
  isAuthenticated: boolean;
  isLoading: boolean; // Initial check loading
  therapistProfile: Therapist | null; // ADDED: Therapist profile data
  isProfileLoading: boolean; // ADDED: Status of profile data fetching
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
  
  // 2. Use the dedicated hook to fetch the profile based on userId
  const { 
    therapist: therapistProfile, 
    isLoading: isProfileLoading, 
  } = useTherapist(userId);


  // Check if user is authenticated on mount
  useEffect(() => {
    const checkAuth = () => {
      const storedUserId = localStorage.getItem("userId");
      const storedJwt = localStorage.getItem("jwt");

      // Only authenticate if we have REAL jwt and userId (not temp ones)
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

      // Store credentials temporarily, DON'T set as authenticated yet
      if (response.access_token && response.user?.id) {
        localStorage.setItem("temp_jwt", response.access_token);
        localStorage.setItem("temp_userId", response.user.id);

        setPendingMFAAuth({
          jwt: response.access_token,
          userId: response.user.id,
        });

        // Return that MFA is required AND return the userId
        return { requiresMFA: true, userId: response.user.id };
      }

      // If no MFA required (shouldn't happen in your case), authenticate immediately
      if (response.access_token) {
        localStorage.setItem("jwt", response.access_token);
      }
      if (response.user?.id) {
        localStorage.setItem("userId", response.user.id);
        setUserId(response.user.id);
        // userId change triggers useTherapist
      }

      return { requiresMFA: false, userId: response.user?.id || null };
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
      // userId change triggers useTherapist
    }
  };

  const signup = async (credentials: PostAuthSignupBody) => {
    try {
      const response = await userSignup(credentials);

      // Store as temp credentials (similar to login) - DON'T authenticate yet
      if (response.access_token && response.user?.id) {
        localStorage.setItem("temp_jwt", response.access_token);
        localStorage.setItem("temp_userId", response.user.id);

        setPendingMFAAuth({
          jwt: response.access_token,
          userId: response.user.id,
        });
      }

      // DON'T set userId or authenticate - wait for MFA
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

    setPendingMFAAuth(null);

    // Clear any remaining cookies 
    document.cookie.split(";").forEach((cookie) => {
      const eqPos = cookie.indexOf("=");
      const name =
        eqPos > -1 ? cookie.substring(0, eqPos).trim() : cookie.trim();
      document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`;
    });

    userLogout();
    setUserId(null); // Setting userId to null automatically resets the profile via useTherapist
    router.push("/login");
  };

  // 3. Update the Context.Provider value
  return (
    <AuthContext.Provider
      value={{
        userId,
        isAuthenticated: !!userId,
        isLoading,
        therapistProfile, // Exposed profile data
        isProfileLoading, // Exposed profile loading state
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