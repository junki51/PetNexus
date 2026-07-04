"use client";

import React from "react";

interface CheckboxProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, "type"> {
  label?: React.ReactNode;
  description?: string;
  id: string;
}

export function Checkbox({ label, description, id, className = "", ...props }: CheckboxProps) {
  return (
    <div className="flex items-start gap-2.5">
      <input
        id={id}
        type="checkbox"
        className={[
          "mt-0.5 w-4 h-4 rounded border-navy-300 text-teal-600 bg-white shrink-0",
          "checked:bg-teal-600 checked:border-teal-600",
          "focus:ring-2 focus:ring-teal-500/30 focus:ring-offset-0",
          "transition-colors duration-150 cursor-pointer",
          className,
        ]
          .filter(Boolean)
          .join(" ")}
        {...props}
      />
      {(label ?? description) && (
        <div className="flex flex-col">
          {label && (
            <label
              htmlFor={id}
              className="text-sm text-navy-700 cursor-pointer leading-5"
            >
              {label}
            </label>
          )}
          {description && (
            <span className="text-xs text-navy-500">{description}</span>
          )}
        </div>
      )}
    </div>
  );
}
