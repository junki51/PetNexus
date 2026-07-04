"use client";

import React from "react";
import {
  Users,
  FileClock,
  CalendarDays,
  UserCheck,
} from "lucide-react";
import { StatCard } from "@/app/components/ui/StatCard";
import { Card, CardHeader, CardTitle, CardBody } from "@/app/components/ui/Card";
import { Badge, StatusBadge } from "@/app/components/ui/Badge";
import { Avatar } from "@/app/components/ui/Avatar";
import { Button } from "@/app/components/ui/Button";
import {
  MOCK_DASHBOARD_STATS,
  MOCK_TODAY_SCHEDULE,
  MOCK_RECENT_ACTIVITY,
} from "@/app/lib/mock-data";
import Link from "next/link";

export default function DashboardPage() {
  const stats = MOCK_DASHBOARD_STATS;
  const today = new Date().toLocaleDateString("en-US", {
    month: "long",
    day: "numeric",
    year: "numeric",
  });

  return (
    <div className="flex flex-col gap-6 max-w-7xl mx-auto">
      {/* Welcome Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-navy-900">
            Good morning, Dr. Emily! 👋
          </h1>
          <p className="text-sm text-navy-500 mt-1">
            Here&apos;s what&apos;s happening at your clinic today.
          </p>
        </div>
        <div className="bg-white border border-navy-200 rounded-lg px-4 py-2 text-xs font-semibold text-navy-600 shadow-sm flex items-center gap-2 self-start sm:self-auto">
          <CalendarDays size={14} className="text-teal-600" />
          <span>{today}</span>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          label="Today's Patients"
          value={stats.todayPatients}
          subtitle={stats.todayPatientsChange}
          icon={<Users size={20} />}
          color="teal"
        />
        <StatCard
          label="Pending Requests"
          value={stats.pendingRequests}
          subtitle="Requires owner approval"
          icon={<FileClock size={20} />}
          color="blue"
          linkLabel="View requests"
          onLinkClick={() => {}}
        />
        <StatCard
          label="Upcoming Appointments"
          value={stats.upcomingAppointments}
          subtitle={stats.nextAppointmentTime}
          icon={<CalendarDays size={20} />}
          color="purple"
        />
        <StatCard
          label="Checked-in Now"
          value={stats.checkedInNow}
          subtitle="Waiting in lobby"
          icon={<UserCheck size={20} />}
          color="green"
          linkLabel="View patients"
          onLinkClick={() => {}}
        />
      </div>

      {/* Main Dashboard Content */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Today's Schedule */}
        <div className="lg:col-span-2 flex flex-col h-full">
          <Card className="flex-1">
            <CardHeader
              action={
                <Link href="/calendar">
                  <Button variant="ghost" size="sm" className="text-teal-600">
                    View full calendar
                  </Button>
                </Link>
              }
            >
              <CardTitle subtitle="Time-based list of scheduled visits for today">
                Today&apos;s Schedule
              </CardTitle>
            </CardHeader>
            <CardBody>
              <div className="flex flex-col divide-y divide-navy-100">
                {MOCK_TODAY_SCHEDULE.map((appt) => (
                  <div
                    key={appt.id}
                    className="py-3.5 flex items-center justify-between gap-4 first:pt-0 last:pb-0"
                  >
                    <div className="flex items-center gap-4">
                      {/* Time Slot */}
                      <span className="text-sm font-semibold text-navy-800 w-20 shrink-0">
                        {appt.time}
                      </span>
                      {/* Pet Avatar & Name */}
                      <Avatar name={appt.petName} size="sm" />
                      <div className="min-w-0">
                        <p className="text-sm font-semibold text-navy-900 flex items-center gap-2">
                          {appt.petName}
                          <span className="text-xs font-normal text-navy-400 capitalize">
                            ({appt.petSpecies})
                          </span>
                        </p>
                        <p className="text-xs text-navy-500 truncate">
                          Owner: {appt.ownerName}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center gap-3">
                      <Badge variant="info">{appt.type}</Badge>
                      <StatusBadge status={appt.status} />
                    </div>
                  </div>
                ))}
              </div>
            </CardBody>
          </Card>
        </div>

        {/* Recent Activity */}
        <Card className="h-full">
          <CardHeader
            action={
              <Button variant="ghost" size="sm" className="text-teal-600">
                View all activity
              </Button>
            }
          >
            <CardTitle subtitle="Real-time log of clinic operations">
              Recent Activity
            </CardTitle>
          </CardHeader>
          <CardBody>
            <div className="flex flex-col gap-4">
              {MOCK_RECENT_ACTIVITY.map((activity) => (
                <div key={activity.id} className="flex gap-3 items-start">
                  {/* Indicator Dot */}
                  <span className="w-1.5 h-1.5 rounded-full bg-teal-500 mt-2 shrink-0" />
                  <div className="min-w-0 flex-1">
                    <p className="text-xs text-navy-800 leading-relaxed">
                      {activity.description}
                    </p>
                    <span className="text-[10px] text-navy-400 font-medium block mt-0.5">
                      {activity.time}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </CardBody>
        </Card>
      </div>
    </div>
  );
}
