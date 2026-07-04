import React from "react";

interface DividerProps {
  label?: string;
  className?: string;
}

export function Divider({ label, className = "" }: DividerProps) {
  if (!label) {
    return (
      <hr className={["border-0 border-t border-navy-200", className].join(" ")} />
    );
  }

  return (
    <div className={["flex items-center gap-3", className].join(" ")}>
      <span className="flex-1 border-t border-navy-200" />
      <span className="text-xs text-navy-400 font-medium whitespace-nowrap">
        {label}
      </span>
      <span className="flex-1 border-t border-navy-200" />
    </div>
  );
}
