"use client";

import React, { useState } from "react";
import { ChevronLeft, ChevronRight, Plus, Calendar as CalendarIcon, Clock } from "lucide-react";
import { Card, CardHeader, CardTitle, CardBody } from "@/app/components/ui/Card";
import { Button } from "@/app/components/ui/Button";
import { MOCK_TODAY_SCHEDULE } from "@/app/lib/mock-data";

export default function CalendarPage() {
  const [currentDate, setCurrentDate] = useState(new Date(2025, 4, 20)); // May 20, 2025

  // Generate calendar days
  const daysInMonth = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 0).getDate();
  const startDayOfWeek = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1).getDay();

  const monthYearStr = currentDate.toLocaleDateString("en-US", {
    month: "long",
    year: "numeric",
  });

  const nextMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 1));
  };

  const prevMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() - 1, 1));
  };

  return (
    <div className="flex flex-col gap-6 max-w-7xl mx-auto">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-navy-900">Calendar</h1>
          <p className="text-sm text-navy-500 mt-1">
            Manage scheduled appointments and veterinarian availability.
          </p>
        </div>
        <Button icon={<Plus size={16} />} className="self-start sm:self-auto">
          Book Appointment
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6 items-start">
        {/* Left column - Calendar grid (lg:col-span-3) */}
        <div className="lg:col-span-3 flex flex-col gap-4">
          <Card className="p-0 overflow-hidden">
            {/* Calendar Controls */}
            <div className="flex items-center justify-between border-b border-navy-200 px-6 py-4 bg-white">
              <h2 className="text-base font-bold text-navy-900 flex items-center gap-2">
                <CalendarIcon size={18} className="text-teal-600" />
                {monthYearStr}
              </h2>
              <div className="flex items-center gap-1.5">
                <button
                  onClick={prevMonth}
                  className="p-1.5 hover:bg-navy-50 rounded-lg border border-navy-200 transition-colors cursor-pointer text-navy-600"
                >
                  <ChevronLeft size={16} />
                </button>
                <button
                  onClick={() => setCurrentDate(new Date(2025, 4, 20))}
                  className="px-3 py-1 text-xs font-semibold hover:bg-navy-50 rounded-lg border border-navy-200 transition-colors cursor-pointer text-navy-600"
                >
                  Today
                </button>
                <button
                  onClick={nextMonth}
                  className="p-1.5 hover:bg-navy-50 rounded-lg border border-navy-200 transition-colors cursor-pointer text-navy-600"
                >
                  <ChevronRight size={16} />
                </button>
              </div>
            </div>

            {/* Days of Week Headers */}
            <div className="grid grid-cols-7 border-b border-navy-200 text-center bg-navy-50/50 py-2.5">
              {["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"].map((day) => (
                <span key={day} className="text-xs font-bold text-navy-500 uppercase tracking-wider">
                  {day}
                </span>
              ))}
            </div>

            {/* Calendar Days Grid */}
            <div className="grid grid-cols-7 bg-navy-100 gap-px">
              {/* Empty starting cells */}
              {Array.from({ length: startDayOfWeek }).map((_, i) => (
                <div key={`empty-${i}`} className="bg-navy-50/20 min-h-[100px] p-2" />
              ))}

              {/* Month cells */}
              {Array.from({ length: daysInMonth }).map((_, i) => {
                const dayNum = i + 1;
                const isSelected = dayNum === 20 && currentDate.getMonth() === 4; // May 20, 2025
                return (
                  <div
                    key={dayNum}
                    className={[
                      "bg-white min-h-[100px] p-2 flex flex-col gap-1 transition-colors hover:bg-navy-50/30",
                      isSelected ? "ring-2 ring-inset ring-teal-500 bg-teal-50/10" : "",
                    ].join(" ")}
                  >
                    <span
                      className={[
                        "text-xs font-semibold w-6 h-6 flex items-center justify-center rounded-full",
                        isSelected
                          ? "bg-teal-600 text-white font-bold"
                          : "text-navy-700",
                      ].join(" ")}
                    >
                      {dayNum}
                    </span>

                    {/* Mock events on selected day */}
                    {isSelected && (
                      <div className="flex flex-col gap-1 overflow-hidden mt-1">
                        <span className="text-[10px] font-bold bg-teal-50 text-teal-700 border border-teal-200 px-1 rounded truncate leading-tight">
                          10:00 AM — Bella (Vacc)
                        </span>
                        <span className="text-[10px] font-bold bg-blue-50 text-blue-700 border border-blue-200 px-1 rounded truncate leading-tight">
                          11:00 AM — Charlie (Vacc)
                        </span>
                        <span className="text-[10px] font-bold bg-purple-50 text-purple-700 border border-purple-200 px-1 rounded truncate leading-tight">
                          11:30 AM — Luna (Follow)
                        </span>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          </Card>
        </div>

        {/* Right column - Event details/summary list (lg:col-span-1) */}
        <Card className="h-full">
          <CardHeader>
            <CardTitle subtitle="Visits scheduled for selected day">
              May 20, 2025
            </CardTitle>
          </CardHeader>
          <CardBody>
            <div className="flex flex-col gap-4">
              {MOCK_TODAY_SCHEDULE.map((appt) => (
                <div
                  key={appt.id}
                  className="p-3 bg-navy-50 rounded-xl border border-navy-200 flex flex-col gap-2 hover:border-navy-300 transition-colors"
                >
                  <div className="flex justify-between items-center">
                    <span className="text-xs font-semibold text-teal-600 flex items-center gap-1">
                      <Clock size={12} />
                      {appt.time}
                    </span>
                    <span className="text-[10px] bg-navy-200 text-navy-700 font-medium px-2 py-0.5 rounded-full capitalize">
                      {appt.petSpecies}
                    </span>
                  </div>
                  <div>
                    <h4 className="text-sm font-bold text-navy-950">
                      {appt.petName}
                    </h4>
                    <p className="text-xs text-navy-500 mt-0.5">
                      Owner: {appt.ownerName}
                    </p>
                  </div>
                  <span className="text-xs text-navy-600 font-medium bg-white border border-navy-200 rounded px-2.5 py-1 self-start">
                    {appt.type}
                  </span>
                </div>
              ))}
            </div>
          </CardBody>
        </Card>
      </div>
    </div>
  );
}
