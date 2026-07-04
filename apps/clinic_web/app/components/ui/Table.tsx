import React from "react";

// ── Table Root ──────────────────────────────────────────────

interface TableProps {
  children: React.ReactNode;
  className?: string;
}

export function Table({ children, className = "" }: TableProps) {
  return (
    <div className="w-full overflow-x-auto rounded-xl border border-navy-200">
      <table className={["w-full text-sm border-collapse", className].join(" ")}>
        {children}
      </table>
    </div>
  );
}

// ── Table Head ──────────────────────────────────────────────

interface TableHeadProps {
  children: React.ReactNode;
  className?: string;
}

export function TableHead({ children, className = "" }: TableHeadProps) {
  return (
    <thead className={["bg-navy-50 border-b border-navy-200", className].join(" ")}>
      {children}
    </thead>
  );
}

// ── Table Body ──────────────────────────────────────────────

interface TableBodyProps {
  children: React.ReactNode;
  className?: string;
}

export function TableBody({ children, className = "" }: TableBodyProps) {
  return (
    <tbody className={["divide-y divide-navy-100 bg-white", className].join(" ")}>
      {children}
    </tbody>
  );
}

// ── Table Row ───────────────────────────────────────────────

interface TableRowProps {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
  hoverable?: boolean;
}

export function TableRow({
  children,
  className = "",
  onClick,
  hoverable = true,
}: TableRowProps) {
  return (
    <tr
      onClick={onClick}
      className={[
        "transition-colors duration-100",
        hoverable ? "hover:bg-navy-50/60" : "",
        onClick ? "cursor-pointer" : "",
        className,
      ]
        .filter(Boolean)
        .join(" ")}
    >
      {children}
    </tr>
  );
}

// ── Table Header Cell ───────────────────────────────────────

interface TableThProps {
  children: React.ReactNode;
  className?: string;
  align?: "left" | "center" | "right";
  width?: string;
}

export function TableTh({
  children,
  className = "",
  align = "left",
  width,
}: TableThProps) {
  const alignClass = {
    left: "text-left",
    center: "text-center",
    right: "text-right",
  }[align];

  return (
    <th
      style={width ? { width } : undefined}
      className={[
        "px-4 py-3 text-xs font-semibold text-navy-500 uppercase tracking-wide whitespace-nowrap",
        alignClass,
        className,
      ]
        .filter(Boolean)
        .join(" ")}
    >
      {children}
    </th>
  );
}

// ── Table Data Cell ─────────────────────────────────────────

interface TableTdProps {
  children: React.ReactNode;
  className?: string;
  align?: "left" | "center" | "right";
}

export function TableTd({
  children,
  className = "",
  align = "left",
}: TableTdProps) {
  const alignClass = {
    left: "text-left",
    center: "text-center",
    right: "text-right",
  }[align];

  return (
    <td
      className={[
        "px-4 py-3 text-navy-700 whitespace-nowrap",
        alignClass,
        className,
      ]
        .filter(Boolean)
        .join(" ")}
    >
      {children}
    </td>
  );
}
