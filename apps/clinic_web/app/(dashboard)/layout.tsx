import React from "react";
import { Sidebar } from "@/app/components/layout/Sidebar";
import { Topbar } from "@/app/components/layout/Topbar";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex h-screen w-screen overflow-hidden bg-page-bg">
      {/* Sidebar (Fixed Width, Sticky Left) */}
      <Sidebar />

      {/* Main Container */}
      <div className="flex-1 flex flex-col min-w-0 overflow-hidden h-full">
        {/* Topbar */}
        <Topbar />

        {/* Content Body (Scrollable) */}
        <main className="flex-1 overflow-y-auto p-4 sm:p-6 lg:p-8 animate-[fade-in_0.2s_ease-out]">
          {children}
        </main>
      </div>
    </div>
  );
}
