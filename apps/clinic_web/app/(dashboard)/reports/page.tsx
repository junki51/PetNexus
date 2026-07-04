"use client";

import React from "react";
import { TrendingUp, Users, Calendar, DollarSign } from "lucide-react";
import { Card, CardHeader, CardTitle, CardBody } from "@/app/components/ui/Card";
import { Select } from "@/app/components/ui/Select";

export default function ReportsPage() {
  return (
    <div className="flex flex-col gap-6 max-w-7xl mx-auto">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-navy-900">Reports & Analytics</h1>
          <p className="text-sm text-navy-500 mt-1">
            Clinic analytics, patient trends, and operational insights.
          </p>
        </div>
        <Select
          options={[
            { value: "7d", label: "Last 7 Days" },
            { value: "30d", label: "Last 30 Days" },
            { value: "ytd", label: "Year to Date" },
          ]}
          value="30d"
          className="w-full sm:w-40 self-start sm:self-auto"
        />
      </div>

      {/* Grid of highlight stat small blocks */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="bg-white border border-navy-200 rounded-xl p-5 shadow-sm flex items-center gap-4">
          <div className="w-12 h-12 rounded-lg bg-teal-50 text-teal-600 flex items-center justify-center shrink-0">
            <Users size={22} />
          </div>
          <div>
            <p className="text-xs font-semibold text-navy-400 uppercase tracking-wide">
              Total Active Patients
            </p>
            <h3 className="text-2xl font-bold text-navy-900 mt-1">1,248</h3>
            <span className="text-[10px] text-emerald-600 font-semibold flex items-center gap-0.5 mt-0.5">
              <TrendingUp size={10} /> +4.2% vs last month
            </span>
          </div>
        </div>

        <div className="bg-white border border-navy-200 rounded-xl p-5 shadow-sm flex items-center gap-4">
          <div className="w-12 h-12 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center shrink-0">
            <Calendar size={22} />
          </div>
          <div>
            <p className="text-xs font-semibold text-navy-400 uppercase tracking-wide">
              Visits This Month
            </p>
            <h3 className="text-2xl font-bold text-navy-900 mt-1">342</h3>
            <span className="text-[10px] text-emerald-600 font-semibold flex items-center gap-0.5 mt-0.5">
              <TrendingUp size={10} /> +12.5% vs last month
            </span>
          </div>
        </div>

        <div className="bg-white border border-navy-200 rounded-xl p-5 shadow-sm flex items-center gap-4">
          <div className="w-12 h-12 rounded-lg bg-purple-50 text-purple-600 flex items-center justify-center shrink-0">
            <DollarSign size={22} />
          </div>
          <div>
            <p className="text-xs font-semibold text-navy-400 uppercase tracking-wide">
              Est. Monthly Revenue
            </p>
            <h3 className="text-2xl font-bold text-navy-900 mt-1">$14,850</h3>
            <span className="text-[10px] text-emerald-600 font-semibold flex items-center gap-0.5 mt-0.5">
              <TrendingUp size={10} /> +8.1% vs last month
            </span>
          </div>
        </div>
      </div>

      {/* Analytics Charts Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Weekly Volume Graph Mockup */}
        <Card>
          <CardHeader>
            <CardTitle subtitle="Daily count of patients checked in">
              Weekly Patient Volume
            </CardTitle>
          </CardHeader>
          <CardBody>
            <div className="h-64 flex items-end gap-3 pt-6 px-4">
              {[
                { day: "Mon", count: 12, height: "h-2/5" },
                { day: "Tue", count: 18, height: "h-3/5" },
                { day: "Wed", count: 24, height: "h-4/5" },
                { day: "Thu", count: 15, height: "h-2.5/5" },
                { day: "Fri", count: 30, height: "h-5/5" },
                { day: "Sat", count: 8, height: "h-1.5/5" },
                { day: "Sun", count: 0, height: "h-0" },
              ].map((bar, idx) => (
                <div key={idx} className="flex-1 flex flex-col items-center gap-2 group">
                  <div className="w-full relative flex items-end h-full">
                    {/* Tooltip on hover */}
                    <span className="absolute bottom-[calc(100%+4px)] left-1/2 -translate-x-1/2 bg-navy-800 text-white text-[10px] font-bold px-1.5 py-0.5 rounded shadow opacity-0 group-hover:opacity-100 transition-opacity">
                      {bar.count}
                    </span>
                    <div
                      className={[
                        "w-full bg-teal-500 hover:bg-teal-600 rounded-t transition-all duration-300",
                        bar.height,
                      ].join(" ")}
                    />
                  </div>
                  <span className="text-[10px] font-bold text-navy-500 uppercase">
                    {bar.day}
                  </span>
                </div>
              ))}
            </div>
          </CardBody>
        </Card>

        {/* Species Distribution Chart Mockup */}
        <Card>
          <CardHeader>
            <CardTitle subtitle="Proportion of clinic visit sessions by species">
              Species Distribution
            </CardTitle>
          </CardHeader>
          <CardBody className="flex flex-col sm:flex-row items-center justify-around gap-6 py-6">
            {/* Pie Chart Representation */}
            <div className="relative w-40 h-40 rounded-full border-8 border-teal-500 flex items-center justify-center shrink-0">
              {/* Inner semi-circle overlay logic for visual presentation */}
              <div className="absolute inset-0 rounded-full border-8 border-navy-800 border-t-transparent border-r-transparent border-b-transparent transform rotate-45" />
              <div className="flex flex-col items-center">
                <span className="text-xl font-extrabold text-navy-900">1,248</span>
                <span className="text-[10px] font-bold text-navy-400 uppercase">
                  Patients
                </span>
              </div>
            </div>

            <div className="flex flex-col gap-3">
              <div className="flex items-center gap-3">
                <span className="w-3.5 h-3.5 rounded bg-teal-500 shrink-0" />
                <div>
                  <p className="text-xs font-bold text-navy-800">Dogs (64%)</p>
                  <p className="text-[10px] text-navy-500">798 active patients</p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <span className="w-3.5 h-3.5 rounded bg-navy-800 shrink-0" />
                <div>
                  <p className="text-xs font-bold text-navy-800">Cats (36%)</p>
                  <p className="text-[10px] text-navy-500">450 active patients</p>
                </div>
              </div>
            </div>
          </CardBody>
        </Card>
      </div>
    </div>
  );
}
