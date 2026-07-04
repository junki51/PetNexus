"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

export default function RegisterRedirectPage() {
  const router = useRouter();

  useEffect(() => {
    // Redirect to login page and default to signup tab (if handled by query param, but we'll just go to login)
    router.replace("/login");
  }, [router]);

  return (
    <div className="flex items-center justify-center min-h-[400px]">
      <div className="animate-spin h-8 w-8 text-teal-600 rounded-full border-4 border-solid border-current border-r-transparent align-[-0.125em] motion-reduce:animate-[spin_1.5s_linear_infinite]" />
    </div>
  );
}
