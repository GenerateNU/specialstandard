"use client";

import { Loader2, LogIn } from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import CustomAlert from "@/components/ui/CustomAlert";
import { Input } from "@/components/ui/input";
import { useAuthContext } from "@/contexts/authContext";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [rememberMe, setRememberMe] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const { login, isAuthenticated } = useAuthContext();
  const router = useRouter();

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      router.push("/students");
    }
  }, [isAuthenticated, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      await login({ email, password, remember_me: rememberMe });
    } catch (err: any) {
      setError(err.response?.data?.message || "Invalid email or password");
    } finally {
      setIsLoading(false);
    }
  };

  const [showError, setShowError] = useState(false);

  useEffect(() => {
    // show the alert whenever `error` becomes non-null
    if (error) setShowError(true);
  }, [error]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-background px-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <Image
            src="/tss.png"
            alt="The Special Standard logo"
            width={180}
            height={38}
            className="mx-auto mb-6"
            priority
          />
          <h1 className="text-3xl font-bold text-primary mb-2">Welcome Back</h1>
          <p className="text-secondary">Sign in to access your account</p>
        </div>

        <div className="bg-card rounded-lg shadow-lg border border-default p-8 flex flex-col gap-2">
          {showError && error && (
            <CustomAlert
              variant="destructive"
              title="Login Failed"
              description={error}
              onClose={() => {
                setShowError(false);
                setError(null);
              }}
            />
          )}

          <form onSubmit={handleSubmit} className="space-y-6 bg">
            <div>
              <label
                htmlFor="email"
                className="block text-sm font-medium text-primary mb-2"
              >
                Email
              </label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                disabled={isLoading}
                placeholder="therapist@example.com"
              />
            </div>

            <div>
              <label
                htmlFor="password"
                className="block text-sm font-medium text-primary mb-2"
              >
                Password
              </label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                disabled={isLoading}
                placeholder="••••••••"
              />
            </div>

            <div className="flex items-center">
              <Checkbox
                id="remember-me"
                checked={rememberMe}
                onCheckedChange={setRememberMe}
                disabled={isLoading}
              />
              <label
                htmlFor="remember-me"
                className="ml-2 block text-sm text-secondary"
              >
                Remember me
              </label>
            </div>

            <Button type="submit" disabled={isLoading} size="long">
              {isLoading ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin" />
                  <span>Signing in...</span>
                </>
              ) : (
                <>
                  <LogIn className="w-5 h-5" />
                  <span>Sign In</span>
                </>
              )}
            </Button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-secondary">
              Don't have an account?{" "}
              <Link href="/signup">
                <Button variant="link" size="sm">
                  Sign up
                </Button>
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
