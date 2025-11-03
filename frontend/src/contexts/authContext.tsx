"use client";
import { useAuth } from "@/hooks/useAuth";
import type {
  PostAuthLoginBody,
  PostAuthSignupBody,
} from "@/lib/api/theSpecialStandardAPI.schemas";
import { useRouter } from "next/navigation";
import { createContext, useContext, useEffect, useState } from "react";

interface AuthContextType {
  userId: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: PostAuthLoginBody) => Promise<{ requiresMFA: boolean }>;
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

  // Check if user is authenticated on mount
  useEffect(() => {
    const checkAuth = () => {
      const isDevMode = false; //process.env.NODE_ENV === "development";

      if (!isDevMode) {
        const mockUserId = "dev-user-123";
        const mockToken = "dev-mock-token";

        if (!localStorage.getItem("userId")) {
          localStorage.setItem("userId", mockUserId);
          localStorage.setItem("jwt", mockToken);
        }

        setUserId(localStorage.getItem("userId"));
        setIsLoading(false);
        return;
      }

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

      // Store credentials temporarily, DON'T set as authenticated yet
      if (response.access_token && response.user?.id) {
        setPendingMFAAuth({
          jwt: response.access_token,
          userId: response.user.id,
        });

        // Return that MFA is required
        return { requiresMFA: true };
      }

      // If no MFA required (shouldn't happen in your case), authenticate immediately
      if (response.access_token) {
        localStorage.setItem("jwt", response.access_token);
      }
      if (response.user?.id) {
        localStorage.setItem("userId", response.user.id);
        setUserId(response.user.id);
      }

      return { requiresMFA: false };
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  };

  // Call this after MFA verification succeeds
  const completeMFALogin = () => {
    if (pendingMFAAuth) {
      localStorage.setItem("jwt", pendingMFAAuth.jwt);
      localStorage.setItem("userId", pendingMFAAuth.userId);
      setUserId(pendingMFAAuth.userId);
      setPendingMFAAuth(null);
    }
  };

  const signup = async (credentials: PostAuthSignupBody) => {
    try {
      const response = await userSignup(credentials);

      if (response.access_token) {
        localStorage.setItem("jwt", response.access_token);
      }
      if (response.user?.id) {
        localStorage.setItem("userId", response.user.id);
        setUserId(response.user.id);
      }

      router.push("/");
    } catch (error) {
      console.error("Signup failed:", error);
      throw error;
    }
  };

  const logout = () => {
    localStorage.removeItem("jwt");
    localStorage.removeItem("userId");
    setPendingMFAAuth(null);

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
