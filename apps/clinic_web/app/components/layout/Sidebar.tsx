"use client";

import React, { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard,
  QrCode,
  Users,
  Calendar,
  FileText,
  BarChart2,
  Settings,
  HelpCircle,
  PawPrint,
  Menu,
  X,
} from "lucide-react";
import { useLanguage } from "@/app/components/LanguageContext";

interface SidebarProps {
  clinicName?: string;
}

export function Sidebar({ clinicName = "Happy Paws Veterinary Clinic" }: SidebarProps) {
  const pathname = usePathname();
  const [mobileOpen, setMobileOpen] = useState(false);
  const { t } = useLanguage();

  function isActive(href: string): boolean {
    if (href === "/(dashboard)") return pathname === "/";
    return pathname.startsWith(href);
  }

  const navItems = [
    { label: t("dashboard"),       href: "/(dashboard)",              icon: <LayoutDashboard size={18} /> },
    { label: t("qr_pet_data"),     href: "/qr-pet-data",              icon: <QrCode size={18} /> },
    { label: t("patients"),        href: "/patients",                 icon: <Users size={18} /> },
    { label: t("calendar"),        href: "/calendar",                 icon: <Calendar size={18} /> },
    { label: t("medical_records"), href: "/medical-records",          icon: <FileText size={18} /> },
    { label: t("reports"),         href: "/reports",                  icon: <BarChart2 size={18} /> },
    { label: t("settings"),        href: "/settings",                 icon: <Settings size={18} /> },
  ];

  const renderSidebarContent = () => (
    <div className="flex flex-col h-full">
      {/* Logo */}
      <div className="flex items-center gap-3 px-5 py-5 border-b border-white/10">
        <div className="w-9 h-9 bg-teal-500 rounded-xl flex items-center justify-center shrink-0 shadow-lg">
          <PawPrint size={20} className="text-white" />
        </div>
        <div>
          <p className="font-bold text-white text-sm leading-tight">PetNexus</p>
          <p className="text-teal-400 text-xs font-medium leading-tight">{t("clinic_platform")}</p>
        </div>
      </div>

      {/* Clinic name */}
      <div className="px-5 py-3 border-b border-white/10">
        <p className="text-navy-400 text-xs font-medium truncate">{clinicName}</p>
      </div>

      {/* Nav items */}
      <nav className="flex-1 overflow-y-auto py-3 px-3">
        <ul className="flex flex-col gap-0.5">
          {navItems.map((item) => {
            const active = isActive(item.href);
            return (
              <li key={item.href}>
                <Link
                  href={item.href === "/(dashboard)" ? "/" : item.href}
                  onClick={() => setMobileOpen(false)}
                  className={[
                    "flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium",
                    "transition-all duration-150 group",
                    active
                      ? "bg-teal-600 text-white shadow-sm"
                      : "text-navy-300 hover:bg-white/10 hover:text-white",
                  ].join(" ")}
                >
                  <span className={active ? "text-white" : "text-navy-400 group-hover:text-white transition-colors"}>
                    {item.icon}
                  </span>
                  {item.label}
                </Link>
              </li>
            );
          })}
        </ul>
      </nav>

      {/* Help & Support */}
      <div className="px-3 pb-5 border-t border-white/10 pt-3">
        <Link
          href="#"
          className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-navy-400 hover:bg-white/10 hover:text-white transition-all duration-150"
        >
          <HelpCircle size={18} />
          {t("help_support")}
        </Link>
      </div>
    </div>
  );

  return (
    <>
      {/* Mobile hamburger button */}
      <button
        className="fixed top-4 left-4 z-50 lg:hidden bg-navy-800 text-white p-2 rounded-lg shadow-lg"
        onClick={() => setMobileOpen((v) => !v)}
        aria-label="Toggle sidebar"
      >
        {mobileOpen ? <X size={20} /> : <Menu size={20} />}
      </button>

      {/* Mobile overlay */}
      {mobileOpen && (
        <div
          className="fixed inset-0 bg-black/40 z-40 lg:hidden"
          onClick={() => setMobileOpen(false)}
        />
      )}

      {/* Mobile sidebar (slide-in) */}
      <aside
        className={[
          "fixed top-0 left-0 h-full w-60 bg-navy-800 z-50 transition-transform duration-300 ease-out lg:hidden",
          mobileOpen ? "translate-x-0" : "-translate-x-full",
        ].join(" ")}
      >
        {renderSidebarContent()}
      </aside>

      {/* Desktop sidebar (always visible) */}
      <aside className="hidden lg:flex flex-col w-60 shrink-0 bg-navy-800 h-screen sticky top-0">
        {renderSidebarContent()}
      </aside>
    </>
  );
}
