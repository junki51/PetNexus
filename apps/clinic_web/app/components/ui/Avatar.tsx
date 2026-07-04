import React from "react";

type AvatarSize = "xs" | "sm" | "md" | "lg" | "xl";

interface AvatarProps {
  src?: string;
  name?: string;
  size?: AvatarSize;
  className?: string;
  shape?: "circle" | "rounded";
}

const sizeClasses: Record<AvatarSize, string> = {
  xs: "w-6 h-6 text-xs",
  sm: "w-8 h-8 text-xs",
  md: "w-10 h-10 text-sm",
  lg: "w-12 h-12 text-base",
  xl: "w-16 h-16 text-lg",
};

// Stable color based on first character of name
const bgColors = [
  "bg-teal-100 text-teal-700",
  "bg-blue-100 text-blue-700",
  "bg-purple-100 text-purple-700",
  "bg-amber-100 text-amber-700",
  "bg-pink-100 text-pink-700",
  "bg-indigo-100 text-indigo-700",
  "bg-emerald-100 text-emerald-700",
  "bg-rose-100 text-rose-700",
];

function getInitials(name: string): string {
  return name
    .split(" ")
    .map((n) => n[0])
    .slice(0, 2)
    .join("")
    .toUpperCase();
}

function getColor(name: string): string {
  const index = name.charCodeAt(0) % bgColors.length;
  return bgColors[index];
}

export function Avatar({
  src,
  name = "",
  size = "md",
  className = "",
  shape = "circle",
}: AvatarProps) {
  const shapeClass = shape === "circle" ? "rounded-full" : "rounded-lg";
  const colorClass = getColor(name);

  if (src) {
    return (
      // eslint-disable-next-line @next/next/no-img-element
      <img
        src={src}
        alt={name}
        className={[
          "object-cover shrink-0",
          sizeClasses[size],
          shapeClass,
          className,
        ].join(" ")}
      />
    );
  }

  return (
    <span
      className={[
        "inline-flex items-center justify-center font-semibold shrink-0",
        sizeClasses[size],
        shapeClass,
        colorClass,
        className,
      ].join(" ")}
      aria-label={name}
    >
      {name ? getInitials(name) : "?"}
    </span>
  );
}
