import React from "react";

interface CardProps {
  children: React.ReactNode;
  className?: string;
  hover?: boolean;
  padding?: "none" | "sm" | "md" | "lg";
}

interface CardHeaderProps {
  children: React.ReactNode;
  className?: string;
  action?: React.ReactNode;
}

interface CardTitleProps {
  children: React.ReactNode;
  className?: string;
  subtitle?: string;
}

interface CardBodyProps {
  children: React.ReactNode;
  className?: string;
}

interface CardFooterProps {
  children: React.ReactNode;
  className?: string;
}

const paddingClasses = {
  none: "",
  sm: "p-4",
  md: "p-5",
  lg: "p-6",
};

export function Card({
  children,
  className = "",
  hover = false,
  padding = "md",
}: CardProps) {
  return (
    <div
      className={[
        "bg-white rounded-xl border border-navy-200 shadow-sm",
        hover ? "hover:shadow-md hover:-translate-y-0.5 transition-all duration-200 cursor-pointer" : "",
        paddingClasses[padding],
        className,
      ]
        .filter(Boolean)
        .join(" ")}
    >
      {children}
    </div>
  );
}

export function CardHeader({ children, className = "", action }: CardHeaderProps) {
  return (
    <div className={["flex items-center justify-between mb-4", className].join(" ")}>
      <div className="flex-1">{children}</div>
      {action && <div className="ml-4 shrink-0">{action}</div>}
    </div>
  );
}

export function CardTitle({ children, className = "", subtitle }: CardTitleProps) {
  return (
    <div>
      <h3 className={["text-base font-semibold text-navy-900", className].join(" ")}>
        {children}
      </h3>
      {subtitle && (
        <p className="text-xs text-navy-500 mt-0.5">{subtitle}</p>
      )}
    </div>
  );
}

export function CardBody({ children, className = "" }: CardBodyProps) {
  return <div className={className}>{children}</div>;
}

export function CardFooter({ children, className = "" }: CardFooterProps) {
  return (
    <div className={["mt-4 pt-4 border-t border-navy-100", className].join(" ")}>
      {children}
    </div>
  );
}
