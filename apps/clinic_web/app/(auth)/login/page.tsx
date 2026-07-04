"use client";

import React, { useState } from "react";
import { PawPrint, Mail, Lock } from "lucide-react";
import { Input } from "@/app/components/ui/Input";
import { Button } from "@/app/components/ui/Button";
import { Checkbox } from "@/app/components/ui/Checkbox";
import { Divider } from "@/app/components/ui/Divider";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function LoginPage() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState<"signin" | "signup">("signin");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [rememberMe, setRememberMe] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    // Simulate login
    setTimeout(() => {
      setLoading(false);
      router.push("/");
    }, 1000);
  };

  return (
    <div className="w-full max-w-4xl bg-white border border-navy-200 rounded-2xl shadow-xl overflow-hidden flex flex-col md:flex-row min-h-140 animate-[fade-in_0.3s_ease-out]">
      {/* Left Panel — Branding & Photo Mock (Teal & Navy Background) */}
      <div className="w-full md:w-1/2 bg-navy-50/50 p-8 md:p-12 flex flex-col justify-between border-b md:border-b-0 md:border-r border-navy-200">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-teal-600 rounded-xl flex items-center justify-center shrink-0 shadow-md">
            <PawPrint className="text-white" size={20} />
          </div>
          <div>
            <h2 className="text-lg font-bold text-navy-900 leading-tight">
              PetNexus
            </h2>
            <p className="text-teal-600 text-xs font-semibold leading-tight">
              Clinic Platform
            </p>
          </div>
        </div>

        <div className="my-8 flex flex-col items-center justify-center text-center">
          {/* Welcome back text */}
          <h1 className="text-2xl font-bold text-navy-800 leading-tight">
            Welcome back!
          </h1>
          <p className="text-sm text-navy-500 mt-2">
            Sign in to manage your clinic, schedule visits, and access pet records.
          </p>

          {/* Dog & Cat Illustration / Placeholder Area */}
          <div className="mt-8 relative w-56 h-56 bg-teal-100/50 rounded-full flex items-center justify-center overflow-hidden border border-teal-200/50">
            {/* Styled Paw SVG background */}
            <PawPrint size={120} className="text-teal-600/10 absolute animate-pulse" />
            <svg
              className="w-48 h-48 text-teal-600/70 relative z-10"
              viewBox="0 0 200 200"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              {/* Dog shape outline simplified */}
              <path
                d="M50 140 C50 120 70 80 90 80 C110 80 120 100 130 90 C140 80 150 60 160 80 C170 100 170 140 160 150 C150 160 60 160 50 140 Z"
                fill="currentColor"
                opacity="0.8"
              />
              {/* Cat shape outline simplified */}
              <path
                d="M100 150 C100 130 115 110 125 110 C135 110 140 120 145 115 C150 110 155 100 160 115 C165 130 165 150 160 155 C155 160 105 160 100 150 Z"
                fill="var(--color-teal-700)"
                opacity="0.9"
              />
            </svg>
          </div>
        </div>

        <div className="text-xs text-navy-400 text-center md:text-left">
          © {new Date().getFullYear()} PetNexus Inc. All rights reserved.
        </div>
      </div>

      {/* Right Panel — Forms (Sign In / Sign Up) */}
      <div className="w-full md:w-1/2 p-8 md:p-12 flex flex-col justify-center">
        {/* Sign In / Create Account Tabs */}
        <div className="flex border-b border-navy-200 mb-8">
          <button
            onClick={() => setActiveTab("signin")}
            className={[
              "flex-1 pb-3 text-sm font-semibold border-b-2 text-center transition-all cursor-pointer",
              activeTab === "signin"
                ? "border-teal-600 text-teal-600"
                : "border-transparent text-navy-500 hover:text-navy-700",
            ].join(" ")}
          >
            Sign In
          </button>
          <button
            onClick={() => setActiveTab("signup")}
            className={[
              "flex-1 pb-3 text-sm font-semibold border-b-2 text-center transition-all cursor-pointer",
              activeTab === "signup"
                ? "border-teal-600 text-teal-600"
                : "border-transparent text-navy-500 hover:text-navy-700",
            ].join(" ")}
          >
            Create Account
          </button>
        </div>

        {activeTab === "signin" ? (
          /* Sign In Form */
          <form onSubmit={handleSubmit} className="flex flex-col gap-5">
            <Input
              label="Username or Email"
              placeholder="Enter username or email"
              type="text"
              required
              prefixIcon={<Mail size={16} />}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />

            <div className="flex flex-col gap-1">
              <Input
                label="Password"
                placeholder="Enter password"
                type="password"
                required
                prefixIcon={<Lock size={16} />}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              <div className="flex items-center justify-between mt-1">
                <Checkbox
                  id="remember-me"
                  label="Remember me"
                  checked={rememberMe}
                  onChange={(e) => setRememberMe(e.target.checked)}
                />
                <Link
                  href="#"
                  className="text-xs font-semibold text-teal-600 hover:text-teal-700"
                >
                  Forgot password?
                </Link>
              </div>
            </div>

            <Button type="submit" loading={loading} fullWidth className="mt-2">
              Sign In
            </Button>

            <Divider label="or continue with" />

            {/* Google Sign In */}
            <button
              type="button"
              onClick={handleSubmit}
              className="w-full h-10 border border-navy-200 rounded-lg flex items-center justify-center gap-3 bg-white text-navy-700 text-sm font-medium hover:bg-navy-50 active:bg-navy-100 transition-colors shadow-sm cursor-pointer"
            >
              <svg className="w-4 h-4 shrink-0" viewBox="0 0 24 24">
                <path
                  fill="#4285F4"
                  d="M23.745 12.27c0-.7-.06-1.4-.19-2.07H12v3.9h6.69c-.29 1.5-.1.14-.14 3.06-.9 2.05-2.6 3.73-4.66 4.67v3.86h7.52c4.4-4.04 6.94-10 6.94-16.73z"
                />
                <path
                  fill="#34A853"
                  d="M12 24c3.24 0 5.97-1.08 7.96-2.91l-7.52-3.86c-2.08 1.4-4.73 2.22-7.96 2.22-6.13 0-11.32-4.14-13.17-9.72H.4v4c3.6 7.15 10.93 11.97 19.32 11.97z"
                />
                <path
                  fill="#FBBC05"
                  d="M6.3 14.13c-.45-1.34-.7-2.77-.7-4.25s.25-2.91.7-4.25v-4H1.36C.49 3.37 0 5.48 0 7.75s.49 4.38 1.36 6.13l4.94-3.75z"
                />
                <path
                  fill="#EA4335"
                  d="M12 4.75c1.77 0 3.35.61 4.6 1.8l3.42-3.42C17.96 1.19 15.24 0 12 0 6.93 0 2.5 4.82.4 11.97l4.94 3.75c1.85-5.58 7.04-9.72 13.17-9.72z"
                />
              </svg>
              Sign in with Google
            </button>

            <p className="text-xs text-navy-500 text-center mt-4">
              Don&apos;t have an account?{" "}
              <button
                type="button"
                onClick={() => setActiveTab("signup")}
                className="font-semibold text-teal-600 hover:text-teal-700 cursor-pointer"
              >
                Register your clinic
              </button>
            </p>
          </form>
        ) : (
          /* Sign Up Form Placeholder */
          <form
            onSubmit={(e) => {
              e.preventDefault();
              setLoading(true);
              setTimeout(() => {
                setLoading(false);
                setActiveTab("signin");
              }, 1000);
            }}
            className="flex flex-col gap-5 animate-[slide-up_0.2s_ease-out]"
          >
            <Input
              label="Clinic Name"
              placeholder="Enter your clinic name"
              type="text"
              required
            />
            <Input
              label="Email Address"
              placeholder="Enter email address"
              type="email"
              required
            />
            <Input
              label="Password"
              placeholder="Create secure password"
              type="password"
              required
            />
            <Checkbox
              id="terms"
              required
              label={
                <span className="text-xs text-navy-500">
                  I agree to the{" "}
                  <Link href="#" className="font-semibold text-teal-600">
                    Terms of Service
                  </Link>{" "}
                  and{" "}
                  <Link href="#" className="font-semibold text-teal-600">
                    Privacy Policy
                  </Link>
                </span>
              }
            />
            <Button type="submit" loading={loading} fullWidth className="mt-2">
              Create Account
            </Button>
            <p className="text-xs text-navy-500 text-center mt-2">
              Already have an account?{" "}
              <button
                type="button"
                onClick={() => setActiveTab("signin")}
                className="font-semibold text-teal-600 hover:text-teal-700 cursor-pointer"
              >
                Sign In
              </button>
            </p>
          </form>
        )}
      </div>
    </div>
  );
}
