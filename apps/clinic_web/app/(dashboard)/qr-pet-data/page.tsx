"use client";

import React, { useState } from "react";
import { QrCode, Keyboard, Camera, AlertCircle, Phone, Search } from "lucide-react";
import { Card, CardHeader, CardTitle, CardBody } from "@/app/components/ui/Card";
import { Tabs, TabList, Tab, TabPanel } from "@/app/components/ui/Tabs";
import { Input } from "@/app/components/ui/Input";
import { Button } from "@/app/components/ui/Button";
import { StatusBadge } from "@/app/components/ui/Badge";
import { Avatar } from "@/app/components/ui/Avatar";
import { MOCK_QR_PET } from "@/app/lib/mock-data";
import Link from "next/link";
import { useLanguage } from "@/app/components/LanguageContext";

export default function QrCheckInPage() {
  const { t } = useLanguage();
  const [code, setCode] = useState("");
  const [phoneOrName, setPhoneOrName] = useState("");
  const [scannedPet, setScannedPet] = useState<typeof MOCK_QR_PET | null>(null);
  const [cameraActive, setCameraActive] = useState(false);
  const [scanning, setScanning] = useState(false);
  const [error, setError] = useState("");

  const handleEnterCode = (e: React.FormEvent) => {
    e.preventDefault();
    if (code.trim().toUpperCase() === "PNX-2034-00087") {
      setScannedPet(MOCK_QR_PET);
      setError("");
    } else {
      setError("Invalid QR Code or Pet ID. Please try again.");
      setScannedPet(null);
    }
  };

  const handleSearchOwner = (e: React.FormEvent) => {
    e.preventDefault();
    const query = phoneOrName.trim();
    // Match phone mock: "081-234-5678" or owner mock: "Sarah Johnson" or "Sarah"
    if (
      query === "081-234-5678" || 
      query === "0812345678" ||
      query.toLowerCase().includes("sarah") ||
      query.toLowerCase().includes("johnson")
    ) {
      setScannedPet(MOCK_QR_PET);
      setError("");
    } else {
      setError("No owner profile matches this name or phone number.");
      setScannedPet(null);
    }
  };

  const handleStartScan = () => {
    setCameraActive(true);
    setScanning(true);
    setError("");

    // Simulate scanning
    setTimeout(() => {
      setScannedPet(MOCK_QR_PET);
      setScanning(false);
      setCameraActive(false);
    }, 2000);
  };

  return (
    <div className="flex flex-col gap-6 max-w-7xl mx-auto">
      <div>
        <h1 className="text-2xl font-bold text-navy-900">{t("qr_title")}</h1>
        <p className="text-sm text-navy-500 mt-1">
          {t("qr_desc")}
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        {/* Left Column — Scanning & Inputs (lg:col-span-7) */}
        <div className="lg:col-span-7">
          <Card className="p-0">
            <Tabs defaultTab="scan">
              <div className="border-b border-navy-200 px-6 pt-4 bg-white">
                <TabList className="w-full flex-wrap mb-4 gap-1 sm:gap-2">
                  <Tab value="scan" className="flex items-center justify-center gap-2">
                    <QrCode size={16} />
                    {t("tab_scan_qr")}
                  </Tab>
                  <Tab value="manual" className="flex items-center justify-center gap-2">
                    <Keyboard size={16} />
                    {t("tab_enter_code")}
                  </Tab>
                  <Tab value="search" className="flex items-center justify-center gap-2">
                    <Search size={16} />
                    {t("tab_search_owner")}
                  </Tab>
                </TabList>
              </div>

              <div className="p-6">
                {/* Scan QR Tab */}
                <TabPanel value="scan">
                  <div className="flex flex-col items-center justify-center py-6">
                    {/* Viewport Frame */}
                    <div className="relative w-64 h-64 border-2 border-navy-300 rounded-xl bg-navy-50 flex items-center justify-center overflow-hidden">
                      {cameraActive ? (
                        <div className="absolute inset-0 flex flex-col items-center justify-center">
                          <Camera className="text-navy-400 w-12 h-12 animate-pulse" />
                          {scanning && (
                            <div className="absolute w-full h-0.5 bg-teal-500 top-1/2 left-0 animate-[bounce_2s_infinite]" />
                          )}
                          <span className="text-xs text-navy-500 mt-2 font-medium">
                            {t("scanning_viewport")}
                          </span>
                        </div>
                      ) : (
                        <div className="flex flex-col items-center justify-center p-6 text-center">
                          <QrCode className="text-navy-300 w-16 h-16 mb-4" />
                          <p className="text-xs font-semibold text-navy-600">
                            {t("position_qr")}
                          </p>
                          <p className="text-[10px] text-navy-400 mt-1">
                            {t("center_to_scan")}
                          </p>
                        </div>
                      )}

                      {/* Corners indicator for QR target */}
                      <div className="absolute top-4 left-4 w-4 h-4 border-t-2 border-l-2 border-teal-500" />
                      <div className="absolute top-4 right-4 w-4 h-4 border-t-2 border-r-2 border-teal-500" />
                      <div className="absolute bottom-4 left-4 w-4 h-4 border-b-2 border-l-2 border-teal-500" />
                      <div className="absolute bottom-4 right-4 w-4 h-4 border-b-2 border-r-2 border-teal-500" />
                    </div>

                    <Button
                      variant="outline"
                      className="mt-6 border-red-200 text-red-600 hover:bg-red-50 hover:border-red-300 animate-pulse"
                      onClick={handleStartScan}
                      disabled={scanning}
                      icon={<Camera size={16} />}
                    >
                      {scanning ? t("loading") : t("turn_on_camera")}
                    </Button>

                    <p className="text-xs text-navy-400 mt-4 text-center">
                      {t("manual_prompt")}{" "}
                      <button
                        type="button"
                        onClick={() => {}}
                        className="text-teal-600 font-semibold cursor-pointer"
                      >
                        {t("manual_link")}
                      </button>
                    </p>
                  </div>
                </TabPanel>

                {/* Manual Code Input Tab */}
                <TabPanel value="manual">
                  <form onSubmit={handleEnterCode} className="max-w-md mx-auto py-8 flex flex-col gap-4">
                    <Input
                      label={t("pnx_id_label")}
                      placeholder="e.g. PNX-2034-00087"
                      value={code}
                      onChange={(e) => setCode(e.target.value)}
                      error={error}
                      required
                    />
                    <Button type="submit" fullWidth>
                      {t("retrieve_data_btn")}
                    </Button>
                    <p className="text-xs text-navy-400 text-center">
                      {t("enter_code_desc")}
                    </p>
                  </form>
                </TabPanel>

                {/* Find User / Phone Tab */}
                <TabPanel value="search">
                  <form onSubmit={handleSearchOwner} className="max-w-md mx-auto py-8 flex flex-col gap-4">
                    <Input
                      label={t("search_owner_label")}
                      placeholder={t("search_owner_placeholder")}
                      value={phoneOrName}
                      onChange={(e) => setPhoneOrName(e.target.value)}
                      error={error}
                      prefixIcon={<Phone size={16} />}
                      required
                    />
                    <Button type="submit" fullWidth>
                      {t("search_owner_btn")}
                    </Button>
                    <p className="text-xs text-navy-400 text-center">
                      Search using registered owner name (e.g. Sarah) or phone (e.g. 081-234-5678).
                    </p>
                  </form>
                </TabPanel>
              </div>
            </Tabs>
          </Card>
        </div>

        {/* Right Column — Loaded Pet Info (lg:col-span-5) */}
        <div className="lg:col-span-5">
          {scannedPet ? (
            <Card className="animate-[slide-up_0.25s_ease-out] bg-white">
              <CardHeader>
                <CardTitle subtitle={t("profile_retrieved")}>
                  {t("checkin_profile")}
                </CardTitle>
              </CardHeader>
              <CardBody>
                {/* Pet Identity Header */}
                <div className="flex items-center gap-4 bg-navy-50/50 p-4 rounded-xl border border-navy-200 mb-6">
                  <Avatar name={scannedPet.name} size="lg" />
                  <div>
                    <h3 className="text-base font-bold text-navy-900">
                      {scannedPet.name}
                    </h3>
                    <p className="text-xs font-semibold text-navy-500 mt-0.5">
                      {scannedPet.breed}
                    </p>
                    <p className="text-[10px] text-teal-600 font-mono font-semibold mt-1">
                      ID: {scannedPet.petNexusId}
                    </p>
                  </div>
                </div>

                {/* Meta details list */}
                <div className="space-y-4 text-sm mb-6">
                  <div className="flex justify-between py-1 border-b border-navy-100">
                    <span className="text-navy-500 font-medium">{t("pet_owner")}</span>
                    <span className="text-navy-800 font-semibold">
                      {scannedPet.ownerName}
                    </span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-navy-100">
                    <span className="text-navy-500 font-medium">{t("pet_phone")}</span>
                    <span className="text-navy-800 font-semibold">
                      {scannedPet.ownerPhone}
                    </span>
                  </div>
                  <div className="flex justify-between py-1 border-b border-navy-100">
                    <span className="text-navy-500 font-medium">{t("last_visit")}</span>
                    <span className="text-navy-800 font-semibold">
                      {scannedPet.birthDate ? "Apr 12, 2025" : "N/A"}
                    </span>
                  </div>
                  <div className="flex justify-between items-center py-1">
                    <span className="text-navy-500 font-medium">{t("pet_status")}</span>
                    <StatusBadge status="checked-in" />
                  </div>
                </div>

                {/* Actions */}
                <div className="grid grid-cols-2 gap-3">
                  <Link href="/patients" className="w-full">
                    <Button variant="outline" className="w-full">
                      {t("view_history_btn")}
                    </Button>
                  </Link>
                  <Link href="/medical-records/new" className="w-full">
                    <Button className="w-full">{t("create_visit_btn")}</Button>
                  </Link>
                </div>
              </CardBody>
            </Card>
          ) : (
            <Card className="border-dashed border-2 border-navy-300 bg-navy-50/20 py-16 flex flex-col items-center justify-center text-center">
              <AlertCircle className="text-navy-400 w-12 h-12 mb-4" />
              <p className="text-sm font-semibold text-navy-600">
                {t("no_pet_checkin_title")}
              </p>
              <p className="text-xs text-navy-400 mt-1 max-w-[240px]">
                {t("no_pet_checkin_desc")}
              </p>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}
