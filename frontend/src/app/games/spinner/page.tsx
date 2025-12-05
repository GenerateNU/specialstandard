import type { Metadata } from "next";
import { Suspense } from "react";
import { SpinnerContent } from "./SpinnerContent";

export const metadata: Metadata = {
  title: "Spinner Game",
  description: "Practice with interactive wheel",
};

function LoadingSpinner() {
  return (
    <div className="min-h-screen bg-background p-8 flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue mx-auto mb-4"></div>
        <p className="text-secondary">Loading...</p>
      </div>
    </div>
  );
}

export default function SpinnerPage() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <SpinnerContent />
    </Suspense>
  );
}
