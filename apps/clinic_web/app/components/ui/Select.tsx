"use client";

import React from "react";
import { ChevronDown } from "lucide-react";
import type { SelectOption } from "@/app/types";

interface SelectProps {
  label?: string;
  value?: string;
  onChange?: (value: string) => void;
  options: SelectOption[];
  placeholder?: string;
  error?: string;
  disabled?: boolean;
  required?: boolean;
  className?: string;
  id?: string;
}

export function Select({
  label,
  value,
  onChange,
  options,
  placeholder = "Select...",
  error,
  disabled = false,
  required = false,
  className = "",
  id,
}: SelectProps) {
  const selectId = id ?? label?.toLowerCase().replace(/\s+/g, "-");

  return (
    <div className="flex flex-col gap-1.5">
      {label && (
        <label htmlFor={selectId} className="text-sm font-medium text-navy-700">
          {label}
        </label>
      )}

      <div className="relative">
        <select
          id={selectId}
          value={value}
          onChange={(e) => onChange?.(e.target.value)}
          disabled={disabled}
          required={required}
          className={[
            "w-full h-10 pl-3 pr-9 rounded-lg border bg-white text-sm appearance-none",
            "transition-colors duration-150 cursor-pointer",
            "focus:outline-none focus:ring-2 focus:ring-teal-500/30 focus:border-teal-500",
            error
              ? "border-red-400 text-navy-900"
              : "border-navy-200 hover:border-navy-300 text-navy-700",
            disabled ? "opacity-50 cursor-not-allowed bg-navy-50" : "",
            className,
          ]
            .filter(Boolean)
            .join(" ")}
        >
          {placeholder && (
            <option value="" disabled={!value}>
              {placeholder}
            </option>
          )}
          {options.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>

        <ChevronDown
          size={16}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-navy-400 pointer-events-none"
        />
      </div>

      {error && <p className="text-xs text-red-500">{error}</p>}
    </div>
  );
}
