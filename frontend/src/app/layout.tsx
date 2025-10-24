// src/app/layout.tsx
import { AppWithMFA } from "@/components/MFA/AppWithMFA";
import { AuthProvider } from "@/contexts/authContext";
import { ReactQueryProvider } from "@/lib/ReactQueryProvider";
import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "The Special Standard",
  description: "Student management system for The Special Standard",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <ReactQueryProvider>
          <AuthProvider>
            <AppWithMFA>{children}</AppWithMFA>
          </AuthProvider>
        </ReactQueryProvider>
      </body>
    </html>
  );
}
