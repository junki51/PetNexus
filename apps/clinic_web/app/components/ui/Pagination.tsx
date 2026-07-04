"use client";

import React from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface PaginationProps {
  page: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  showInfo?: boolean;
  total?: number;
  limit?: number;
}

export function Pagination({
  page,
  totalPages,
  onPageChange,
  showInfo = true,
  total,
  limit,
}: PaginationProps) {
  const from = total && limit ? (page - 1) * limit + 1 : undefined;
  const to = total && limit ? Math.min(page * limit, total) : undefined;

  function getPageNumbers(): (number | "...")[] {
    if (totalPages <= 7) {
      return Array.from({ length: totalPages }, (_, i) => i + 1);
    }
    const pages: (number | "...")[] = [1];
    if (page > 3) pages.push("...");
    const start = Math.max(2, page - 1);
    const end = Math.min(totalPages - 1, page + 1);
    for (let i = start; i <= end; i++) pages.push(i);
    if (page < totalPages - 2) pages.push("...");
    pages.push(totalPages);
    return pages;
  }

  return (
    <div className="flex items-center justify-between gap-4 flex-wrap">
      {showInfo && total !== undefined && from !== undefined && to !== undefined ? (
        <p className="text-sm text-navy-500">
          Showing <span className="font-medium text-navy-700">{from}–{to}</span> of{" "}
          <span className="font-medium text-navy-700">{total}</span> results
        </p>
      ) : showInfo ? (
        <p className="text-sm text-navy-500">
          Page <span className="font-medium text-navy-700">{page}</span> of{" "}
          <span className="font-medium text-navy-700">{totalPages}</span>
        </p>
      ) : (
        <div />
      )}

      <div className="flex items-center gap-1">
        <button
          onClick={() => onPageChange(page - 1)}
          disabled={page <= 1}
          className="flex items-center justify-center w-8 h-8 rounded-lg border border-navy-200 text-navy-500 hover:bg-navy-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          aria-label="Previous page"
        >
          <ChevronLeft size={16} />
        </button>

        {getPageNumbers().map((p, i) =>
          p === "..." ? (
            <span key={`ellipsis-${i}`} className="w-8 h-8 flex items-center justify-center text-navy-400 text-sm">
              …
            </span>
          ) : (
            <button
              key={p}
              onClick={() => onPageChange(p as number)}
              className={[
                "w-8 h-8 rounded-lg text-sm font-medium transition-colors",
                page === p
                  ? "bg-teal-600 text-white border border-teal-600"
                  : "border border-navy-200 text-navy-600 hover:bg-navy-50",
              ].join(" ")}
              aria-current={page === p ? "page" : undefined}
            >
              {p}
            </button>
          )
        )}

        <button
          onClick={() => onPageChange(page + 1)}
          disabled={page >= totalPages}
          className="flex items-center justify-center w-8 h-8 rounded-lg border border-navy-200 text-navy-500 hover:bg-navy-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          aria-label="Next page"
        >
          <ChevronRight size={16} />
        </button>
      </div>
    </div>
  );
}
