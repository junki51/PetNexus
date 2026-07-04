"use client";

import React, { createContext, useContext } from "react";

// ── Context ────────────────────────────────────────────────

interface TabsContextValue {
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

const TabsContext = createContext<TabsContextValue | null>(null);

function useTabsContext() {
  const ctx = useContext(TabsContext);
  if (!ctx) throw new Error("Tabs subcomponent used outside of <Tabs>");
  return ctx;
}

// ── Tabs Root ───────────────────────────────────────────────

interface TabsProps {
  defaultTab: string;
  children: React.ReactNode;
  className?: string;
}

export function Tabs({ defaultTab, children, className = "" }: TabsProps) {
  const [activeTab, setActiveTab] = React.useState(defaultTab);

  return (
    <TabsContext.Provider value={{ activeTab, setActiveTab }}>
      <div className={className}>{children}</div>
    </TabsContext.Provider>
  );
}

// ── Tab List ────────────────────────────────────────────────

interface TabListProps {
  children: React.ReactNode;
  className?: string;
}

export function TabList({ children, className = "" }: TabListProps) {
  return (
    <div
      role="tablist"
      className={[
        "flex gap-1 p-1 bg-navy-100 rounded-lg",
        className,
      ].join(" ")}
    >
      {children}
    </div>
  );
}

// ── Tab ─────────────────────────────────────────────────────

interface TabProps {
  value: string;
  children: React.ReactNode;
  className?: string;
}

export function Tab({ value, children, className = "" }: TabProps) {
  const { activeTab, setActiveTab } = useTabsContext();
  const isActive = activeTab === value;

  return (
    <button
      role="tab"
      aria-selected={isActive}
      onClick={() => setActiveTab(value)}
      className={[
        "flex-1 px-4 py-2 text-sm font-medium rounded-md transition-all duration-150",
        isActive
          ? "bg-white text-navy-900 shadow-sm"
          : "text-navy-500 hover:text-navy-700",
        className,
      ]
        .filter(Boolean)
        .join(" ")}
    >
      {children}
    </button>
  );
}

// ── Tab Panel ────────────────────────────────────────────────

interface TabPanelProps {
  value: string;
  children: React.ReactNode;
  className?: string;
}

export function TabPanel({ value, children, className = "" }: TabPanelProps) {
  const { activeTab } = useTabsContext();

  if (activeTab !== value) return null;

  return (
    <div
      role="tabpanel"
      className={["animate-[fade-in_0.15s_ease-out]", className].join(" ")}
    >
      {children}
    </div>
  );
}
