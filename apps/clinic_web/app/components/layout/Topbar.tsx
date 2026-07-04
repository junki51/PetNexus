"use client";

import React, { useState } from "react";
import { Bell, ChevronDown, LogOut, User as UserIcon } from "lucide-react";
import { Avatar } from "@/app/components/ui/Avatar";
import { MOCK_CURRENT_USER } from "@/app/lib/mock-data";
import Link from "next/link";
import { useLanguage } from "@/app/components/LanguageContext";

interface TopbarProps {
  clinicName?: string;
  onLogout?: () => void;
}

export function Topbar({
  clinicName = "Happy Paws Veterinary Clinic",
  onLogout,
}: TopbarProps) {
  const [profileOpen, setProfileOpen] = useState(false);
  const [clinicOpen, setClinicOpen] = useState(false);
  const { lang, setLang, t } = useLanguage();

  return (
    <header className="h-15 border-b border-navy-200 bg-white flex items-center justify-between px-6 z-30 sticky top-0">
      {/* Clinic Selector (Left) */}
      <div className="relative">
        <button
          onClick={() => setClinicOpen(!clinicOpen)}
          className="flex items-center gap-1.5 text-navy-800 font-medium text-sm hover:text-navy-950 transition-colors py-2 cursor-pointer"
        >
          <span>{clinicName}</span>
          <ChevronDown size={14} className="text-navy-400" />
        </button>

        {clinicOpen && (
          <>
            <div
              className="fixed inset-0 z-30"
              onClick={() => setClinicOpen(false)}
            />
            <div className="absolute left-0 mt-1 w-56 bg-white border border-navy-200 rounded-lg shadow-lg py-1 z-40 animate-[scale-in_0.1s_ease-out]">
              <button
                className="w-full text-left px-4 py-2 text-xs font-semibold text-navy-400 uppercase tracking-wide"
                disabled
              >
                {t("switch_clinic")}
              </button>
              <button
                className="w-full text-left px-4 py-2 text-sm text-navy-700 bg-navy-50 font-medium border-l-2 border-teal-600"
                onClick={() => setClinicOpen(false)}
              >
                {clinicName}
              </button>
              <button
                className="w-full text-left px-4 py-2 text-sm text-navy-600 hover:bg-navy-50"
                onClick={() => setClinicOpen(false)}
              >
                {t("register_new_clinic")}
              </button>
            </div>
          </>
        )}
      </div>

      {/* User Actions (Right) */}
      <div className="flex items-center gap-4">
        {/* Language Switcher */}
        <div className="flex items-center bg-navy-100 rounded-lg p-0.5 border border-navy-200 select-none">
          <button
            onClick={() => setLang("th")}
            className={`px-2 py-0.5 text-xs font-semibold rounded cursor-pointer transition-all duration-150 ${
              lang === "th"
                ? "bg-white text-teal-600 shadow-sm"
                : "text-navy-500 hover:text-navy-800"
            }`}
          >
            ไทย
          </button>
          <button
            onClick={() => setLang("en")}
            className={`px-2 py-0.5 text-xs font-semibold rounded cursor-pointer transition-all duration-150 ${
              lang === "en"
                ? "bg-white text-teal-600 shadow-sm"
                : "text-navy-500 hover:text-navy-800"
            }`}
          >
            EN
          </button>
        </div>

        {/* Notification Bell */}
        <button className="w-8 h-8 rounded-lg flex items-center justify-center text-navy-500 hover:bg-navy-100 hover:text-navy-700 transition-colors relative cursor-pointer">
          <Bell size={18} />
          <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full border border-white" />
        </button>

        {/* Vertical Separator */}
        <div className="h-6 w-px bg-navy-200" />

        {/* Profile Dropdown */}
        <div className="relative">
          <button
            onClick={() => setProfileOpen(!profileOpen)}
            className="flex items-center gap-2 py-1.5 cursor-pointer group"
          >
            <Avatar
              name={MOCK_CURRENT_USER.name}
              size="sm"
            />
            <div className="text-left hidden sm:block">
              <p className="text-xs font-semibold text-navy-800 leading-tight group-hover:text-navy-950 transition-colors">
                {MOCK_CURRENT_USER.name}
              </p>
              <p className="text-[10px] text-navy-500 font-medium leading-tight capitalize">
                {MOCK_CURRENT_USER.role}
              </p>
            </div>
            <ChevronDown size={14} className="text-navy-400 group-hover:text-navy-600 transition-colors" />
          </button>

          {profileOpen && (
            <>
              <div
                className="fixed inset-0 z-30"
                onClick={() => setProfileOpen(false)}
              />
              <div className="absolute right-0 mt-1 w-48 bg-white border border-navy-200 rounded-lg shadow-lg py-1 z-40 animate-[scale-in_0.1s_ease-out]">
                <div className="px-4 py-2 border-b border-navy-100">
                  <p className="text-xs font-medium text-navy-400">{t("signed_in_as")}</p>
                  <p className="text-sm font-semibold text-navy-800 truncate">
                    emily.carter@petnexus.com
                  </p>
                </div>
                <Link
                  href="/settings"
                  className="flex items-center gap-2 px-4 py-2 text-sm text-navy-700 hover:bg-navy-50 transition-colors"
                  onClick={() => setProfileOpen(false)}
                >
                  <UserIcon size={16} className="text-navy-400" />
                  {t("my_profile")}
                </Link>
                <Link
                  href="/login"
                  className="flex items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors border-t border-navy-100"
                  onClick={() => {
                    setProfileOpen(false);
                    onLogout?.();
                  }}
                >
                  <LogOut size={16} />
                  {t("sign_out")}
                </Link>
              </div>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
