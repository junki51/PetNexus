"use client";

import React, { useState, useMemo } from "react";
import { Plus, Eye, FileText, PawPrint } from "lucide-react";
import { Card } from "@/app/components/ui/Card";
import { Table, TableHead, TableBody, TableRow, TableTh, TableTd } from "@/app/components/ui/Table";
import { Select } from "@/app/components/ui/Select";
import { Button } from "@/app/components/ui/Button";
import { SearchInput } from "@/app/components/ui/SearchInput";
import { StatusBadge } from "@/app/components/ui/Badge";
import { Avatar } from "@/app/components/ui/Avatar";
import { Pagination } from "@/app/components/ui/Pagination";
import { MOCK_PATIENTS } from "@/app/lib/mock-data";
import type { SelectOption } from "@/app/types";
import Link from "next/link";

const SPECIES_OPTIONS: SelectOption[] = [
  { value: "all", label: "All Species" },
  { value: "dog", label: "Dogs" },
  { value: "cat", label: "Cats" },
];

const STATUS_OPTIONS: SelectOption[] = [
  { value: "all", label: "All Status" },
  { value: "approved", label: "Approved" },
  { value: "pending", label: "Pending" },
  { value: "rejected", label: "Rejected" },
  { value: "revoked", label: "Revoked" },
];

export default function PatientsPage() {
  const [search, setSearch] = useState("");
  const [species, setSpecies] = useState("all");
  const [status, setStatus] = useState("all");
  const [page, setPage] = useState(1);
  const limit = 6;

  // Filter logic
  const filteredPatients = useMemo(() => {
    return MOCK_PATIENTS.filter((patient) => {
      const matchSearch =
        patient.name.toLowerCase().includes(search.toLowerCase()) ||
        patient.ownerName.toLowerCase().includes(search.toLowerCase());
      const matchSpecies = species === "all" || patient.species === species;
      const matchStatus = status === "all" || patient.status === status;
      return matchSearch && matchSpecies && matchStatus;
    });
  }, [search, species, status]);

  // Paginated patients
  const paginatedPatients = useMemo(() => {
    const start = (page - 1) * limit;
    return filteredPatients.slice(start, start + limit);
  }, [filteredPatients, page]);

  const totalPages = Math.ceil(filteredPatients.length / limit) || 1;

  return (
    <div className="flex flex-col gap-6 max-w-7xl mx-auto">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-navy-900">Patients</h1>
          <p className="text-sm text-navy-500 mt-1">
            Manage and view all your clinic patients.
          </p>
        </div>
        <Button icon={<Plus size={16} />} className="self-start sm:self-auto">
          Add Patient
        </Button>
      </div>

      {/* Filters Bar Card */}
      <Card padding="sm" className="bg-white">
        <div className="flex flex-col sm:flex-row gap-4 items-center justify-between">
          <div className="w-full sm:max-w-md shrink-0">
            <SearchInput
              value={search}
              onChange={(val) => {
                setSearch(val);
                setPage(1);
              }}
              placeholder="Search by pet or owner name..."
            />
          </div>

          <div className="flex flex-col sm:flex-row gap-3 w-full sm:w-auto items-stretch sm:items-center">
            <Select
              options={SPECIES_OPTIONS}
              value={species}
              onChange={(val) => {
                setSpecies(val);
                setPage(1);
              }}
              placeholder="Species"
              className="w-full sm:w-40"
            />
            <Select
              options={STATUS_OPTIONS}
              value={status}
              onChange={(val) => {
                setStatus(val);
                setPage(1);
              }}
              placeholder="Status"
              className="w-full sm:w-40"
            />
          </div>
        </div>
      </Card>

      {/* Patient Table Card */}
      <Card padding="none" className="overflow-hidden">
        <Table>
          <TableHead>
            <TableRow hoverable={false}>
              <TableTh>Pet Name</TableTh>
              <TableTh>Species</TableTh>
              <TableTh>Owner</TableTh>
              <TableTh>Status</TableTh>
              <TableTh>Last Visit</TableTh>
              <TableTh align="center" width="80px">
                Actions
              </TableTh>
            </TableRow>
          </TableHead>
          <TableBody>
            {paginatedPatients.length > 0 ? (
              paginatedPatients.map((patient) => (
                <TableRow key={patient.id}>
                  {/* Pet Name (Avatar + text) */}
                  <TableTd>
                    <div className="flex items-center gap-3">
                      <Avatar name={patient.name} size="sm" />
                      <div>
                        <p className="text-sm font-semibold text-navy-900 leading-tight">
                          {patient.name}
                        </p>
                        <p className="text-[10px] text-navy-400 font-medium">
                          {patient.breed}
                        </p>
                      </div>
                    </div>
                  </TableTd>

                  {/* Species (Capitalized) */}
                  <TableTd className="capitalize">{patient.species}</TableTd>

                  {/* Owner */}
                  <TableTd>
                    <div>
                      <p className="text-sm font-medium text-navy-800 leading-tight">
                        {patient.ownerName}
                      </p>
                      <p className="text-[10px] text-navy-400 font-mono">
                        {patient.ownerPhone}
                      </p>
                    </div>
                  </TableTd>

                  {/* Status */}
                  <TableTd>
                    <StatusBadge status={patient.status} />
                  </TableTd>

                  {/* Last Visit */}
                  <TableTd className="text-navy-600 font-medium">
                    {patient.lastVisit || "No visits yet"}
                  </TableTd>

                  {/* Actions */}
                  <TableTd align="center">
                    <div className="flex items-center justify-center gap-1.5">
                      <Link
                        href={`/medical-records/new`}
                        title="New Medical Record"
                        className="p-1 text-navy-500 hover:text-teal-600 rounded hover:bg-navy-100 transition-colors"
                      >
                        <FileText size={16} />
                      </Link>
                      <button
                        title="View Profile"
                        className="p-1 text-navy-500 hover:text-navy-700 rounded hover:bg-navy-100 transition-colors cursor-pointer"
                      >
                        <Eye size={16} />
                      </button>
                    </div>
                  </TableTd>
                </TableRow>
              ))
            ) : (
              <TableRow hoverable={false}>
                <td
                  colSpan={6}
                  className="py-16 text-center text-navy-400 text-sm font-medium"
                >
                  <PawPrint size={40} className="text-navy-300 mx-auto mb-3" />
                  No patients match the filters.
                </td>
              </TableRow>
            )}
          </TableBody>
        </Table>

        {/* Table Footer with Pagination */}
        <div className="px-6 py-4 bg-navy-50 border-t border-navy-200">
          <Pagination
            page={page}
            totalPages={totalPages}
            onPageChange={setPage}
            total={filteredPatients.length}
            limit={limit}
          />
        </div>
      </Card>
    </div>
  );
}
