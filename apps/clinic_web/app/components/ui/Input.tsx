"use client";

import React, { useState } from "react";
import { Eye, EyeOff } from "lucide-react";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
  prefixIcon?: React.ReactNode;
  suffixIcon?: React.ReactNode;
}

export function Input({
  label,
  error,
  helperText,
  prefixIcon,
  suffixIcon,
  type = "text",
  className = "",
  id,
  ...props
}: InputProps) {
  const [showPassword, setShowPassword] = useState(false);
  const inputId = id ?? label?.toLowerCase().replace(/\s+/g, "-");
  const isPassword = type === "password";
  const resolvedType = isPassword ? (showPassword ? "text" : "password") : type;

  return (
    <div className="flex flex-col gap-1.5">
      {label && (
        <label
          htmlFor={inputId}
          className="text-sm font-medium text-navy-700"
        >
          {label}
        </label>
      )}

      <div className="relative flex items-center">
        {prefixIcon && (
          <span className="absolute left-3 text-navy-400 pointer-events-none">
            {prefixIcon}
          </span>
        )}

        <input
          id={inputId}
          type={resolvedType}
          className={[
            "w-full h-10 rounded-lg border bg-white text-navy-900 text-sm",
            "placeholder:text-navy-400",
            "transition-colors duration-150",
            "focus:outline-none focus:ring-2 focus:ring-teal-500/30 focus:border-teal-500",
            error
              ? "border-red-400 focus:ring-red-300/30 focus:border-red-500"
              : "border-navy-200 hover:border-navy-300",
            prefixIcon ? "pl-10" : "pl-3",
            isPassword || suffixIcon ? "pr-10" : "pr-3",
            className,
          ]
            .filter(Boolean)
            .join(" ")}
          {...props}
        />

        {isPassword && (
          <button
            type="button"
            onClick={() => setShowPassword((v) => !v)}
            className="absolute right-3 text-navy-400 hover:text-navy-600 transition-colors"
            tabIndex={-1}
            aria-label={showPassword ? "Hide password" : "Show password"}
          >
            {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
          </button>
        )}

        {!isPassword && suffixIcon && (
          <span className="absolute right-3 text-navy-400 pointer-events-none">
            {suffixIcon}
          </span>
        )}
      </div>

      {error && (
        <p className="text-xs text-red-500">{error}</p>
      )}
      {!error && helperText && (
        <p className="text-xs text-navy-500">{helperText}</p>
      )}
    </div>
  );
}
