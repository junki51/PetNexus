import React from "react";

type StatCardColor = "teal" | "blue" | "purple" | "green" | "amber" | "rose";

interface StatCardProps {
  label: string;
  value: number | string;
  subtitle?: string;
  icon: React.ReactNode;
  color?: StatCardColor;
  linkLabel?: string;
  onLinkClick?: () => void;
}

const colorClasses: Record<StatCardColor, { bg: string; text: string }> = {
  teal:   { bg: "bg-teal-50",   text: "text-teal-600" },
  blue:   { bg: "bg-blue-50",   text: "text-blue-600" },
  purple: { bg: "bg-purple-50", text: "text-purple-600" },
  green:  { bg: "bg-emerald-50", text: "text-emerald-600" },
  amber:  { bg: "bg-amber-50",  text: "text-amber-600" },
  rose:   { bg: "bg-rose-50",   text: "text-rose-600" },
};

export function StatCard({
  label,
  value,
  subtitle,
  icon,
  color = "teal",
  linkLabel,
  onLinkClick,
}: StatCardProps) {
  const { bg, text } = colorClasses[color];

  return (
    <div className="bg-white rounded-xl border border-navy-200 shadow-sm p-5 flex flex-col gap-3 hover:shadow-md transition-shadow duration-200">
      <div className="flex items-start justify-between">
        <span className="text-sm font-medium text-navy-500">{label}</span>
        <span className={["w-10 h-10 rounded-lg flex items-center justify-center shrink-0", bg, text].join(" ")}>
          {icon}
        </span>
      </div>

      <div>
        <span className="text-3xl font-bold text-navy-900 leading-none">{value}</span>
        {subtitle && (
          <p className="text-xs text-navy-500 mt-1">{subtitle}</p>
        )}
      </div>

      {linkLabel && (
        <button
          onClick={onLinkClick}
          className="text-xs text-teal-600 hover:text-teal-700 font-medium flex items-center gap-1 transition-colors mt-auto"
        >
          {linkLabel} →
        </button>
      )}
    </div>
  );
}
