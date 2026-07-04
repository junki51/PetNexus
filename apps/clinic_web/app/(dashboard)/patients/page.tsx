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
import { useLanguage } from "@/app/components/LanguageContext";

export default function PatientsPage() {
  const { t } = useLanguage();
  const [search, setSearch] = useState("");
  const [species, setSpecies] = useState("all");
  const [status, setStatus] = useState("all");
  const [page, setPage] = useState(1);
  const limit = 6;

  const SPECIES_OPTIONS: SelectOption[] = useMemo(() => [
    { value: "all", label: t("all_species") },
    { value: "dog", label: "Dogs" },
    { value: "cat", label: "Cats" },
  ], [t]);

  const STATUS_OPTIONS: SelectOption[] = useMemo(() => [
    { value: "all", label: t("all_status") },
    { value: "approved", label: "Approved" },
    { value: "pending", label: "Pending" },
    { value: "rejected", label: "Rejected" },
    { value: "revoked", label: "Revoked" },
  ], [t]);

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
          <h1 className="text-2xl font-bold text-navy-900">{t("patients_title")}</h1>
          <p className="text-sm text-navy-500 mt-1">
            {t("patients_desc")}
          </p>
        </div>
        <Button icon={<Plus size={16} />} className="self-start sm:self-auto cursor-pointer">
          {t("add_patient")}
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
              placeholder={t("search_placeholder")}
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
      <Card padding="none" className="overflow-hidden bg-white">
        <Table>
          <TableHead>
            <TableRow hoverable={false}>
              <TableTh>{t("pet_name_col")}</TableTh>
              <TableTh>{t("species_col")}</TableTh>
              <TableTh>{t("owner_col")}</TableTh>
              <TableTh>{t("status_col")}</TableTh>
              <TableTh>{t("last_visit")}</TableTh>
              <TableTh align="center" width="80px">
                {t("actions_col")}
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
                  {t("no_patients_found")}
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
            showInfo={false}
          />
          {filteredPatients.length > 0 && (
            <p className="text-xs text-navy-400 mt-2 text-center sm:text-left">
              {t("showing_results", {
                from: (page - 1) * limit + 1,
                to: Math.min(page * limit, filteredPatients.length),
                total: filteredPatients.length,
              })}
            </p>
          )}
        </div>
      </Card>
    </div>
  );
}
