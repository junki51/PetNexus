"use client";

import React, { useState, useMemo } from "react";
import { Plus, FileText, PawPrint } from "lucide-react";
import { Card } from "@/app/components/ui/Card";
import { Table, TableHead, TableBody, TableRow, TableTh, TableTd } from "@/app/components/ui/Table";
import { Button } from "@/app/components/ui/Button";
import { SearchInput } from "@/app/components/ui/SearchInput";
import { Avatar } from "@/app/components/ui/Avatar";
import { Badge } from "@/app/components/ui/Badge";
import { Pagination } from "@/app/components/ui/Pagination";
import Link from "next/link";
import { useLanguage } from "@/app/components/LanguageContext";

export default function MedicalRecordsListPage() {
  const { t } = useLanguage();
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const limit = 5;

  // Derive mock medical records from patients for display
  const records = useMemo(() => {
    return [
      {
        id: "rec-001",
        date: "May 20, 2025",
        petName: "Bella",
        species: "dog",
        vetName: "Dr. Emily Carter",
        type: "Vaccination",
        complaint: "Routine 5-in-1 vaccine booster shot.",
      },
      {
        id: "rec-002",
        date: "May 23, 2025",
        petName: "Billie",
        species: "dog",
        vetName: "Dr. Emily Carter",
        type: "Consultation",
        complaint: "Mild skin itching and red spots on stomach area.",
      },
      {
        id: "rec-003",
        date: "May 24, 2025",
        petName: "Neo",
        species: "dog",
        vetName: "Dr. James Wilson",
        type: "Vaccination",
        complaint: "Rabies vaccine annual booster.",
      },
      {
        id: "rec-004",
        date: "May 19, 2025",
        petName: "Latte",
        species: "cat",
        vetName: "Dr. Emily Carter",
        type: "Follow-up",
        complaint: "Check post-operation healing status. Stitches look great.",
      },
      {
        id: "rec-005",
        date: "May 17, 2025",
        petName: "Luna",
        species: "cat",
        vetName: "Dr. James Wilson",
        type: "Emergency",
        complaint: "Accidental ingestion of chocolate wrapper. Induced vomiting.",
      },
    ];
  }, []);

  const filteredRecords = useMemo(() => {
    return records.filter((rec) =>
      rec.petName.toLowerCase().includes(search.toLowerCase()) ||
      rec.type.toLowerCase().includes(search.toLowerCase()) ||
      rec.complaint.toLowerCase().includes(search.toLowerCase())
    );
  }, [records, search]);

  const paginatedRecords = useMemo(() => {
    const start = (page - 1) * limit;
    return filteredRecords.slice(start, start + limit);
  }, [filteredRecords, page]);

  const totalPages = Math.ceil(filteredRecords.length / limit) || 1;

  return (
    <div className="flex flex-col gap-6 max-w-7xl mx-auto">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-navy-900">{t("medical_records")}</h1>
          <p className="text-sm text-navy-500 mt-1">
            Browse and view historical medical records for checked-in patients.
          </p>
        </div>
        <Link href="/medical-records/new" className="self-start sm:self-auto">
          <Button icon={<Plus size={16} />} className="cursor-pointer">{t("new_record_title")}</Button>
        </Link>
      </div>

      {/* Filter Row Card */}
      <Card padding="sm" className="bg-white">
        <div className="flex items-center">
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
        </div>
      </Card>

      {/* Records Table Card */}
      <Card padding="none" className="overflow-hidden bg-white">
        <Table>
          <TableHead>
            <TableRow hoverable={false}>
              <TableTh>{t("visit_date")}</TableTh>
              <TableTh>{t("pet_name_col")}</TableTh>
              <TableTh>{t("veterinarian")}</TableTh>
              <TableTh>{t("visit_type")}</TableTh>
              <TableTh>Symptoms / Diagnosis Excerpt</TableTh>
              <TableTh align="center" width="80px">
                View
              </TableTh>
            </TableRow>
          </TableHead>
          <TableBody>
            {paginatedRecords.length > 0 ? (
              paginatedRecords.map((rec) => (
                <TableRow key={rec.id}>
                  {/* Date */}
                  <TableTd className="font-semibold text-navy-800">{rec.date}</TableTd>

                  {/* Pet Name */}
                  <TableTd>
                    <div className="flex items-center gap-3">
                      <Avatar name={rec.petName} size="sm" />
                      <div>
                        <p className="text-sm font-semibold text-navy-900 leading-tight">
                          {rec.petName}
                        </p>
                        <p className="text-[10px] text-navy-400 font-medium capitalize">
                          {rec.species}
                        </p>
                      </div>
                    </div>
                  </TableTd>

                  {/* Veterinarian */}
                  <TableTd className="text-navy-700 font-medium">{rec.vetName}</TableTd>

                  {/* Visit Type */}
                  <TableTd>
                    <Badge variant="info">{rec.type}</Badge>
                  </TableTd>

                  {/* Excerpt */}
                  <TableTd className="text-navy-500 max-w-xs truncate font-medium">
                    {rec.complaint}
                  </TableTd>

                  {/* Action */}
                  <TableTd align="center">
                    <button
                      title="View Full Record"
                      className="p-1 text-navy-500 hover:text-teal-600 rounded hover:bg-navy-100 transition-colors cursor-pointer"
                    >
                      <FileText size={16} />
                    </button>
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
                  No medical records found matching search filters.
                </td>
              </TableRow>
            )}
          </TableBody>
        </Table>

        {/* Footer with Pagination */}
        <div className="px-6 py-4 bg-navy-50 border-t border-navy-200">
          <Pagination
            page={page}
            totalPages={totalPages}
            onPageChange={setPage}
            total={filteredRecords.length}
            limit={limit}
            showInfo={false}
          />
          {filteredRecords.length > 0 && (
            <p className="text-xs text-navy-400 mt-2 text-center sm:text-left">
              {t("showing_results", {
                from: (page - 1) * limit + 1,
                to: Math.min(page * limit, filteredRecords.length),
                total: filteredRecords.length,
              })}
            </p>
          )}
        </div>
      </Card>
    </div>
  );
}
