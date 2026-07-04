"use client";

import React from "react";
import { Search, X } from "lucide-react";

interface SearchInputProps {
  value?: string;
  onChange?: (value: string) => void;
  placeholder?: string;
  className?: string;
  id?: string;
}

export function SearchInput({
  value = "",
  onChange,
  placeholder = "Search...",
  className = "",
  id = "search",
}: SearchInputProps) {
  return (
    <div className={["relative flex items-center", className].join(" ")}>
      <Search
        size={16}
        className="absolute left-3 text-navy-400 pointer-events-none shrink-0"
      />
      <input
        id={id}
        type="search"
        value={value}
        onChange={(e) => onChange?.(e.target.value)}
        placeholder={placeholder}
        className={[
          "h-9 pl-9 pr-8 rounded-lg border border-navy-200 bg-white text-sm text-navy-800",
          "placeholder:text-navy-400",
          "hover:border-navy-300 focus:outline-none focus:ring-2 focus:ring-teal-500/30 focus:border-teal-500",
          "transition-colors duration-150 w-full",
          "[&::-webkit-search-cancel-button]:hidden",
        ].join(" ")}
      />
      {value && (
        <button
          onClick={() => onChange?.("")}
          className="absolute right-2.5 text-navy-400 hover:text-navy-600 transition-colors"
          aria-label="Clear search"
        >
          <X size={14} />
        </button>
      )}
    </div>
  );
}
