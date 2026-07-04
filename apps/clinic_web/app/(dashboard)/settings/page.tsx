"use client";

import React, { useState } from "react";
import { Save, UserPlus, Bell, Building, ShieldCheck } from "lucide-react";
import { Card } from "@/app/components/ui/Card";
import { Tabs, TabList, Tab, TabPanel } from "@/app/components/ui/Tabs";
import { Input } from "@/app/components/ui/Input";
import { Button } from "@/app/components/ui/Button";
import { Checkbox } from "@/app/components/ui/Checkbox";
import { Avatar } from "@/app/components/ui/Avatar";
import { MOCK_CLINIC } from "@/app/lib/mock-data";

export default function SettingsPage() {
  const [clinicName, setClinicName] = useState(MOCK_CLINIC.name);
  const [clinicPhone, setClinicPhone] = useState("02-123-4567");
  const [clinicEmail, setClinicEmail] = useState("contact@happypaws.com");
  const [clinicAddress, setClinicAddress] = useState("123 Pet Street, Bangkok, 10500");

  const [loading, setLoading] = useState(false);

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
    }, 1000);
  };

  return (
    <div className="flex flex-col gap-6 max-w-7xl mx-auto">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-navy-900">Settings</h1>
        <p className="text-sm text-navy-500 mt-1">
          Manage clinic settings, staff members, and system preferences.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        {/* Settings Navigation and Panels inside a Card (lg:col-span-12) */}
        <div className="lg:col-span-12">
          <Card className="p-0 overflow-hidden">
            <Tabs defaultTab="clinic">
              <div className="border-b border-navy-200 px-6 pt-4 bg-white">
                <TabList className="w-full max-w-md mb-4">
                  <Tab value="clinic" className="flex items-center justify-center gap-2">
                    <Building size={16} />
                    Clinic Profile
                  </Tab>
                  <Tab value="staff" className="flex items-center justify-center gap-2">
                    <ShieldCheck size={16} />
                    Staff Members
                  </Tab>
                  <Tab value="notifications" className="flex items-center justify-center gap-2">
                    <Bell size={16} />
                    System Rules
                  </Tab>
                </TabList>
              </div>

              <div className="p-6 bg-white">
                {/* Clinic Profile Tab */}
                <TabPanel value="clinic">
                  <form onSubmit={handleSave} className="max-w-2xl flex flex-col gap-6">
                    <div>
                      <h3 className="text-sm font-bold text-navy-800 uppercase tracking-wider mb-4">
                        Clinic Details
                      </h3>
                      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div className="sm:col-span-2">
                          <Input
                            label="Clinic Name"
                            value={clinicName}
                            onChange={(e) => setClinicName(e.target.value)}
                            required
                          />
                        </div>
                        <Input
                          label="Phone Number"
                          value={clinicPhone}
                          onChange={(e) => setClinicPhone(e.target.value)}
                          required
                        />
                        <Input
                          label="Email Address"
                          type="email"
                          value={clinicEmail}
                          onChange={(e) => setClinicEmail(e.target.value)}
                          required
                        />
                        <div className="sm:col-span-2">
                          <Input
                            label="Full Address"
                            value={clinicAddress}
                            onChange={(e) => setClinicAddress(e.target.value)}
                            required
                          />
                        </div>
                      </div>
                    </div>

                    <div className="flex justify-end gap-3 border-t border-navy-100 pt-4">
                      <Button type="submit" loading={loading} icon={<Save size={16} />}>
                        Save Preferences
                      </Button>
                    </div>
                  </form>
                </TabPanel>

                {/* Staff Members Tab */}
                <TabPanel value="staff">
                  <div className="flex flex-col gap-6">
                    <div className="flex items-center justify-between">
                      <h3 className="text-sm font-bold text-navy-800 uppercase tracking-wider">
                        Active staff members
                      </h3>
                      <Button size="sm" icon={<UserPlus size={14} />}>
                        Add Staff
                      </Button>
                    </div>

                    <div className="flex flex-col divide-y divide-navy-100 border border-navy-200 rounded-xl overflow-hidden bg-white">
                      {[
                        {
                          name: "Dr. Emily Carter",
                          role: "Chief Veterinarian",
                          email: "emily.carter@petnexus.com",
                          license: "VET-8839-2021",
                          active: true,
                        },
                        {
                          name: "Dr. James Wilson",
                          role: "Veterinarian",
                          email: "james.wilson@happypaws.com",
                          license: "VET-9402-2023",
                          active: true,
                        },
                        {
                          name: "Sarah Jenkins",
                          role: "Clinic Assistant",
                          email: "sarah.j@happypaws.com",
                          license: "N/A",
                          active: true,
                        },
                      ].map((member, idx) => (
                        <div
                          key={idx}
                          className="p-4 flex items-center justify-between gap-4 hover:bg-navy-50/50 transition-colors"
                        >
                          <div className="flex items-center gap-4">
                            <Avatar name={member.name} size="md" />
                            <div>
                              <p className="text-sm font-bold text-navy-950">
                                {member.name}
                              </p>
                              <p className="text-xs text-teal-600 font-semibold mt-0.5">
                                {member.role}
                              </p>
                            </div>
                          </div>

                          <div className="hidden sm:flex flex-col text-right">
                            <p className="text-xs text-navy-600 font-medium">
                              {member.email}
                            </p>
                            <p className="text-[10px] text-navy-400 font-mono mt-0.5">
                              License: {member.license}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                </TabPanel>

                {/* Notifications & System Rules Tab */}
                <TabPanel value="notifications">
                  <div className="max-w-2xl flex flex-col gap-6">
                    <div>
                      <h3 className="text-sm font-bold text-navy-800 uppercase tracking-wider mb-4">
                        Data Sharing Preferences
                      </h3>
                      <div className="flex flex-col gap-4">
                        <Checkbox
                          id="auto-request"
                          label="Auto-request QR records access"
                          description="Automatically send an authorization request when a new pet QR is scanned."
                          defaultChecked
                        />
                        <Checkbox
                          id="strict-auth"
                          label="Strict visit verification policy"
                          description="Require digital signatures for all medical treatments entered into records."
                          defaultChecked
                        />
                        <Checkbox
                          id="notify-email"
                          label="Email notifications for authorization status"
                          description="Receive an email summary whenever an owner approves or revokes record sharing permissions."
                        />
                      </div>
                    </div>
                  </div>
                </TabPanel>
              </div>
            </Tabs>
          </Card>
        </div>
      </div>
    </div>
  );
}
