import React from "react";

type BadgeVariant = "success" | "warning" | "danger" | "info" | "neutral" | "primary";

interface BadgeProps {
  children: React.ReactNode;
  variant?: BadgeVariant;
  dot?: boolean;
  className?: string;
}

const variantClasses: Record<BadgeVariant, string> = {
  success: "bg-emerald-50 text-emerald-700 border border-emerald-200",
  warning: "bg-amber-50 text-amber-700 border border-amber-200",
  danger:  "bg-red-50 text-red-700 border border-red-200",
  info:    "bg-blue-50 text-blue-700 border border-blue-200",
  neutral: "bg-navy-100 text-navy-600 border border-navy-200",
  primary: "bg-teal-50 text-teal-700 border border-teal-200",
};

const dotClasses: Record<BadgeVariant, string> = {
  success: "bg-emerald-500",
  warning: "bg-amber-500",
  danger:  "bg-red-500",
  info:    "bg-blue-500",
  neutral: "bg-navy-400",
  primary: "bg-teal-500",
};

export function Badge({
  children,
  variant = "neutral",
  dot = false,
  className = "",
}: BadgeProps) {
  return (
    <span
      className={[
        "inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium",
        variantClasses[variant],
        className,
      ]
        .filter(Boolean)
        .join(" ")}
    >
      {dot && (
        <span className={["w-1.5 h-1.5 rounded-full shrink-0", dotClasses[variant]].join(" ")} />
      )}
      {children}
    </span>
  );
}

// Convenience helpers
export function StatusBadge({ status }: { status: string }) {
  const map: Record<string, { label: string; variant: BadgeVariant }> = {
    approved:    { label: "Approved", variant: "success" },
    pending:     { label: "Pending", variant: "warning" },
    rejected:    { label: "Rejected", variant: "danger" },
    revoked:     { label: "Revoked", variant: "neutral" },
    "checked-in":{ label: "Checked In", variant: "primary" },
    "in-progress":{ label: "In Progress", variant: "info" },
    done:        { label: "Done", variant: "success" },
    scheduled:   { label: "Scheduled", variant: "neutral" },
    cancelled:   { label: "Cancelled", variant: "danger" },
  };

  const config = map[status] ?? { label: status, variant: "neutral" as BadgeVariant };

  return (
    <Badge variant={config.variant} dot>
      {config.label}
    </Badge>
  );
}
